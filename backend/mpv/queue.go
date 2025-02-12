package mpv

import (
	"fmt"
)

func (mpv *MPV) Enqueue(path string) {
	mpv.Mu.Lock()
	defer mpv.Mu.Unlock()
	mpv.Queue = append(mpv.Queue, path)
}

// Delete removes an item from the queue by index.
func (mpv *MPV) Dequeue(i int) error {
	mpv.Mu.Lock()
	defer mpv.Mu.Unlock()
	if i < 0 || i >= len(mpv.Queue) {
		return fmt.Errorf("index out of range")
	}
	mpv.Queue = append(mpv.Queue[:i], mpv.Queue[i+1:]...)
	return nil
}

func (mpv *MPV) ClearQueue() {
	mpv.Mu.Lock()
	defer mpv.Mu.Unlock()
	mpv.Queue = make([]string, 0)
}
