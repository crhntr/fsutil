package filesystem

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/go-git/go-billy/v5"
)

// Copy recursively from src to root
func Copy(dst, src billy.Filesystem) error {
	return cpDir(src.Root(), dst.Root(), src, dst)
}

func cpFile(srcPath, dstPath string, src, dst billy.Filesystem) error {
	var (
		err     error
		srcFD   billy.File
		dstFD   billy.File
		srcInfo os.FileInfo
	)

	if srcFD, err = src.Open(srcPath); err != nil {
		return err
	}
	defer func() {
		_ = srcFD.Close()
	}()

	if srcInfo, err = src.Stat(srcPath); err != nil {
		return err
	}

	if dstFD, err = dst.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcInfo.Mode()); err != nil {
		return err
	}
	defer func() {
		_ = dstFD.Close()
	}()

	if _, err = io.Copy(dstFD, srcFD); err != nil {
		return err
	}

	return nil
}

func cpDir(srcPath, dstPath string, src, dst billy.Filesystem) error {
	var (
		err     error
		fds     []os.FileInfo
		srcInfo os.FileInfo
	)

	if srcInfo, err = src.Stat(srcPath); err != nil {
		return err
	}

	if err = dst.MkdirAll(dstPath, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = src.ReadDir(srcPath); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFP := path.Join(srcPath, fd.Name())
		dstFP := path.Join(dstPath, fd.Name())

		if fd.IsDir() {
			if err = cpDir(srcFP, dstFP, src, dst); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = cpFile(srcFP, dstFP, src, dst); err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}
