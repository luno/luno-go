package clientgen

import (
	"bytes"
	"path"
	"path/filepath"
)

type File struct {
	relPath string
	name    string
	buf     *bytes.Buffer
}

func NewFile(relPath string) File {
	return File{
		relPath: filepath.Dir(relPath),
		name:    filepath.Base(relPath),
		buf:     bytes.NewBuffer(nil),
	}
}

func (f File) RelPath() string {
	return path.Join(f.relPath, f.name)
}

func (f File) Write(b []byte) (int, error) {
	return f.buf.Write(b)
}

func (f File) Bytes() []byte {
	return f.buf.Bytes()
}
