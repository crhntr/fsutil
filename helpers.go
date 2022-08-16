package fsutil

import (
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
)

// ReadFile is based on the standard library function fs.ReadFile
func ReadFile(dir billy.Basic, name string) ([]byte, error) {
	f, err := dir.Open(name)
	if err != nil {
		return nil, err
	}
	defer closeAndIgnoreError(f)
	return io.ReadAll(f)
}

// WriteFile is based on the standard library function os.WriteFile but adds a dir parameter
func WriteFile(dir billy.Basic, name string, data []byte, perm os.FileMode) error {
	f, err := dir.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer closeAndIgnoreError(f)
	_, err = f.Write(data)
	return err
}

func closeAndIgnoreError(c io.Closer) {
	_ = c.Close()
}
