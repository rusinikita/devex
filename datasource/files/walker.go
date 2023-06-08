package files

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Package string
	Name    string
	Lines   uint32
	Symbols uint32
	Tags    map[string]uint32
	Imports []string
}

var Tags = []string{"todo", "fix", "note", "nolint", "billing", "money", "order", "pylint: disable"}

func Extract(_ context.Context, rootPath string, c chan<- File) error {
	defer close(c)

	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fPackage := strings.TrimPrefix(filepath.ToSlash(filepath.Dir(path)), filepath.ToSlash(rootPath))
		fPackage = strings.TrimPrefix(fPackage, "/")

		fName := filepath.Base(path)

		if strings.HasPrefix(fPackage, ".") || strings.HasPrefix(fName, ".") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		lines := strings.Split(content, "\n")

		f := File{
			Package: fPackage,
			Name:    fName,
			Lines:   uint32(len(lines)),
			Symbols: uint32(len(content)),
			Imports: extractImports(lines),
			Tags:    extractTags(content),
		}

		c <- f

		return err
	})
}

func extractTags(content string) (tags map[string]uint32) {
	tags = map[string]uint32{}
	content = strings.ToLower(content)

	for _, tag := range Tags {
		count := strings.Count(content, tag)
		if count == 0 {
			continue
		}

		tags[tag] = uint32(count)
	}

	return tags
}

func extractImports(contentLines []string) (paths []string) {
	inImportClosure := false
	for _, line := range contentLines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "import (") {
			inImportClosure = true
			continue
		}

		if inImportClosure && line == ")" {
			inImportClosure = false
			continue
		}

		if inImportClosure || strings.HasPrefix(line, "import") || strings.HasPrefix(line, "from") {
			i := strings.TrimPrefix(line, "import")

			split := strings.Split(i, " ")
			i = split[0]
			if len(split) > 1 {
				i = split[1]
			}

			i = strings.TrimSpace(i)
			i = strings.Trim(i, "\"")

			paths = append(paths, i)
		}
	}

	return paths
}
