package xfs

import (
	"io/fs"
	"path/filepath"
)

type prefixFS struct {
	prefix string
	fs     fs.FS
}

func (f *prefixFS) Open(name string) (fs.File, error) {
	name = filepath.Join(f.prefix, name)
	return f.fs.Open(name)
}

// PrefixFS delegates Open(filename) to fs.Open(prefix+"/"+filename)
func PrefixFS(prefix string, fs fs.FS) fs.FS {
	return &prefixFS{prefix: prefix, fs: fs}
}
