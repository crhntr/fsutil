package fsutil

import (
	"path/filepath"
)

// FileNames returns name of non directory files
// For now this includes symlinks... this may change.
func FileNames(repo DirReader) ([]string, error) {
	var walk func(*[]string, DirReader, string) error
	walk = func(paths *[]string, dr DirReader, dir string) error {
		files, err := dr.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, stat := range files {
			switch stat.Name() {
			case ".git", "vendor", "node_modules":
				continue
			}
			fp := filepath.Join(dir, stat.Name())
			if !stat.IsDir() {
				*paths = append(*paths, fp)
				continue
			}
			if err := walk(paths, dr, fp); err != nil {
				return err
			}
		}
		return nil
	}

	var paths []string
	if err := walk(&paths, repo, ""); err != nil {
		return nil, err
	}

	return paths, nil
}
