package gnuplex

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
	var reterr error
	mediadirsFromDB, err := gnuplex.NewDB.GetMediaDirs()
	if err != nil {
		return err
	}
	mediaItemsFromDB, err := gnuplex.NewDB.GetMediaItems()
	if err != nil {
		return err
	}
	fileExts, err := gnuplex.NewDB.GetFileExts()
	if err != nil {
		return err
	}
	fileExtHash := make(map[string]bool)
	medialist := make(map[string](*dBFSDiff), len(mediadirsFromDB))
	for _, mediaItem := range mediaItemsFromDB {
		medialist[mediaItem.Path] = &dBFSDiff{
			inDB: true,
			inFS: false,
		}
	}
	for _, ext := range fileExts {
		fileExtHash[ext.Extension] = true
	}
	for _, mediadirFromDB := range mediadirsFromDB {
		dir, err := os.Stat(mediadirFromDB.Path)
		if (err == nil) && dir.IsDir() {
			err = filepath.WalkDir(mediadirFromDB.Path, func(path string, entry fs.DirEntry, err error) error {
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
						return gnuplex.NewDB.AddMediaItem(path)
					} else if medialist[path].inDB {
						medialist[path].inFS = true
						return nil
					} else {
						medialist[path] = &dBFSDiff{inDB: false, inFS: true}
						return gnuplex.NewDB.AddMediaItem(path)
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
	for _, mediaItem := range mediaItemsFromDB {
		if !medialist[mediaItem.Path].inFS {
			gnuplex.NewDB.DeleteMediaItem(&mediaItem)
		}
	}
	return reterr
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

func (gnuplex *GNUPlex) ReplaceQueueAndPlay(id models.MediaItemId) error {
	var mediaItem *models.MediaItem
	if err := gnuplex.NewDB.ORM.First(&mediaItem, id).Error; err != nil {
		return err
	}
	if mediaItem != nil {
		gnuplex.PlayQueue = []*models.MediaItem{mediaItem}
	}
	if err := gnuplex.MPV.ReplaceQueueAndPlay(mediaItem.Path); err != nil {
		return err
	}
	if err := gnuplex.NewDB.UpdateLastPlayed(mediaItem); err != nil {
		return err
	}
	return nil
}

func (gnuplex *GNUPlex) QueueLast(id models.MediaItemId) *models.MediaItem {
	var mediaItem *models.MediaItem
	gnuplex.NewDB.ORM.First(&mediaItem, id)
	if mediaItem != nil {
		gnuplex.PlayQueue = append(gnuplex.PlayQueue, mediaItem)
	}
	gnuplex.MPV.QueueMedia(mediaItem.Path)
	return mediaItem
}
