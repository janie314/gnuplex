package gnuplex

import (
	"time"

	"github.com/asticode/go-astisub"
)

type MPVCommander interface {
	GetSubDelay() (float64, error)
	SetSubDelay(delay float64) error
	GetCurrentSubFilename() (string, error)
	SubReload() error
}

type GNUPlexForTesting struct {
	GNUPlex
	MPV *MockMPV
}

func (g *GNUPlexForTesting) SaveSubDelay() error {
	delay, err := g.MPV.GetSubDelay()
	if err != nil {
		return err
	}
	if delay == 0 {
		return nil
	}

	filename, err := g.MPV.GetCurrentSubFilename()
	if err != nil {
		return err
	}
	if filename == "" {
		return nil
	}

	subs, err := astisub.OpenFile(filename)
	if err != nil {
		return err
	}

	subs.Add(time.Duration(delay * 1e9))

	if err := subs.Write(filename); err != nil {
		return err
	}

	if err := g.MPV.SubReload(); err != nil {
		return err
	}

	return g.MPV.SetSubDelay(0)
}

type MockMPV struct {
	SubDelayValue         float64
	SubDelayError         error
	SubFilenameValue      string
	SubFilenameError      error
	SubReloadError        error
	SetSubDelayError      error
	SetSubDelayCalledWith float64
	SubReloadCalled       bool
}

func (m *MockMPV) GetSubDelay() (float64, error) {
	return m.SubDelayValue, m.SubDelayError
}

func (m *MockMPV) SetSubDelay(delay float64) error {
	m.SetSubDelayCalledWith = delay
	return m.SetSubDelayError
}

func (m *MockMPV) GetCurrentSubFilename() (string, error) {
	return m.SubFilenameValue, m.SubFilenameError
}

func (m *MockMPV) SubReload() error {
	m.SubReloadCalled = true
	return m.SubReloadError
}
