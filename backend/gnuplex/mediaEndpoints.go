package gnuplex

import (
	"gnuplex/models"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

// Updates GNUPlex's stored MediaItems, using its configured MediaDirs.
func (gnuplex *GNUPlex) ScanLib() error {
	// Grab MediaDirs, FileExts from the database
	mediaDirs, err := gnuplex.DB.GetMediaDirs()
	if err != nil {
		return err
	}
	fileExts, err := gnuplex.DB.GetFileExts()
	if err != nil {
		return err
	}
	fileExtH := make(map[string]bool)
	for _, fileExt := range fileExts {
		fileExtH[fileExt.Extension] = true
	}
	// Add new stuff
	lastScanUUID := uuid.New().String()
	var batch []models.MediaItem
	for i, mediaDir := range mediaDirs {
		if (i%100 == 0) && (i != 0) {
			if err = gnuplex.processScanLibBatch(batch, lastScanUUID); err != nil {
				return err
			}
		}
		dir, err := os.Stat(mediaDir.Path)
		if (err == nil) && dir.IsDir() {
			if err = filepath.WalkDir(mediaDir.Path, func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				} else if !entry.IsDir() {
					ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
					if _, match := fileExtH[ext]; !match {
						batch = append(batch, models.MediaItem{Path: path, Type: models.File, LastScanUUID: lastScanUUID})
					}
					return nil
				} else {
					return nil
				}
			}); err != nil {
				return err
			}
		} else if err != nil {
			log.Println("skipping", dir, "- could not stat this directory")
		}
	}
	if len(batch) != 0 {
		if err = gnuplex.processScanLibBatch(batch, lastScanUUID); err != nil {
			return err
		}
	}
	return gnuplex.DB.DeleteMediaItemFilesNotMatchingUUID(lastScanUUID)
}

func (gnuplex *GNUPlex) processScanLibBatch(batch []models.MediaItem, lastScanUUID string) error {
	return gnuplex.DB.ORM.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "path"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"last_scan_uuid": lastScanUUID}),
		}).
		CreateInBatches(batch, 100).Error
}

// Returns the currently playing MediaItem
func (gnuplex *GNUPlex) GetNowPlaying() (*models.MediaItem, error) {
	path, err := gnuplex.MPV.GetCurrentFilepath()
	if err != nil {
		return nil, err
	}
	return gnuplex.DB.GetMediaItemByPath(path)
}

// Replace the current queue with a MediaItem from the library (by ID).
func (gnuplex *GNUPlex) ReplaceQueueAndPlay(id models.MediaItemId) error {
	var mediaItem *models.MediaItem
	if err := gnuplex.DB.ORM.First(&mediaItem, id).Error; err != nil {
		return err
	}
	if mediaItem != nil {
		gnuplex.PlayQueue = []*models.MediaItem{mediaItem}
	}
	if err := gnuplex.MPV.ReplaceQueueAndPlay(mediaItem.Path); err != nil {
		return err
	}
	if err := gnuplex.DB.UpdateLastPlayed(mediaItem); err != nil {
		return err
	}
	return nil
}

// Replace the current queue with a MediaItem from the library (by path/URL).
func (gnuplex *GNUPlex) ReplaceQueueAndPlayByPath(path string) error {
	var mediaItem *models.MediaItem
	if err := gnuplex.DB.ORM.First(&mediaItem, "path = ?", path).Error; err != nil {
		return err
	}
	if mediaItem != nil {
		gnuplex.PlayQueue = []*models.MediaItem{mediaItem}
	}
	if err := gnuplex.MPV.ReplaceQueueAndPlay(mediaItem.Path); err != nil {
		return err
	}
	if err := gnuplex.DB.UpdateLastPlayed(mediaItem); err != nil {
		return err
	}
	return nil
}

// Replace the current queue with a cast URL, without adding that URL to the library.
func (gnuplex *GNUPlex) ReplaceQueueAndCastTempUrl(url string) error {
	return gnuplex.MPV.ReplaceQueueAndPlay(url)
}

// Add a MediaItem to the end of the current queue.
func (gnuplex *GNUPlex) QueueLast(id models.MediaItemId) *models.MediaItem {
	var mediaItem *models.MediaItem
	gnuplex.DB.ORM.First(&mediaItem, id)
	if mediaItem != nil {
		gnuplex.PlayQueue = append(gnuplex.PlayQueue, mediaItem)
	}
	gnuplex.MPV.QueueMedia(mediaItem.Path)
	return mediaItem
}

// Cast a URL to the media player. `temp` determines whether or not it should be added to your library.
func (gnuplex *GNUPlex) Cast(url string, temp bool) error {
	if temp {
		return gnuplex.ReplaceQueueAndCastTempUrl(url)
	} else if err := gnuplex.DB.AddMediaItemURL(url); err != nil {
		return err
	} else {
		return gnuplex.ReplaceQueueAndPlayByPath(url)
	}
}

// Cycle subtitle track.
func (gnuplex *GNUPlex) GetSubs() ([]models.Track, error) {
	tracks, err := gnuplex.MPV.GetTracks()
	if err != nil {
		return nil, err
	}
	var res []models.Track
	for _, track := range tracks {
		if track.Type == "sub" {
			res = append(res, track)
		}
	}
	return res, nil
}

// Set subtitle visibility.
func (gnuplex *GNUPlex) SetSubVisibility(visible bool) error {
	return gnuplex.MPV.SetSubVisibility(visible)
}

// Set subtitle track.
func (gnuplex *GNUPlex) SetSubTrack(trackID int64) error {
	return gnuplex.MPV.SetSubTrack(trackID)
}
