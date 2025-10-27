package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
)

type FS struct {
	root  string
	cache map[string]bool
}

func NewFS(path string) (*FS, error) {
	fmt.Println(os.Args[0])
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	f := &FS{
		root:  path,
		cache: make(map[string]bool),
	}

	if err := f.cached(); err != nil {
		return nil, err
	}

	return f, nil
}

func (fs *FS) cached() error {
	names, err := fs.ReadNames()
	if err != nil {
		return err
	}

	for _, name := range names {
		fs.cache[name] = true
	}
	return nil
}

func (fs *FS) Open(name string) (fs.File, error) {
	return os.OpenFile(path.Join(fs.root, name), os.O_RDONLY, os.ModePerm)
}

func (fs *FS) CheckExist(name string) bool {
	if exist, ok := fs.cache[name]; ok {
		return exist
	}

	_, err := os.Stat(path.Join(fs.root, name))
	exist := err == nil || os.IsExist(err)

	fs.cache[name] = exist

	return exist
}

func (fs *FS) Save(data io.Reader, name string) error {
	f, err := os.OpenFile(path.Join(fs.root, name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return err
	}

	fs.cache[name] = true

	return nil
}

func (fs *FS) Remove(name string) error {
	if err := os.Remove(path.Join(fs.root, name)); err != nil {
		return err
	}
	fs.cache[name] = false
	return nil
}

func (fs *FS) Count() (int, error) {
	names, err := fs.ReadNames()
	if err != nil {
		return 0, err
	}

	return len(names), nil
}

func (fs *FS) ReadNames() ([]string, error) {
	f, err := os.OpenFile(fs.root, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}
