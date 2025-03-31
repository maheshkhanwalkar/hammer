package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFile(t *testing.T) {
	source := create(t.TempDir(), "test.c")[0]

	f, err := NewFile(source)
	if err != nil {
		t.Fatal(err)
	}

	if f.Source == nil {
		t.Fatal("expected non-nil source file")
	}

	if f.Object != nil {
		// We expect f.Object to be nil, because it doesn't exist
		t.Fatal("expected nil object file")
	}

	_ = f.Close()
}

func TestFile_ShouldCompile_IfObjectFileDoesNotExist_Then_ReturnTrue(t *testing.T) {
	source := create(t.TempDir(), "test.c")[0]

	f, _ := NewFile(source)
	if !f.ShouldCompile() {
		t.Fatal("expected compilation to be necessary, as object file does not exist")
	}
}

func TestFile_ShouldCompile_IfObjectFileIsNewer_Then_ReturnFalse(t *testing.T) {
	source := create(t.TempDir(), "test.c", "test.o")[0]

	f, _ := NewFile(source)
	if f.ShouldCompile() {
		t.Fatal("expected compilation to be not be necessary, as object file is newer")
	}
}

func TestFile_ShouldCompile_IfObjectFileIsOlder_Then_ReturnTrue(t *testing.T) {
	// Create object first, so it's older
	source := create(t.TempDir(), "test.o", "test.c")[1]

	f, _ := NewFile(source)
	if !f.ShouldCompile() {
		t.Fatal("expected compilation to be be necessary, as object file is older")
	}
}

// create the files in the order specified in names in the directory, returning the full paths
// of the files created
func create(dir string, names ...string) []string {
	res := make([]string, 0, len(names))

	for i := range names {
		name := names[i]
		fullPath := filepath.Join(dir, name)

		_, _ = os.Create(fullPath)
		res = append(res, fullPath)
	}

	return res
}
