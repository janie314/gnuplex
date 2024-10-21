package gnuplex

import (
	"errors"
	"gnuplex-backend/models"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (gnuplex *GNUPlex) ScanLib() error {
	/*
	 * grab MediaDirs, MediaItems, FileExts from the database
	 */
	mediaDirs, err := gnuplex.NewDB.GetMediaDirs()
	if err != nil {
		return err
	}
	mediaItems, err := gnuplex.NewDB.GetMediaItems()
	if err != nil {
		return err
	}
	mediaItemH := make(map[string]models.MediaItem)
	for _, mediaItem := range mediaItems {
		mediaItemH[mediaItem.Path] = mediaItem
	}
	fileExts, err := gnuplex.NewDB.GetFileExts()
	if err != nil {
		return err
	}
	fileExtH := make(map[string]bool)
	for _, fileExt := range fileExts {
		fileExtH[fileExt.Extension] = true
	}
	/*
	 * add new stuff
	 */
	for _, mediaDir := range mediaDirs {
		dir, err := os.Stat(mediaDir.Path)
		if (err == nil) && dir.IsDir() {
			return filepath.WalkDir(mediaDir.Path, func(path string, entry fs.DirEntry, err error) error {
				log.Println("path", path)
				if err != nil {
					log.Println("Walkdir prob: ", mediaDir)
					return err
				} else if !entry.IsDir() {
					ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
					log.Println("ext", ext)
					if _, ok := fileExtH[ext]; ok {
						return gnuplex.NewDB.DeleteMediaItemByPath(path)
					}
					return gnuplex.NewDB.AddMediaItem(path)
				} else {
					return nil
				}
			})

		}
	}
	/*
	 * remove stuff that no longer exists
	 */
	for _, mediaItem := range mediaItems {
		_, err := os.Stat(mediaItem.Path)
		if err != nil {
			err = gnuplex.NewDB.DeleteMediaItemByPath(mediaItem.Path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gnuplex *GNUPlex) SetFileExts(file_exts []string) error {
	// TODO
	return nil
}

func (gnuplex *GNUPlex) Last25() []string {
	// TODO return nil
	return nil
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
