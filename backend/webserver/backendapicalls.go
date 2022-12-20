package webserver

import (
	"fmt"
	"gnuplex-backend/liteDB"
	"io/fs"
	"os"
	"path/filepath"
)

func ScanLib(db *liteDB.LiteDB) error {
	db.Mu.Lock()
	fmt.Println("got a lock")
	defer db.Mu.Unlock()
	defer fmt.Println("rem a lock")
	var reterr error
	mediadirs := GetMediadirs(db, true)
	for _, mediadir := range mediadirs {
		dir, err := os.Stat(mediadir)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediadir, func(path string, entry fs.DirEntry, err error) error {
				if err == nil && (!entry.IsDir()) {
					return AddMedia(db, path, true)
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
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.srt';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.txt';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.jpg';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.jpeg';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.torrent';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.ico';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.docx';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.pdf';`)
	db.SqliteConn.Exec(`delete from medialist where filepath like '%.png';`)
	return reterr
}

func AddHist(liteDB *liteDB.LiteDB, mediafile string) error {
	liteDB.Mu.Lock()
	fmt.Println("got b lock")
	defer liteDB.Mu.Unlock()
	defer fmt.Println("rem b lock")
	_, err := liteDB.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "AddHist err", err)
	}
	return err
}

func AddMedia(db *liteDB.LiteDB, mediafile string, ignorelock bool) error {
	if !ignorelock {
		db.Mu.Lock()
		fmt.Println("got c lock")
		defer db.Mu.Unlock()
		defer fmt.Println("rem c lock")
	} else {
		fmt.Println("ignoring c lock")
	}
	_, err := db.SqliteConn.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "AddMedia err", err)
	}
	return err
}

func GetMediadirs(db *liteDB.LiteDB, ignorelock bool) []string {
	if !ignorelock {
		db.Mu.Lock()
		fmt.Println("got d lock")
		defer db.Mu.Unlock()
		defer fmt.Println("rem d lock")
	} else {
		fmt.Println("ignoring d lock")
	}
	rows, err := db.SqliteConn.Query("select filepath from mediadirs order by filepath;")
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

func SetMediadirs(db *liteDB.LiteDB, mediadirs []string) error {
	db.Mu.Lock()
	fmt.Println("got e lock")
	defer db.Mu.Unlock()
	defer fmt.Println("rem e lock")
	var err error
	db.SqliteConn.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := db.SqliteConn.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "AddMediadir err", err)
		}
	}
	return err
}

func GetMedialib(db *liteDB.LiteDB) []string {
	db.Mu.Lock()
	fmt.Println("got f lock")
	defer db.Mu.Unlock()
	defer fmt.Println("rem f lock")
	rows, err := db.SqliteConn.Query("select filepath from medialist order by filepath;")
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

func Last25(db *liteDB.LiteDB) []string {
	db.Mu.Lock()
	fmt.Println("got g lock")
	defer db.Mu.Unlock()
	defer fmt.Println("rem g lock")
	rows, err := db.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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
