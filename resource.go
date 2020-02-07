package gapt

import (
	"html/template"
	"io"
)

type Folder struct {
	Name    string // like templates
	Files   []*File
	Folders []*Folder
}

type File struct {
	Name     string             // like users.gohtml
	Cached   []byte             //only if cached?
	template *template.Template // only if template?

	Blob Blob
}

// A Blob can be
//  a local file in the source project
//  a local file in a configuration folder
//  an in-memory file
//  a cached file from a local file
//  an embedded gzipped stream
//  a virtual blob delegating and trying a bunch of other blobs and fails if none exists
type Blob interface {
	Read() (io.Reader, error)
	Write() (io.Writer, error)
	Version() [32]byte // version id e.g. a hash
}

type Gzipped interface {
	ReadRaw() (io.Reader, error)
}
