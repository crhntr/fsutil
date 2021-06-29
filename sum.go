package fsutil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/go-git/go-billy/v5"
)

func SumFileNameAndContents(fs billy.Filesystem) (string, error) {
	var (
		err error
		fds []os.FileInfo
	)

	var recursiveSum func(fs billy.Filesystem, path string) ([]byte, error)
	recursiveSum = func(fs billy.Filesystem, p string) ([]byte, error) {
		if fds, err = fs.ReadDir(p); err != nil {
			return nil, err
		}

		sort.Sort(sorter{
			len(fds),
			func(i, j int) { fds[i], fds[j] = fds[j], fds[i] },
			func(i, j int) bool {
				return strings.Compare(fds[i].Name(), fds[j].Name()) > 0
			},
		})

		var sums []byte
		for _, fd := range fds {
			nm := path.Join(p, fd.Name())
			if fd.IsDir() {
				sum, err := recursiveSum(fs, nm)
				if err != nil {
					return nil, err
				}
				sums = append(sums, sum...)
			} else {
				h := sha256.New()
				f, err := fs.Open(nm)
				if err != nil {
					return nil, err
				}
				if _, err := h.Write([]byte(nm)); err != nil {
					return nil, err
				}
				if _, err := io.Copy(h, f); err != nil {
					return nil, err
				}
				sums = append(sums, h.Sum(nil)...)
			}
		}
		sm := sha256.Sum256(sums)

		return sm[:], nil
	}

	buf, err := recursiveSum(fs, "/")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

// srt "generic"
// https://medium.com/capital-one-tech/closures-are-the-generics-for-go-cb32021fb5b5
type sorter struct {
	len  int
	swap func(i, j int)
	less func(i, j int) bool
}

func (x sorter) Len() int           { return x.len }
func (x sorter) Swap(i, j int)      { x.swap(i, j) }
func (x sorter) Less(i, j int) bool { return x.less(i, j) }
