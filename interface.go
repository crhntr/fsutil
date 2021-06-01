package filesystem

import "os"

type DirReader interface {
	ReadDir(path string) ([]os.FileInfo, error)
}
