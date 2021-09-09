package main

import (
	"fmt"
	"os"
	"strings"
)

// nextFilename returns the next available filename,
// e.g. if "file.zip" exists "file (2).zip" will be used, and so on.
func nextFilename(directory, filename string) string {
	os.Chdir(directory)
	if _, err := os.Stat(filename); err != nil {
		return filename
	}
	parts := strings.SplitN(filename, ".", 2)
	name := parts[0]
	ext := parts[1]
	for i := 1; ; i++ {
		filename := fmt.Sprintf("%s (%d).%s", name, i, ext)
		if _, err := os.Stat(filename); err != nil {
			return filename
		}
	}
}
