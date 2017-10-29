// mirrorfs implementation
package mfs

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var _ fs.Node = (*Dir)(nil)
var _ fs.NodeCreater = (*Dir)(nil)
var _ fs.NodeMkdirer = (*Dir)(nil)
var _ fs.NodeRemover = (*Dir)(nil)
var _ fs.NodeRenamer = (*Dir)(nil)
var _ fs.NodeStringLookuper = (*Dir)(nil)

// 表示文件系统中的目录节点
// Dir既是fs.Node也是fs.Handle
type Dir struct {
	sync.RWMutex
	attr fuse.Attr

	path string
	fs   *MirrorFS
}

func (d *Dir) Attr(ctx context.Context, o *fuse.Attr) error {
	d.RLock()
	*o = d.attr
	d.RUnlock()

	return nil
}

// 通过文件名查找当前Dir对应的目录下的文件结点
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	d.RLock()
	defer d.RUnlock()

	path := filepath.Join(d.path, name)
	stats, err := os.Stat(path)
	if err != nil {
		//The real file does not exists.
		log.Println("Lookup ERR:", err)
		return nil, fuse.ENOENT
	}

	switch {
	case stats.IsDir():
		return d.fs.newDir(path, stats.Mode()), nil //生成新的目录结点

	case stats.Mode().IsRegular():
		return d.fs.newFile(path, stats.Mode()), nil //生成新的文件结点

	default:
		panic("Unknown type in filesystem")
	}
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	d.RLock()
	defer d.RUnlock()

	var out []fuse.Dirent
	files, err := ioutil.ReadDir(d.path)
	if err != nil {
		log.Println("ReadDirAll ERR:", err)
		return nil, err
	}

	for _, node := range files {
		de := fuse.Dirent{Name: node.Name()}
		if node.IsDir() {
			de.Type = fuse.DT_Dir
		}
		if node.Mode().IsRegular() {
			de.Type = fuse.DT_File
		}

		out = append(out, de)
	}

	return out, nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	d.Lock()
	defer d.Unlock()

	if exists := d.exists(req.Name); exists {
		log.Println("Mkdir ERR: EEXIST")
		return nil, fuse.EEXIST
	}

	path := filepath.Join(d.path, req.Name)
	n := d.fs.newDir(path, req.Mode)

	if err := os.Mkdir(path, req.Mode); err != nil {
		log.Println("Mkdir ERR:", err)
		return nil, err
	}

	return n, nil
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	d.Lock()
	defer d.Unlock()

	if exists := d.exists(req.Name); exists {
		log.Println("Create open ERR: EEXIST")
		return nil, nil, fuse.EEXIST
	}

	path := filepath.Join(d.path, req.Name)
	fHandler, err := os.OpenFile(path, int(req.Flags), req.Mode)
	if err != nil {
		log.Println("Create open ERR:", err)
		return nil, nil, err
	}

	n := d.fs.newFile(path, req.Mode)
	n.handler = fHandler

	resp.Attr = n.attr

	return n, n, nil
}

func (d *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	nd := newDir.(*Dir)
	log.Println(req)
	if d.attr.Inode == nd.attr.Inode {
		d.Lock()
		defer d.Unlock()
	} else if d.attr.Inode < nd.attr.Inode {
		d.Lock()
		defer d.Unlock()
		nd.Lock()
		defer nd.Unlock()
	} else {
		nd.Lock()
		defer nd.Unlock()
		d.Lock()
		defer d.Unlock()
	}

	if exists := d.exists(req.OldName); !exists {
		log.Println("Rename ERR: ENOENT")
		return fuse.ENOENT
	}

	oldPath := filepath.Join(d.path, req.OldName)
	newPath := filepath.Join(nd.path, req.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		log.Println("Rename ERR:", err)
		return err
	}

	return nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	d.Lock()
	defer d.Unlock()
	log.Println(req, filepath.Base(d.path), req.Name)

	if exists := d.exists(req.Name); !exists {
		log.Println("Remove ERR: ENOENT")
		return fuse.ENOENT
	}

	path := filepath.Join(d.path, req.Name)
	if err := os.Remove(path); err != nil {
		log.Println("Remove ERR:", err)
		return err
	}
	return nil
}

func (d *Dir) exists(name string) bool {
	path := filepath.Join(d.path, name)
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}
