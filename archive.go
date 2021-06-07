package main

import (
	"archive/zip"
	"os"
	"sync"
)

var zipMux sync.Mutex

// archiveWriter is an archive writer interface
type archiveWriter interface {
	add(string, []byte) error
	close() error
}

// archive is a file-baesd implementation of archiveWriter
type archive struct {
	file *os.File
	zw   *zip.Writer
}

// newArchiveWriter creates a new file-based archive writer
func newArchiveWriter(name string) (archiveWriter, error) {
	f, err := os.Create(name)
	if err != nil {
		return archive{}, err
	}
	zw := zip.NewWriter(f)
	return archive{
		file: f,
		zw:   zw,
	}, nil
}

// close closes the zip writer and file
func (a archive) close() error {
	err := a.zw.Close()
	if err != nil {
		return err
	}
	return a.file.Close()
}

// add adds a file and content to the zip archive
func (a archive) add(name string, content []byte) error {
	zipMux.Lock()
	defer zipMux.Unlock()
	f, err := a.zw.Create(name)
	if err != nil {
		return nil
	}
	_, err = f.Write(content)
	return err
}
