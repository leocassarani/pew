package store

import (
	"os"
	"path"
)

type FileStore struct {
	Root string
}

func (fs FileStore) Create(cmd string) (*File, error) {
	if err := fs.makePewDir(); err != nil {
		return nil, err
	}

	filepath := fs.pewFilePath(cmd)
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return &File{fd: file}, nil
}

func (fs FileStore) makePewDir() error {
	dirpath := path.Join(fs.Root, ".pew")
	return os.MkdirAll(dirpath, os.ModePerm)
}

func (fs FileStore) pewFilePath(cmd string) string {
	_, binary := path.Split(cmd)
	filename := binary + ".csv"
	return path.Join(fs.Root, ".pew", filename)
}

type File struct {
	fd   *os.File
	path string
}

func (f *File) Close() error {
	return f.fd.Close()
}

func (f *File) Write(b []byte) (int, error) {
	return f.fd.Write(b)
}

func (f *File) Name() string {
	return f.fd.Name()
}
