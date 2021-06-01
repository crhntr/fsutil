package filesystem

import (
	"os"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
)

func TestWriteFilesystemToTarball(t *testing.T) {
	t.SkipNow()

	fs := memfs.New()

	f, err := fs.Create("./hello.go")
	if err != nil {
		t.Errorf("it should not error: %q", err)
	}

	_, _ = f.Write([]byte("package main\n\nfunc main() {}\n"))
	_ = f.Close()

	if err := os.MkdirAll("test_data", 0777); err != nil {
		t.Fatal("could not create test_data directory")
	}

	img, err := os.Create("test_data/image.tar")
	if err != nil {
		t.Fatal("could not create test_data directory")
	}
	defer func() {
		_ = img.Close()
	}()

	if err := WriteToTarball(img, fs); err != nil {
		t.Errorf("it should not error: %q", err)
	}

	// os.RemoveAll("test_data")
}
