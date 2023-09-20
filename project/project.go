package project

type ID uint64

// TODO packages prioritizing based on business value
// type Package struct {
// 	ID        ID
// 	ProjectID ID
// 	Path      string
// 	Priority  Priority
// 	Present   bool
// }
//
// type Priority string
//
// const (
// 	Regular    Priority = ""
// 	Vital               = "vital"
// 	Money               = "money"
// 	Critical            = "critical"
// 	Deprecated          = "deprecated"
// )

// TODO next
// Code version.
// When loading new data from a more recent commit,
// old data is marked as irrelevant
// type Revision struct {
// 	Hash string
// }

// TODO future UI
// DataFetchJob contains project data collection job state
// type DataFetchJob struct {
// 	ProjectID  ID
// 	DataSource string
// 	CreatedAt  time.Time
// 	FinishedAt *time.Time
// }
