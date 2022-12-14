package sqliteconn

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func ScanLib(db *sql.DB) error {
	var err error
	var dir string
	rows, err := db.Query("select distinct filepath from mediadirs;")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Scanlib query problem")
		return err
	}
	defer db.Exec(`delete from medialist where filepath like '%.srt';`)
	defer db.Exec(`delete from medialist where filepath like '%.txt';`)
	defer db.Exec(`delete from medialist where filepath like '%.jpg';`)
	defer db.Exec(`delete from medialist where filepath like '%.docx';`)
	defer db.Exec(`delete from medialist where filepath like '%.pdf';`)
	for rows.Next() {
		err := rows.Scan(&dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			defer filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
				if !entry.IsDir() {
					return AddMedia(db, path)
				} else {
					return nil
				}
			})
		}
	}
	return err
}

func AddHist(db *sql.DB, mediafile string) error {
	_, err := db.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return err
}

func AddMedia(db *sql.DB, mediafile string) error {
	_, err := db.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return err
}
