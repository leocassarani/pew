package store

import (
	"os"
	"path"
	"strconv"
	"time"
)

func Create() (*File, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err = makePewDir(pwd); err != nil {
		return nil, err
	}

	filepath := pewFilePath(pwd)
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return &File{fd: file}, nil
}

func makePewDir(dir string) error {
	dirpath := path.Join(dir, ".pew")
	return os.MkdirAll(dirpath, os.ModePerm)
}

func pewFilePath(dir string) string {
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	return path.Join(dir, ".pew", timeStr+".csv")
}

type File struct {
	fd *os.File
}

func (f *File) Close() error {
	return f.fd.Close()
}

func (f *File) Write(b []byte) (int, error) {
	return f.fd.Write(b)
}
