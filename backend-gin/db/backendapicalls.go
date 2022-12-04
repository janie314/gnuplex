package db

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func ScanDir(dir string) {
	filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			fmt.Println(path)
		}
		return err
	})
}
