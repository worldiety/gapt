package gapt

import (
	"github.com/worldiety/gapt/unix"
	"net/http"
	"os"
	"time"
)

var _ http.FileSystem = (*Filesystem)(nil)

type Filesystem struct {
}

func (f *Filesystem) Open(name string) (http.File, error) {
	panic("implement me")
}

var _ http.File = (*File)(nil)

type File struct {
	name    string
	dir     bool
	lastMod unix.Time
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	panic("implement me")
}

func (f *File) Mode() os.FileMode {
	if f.dir {
		return os.ModeDir
	} else {
		return 0
	}
}

func (f *File) ModTime() time.Time {
	return f.lastMod.Time()
}

func (f *File) IsDir() bool {
	return f.dir
}

func (f *File) Sys() interface{} {
	return nil
}

func (f *File) Close() error {
	panic("implement me")
}

func (f *File) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	panic("implement me")
}

func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}
