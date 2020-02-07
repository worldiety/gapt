package gapt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resource represents an identified resource at generation time
type fResource struct {
	Src  string // Src contains the local absolute filename
	Path string // Path contains the fully qualified resource name, relative to the module root
}

// A Project represents the go project whose resources needs to be processed.
type Project struct {
	root      string
	pkg       *Package
	config    *Config
	resources []Resource
}

// NewProject parses the given directory as a go module.
func NewProject(dir string) (*Project, error) {
	pkg, err := GoList(dir)
	if err != nil {
		return nil, fmt.Errorf("not a go module: %s: %w", dir, err)
	}

}

// Collect
func (b *Project) Collect(dir string) {

}

// collect scans all files and only includes those which are candidates for embedding.
// Ignored files (case insensitive) are
//   Makefile
//   LICENSE
//   *.md
//   *.go
//   .*
//    *.mod
//   *.sum
//   build
func collect(root string) ([]string, error) {
	var res []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") || info.Name() == "build" {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Mode().IsRegular() {
			lcase := strings.ToLower(info.Name())
			allowed := true
			for _, ext := range ignoreFileExt {
				if strings.HasSuffix(lcase, ext) {
					allowed = false
					break
				}
			}
			if allowed {
				res = append(res, path)
			}
		}
		return nil
	})
	return res, err
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(dir)
	}
	files, err := collect(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file[len(dir):])
	}

	pkg, err := GoList(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println(pkg.String())
}
