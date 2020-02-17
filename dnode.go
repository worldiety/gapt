package gapt

import (
	"fmt"
	"github.com/orcaman/writerseeker"
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

type Stream interface {
	io.Writer
	io.Reader
	io.Seeker
	io.Closer
}

// ErrAlreadyExists is a sentinel error to show that something already exists but was not expected
var ErrAlreadyExists = fmt.Errorf("already exists")

// ErrIsDir indicates that a directory was found where a file has been expected
var ErrIsDir = fmt.Errorf("is directory")

// ErrAlreadyClosed complaints about something closed, like a reader
var ErrAlreadyClosed = fmt.Errorf("already closed")

// ErrReadOnly complaints about something is read only, like a writer
var ErrReadOnly = fmt.Errorf("read only")

// dnType defines which kind of dnode we have
type dnType int

const (
	dnTypeBlob      dnType = iota
	dnTypeDirectory        // dnEnc, size, data and hash are meaningless
)

// dnEnc defines how data is encoded within a dnode
type dnEnc int

const (
	dnEncRaw dnEnc = iota
	dnEncGzip
)

// A dnode represents a "data node" which may refer to an in-memory or to a local file.
type dnode struct {
	name     string            // name of the resource, not fully qualified
	dnType   dnType            // dnType tells us if it is a regular file or a folder
	dnEnc    dnEnc             // dnEnc tells us if the resource is already compressed or if it is a directory
	size     int               // size is the raw length in bytes and if encoded as raw, equal to len(data)
	data     []byte            // data is the actual blob, encoded as told by dnEnc
	hash     [32]byte          // hash is the sha512-256 hash of the uncompressed data
	lastMod  int64             // lastMod is unix time in milliseconds since epoch
	children map[string]*dnode // children contains the list of all referring files and folders. Only valid if dnTypeDirectory
	overlay  string            // overlay is either empty or an absolute local file path
	mutex    sync.RWMutex
}

// IsFile returns true if this is a file
func (n *dnode) IsFile() bool {
	return n.dnType == dnTypeBlob
}

// IsDir returns true if this is a folder
func (n *dnode) IsDir() bool {
	return n.dnType == dnTypeDirectory
}

// Child returns the named child data node or nil
func (n *dnode) Child(name string) *dnode {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if n.children == nil {
		return nil
	}
	return n.children[name]
}

// MkDir returns either an existing directory or creates a new one. If a file with that name exists, an ErrAlreadyExists
// is returned
func (n *dnode) MkDir(name string) (*dnode, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	node := n.Child(name)
	if node == nil {
		node = &dnode{name: name, dnType: dnTypeDirectory}
		n.children[name] = node
		n.lastMod = nowAsUnixMilliseconds()
	}
	if node.dnType != dnTypeDirectory {
		return node, ErrAlreadyExists
	}
	return node, nil
}

// Delete removes the named node and returns it. If no such node exists, returns nil.
func (n *dnode) Delete(name string) *dnode {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	node := n.Child(name)
	if node != nil {
		delete(n.children, name)
		n.lastMod = nowAsUnixMilliseconds()
	}
	return node
}

// MkFile returns either an existing file or creates a new one. If a dir with that name exists, an ErrAlreadyExists
// is returned.
func (n *dnode) MkFile(name string) (*dnode, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	node := n.Child(name)
	if node == nil {
		node = &dnode{name: name, dnType: dnTypeBlob, lastMod: nowAsUnixMilliseconds(), hash: hash(nil)}
		n.children[name] = node
	}
	if node.dnType != dnTypeDirectory {
		return node, ErrAlreadyExists
	}
	return node, nil
}

// Names returns a sorted list of file names
func (n *dnode) Names() []string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	tmp := make([]string, 0, len(n.children))
	for k := range n.children {
		tmp = append(tmp, k)
	}
	sort.Strings(tmp)
	return tmp
}

// Open returns a stream to read or write the file. If it is a directory ErrIsDir is returned
func (n *dnode) Open(readOnly bool) (Stream, error) {
	if readOnly {
		n.mutex.RLock()
	} else {
		n.mutex.Lock()
	}
	if !n.IsFile() {
		if readOnly {
			n.mutex.RUnlock()
		} else {
			n.mutex.Unlock()
		}
		return nil, ErrIsDir
	}

	return newDnodeStream(n, readOnly), nil
}

type dnodeStream struct {
	dnode    *dnode
	readOnly bool
	closed   int32
	delegate *writerseeker.WriterSeeker
}

func newDnodeStream(parent *dnode, readOnly bool) *dnodeStream {
	s := &dnodeStream{
		dnode:    parent,
		readOnly: readOnly,
		delegate: &writerseeker.WriterSeeker{},
	}
	_, err := s.delegate.Write(parent.data)
	assertNil(err)
	_, err = s.delegate.Seek(0, io.SeekStart)
	assertNil(err)

	return s
}

func (s *dnodeStream) Write(p []byte) (n int, err error) {
	if s.closed != 0 {
		return 0, ErrAlreadyClosed
	}
	if s.readOnly {
		return 0, ErrReadOnly
	}
	return s.delegate.Write(p)
}

func (s *dnodeStream) Read(p []byte) (n int, err error) {
	if s.closed != 0 {
		return 0, ErrAlreadyClosed
	}
	//TODO the writer seeker is not implemented correctly, reading and seeking will break badly
}

func (s *dnodeStream) Seek(offset int64, whence int) (int64, error) {
	if s.closed != 0 {
		return 0, ErrAlreadyClosed
	}
}

// Close is idempotent
func (s *dnodeStream) Close() error {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return ErrAlreadyClosed
	}
	if s.readOnly {
		s.dnode.mutex.RUnlock()
	} else {
		s.dnode.mutex.Unlock()
	}
	return nil
}
