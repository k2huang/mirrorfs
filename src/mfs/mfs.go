// mirrorfs implementation
package mfs

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// 表示根文件系统
type MirrorFS struct {
	root   *Dir // 根目录
	nodeId uint64
	path   string //mirror path
}

// Compile-time interface checks.
var _ fs.FS = (*MirrorFS)(nil)
var _ fs.FSStatfser = (*MirrorFS)(nil)

const DEF_MODE = os.FileMode(0777)

func NewMirrorFS(path string) *MirrorFS {
	fs := &MirrorFS{
		path: path,
	}
	fs.root = fs.newDir(path, os.ModeDir|DEF_MODE)
	if fs.root.attr.Inode != 1 {
		panic("Root node should have been assigned id 1")
	}
	return fs
}

func (m *MirrorFS) nextId() uint64 {
	return atomic.AddUint64(&m.nodeId, 1)
}

func (m *MirrorFS) newDir(path string, mode os.FileMode) *Dir {
	n := time.Now()
	return &Dir{
		attr: fuse.Attr{
			Inode:  m.nextId(),
			Atime:  n,
			Mtime:  n,
			Ctime:  n,
			Crtime: n,
			Mode:   os.ModeDir | mode,
		},
		path: path,
		fs:   m,
	}
}

func (m *MirrorFS) newFile(path string, mode os.FileMode) *File {
	n := time.Now()
	return &File{
		attr: fuse.Attr{
			Inode:  m.nextId(),
			Atime:  n,
			Mtime:  n,
			Ctime:  n,
			Crtime: n,
			Mode:   mode,
		},
		path: path,
	}
}

func (f *MirrorFS) Root() (fs.Node, error) {
	return f.root, nil
}

func (f *MirrorFS) Statfs(ctx context.Context, req *fuse.StatfsRequest, res *fuse.StatfsResponse) error {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(f.path, &s)
	if err != nil {
		log.Println("DRIVE | Statfs syscall failed:", err)
		return err
	}

	res.Blocks = s.Blocks
	res.Bfree = s.Bfree
	res.Bavail = s.Bavail
	res.Ffree = s.Ffree
	res.Bsize = uint32(s.Bsize)

	return nil
}
