package archive

import (
	"archive/zip"
	"os"
	"sync"
)

var zipMux sync.Mutex

// Writer is an archive writer interface
type Writer interface {
	Add(string, []byte) error
	Close() error
}

// FileWriter is a file-based implementation of archiveWriter
type FileWriter struct {
	file *os.File
	zw   *zip.Writer
}

// NewWriter creates a new file-based archive writer
func NewWriter(name string) (Writer, error) {
	f, err := os.Create(name)
	if err != nil {
		return FileWriter{}, err
	}
	zw := zip.NewWriter(f)
	return FileWriter{
		file: f,
		zw:   zw,
	}, nil
}

// Close closes the zip writer and file
func (a FileWriter) Close() error {
	err := a.zw.Close()
	if err != nil {
		return err
	}
	return a.file.Close()
}

// Add adds a file and content to the zip archive
func (a FileWriter) Add(name string, content []byte) error {
	zipMux.Lock()
	defer zipMux.Unlock()
	f, err := a.zw.Create(name)
	if err != nil {
		return nil
	}
	_, err = f.Write(content)
	return err
}
