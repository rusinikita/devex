package dashboard

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/rusinikita/devex/dao"
	"github.com/rusinikita/devex/project"
)

const name = ", name"

func GitChangesTop(db *gorm.DB, filesMode bool, projects []project.ID, filesFilter string) (result Values, err error) {
	grouping := "alias, package"
	filter := filesFilter

	if filesMode {
		grouping += name
		filter = "and f.present > 0\n" + filter
	}

	sqlBars := getSQLBars()
	sqlBars = fmt.Sprintf(sqlBars, grouping, filter)

	err = db.Raw(sqlBars, projects).Scan(&result).Error

	return result, err
}

func getSQLBars() string {
	return `
	with fcm as (select %[1]s, date("time", 'start of month') as month, sum(rows_added + rows_removed) as line_changes
		from git_changes as ch
		join files f on ch.file = f.id
		join projects p on f.project = p.id
		where f.project in ?
		   %[2]s
		   and time > date('now', '-48 month')
		group by %[1]s, date("time", 'start of month'))
	select %[1]s, count(*), sum(line_changes), avg(line_changes) as value
	from fcm group by %[1]s
	having count(*) > 3
	order by avg(line_changes) desc
	limit 100
`
}

func GitChangesData(db *gorm.DB, filesMode bool, projects []project.ID, bars Values) (result Values, err error) {
	// Future: months/weeks selector

	barStrings := bars.BarNames()

	grouping := "alias, package"
	barFilter := "alias || '/' || package"

	if filesMode {
		grouping += name
		barFilter += " || '/' || name"
	}

	sql := getSQLChangesData()
	sql = fmt.Sprintf(sql, grouping, barFilter)

	err = db.Raw(sql, projects, barStrings).Scan(&result).Error

	return result, err
}

func getSQLChangesData() string {
	return `
	select %[1]s, date("time", 'start of month') as 'time', sum(rows_added + rows_removed) as value
	from git_changes as ch
	join files f on ch.file = f.id
	join projects p on f.project = p.id
	where f.project in ?
		and %[2]s in ?
		and time > date('now', '-24 month')
	group by %[1]s, date("time", 'start of month')
`
}

func FileSizes(db *gorm.DB, projects []project.ID, filesFilter string) (result Values, err error) {
	err = db.Model(dao.File{}).
		Select("alias", "package", "name", "lines as value").
		Joins("join projects p on p.id = files.project").
		Where("present > 0 and project in ?"+filesFilter, projects).
		Scan(&result).
		Error

	return result, err
}

func Contribution(db *gorm.DB, filesMode bool, projects []project.ID, filesFilter string) (result Values, err error) {
	grouping := "package"
	if filesMode {
		grouping += name
	}

	err = db.Model(dao.GitChange{}).
		Select("alias", grouping, "author", "sum(rows_added+rows_removed) as value").
		Joins("join git_commits c on c.id = git_changes.'commit'").
		Joins("join files f on f.id = git_changes.file").
		Joins("join projects p on p.id = f.project").
		Where("git_changes.time > date('now', '-12 month') and f.project in ?"+filesFilter, projects).
		Group("alias, author, " + grouping).
		Having("sum(rows_added+rows_removed) > 300").
		Scan(&result).
		Error

	return result, err
}

// CommitMessages TODO contribution pace. velocity per month
func CommitMessages(db *gorm.DB, filesMode bool, projects []project.ID, filter string) (result Values, err error) {
	grouping := "package"
	if filesMode {
		grouping += name
	}

	err = db.Model(dao.GitChange{}).
		Select("alias", grouping, "count(*) as value").
		Joins("join git_commits c on c.id = git_changes.'commit'").
		Joins("join files f on f.id = git_changes.file").
		Joins("join projects p on p.id = f.project").
		Where("git_changes.time > date('now', '-24 month') and f.present > 0 and f.project in ?"+filter, projects).
		Group("alias, " + grouping).
		Having("count(*) > 0").
		Order("count(*) desc").
		Limit(40).
		Scan(&result).
		Error

	return result, err
}

func FileTags(db *gorm.DB, projects []project.ID, filesFilter, tagsFilter string) (result Values, err error) {
	err = db.Model(dao.File{}).
		Select("alias", "package", "name", "tags").
		Joins("join projects p on p.id = files.project").
		Where("present > 0 and project in ?"+filesFilter+tagsFilter, projects).
		Scan(&result).Error

	return result, err
}

func Imports(db *gorm.DB, filesMode bool, projects []project.ID, filesFilter string) (result AllImports, err error) {
	grouping := "package"
	if filesMode {
		grouping += name
	}

	err = db.Model(dao.File{}).
		Select("alias", grouping, "lines", "imports").
		Joins("join projects p on p.id = files.project").
		Find(&result, "present > 0 and project in ?"+filesFilter, projects).
		Error

	return result, err
}
