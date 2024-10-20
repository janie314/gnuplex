package server

import (
	"errors"
	"fmt"
	"gnuplex-backend/models"
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

func (gnuplex *GNUPlex) ScanLib() error {
	gnuplex.DB.Mu.Lock()
	log.Println("Got ScanLib lock")
	defer gnuplex.DB.Mu.Unlock()
	defer log.Println("Rem ScanLib lock")
	var reterr error
	mediadirsFromDB := gnuplex.GetMediadirs(true)
	medialistFromDB := gnuplex.GetMedialib(true)
	fileExts := gnuplex.GetFileExts(true)
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
						return gnuplex.AddMedia(path, true)
					} else if medialist[path].inDB {
						medialist[path].inFS = true
						return nil
					} else {
						medialist[path] = &dBFSDiff{inDB: false, inFS: true}
						return gnuplex.AddMedia(path, true)
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
			gnuplex.DB.SqliteConn.Exec(`delete from medialist where filepath = ?;`, path)
		}
	}
	return reterr
}

func (gnuplex *GNUPlex) AddHist(mediafile string) error {
	gnuplex.DB.Mu.Lock()
	log.Println("Got AddHist lock")
	defer gnuplex.DB.Mu.Unlock()
	defer log.Println("Rem AddHist lock")
	_, err := gnuplex.DB.SqliteConn.Exec("insert into history (mediafile) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddHist err", err)
	}
	return err
}

func (gnuplex *GNUPlex) AddMedia(mediafile string, ignorelock bool) error {
	if !ignorelock {
		gnuplex.DB.Mu.Lock()
		log.Println("Got AddMedia lock")
		defer gnuplex.DB.Mu.Unlock()
		defer log.Println("Rem AddMedia lock")
	} else {
		log.Println("Ignoring AddMedia lock")
	}
	_, err := gnuplex.DB.SqliteConn.Exec("insert or replace into medialist (filepath) values (?);", mediafile)
	if err != nil {
		log.Println("Error: AddMedia err", err)
	}
	return err
}

func (gnuplex *GNUPlex) GetMediadirs(ignorelock bool) []string {
	if !ignorelock {
		gnuplex.DB.Mu.Lock()
		log.Println("Got GetMediadirs lock")
		defer gnuplex.DB.Mu.Unlock()
		defer log.Println("Rem GetMediadirs lock")
	} else {
		log.Println("Ignoring GetMediadirs lock")
	}
	rows, err := gnuplex.DB.SqliteConn.Query("select filepath from mediadirs order by filepath;")
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

func (gnuplex *GNUPlex) SetMediadirs(mediadirs []string) error {
	gnuplex.DB.Mu.Lock()
	log.Println("Got SetMediadirs lock")
	defer gnuplex.DB.Mu.Unlock()
	defer log.Println("Rem SetMediadirs lock")
	var err error
	gnuplex.DB.SqliteConn.Exec("delete from mediadirs;")
	for _, mediafile := range mediadirs {
		_, err := gnuplex.DB.SqliteConn.Exec("insert or ignore into mediadirs (filepath) values (?);", mediafile)
		if err != nil {
			log.Println("Error: AddMediadir", err)
		}
	}
	return err
}

func (gnuplex *GNUPlex) GetFileExts(ignorelock bool) []string {
	if !ignorelock {
		gnuplex.DB.Mu.Lock()
		log.Println("Got GetFileExts lock")
		defer gnuplex.DB.Mu.Unlock()
		defer log.Println("Rem GetFileExtslock")
	} else {
		log.Println("Ignoring GetFileExts lock")
	}
	rows, err := gnuplex.DB.SqliteConn.Query("select (ext) from file_exts order by ext ;")
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

func (gnuplex *GNUPlex) SetFileExts(file_exts []string) error {
	gnuplex.DB.Mu.Lock()
	log.Println("Got SetFileExtslock")
	defer gnuplex.DB.Mu.Unlock()
	defer log.Println("Rem SetFileExtslock")
	var err error
	gnuplex.DB.SqliteConn.Exec("delete from file_exts;")
	for _, ext := range file_exts {
		_, err := gnuplex.DB.SqliteConn.Exec("insert or ignore into file_exts (ext, exclude) values (?, 1);", strings.ToLower(ext))
		if err != nil {
			log.Println("Error: SetFileExts", err)
		}
	}
	return err
}

func (gnuplex *GNUPlex) GetMedialib(ignorelock bool) []string {
	if !ignorelock {
		gnuplex.DB.Mu.Lock()
		log.Println("Got GetMedialib lock")
		defer gnuplex.DB.Mu.Unlock()
		defer log.Println("Rem GetMedialib lock")
	} else {
		log.Println("Ignoring GetMedialib lock")
	}
	rows, err := gnuplex.DB.SqliteConn.Query("select filepath from medialist order by filepath;")
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

func (gnuplex *GNUPlex) Last25() []string {
	gnuplex.DB.Mu.Lock()
	log.Println("Got Last25 lock")
	defer gnuplex.DB.Mu.Unlock()
	defer fmt.Println("Rem Last25 lock")
	rows, err := gnuplex.DB.SqliteConn.Query("select distinct mediafile from history order by id desc limit 25;")
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

func (gnuplex *GNUPlex) NowPlaying() (*models.MediaItem, error) {
	if len(gnuplex.PlayQueue) == 0 {
		return nil, errors.New("PlayQueue is empty at the moment")
	}
	return gnuplex.PlayQueue[0], nil
}

func (gnuplex *GNUPlex) Queue(id models.MediaItemId) *models.MediaItem {
	var mediaItem *models.MediaItem
	gnuplex.NewDB.First(&mediaItem, id)
	if mediaItem != nil {
		gnuplex.PlayQueue = append(gnuplex.PlayQueue, mediaItem)
	}
	return mediaItem
}
