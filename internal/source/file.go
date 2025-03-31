package source

import (
	"errors"
	"os"
	"strings"
)

// File is the lowest level (leaf) construct in the build system
type File struct {
	Source *os.File
	Object *os.File
}

// NewFile creates a new file from the source path
func NewFile(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dotPos := strings.LastIndex(file.Name(), ".")
	objPath := path[:dotPos] + ".o"

	objFile, err := os.Open(objPath)
	if err != nil {
		objFile = nil
	}

	return &File{Source: file, Object: objFile}, nil
}

// ShouldCompile determines whether the corresponding object file
// is up to date with the source -- that is, whether we should (re)compile
func (f *File) ShouldCompile() bool {
	if f.Object == nil {
		return true
	}
	objInfo, _ := f.Object.Stat()
	srcInfo, _ := f.Source.Stat()

	return objInfo.ModTime().Before(srcInfo.ModTime())
}

// Close the file
func (f *File) Close() error {
	allErrors := make([]error, 0)
	allErrors = append(allErrors, f.Source.Close())

	if f.Object != nil {
		allErrors = append(allErrors, f.Object.Close())
	}

	return errors.Join(allErrors...)
}
