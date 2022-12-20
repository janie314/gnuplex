package sqliteconn

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func ScanLib(db *sql.DB, mu *sync.Mutex) error {
	mu.Lock()
	fmt.Println("got a lock")
	defer mu.Unlock()
	defer fmt.Println("rem a lock")
	var reterr error
	mediadirs := GetMediadirs(db, mu, true)
	for _, mediadir := range mediadirs {
		dir, err := os.Stat(mediadir)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediadir, func(path string, entry fs.DirEntry, err error) error {
				if err == nil && (!entry.IsDir()) {
					return AddMedia(db, mu, path, true)
				} else if err != nil {
					fmt.Fprintln(os.Stderr, "Walkdir prob: ", mediadir)
					return err
				}
				return nil
			})
			if err != nil {
				reterr = err
			}
		} else {
			fmt.Fprintln(os.Stderr, "Bad mediadir: ", mediadir)
			reterr = err
		}
	}
	db.Exec(`delete from medialist where filepath like '%.srt';`)
	db.Exec(`delete from medialist where filepath like '%.txt';`)
	db.Exec(`delete from medialist where filepath like '%.jpg';`)
	db.Exec(`delete from medialist where filepath like '%.docx';`)
	db.Exec(`delete from medialist where filepath like '%.pdf';`)
	return reterr
}

func AddHist(db *sql.DB, mu *sync.Mutex, mediafile string) error {
	mu.Lock()
	fmt.Println("got b lock")
	defer mu.Unlock()
	defer fmt.Println("rem b lock")
	_, err := db.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "AddHist err", err)
	}
	return err
}

func AddMedia(db *sql.DB, mu *sync.Mutex, mediafile string, ignorelock bool) error {
	if !ignorelock {
		mu.Lock()
		fmt.Println("got c lock")
		defer mu.Unlock()
		defer fmt.Println("rem c lock")
	} else {
		fmt.Println("ignoring c lock")
	}
	_, err := db.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "AddMedia err", err)
	}
	return err
}

func GetMediadirs(db *sql.DB, mu *sync.Mutex, ignorelock bool) []string {
	if !ignorelock {
		mu.Lock()
		fmt.Println("got d lock")
		defer mu.Unlock()
		defer fmt.Println("rem d lock")
	} else {
		fmt.Println("ignoring d lock")
	}
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

func SetMediadirs(db *sql.DB, mu *sync.Mutex, mediadirs []string) error {
	mu.Lock()
	fmt.Println("got e lock")
	defer mu.Unlock()
	defer fmt.Println("rem e lock")
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

func GetMedialib(db *sql.DB, mu *sync.Mutex) []string {
	mu.Lock()
	fmt.Println("got f lock")
	defer mu.Unlock()
	defer fmt.Println("rem f lock")
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

func Last25(db *sql.DB, mu *sync.Mutex) []string {
	mu.Lock()
	fmt.Println("got g lock")
	defer mu.Unlock()
	defer fmt.Println("rem g lock")
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
