package filesystem

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

func WriteToTarball(wr io.Writer, fs billy.Filesystem) (deferErr error) {
	tw := tar.NewWriter(wr)

	defer func() {
		deferErr = tw.Close()
	}()

	return Walk(fs, "./", func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel("./", name)
		if err != nil {
			return err
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := fs.Open(name)
		if err != nil {
			return err
		}
		defer func() {
			_ = f.Close()
		}()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})
}

func ReadFromTarball(r io.Reader, fs billy.Filesystem) error {
	fsTar := tar.NewReader(r)

	for {
		header, err := fsTar.Next()
		if err == io.EOF {
			break // End of archive
		}
		if header == nil {
			continue
		}

		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeDir {
			if _, err := os.Stat(header.Name); os.IsNotExist(err) {
				if err = fs.MkdirAll(path.Clean(header.Name), os.FileMode(header.Mode)); err != nil {
					return err
				}
			}
			continue
		}

		if err := func() error {
			srcFile, err := fs.OpenFile(header.Name, os.O_RDWR|os.O_CREATE, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer func() {
				_ = srcFile.Close()
			}()

			_, err = io.Copy(srcFile, fsTar)
			if err != nil {
				return err
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	return nil
}
