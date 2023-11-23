package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Walk a directory tree returning a of files with a given extension
func Walk(rootpath string, filter string) []string {

	var results []string

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {

		// Apple cache data is never useful
		basename := filepath.Base(path)
		if strings.HasPrefix(basename, ".") {
			return nil
		}

		if filepath.Ext(path) == filter {
			println(".. ", len(results), " - ", path)
			results = append(results, path)
		}

		return nil // No error...
	})
	if err != nil {
		fmt.Printf("Walk error [%v]\n", err)
	}

	return results

}
