package filesystem

import (
	"path"
	"strings"
)

func HasDotPrefixedSegment(fp string) bool {
	dir := fp
	for {
		var file string
		dir, file = path.Split(dir)

		if strings.HasPrefix(file, ".") {
			return true
		}

		if dir == "/" || dir == "" || file == "" {
			break
		}
	}

	return false
}
