package pp

import (
	"io/fs"
	"path/filepath"
)

type FileMeta struct {
	Name string
	Size int64
}

func BuildFileCatalog(folder string) ([]FileMeta, error) {
	var files []FileMeta

	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		files = append(files, FileMeta{
			Name: d.Name(),
			Size: info.Size(),
		})
		return nil
	})

	return files, err
}
