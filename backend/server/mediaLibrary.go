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

func (srv *Server) ScanLib(ignorelock bool) error {
	srv.DB.Lock("ScanLib", ignorelock)
	defer srv.DB.Unlock("ScanLib", ignorelock)
	var reterr error
	mediadirsFromDB := srv.GetMediadirs(true)
	medialistFromDB := srv.GetMedialib(true)
	fileExts := srv.GetFileExts(true)
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
						return srv.AddMedia(path, true)
					} else if medialist[path].inDB {
						medialist[path].inFS = true
						return nil
					} else {
						medialist[path] = &dBFSDiff{inDB: false, inFS: true}
						return srv.AddMedia(path, true)
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
			srv.DB.SqliteConn.Exec(`delete from medialist where filepath = ?;`, path)
		}
	}
	return reterr
}

func (srv *Server) GetMediadirs(ignorelock bool) []string {
	srv.DB.Lock("GetMediadirs", ignorelock)
	defer srv.DB.Unlock("GetMediadirs", ignorelock)
	rows, err := srv.DB.SqliteConn.Query("select filepath from mediadirs order by filepath;")
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

func (srv *Server) SetMediadirs(mediadirs []string, ignorelock bool) error {
	srv.DB.Lock("SetMediadirs", ignorelock)
	defer srv.DB.Unlock("SetMediadirs", ignorelock)
	var err error
	srv.DB.SqliteConn.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := srv.DB.SqliteConn.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			log.Println("Error: AddMediadir", err)
		}
	}
	return err
}

func (srv *Server) GetFileExts(ignorelock bool) []string {
	srv.DB.Lock("GetFileExts", ignorelock)
	defer srv.DB.Unlock("GetFileExts", ignorelock)
	rows, err := srv.DB.SqliteConn.Query("select (ext) from file_exts order by ext ;")
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

func (srv *Server) SetFileExts(file_exts []string, ignorelock bool) error {
	srv.DB.Lock("SetFileExts", ignorelock)
	defer srv.DB.Unlock("SetFileExts", ignorelock)
	var err error
	srv.DB.SqliteConn.Exec("delete from file_exts;")
	for _, ext := range file_exts {
		_, err := srv.DB.SqliteConn.Exec("insert or ignore into file_exts (ext, exclude) values (?, 1);", strings.ToLower(ext))
		if err != nil {
			log.Println("Error: SetFileExts", err)
		}
	}
	return err
}

func (srv *Server) GetMedialib(ignorelock bool) []string {
	srv.DB.Lock("GetMedialib", ignorelock)
	defer srv.DB.Unlock("GetMedialib", ignorelock)
	rows, err := srv.DB.SqliteConn.Query("select filepath from medialist order by filepath;")
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

func (srv *Server) AddMedia(mediafile string, ignorelock bool) error {
	srv.DB.Lock("AddMedia", ignorelock)
	defer srv.DB.Unlock("AddMedia", ignorelock)
	_, err := srv.DB.SqliteConn.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddMedia err", err)
	}
	return err
}
