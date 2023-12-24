package server

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type dBFSDiff struct {
	inDB bool
	inFS bool
}

func (server *Server) ScanLib() error {
	server.DB.Mu.Lock()
	log.Println("Got ScanLib lock")
	defer server.DB.Mu.Unlock()
	defer log.Println("Rem ScanLib lock")
	var reterr error
	mediadirsFromDB := server.GetMediadirs(true)
	medialistFromDB := server.GetMedialib(true)
	fileExts := server.GetFileExts(true)
	fileExtHash := make(map[string]bool)
	medialist := make(map[string](*dBFSDiff), len(mediadirsFromDB))
	for _, path := range medialistFromDB {
		medialist[path] = &dBFSDiff{
			inDB: true,
			inFS: false,
		}
	}
	for _, ext := range fileExts {
		fileExtHash[ext] = true
	}
	for _, mediadirFromDB := range mediadirsFromDB {
		dir, err := os.Stat(mediadirFromDB)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediadirFromDB, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					fmt.Fprintln(os.Stderr, "Walkdir prob: ", mediadirFromDB)
					return err
				} else if !entry.IsDir() {
					pathLC := strings.ToLower(path)
					ext := pathLC[strings.LastIndex(path, ".")+1:]
					if fileExtHash[ext] || fileExtHash["."+ext] {
						return nil
					} else if medialist[path] == nil {
						medialist[path] = &dBFSDiff{inDB: false, inFS: true}
						return server.AddMedia(path, true)
					} else if medialist[path].inDB {
						medialist[path].inFS = true
						return nil
					} else {
						medialist[path] = &dBFSDiff{inDB: false, inFS: true}
						return server.AddMedia(path, true)
					}
				} else {
					return nil
				}
			})
			if err != nil {
				reterr = err
			}
		} else {
			log.Println("Error: Bad mediadir: ", mediadirFromDB)
			reterr = err
		}
	}
	for _, path := range medialistFromDB {
		if !medialist[path].inFS {
			server.DB.SqliteConn.Exec(`delete from medialist where filepath = ?;`, path)
		}
	}
	return reterr
}

func (server *Server) AddHist(mediafile string) error {
	server.DB.Mu.Lock()
	log.Println("Got AddHist lock")
	defer server.DB.Mu.Unlock()
	defer log.Println("Rem AddHist lock")
	_, err := server.DB.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddHist err", err)
	}
	return err
}

func (server *Server) AddMedia(mediafile string, ignorelock bool) error {
	if !ignorelock {
		server.DB.Mu.Lock()
		log.Println("Got AddMedia lock")
		defer server.DB.Mu.Unlock()
		defer log.Println("Rem AddMedia lock")
	} else {
		log.Println("Ignoring AddMedia lock")
	}
	_, err := server.DB.SqliteConn.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddMedia err", err)
	}
	return err
}

func (server *Server) GetMediadirs(ignorelock bool) []string {
	if !ignorelock {
		server.DB.Mu.Lock()
		log.Println("Got GetMediadirs lock")
		defer server.DB.Mu.Unlock()
		defer log.Println("Rem GetMediadirs lock")
	} else {
		log.Println("Ignoring GetMediadirs lock")
	}
	rows, err := server.DB.SqliteConn.Query("select filepath from mediadirs order by filepath;")
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

func (server *Server) SetMediadirs(mediadirs []string) error {
	server.DB.Mu.Lock()
	log.Println("Got SetMediadirs lock")
	defer server.DB.Mu.Unlock()
	defer log.Println("Rem SetMediadirs lock")
	var err error
	server.DB.SqliteConn.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := server.DB.SqliteConn.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			log.Println("Error: AddMediadir", err)
		}
	}
	return err
}

func (server *Server) GetFileExts(ignorelock bool) []string {
	if !ignorelock {
		server.DB.Mu.Lock()
		log.Println("Got GetFileExts lock")
		defer server.DB.Mu.Unlock()
		defer log.Println("Rem GetFileExtslock")
	} else {
		log.Println("Ignoring GetFileExts lock")
	}
	rows, err := server.DB.SqliteConn.Query("select (ext) from file_exts order by ext ;")
	if err != nil {
		log.Println("Error: GetFileExts: ", err)
		return []string{}
	}
	res := make([]string, 10000)
	str := ""
	i := 0
	for rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			fmt.Println("Error: GetFileExts:", err)
		} else if i < len(res) {
			res[i] = str
			i++
		} else {
			res = append(res, str)
		}
	}
	return res[:i]
}

func (oc *Server) SetFileExts(file_exts []string) error {
	oc.DB.Mu.Lock()
	log.Println("Got SetFileExtslock")
	defer oc.DB.Mu.Unlock()
	defer log.Println("Rem SetFileExtslock")
	var err error
	oc.DB.SqliteConn.Exec("delete from file_exts;")
	for _, ext := range file_exts {
		_, err := oc.DB.SqliteConn.Exec("insert or ignore into file_exts (ext, exclude) values (?, 1);", strings.ToLower(ext))
		if err != nil {
			log.Println("Error: SetFileExts", err)
		}
	}
	return err
}

func (oc *Server) GetMedialib(ignorelock bool) []string {
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

func (oc *Server) Last25() []string {
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
