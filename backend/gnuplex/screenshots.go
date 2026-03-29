package gnuplex

import (
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Screenshot struct {
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	ModifiedAt time.Time `json:"modified_at"`
}

func sanitizeScreenshotSource(filename string) string {
	sanitized := strings.ReplaceAll(filename, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	if sanitized == "" {
		return "unknown"
	}
	return sanitized
}

func buildScreenshotFilename(filename string, now time.Time) string {
	return now.Format("2006-01-02-15:04:05") + "." + sanitizeScreenshotSource(filename) + ".png"
}

func screenshotURL(name string) string {
	return "/screenshots/" + url.PathEscape(name)
}

func (gnuplex *GNUPlex) TakeScreenshot() (*Screenshot, error) {
	currentFilename, err := gnuplex.MPV.GetCurrentFilename()
	if err != nil {
		return nil, err
	}

	name := buildScreenshotFilename(currentFilename, time.Now())
	path := filepath.Join(gnuplex.ScreenshotsDir, name)
	if err := gnuplex.MPV.ScreenshotToFile(path); err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return &Screenshot{
		Name:       name,
		URL:        screenshotURL(name),
		ModifiedAt: info.ModTime(),
	}, nil
}

func (gnuplex *GNUPlex) ListRecentScreenshots(limit int) ([]Screenshot, error) {
	entries, err := os.ReadDir(gnuplex.ScreenshotsDir)
	if err != nil {
		return nil, err
	}

	screenshots := make([]Screenshot, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		screenshots = append(screenshots, Screenshot{
			Name:       entry.Name(),
			URL:        screenshotURL(entry.Name()),
			ModifiedAt: info.ModTime(),
		})
	}

	sort.Slice(screenshots, func(i, j int) bool {
		return screenshots[i].ModifiedAt.After(screenshots[j].ModifiedAt)
	})

	if limit > 0 && len(screenshots) > limit {
		screenshots = screenshots[:limit]
	}

	return screenshots, nil
}
