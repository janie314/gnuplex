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
		fmt.Fprintln(os.Stderr, "AddHist err", err)
	}
	return err
}

func AddMedia(db *sql.DB, mediafile string) error {
	_, err := db.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "AddMedia err", err)
	}
	return err
}

func GetMediadirs(db *sql.DB) []string {
	rows, err := db.Query("select filepath from mediadirs order by filepath;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return []string{}
	}
	// TODO: append or [i]
	res := make([]string, 10000)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			res[i] = str
			i++
		}
	}
	return res[:i]
}

func SetMediadirs(db *sql.DB, mediadirs []string) error {
	var err error
	db.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := db.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "AddMediadir err", err)
		}
	}
	return err
}

func GetMedialib(db *sql.DB) []string {
	rows, err := db.Query("select filepath from medialist order by filepath;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return []string{}
	}
	// TODO: append or [i]
	res := make([]string, 131072)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			res[i] = str
			i++
		}
	}
	return res[:i]
}

func Last25(db *sql.DB) []string {
	rows, err := db.Query("select distinct mediafile from history order by id desc limit 25;")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return []string{}
	}
	res := make([]string, 16384)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			res[i] = str
			i++
		}
	}
	return res[:i]
}
