package sqliteconn

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
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

func AddHist(db *sql.DB, mediafile string) {
	_, err := db.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
