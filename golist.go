// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gapt

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GoList performs 'go list -json' within the given directory. There is no public API for this and
// it is recommended to use the installed go version, instead of bringing our own.
func GoList(dir string) (*Package, error) {
	cmd := exec.Command("go", "list", "-json")
	cmd.Dir = dir
	res, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to 'go list -json': %w", err)
	}

	p := &Package{}
	err = json.Unmarshal(res, p)

	if err != nil {
		return nil, fmt.Errorf("failed to parse output from 'go list -json': %w", err)
	}

	return p, nil
}

// Package contains all those info from 'go list -json'
type Package struct {
	Dir           string   `json:",omitempty"` // directory containing package sources
	ImportPath    string   `json:",omitempty"` // import path of package in dir
	ImportComment string   `json:",omitempty"` // path in import comment on package statement
	Name          string   `json:",omitempty"` // package name
	Doc           string   `json:",omitempty"` // package documentation string
	Target        string   `json:",omitempty"` // installed target for this package (may be executable)
	Shlib         string   `json:",omitempty"` // the shared library that contains this package (only set when -linkshared)
	Root          string   `json:",omitempty"` // Go root, Go path dir, or module root dir containing this package
	ConflictDir   string   `json:",omitempty"` // Dir is hidden by this other directory
	ForTest       string   `json:",omitempty"` // package is only for use in named test
	Export        string   `json:",omitempty"` // file containing export data (set by go list -export)
	Module        *Module  `json:",omitempty"` // info about package's module, if any
	Match         []string `json:",omitempty"` // command-line patterns matching this package
	Goroot        bool     `json:",omitempty"` // is this package found in the Go root?
	Standard      bool     `json:",omitempty"` // is this package part of the standard Go library?
	DepOnly       bool     `json:",omitempty"` // package is only as a dependency, not explicitly listed
	BinaryOnly    bool     `json:",omitempty"` // package cannot be recompiled
	Incomplete    bool     `json:",omitempty"` // was there an error loading this package or dependencies?

	// Stale and StaleReason remain here *only* for the list command.
	// They are only initialized in preparation for list execution.
	// The regular build determines staleness on the fly during action execution.
	Stale       bool   `json:",omitempty"` // would 'go install' do anything for this package?
	StaleReason string `json:",omitempty"` // why is Stale true?

	// Source files
	GoFiles         []string `json:",omitempty"` // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles        []string `json:",omitempty"` // .go source files that import "C"
	CompiledGoFiles []string `json:",omitempty"` // .go output from running cgo on CgoFiles
	IgnoredGoFiles  []string `json:",omitempty"` // .go source files ignored due to build constraints
	CFiles          []string `json:",omitempty"` // .c source files
	CXXFiles        []string `json:",omitempty"` // .cc, .cpp and .cxx source files
	MFiles          []string `json:",omitempty"` // .m source files
	HFiles          []string `json:",omitempty"` // .h, .hh, .hpp and .hxx source files
	FFiles          []string `json:",omitempty"` // .f, .F, .for and .f90 Fortran source files
	SFiles          []string `json:",omitempty"` // .s source files
	SwigFiles       []string `json:",omitempty"` // .swig files
	SwigCXXFiles    []string `json:",omitempty"` // .swigcxx files
	SysoFiles       []string `json:",omitempty"` // .syso system object files added to package

	// Cgo directives
	CgoCFLAGS    []string `json:",omitempty"` // cgo: flags for C compiler
	CgoCPPFLAGS  []string `json:",omitempty"` // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string `json:",omitempty"` // cgo: flags for C++ compiler
	CgoFFLAGS    []string `json:",omitempty"` // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string `json:",omitempty"` // cgo: flags for linker
	CgoPkgConfig []string `json:",omitempty"` // cgo: pkg-config names

	// Dependency information
	Imports   []string          `json:",omitempty"` // import paths used by this package
	ImportMap map[string]string `json:",omitempty"` // map from source import to ImportPath (identity entries omitted)
	Deps      []string          `json:",omitempty"` // all (recursively) imported dependencies

	// Error information
	// Incomplete is above, packed into the other bools
	Error      *PackageError   `json:",omitempty"` // error loading this package (not dependencies)
	DepsErrors []*PackageError `json:",omitempty"` // errors loading dependencies

	// Test information
	// If you add to this list you MUST add to p.AllFiles (below) too.
	// Otherwise file name security lists will not apply to any new additions.
	TestGoFiles  []string `json:",omitempty"` // _test.go files in package
	TestImports  []string `json:",omitempty"` // imports from TestGoFiles
	XTestGoFiles []string `json:",omitempty"` // _test.go files outside package
	XTestImports []string `json:",omitempty"` // imports from XTestGoFiles
}

func (p *Package) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// A Module contains all those info from 'go list -json'
type Module struct {
	Path      string       `json:",omitempty"` // module path
	Version   string       `json:",omitempty"` // module version
	Versions  []string     `json:",omitempty"` // available module versions
	Replace   *Module      `json:",omitempty"` // replaced by this module
	Time      *time.Time   `json:",omitempty"` // time version was created
	Update    *Module      `json:",omitempty"` // available update (with -u)
	Main      bool         `json:",omitempty"` // is this the main module?
	Indirect  bool         `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string       `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string       `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string       `json:",omitempty"` // go version used in module
	Error     *ModuleError `json:",omitempty"` // error loading module
}

// A ModuleError describes an error loading information about a module.
type ModuleError struct {
	Err string // error text
}

func (m *Module) String() string {
	s := m.Path
	if m.Version != "" {
		s += " " + m.Version
		if m.Update != nil {
			s += " [" + m.Update.Version + "]"
		}
	}

	if m.Replace != nil {
		s += " => " + m.Replace.Path
		if m.Replace.Version != "" {
			s += " " + m.Replace.Version
			if m.Replace.Update != nil {
				s += " [" + m.Replace.Update.Version + "]"
			}
		}
	}

	return s
}

// A PackageError describes an error loading information about a package.
type PackageError struct {
	ImportStack   []string // shortest path from package named on command line to this one
	Pos           string   // position of error
	Err           string   // the error itself
	IsImportCycle bool     `json:"-"` // the error is an import cycle
	Hard          bool     `json:"-"` // whether the error is soft or hard; soft errors are ignored in some places
}

func (p *PackageError) Error() string {
	// Import cycles deserve special treatment.
	if p.IsImportCycle {
		return fmt.Sprintf("%s\npackage %s\n", p.Err, strings.Join(p.ImportStack, "\n\timports "))
	}

	if p.Pos != "" {
		// Omit import stack. The full path to the file where the error
		// is the most important thing.
		return p.Pos + ": " + p.Err
	}

	if len(p.ImportStack) == 0 {
		return p.Err
	}

	return "package " + strings.Join(p.ImportStack, "\n\timports ") + ": " + p.Err
}

// An ImportStack is a stack of import paths, possibly with the suffix " (test)" appended.
// The import path of a test package is the import path of the corresponding
// non-test package with the suffix "_test" added.
type ImportStack []string
