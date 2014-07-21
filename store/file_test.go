package store

import (
	"os"
	"path"
	"testing"
)

func TestFilePath(t *testing.T) {
	tmp := os.TempDir()
	store := FileStore{Root: tmp}

	file, err := store.Create("/usr/bin/sleep")
	if err != nil {
		t.Fatal(err)
	}

	expected := path.Join(tmp, ".pew", "sleep.csv")
	if file.Name() != expected {
		t.Fatalf("unexpected path: %v", file.Name())
	}
}
