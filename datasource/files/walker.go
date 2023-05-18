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
}

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

		f := File{
			Package: fPackage,
			Name:    fName,
			Lines:   uint32(len(strings.Split(content, "\n"))),
			Symbols: uint32(len(content)),
		}

		c <- f

		return err
	})
}
