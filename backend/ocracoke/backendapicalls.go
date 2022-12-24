package ocracoke

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func (oc *Ocracoke) ScanLib() error {
	oc.DB.Mu.Lock()
	log.Println("Got ScanLib lock")
	defer oc.DB.Mu.Unlock()
	defer log.Println("Rem ScanLib lock")
	var reterr error
	mediadirs := oc.GetMediadirs(true)
	medialist := oc.GetMedialib(true)
	mediadirHash := make(map[string]bool, len(mediadirs))
	for _, path := range medialist {
		mediadirHash[path] = true
	}
	for _, mediadir := range mediadirs {
		dir, err := os.Stat(mediadir)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediadir, func(path string, entry fs.DirEntry, err error) error {
				if err == nil && (!entry.IsDir()) && (!mediadirHash[path]) {
					return oc.AddMedia(path, true)
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
			log.Println("Error: Bad mediadir: ", mediadir)
			reterr = err
		}
	}
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.srt';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.txt';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.jpg';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.jpeg';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.torrent';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.ico';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.docx';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.pdf';`)
	oc.DB.SqliteConn.Exec(`delete from medialist where filepath like '%.png';`)
	return reterr
}

func (oc *Ocracoke) AddHist(mediafile string) error {
	oc.DB.Mu.Lock()
	log.Println("Got AddHist lock")
	defer oc.DB.Mu.Unlock()
	defer log.Println("Rem AddHist lock")
	_, err := oc.DB.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddHist err", err)
	}
	return err
}

func (oc *Ocracoke) AddMedia(mediafile string, ignorelock bool) error {
	if !ignorelock {
		oc.DB.Mu.Lock()
		log.Println("Got AddMedia lock")
		defer oc.DB.Mu.Unlock()
		defer log.Println("Rem AddMedia lock")
	} else {
		log.Println("Ignoring AddMedia lock")
	}
	_, err := oc.DB.SqliteConn.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddMedia err", err)
	}
	return err
}

func (oc *Ocracoke) GetMediadirs(ignorelock bool) []string {
	if !ignorelock {
		oc.DB.Mu.Lock()
		log.Println("Got GetMediadirs lock")
		defer oc.DB.Mu.Unlock()
		defer log.Println("Rem GetMediadirs lock")
	} else {
		log.Println("Ignoring GetMediadirs lock")
	}
	rows, err := oc.DB.SqliteConn.Query("select filepath from mediadirs order by filepath;")
	if err != nil {
		log.Println("Error: GetMediadirs: ", err)
		return []string{}
	}
	res := make([]string, 10000)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Println("Error: GetMediadirs:", err)
		} else if i < len(res) {
			res[i] = str
			i++
		} else {
			res = append(res, str)
		}
	}
	return res[:i]
}

func (oc *Ocracoke) SetMediadirs(mediadirs []string) error {
	oc.DB.Mu.Lock()
	log.Println("Got SetMediadirs lock")
	defer oc.DB.Mu.Unlock()
	defer log.Println("Rem SetMediadirs lock")
	var err error
	oc.DB.SqliteConn.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := oc.DB.SqliteConn.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			log.Println("Error: AddMediadir", err)
		}
	}
	return err
}

func (oc *Ocracoke) GetMedialib(ignorelock bool) []string {
	if !ignorelock {
		oc.DB.Mu.Lock()
		log.Println("Got GetMedialib lock")
		defer oc.DB.Mu.Unlock()
		defer log.Println("Rem GetMedialib lock")
	} else {
		log.Println("Ignoring GetMedialib lock")
	}
	rows, err := oc.DB.SqliteConn.Query("select filepath from medialist order by filepath;")
	if err != nil {
		log.Println("Error: GetMedialib: ", err)
		return []string{}
	}
	res := make([]string, 131072)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Println("Error: GetMedialib:", err)
		} else if i < len(res) {
			res[i] = str
			i++
		} else {
			res = append(res, str)
		}
	}
	return res[:i]
}

func (oc *Ocracoke) Last25() []string {
	oc.DB.Mu.Lock()
	log.Println("Got Last25 lock")
	defer oc.DB.Mu.Unlock()
	defer fmt.Println("Rem Last25 lock")
	rows, err := oc.DB.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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
			fmt.Println("Error: Last25:", err)
		} else if i < len(res) {
			res[i] = str
			i++
		} else {
			res = append(res, str)
		}
	}
	return res[:i]
}
