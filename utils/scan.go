package utils

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

type visitor struct {
	results []string
}

func (v *visitor) scan(filepath string, ext string) error {
	ls, err := ioutil.ReadDir(filepath)
	if err != nil {
		return fmt.Errorf("cannot read the path %s: %s", filepath, err)
	}

	for _, entry := range ls {
		fullpath := path.Join(filepath, entry.Name())

		if entry.IsDir() {
			if v.validDir(entry.Name()) {
				if err := v.scan(fullpath, ext); err != nil {
					return err
				}
			}
		} else if strings.HasSuffix(entry.Name(), ext) {
			v.results = append(v.results, fullpath)
		}
	}

	return nil
}

// Returns true if the directory name is worth scanning.
func (v *visitor) validDir(name string) bool {
	return name != ".svn" && name != ".hg" && name != ".git"
}

// Scans folder recursively search for files with the ext
// extension and returns the whole list.
func Scan(folder string, ext string) ([]string, error) {
	v := &visitor{[]string{}}
	if err := v.scan(folder, ext); err != nil {
		return nil, err
	}

	return v.results, nil
}