package sourcemap

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Path     string
	Filename string
}

func Discover(dir string) ([]File, error) {
	var maps []File

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "node_modules" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".map") {
			maps = append(maps, File{
				Path:     path,
				Filename: info.Name(),
			})
		}

		return nil
	})

	return maps, err
}
