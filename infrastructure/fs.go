package infrastructure

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/prskr/git-age/core/ports"
)

var _ ports.ReadWriteFS = (*ReadWriteDirFS)(nil)

func NewReadWriteDirFS(rootPath string) *ReadWriteDirFS {
	return &ReadWriteDirFS{
		rootPath:   rootPath,
		underlying: os.DirFS(rootPath),
	}
}

type ReadWriteDirFS struct {
	rootPath   string
	underlying fs.FS
}

func (f ReadWriteDirFS) Rename(oldPath, newPath string) error {
	return os.Rename(filepath.Join(f.rootPath, oldPath), filepath.Join(f.rootPath, newPath))
}

func (f ReadWriteDirFS) Remove(filePath string) error {
	return os.Remove(filepath.Join(f.rootPath, filePath))
}

func (f ReadWriteDirFS) TempFile(dir, pattern string) (ports.ReadWriteFile, error) {
	tmpFile, err := os.CreateTemp(filepath.Join(f.rootPath, dir), pattern)
	if err != nil {
		return nil, err
	}

	return readWriteOsFile{File: tmpFile, fsRoot: f.rootPath}, nil
}

func (f ReadWriteDirFS) OpenRW(filePath string) (ports.ReadWriteFile, error) {
	file, err := os.OpenFile(filepath.Join(f.rootPath, filePath), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return readWriteOsFile{File: file, fsRoot: f.rootPath}, nil
}

func (f ReadWriteDirFS) Open(name string) (fs.File, error) {
	return f.underlying.Open(name)
}

var _ ports.ReadWriteFile = (*readWriteOsFile)(nil)

type readWriteOsFile struct {
	fsRoot string
	*os.File
}

func (r readWriteOsFile) Name() string {
	name := r.File.Name()

	var err error
	name, err = filepath.Rel(r.fsRoot, name)

	if err != nil {
		panic(fmt.Sprintf("failed to get relative path: %s", err.Error()))
	}

	return name
}
