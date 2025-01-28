package gowork

import (
	"os"
	"path"

	"golang.org/x/mod/modfile"
)

// Module represents a Go module.
type Module struct {
	// Name is a module name.
	Name string
	// Path is a module file path.
	Path string
}

// Open will parse a module file.
func (m *Module) Open() (*modfile.File, error) {
	gomodPath := path.Join(m.Path, "go.mod")
	bytes, ferr := os.ReadFile(gomodPath)
	if ferr != nil {
		return nil, ferr
	}
	return modfile.ParseLax(gomodPath, bytes, nil)
}
