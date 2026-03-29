package gnuplex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/asticode/go-astisub"
)

func TestSaveSubDelay_DelayZero(t *testing.T) {
	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue: 0,
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSaveSubDelay_NoSubFilename(t *testing.T) {
	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    1.5,
			SubFilenameValue: "",
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSaveSubDelay_GetSubDelayError(t *testing.T) {
	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayError:    os.ErrNotExist,
			SubDelayValue:    1.5,
			SubFilenameValue: "test.srt",
		},
	}

	err := g.SaveSubDelay()
	if err != os.ErrNotExist {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestSaveSubDelay_GetSubFilenameError(t *testing.T) {
	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    1.5,
			SubFilenameValue: "test.srt",
			SubFilenameError: os.ErrNotExist,
		},
	}

	err := g.SaveSubDelay()
	if err != os.ErrNotExist {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestSaveSubDelay_Success(t *testing.T) {
	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")

	subtitles := astisub.NewSubtitles()
	subtitles.Items = append(subtitles.Items, &astisub.Item{
		StartAt: 3000000000,
		EndAt:   10000000000,
		Lines:   []astisub.Line{{Items: []astisub.LineItem{{Text: "Hello"}}}},
	})
	if err := subtitles.Write(srtPath); err != nil {
		t.Fatalf("failed to create test SRT: %v", err)
	}

	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    2.0,
			SubFilenameValue: srtPath,
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if g.MPV.SetSubDelayCalledWith != 0 {
		t.Errorf("expected SetSubDelay(0), got SetSubDelay(%v)", g.MPV.SetSubDelayCalledWith)
	}

	if !g.MPV.SubReloadCalled {
		t.Error("expected SubReload to be called")
	}

	subs, err := astisub.OpenFile(srtPath)
	if err != nil {
		t.Fatalf("failed to reopen SRT: %v", err)
	}

	expectedStart := int64(5000000000)
	if subs.Items[0].StartAt != time.Duration(expectedStart) {
		t.Errorf("expected subtitle start at %d, got %d", expectedStart, subs.Items[0].StartAt)
	}
}

func TestSaveSubDelay_NegativeDelay(t *testing.T) {
	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")

	subtitles := astisub.NewSubtitles()
	subtitles.Items = append(subtitles.Items, &astisub.Item{
		StartAt: 10000000000,
		EndAt:   20000000000,
		Lines:   []astisub.Line{{Items: []astisub.LineItem{{Text: "Hello"}}}},
	})
	if err := subtitles.Write(srtPath); err != nil {
		t.Fatalf("failed to create test SRT: %v", err)
	}

	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    -1.5,
			SubFilenameValue: srtPath,
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	subs, err := astisub.OpenFile(srtPath)
	if err != nil {
		t.Fatalf("failed to reopen SRT: %v", err)
	}

	expectedStart := int64(8500000000)
	if subs.Items[0].StartAt != time.Duration(expectedStart) {
		t.Errorf("expected subtitle start at %d, got %d", expectedStart, subs.Items[0].StartAt)
	}
}

func TestSaveSubDelay_SubReloadError(t *testing.T) {
	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "test.srt")

	subtitles := astisub.NewSubtitles()
	subtitles.Items = append(subtitles.Items, &astisub.Item{
		StartAt: 0,
		EndAt:   10000000000,
		Lines:   []astisub.Line{{Items: []astisub.LineItem{{Text: "Hello"}}}},
	})
	if err := subtitles.Write(srtPath); err != nil {
		t.Fatalf("failed to create test SRT: %v", err)
	}

	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    1.0,
			SubFilenameValue: srtPath,
			SubReloadError:   os.ErrPermission,
		},
	}

	err := g.SaveSubDelay()
	if err != os.ErrPermission {
		t.Errorf("expected os.ErrPermission, got %v", err)
	}
}

func copySRT(t *testing.T, src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("failed to read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", dst, err)
	}
}

func TestSaveSubDelay_CasablancaShiftPlus5(t *testing.T) {
	src := "testdata/casablanca.srt"
	if _, err := os.Stat(src); os.IsNotExist(err) {
		t.Skip("testdata/casablanca.srt not found")
	}

	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "casablanca.srt")
	copySRT(t, src, srtPath)

	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    5.0,
			SubFilenameValue: srtPath,
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Fatalf("SaveSubDelay failed: %v", err)
	}

	subs, err := astisub.OpenFile(srtPath)
	if err != nil {
		t.Fatalf("failed to open SRT: %v", err)
	}

	firstStart := subs.Items[0].StartAt
	originalStart := 71*time.Second + 570*time.Millisecond
	expectedStart := originalStart + 5*time.Second

	if firstStart != expectedStart {
		t.Errorf("expected first subtitle start at %v, got %v", expectedStart, firstStart)
	}

	t.Logf("First subtitle shifted from %v to %v", originalStart, firstStart)
}

func TestSaveSubDelay_CasablancaShiftMinus5(t *testing.T) {
	src := "testdata/casablanca.srt"
	if _, err := os.Stat(src); os.IsNotExist(err) {
		t.Skip("testdata/casablanca.srt not found")
	}

	tmpDir := t.TempDir()
	srtPath := filepath.Join(tmpDir, "casablanca.srt")
	copySRT(t, src, srtPath)

	g := &GNUPlexForTesting{
		MPV: &MockMPV{
			SubDelayValue:    -5.0,
			SubFilenameValue: srtPath,
		},
	}

	err := g.SaveSubDelay()
	if err != nil {
		t.Fatalf("SaveSubDelay failed: %v", err)
	}

	subs, err := astisub.OpenFile(srtPath)
	if err != nil {
		t.Fatalf("failed to open SRT: %v", err)
	}

	firstStart := subs.Items[0].StartAt
	originalStart := 71*time.Second + 570*time.Millisecond
	expectedStart := originalStart - 5*time.Second

	if firstStart != expectedStart {
		t.Errorf("expected first subtitle start at %v, got %v", expectedStart, firstStart)
	}

	t.Logf("First subtitle shifted from %v to %v", originalStart, firstStart)
}

func TestBuildScreenshotFilename(t *testing.T) {
	now := time.Date(2026, 3, 29, 14, 5, 6, 0, time.UTC)
	got := buildScreenshotFilename("/media/movies/Alien (1979).mkv", now)
	want := "2026-03-29-14:05:06._media_movies_Alien (1979).mkv.png"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestListRecentScreenshots(t *testing.T) {
	tmpDir := t.TempDir()
	g := &GNUPlex{ScreenshotsDir: tmpDir}

	olderName := "older.png"
	newerName := "newer.png"
	notScreenshot := "ignore.txt"

	for _, name := range []string{olderName, newerName, notScreenshot} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("x"), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
	}

	oldTime := time.Date(2026, 3, 29, 12, 0, 0, 0, time.UTC)
	newTime := time.Date(2026, 3, 29, 13, 0, 0, 0, time.UTC)
	if err := os.Chtimes(filepath.Join(tmpDir, olderName), oldTime, oldTime); err != nil {
		t.Fatalf("failed to set mtime on older screenshot: %v", err)
	}
	if err := os.Chtimes(filepath.Join(tmpDir, newerName), newTime, newTime); err != nil {
		t.Fatalf("failed to set mtime on newer screenshot: %v", err)
	}

	res, err := g.ListRecentScreenshots(10)
	if err != nil {
		t.Fatalf("ListRecentScreenshots failed: %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 screenshots, got %d", len(res))
	}
	if res[0].Name != newerName || res[1].Name != olderName {
		t.Fatalf("unexpected screenshot order: %+v", res)
	}
	if !strings.HasPrefix(res[0].URL, "/screenshots/") {
		t.Fatalf("expected screenshot URL to use /screenshots/, got %q", res[0].URL)
	}
}
