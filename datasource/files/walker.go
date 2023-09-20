package files

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Tags    map[string]uint32
	Package string
	Name    string
	Imports []string
	Lines   uint32
	Symbols uint32
}

// fileDto TODO please give me normal name =)
type fileDto struct {
	Package string
	Name    string
	Path    string
}

var Tags = []string{
	"todo", "fix", "note", "nolint", "pylint: disable",
	"billing", "money", "order",
}

func Extract(_ context.Context, rootPath string, c chan<- File) error {
	defer close(c)

	pathes, err := getPathes(rootPath)

	if err != nil {
		return err
	}

	return writeFileToChannel(c, pathes)
}

func writeFileToChannel(c chan<- File, pathes []fileDto) error {
	for _, dto := range pathes {
		data, err := os.ReadFile(dto.Path)
		if err != nil {
			return err
		}

		// skip not code files
		if !strings.HasPrefix(http.DetectContentType(data), "text/") {
			continue
		}

		f := getFile(dto, data)

		c <- f
	}

	return nil
}

func getFile(dto fileDto, data []byte) File {
	content := string(data)
	lines := strings.Split(content, "\n")

	return File{
		Package: dto.Package,
		Name:    dto.Name,
		Lines:   uint32(len(lines)),
		Symbols: uint32(len(content)),
		Imports: extractImports(lines),
		Tags:    extractTags(content),
	}
}

func getPathes(rootPath string) ([]fileDto, error) {
	var pathes []fileDto
	err := filepath.WalkDir(
		rootPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			dto := getFileDto(rootPath, path)

			if strings.HasPrefix(dto.Package, ".") || strings.HasPrefix(dto.Name, ".") {
				return nil
			}

			pathes = append(pathes, dto)

			return nil
		})

	return pathes, err
}

func getFileDto(rootPath string, path string) fileDto {
	pathWithSlash := filepath.ToSlash(filepath.Dir(path))
	prefix := filepath.ToSlash(rootPath)
	fPackage := strings.TrimPrefix(pathWithSlash, prefix)
	fPackage = strings.TrimPrefix(fPackage, "/")

	fName := filepath.Base(path)

	return fileDto{
		Package: fPackage,
		Path:    path,
		Name:    fName,
	}
}

func extractTags(content string) (tags map[string]uint32) {
	tags = map[string]uint32{}
	lowerContent := strings.ToLower(content)

	for _, tag := range Tags {
		count := strings.Count(lowerContent, tag)
		if count == 0 {
			continue
		}

		tags[tag] = uint32(count)
	}

	return tags
}

func extractImports(contentLines []string) (paths []string) {
	lines := filterLines(contentLines)

	return appendLinesToPaths(lines, paths)
}

func appendLinesToPaths(lines []string, paths []string) []string {
	results := paths

	for _, line := range lines {
		result := getResultForLine(line)

		results = append(results, result)
	}

	return results
}

func getResultForLine(line string) string {
	result := strings.TrimPrefix(line, "import")

	split := strings.Split(result, " ")
	result = split[0]
	if len(split) > 1 {
		result = split[1]
	}

	result = strings.TrimSpace(result)
	return strings.Trim(result, "\"")
}

func filterLines(contentLines []string) []string {
	var lines []string
	inImportClosure := false
	skipLine := false

	for _, line := range contentLines {
		lines, skipLine, inImportClosure = appendLine(lines, line, skipLine, inImportClosure)
	}

	return lines
}

func appendLine(lines []string, line string, skipLine bool, inImportClosure bool) ([]string, bool, bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return lines, skipLine, inImportClosure
	}

	skipLine, inImportClosure = isSkipLine(line, inImportClosure)

	if skipLine {
		return lines, skipLine, inImportClosure
	}

	importLine := clearImport(line)
	return append(lines, importLine), skipLine, inImportClosure
}

func isSkipLine(
	line string,
	inClosure bool,
) (skipLine bool, inImportClosure bool) {
	if strings.HasPrefix(line, "import (") {
		return true, true
	}

	if inClosure && line == ")" {
		return true, false
	}

	if !isImport(line, inClosure) {
		return true, inClosure
	}

	return false, inClosure
}

func clearImport(line string) string {
	// It looks somewhat redundant, because... keep cleaning it anyway
	i := strings.TrimPrefix(line, "import")

	split := strings.Split(i, " ")
	i = split[0]
	if len(split) > 1 {
		i = split[1]
	}

	i = strings.TrimSpace(i)
	i = strings.Trim(i, "\"")

	return i
}

func isImport(line string, inImportClosure bool) bool {
	return inImportClosure ||
		strings.HasPrefix(line, "import") ||
		strings.HasPrefix(line, "from")
}
