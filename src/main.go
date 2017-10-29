package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"mfs"
)

var (
	debug  = flag.Bool("debug", false, "enable debug log messages to stderr")
	mirror = flag.String("mirror", "", "path to mirror contents (required)")
	mount  = flag.String("mount", "", "path to mount volume (required)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func debugLog(msg interface{}) {
	fmt.Printf("%s", msg)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *mount == "" || *mirror == "" {
		usage()
		os.Exit(2)
	}

	c, err := fuse.Mount(
		*mount,
		fuse.FSName("mirrorfs"),
		fuse.Subtype("mirrorfs"),
		fuse.VolumeName("Mirror FS"),
		// fuse.LocalVolume(),
		// fuse.AllowOther(),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cfg := &fs.Config{}
	if *debug {
		cfg.Debug = debugLog
	}
	srv := fs.New(c, cfg)
	filesys := mfs.NewMirrorFS(*mirror)

	if err := srv.Serve(filesys); err != nil {
		log.Fatal(err)
	}

	// Check if the mount process has an error to report.
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
