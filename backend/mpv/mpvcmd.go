package mpv

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gnuplex/models"
	"log"
	"os"
	"strings"
)

/*
 * Types
 */
type IMPVQuery struct {
	Command []interface{} `json:"command"`
}

type IMPVQueryString struct {
	Command []string `json:"command"`
}

type IMPVResponseBool struct {
	Data bool `json:"data"`
}

type IMPVResponseString struct {
	Data string `json:"data"`
}

type IMPVResponseInt struct {
	Data int `json:"data"`
}

type MPVGetResult[T string | int | float64 | []models.Track] struct {
	Data      T      `json:"data"`
	RequestId int    `json:"request_id"`
	Error     string `json:"error"`
}

type MPVSetResult struct {
	RequestId int    `json:"request_id"`
	Error     string `json:"error"`
}

func processMPVGetResult[T string | int | float64 | []models.Track](resBytes []byte) (T, error) {
	var res MPVGetResult[T]
	var defaultVal T
	err := json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("mpv result error", err)
		return defaultVal, err
	} else if res.Error != "success" {
		log.Println("mpv result error", err)
		return defaultVal, errors.New(res.Error)
	}
	return res.Data, nil

}

func processMPVSetResult(resBytes []byte) error {
	var res MPVSetResult
	err := json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("mpv result error", err)
		return err
	} else if res.Error != "success" {
		log.Println("mpv result error", err)
		return errors.New(res.Error)
	}
	return nil

}

func (mpv *MPV) unixMsg(msg []byte) []byte {
	mpv.Mu.Lock()
	defer mpv.Mu.Unlock()
	_, err := mpv.Conn.Write(append(msg, '\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(mpv.Conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "request_id") {
			return []byte(line)
		}
	}
	return []byte{}
}

func (mpv *MPV) GetCmd(cmd []string) []byte {
	query := IMPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return mpv.unixMsg(jsonData)
}

func (mpv *MPV) SetCmd(cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return mpv.unixMsg(jsonData)
}

/*
 * MPV command public fxns
 */
func (mpv *MPV) Play() error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "pause", false}),
	)
}

func (mpv *MPV) Pause() error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "pause", true}),
	)
}

func (mpv *MPV) GetNowPlaying() (string, error) {
	res := mpv.GetCmd([]string{"get_property", "path"})
	return processMPVGetResult[string](res)
}

func (mpv *MPV) SetNowPlaying(filepath string) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"loadfile", filepath}),
	)
}

func (mpv *MPV) ReplaceQueueAndPlay(filepath string) error {
	res := mpv.SetCmd([]interface{}{"loadfile", filepath})
	return processMPVSetResult(res)
}

func (mpv *MPV) QueueMedia(filepath string) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"loadfile", filepath, "append-play"}),
	)
}

func (mpv *MPV) GetVol() (int, error) {
	resBytes := mpv.GetCmd([]string{"get_property", "volume"})
	n, err := processMPVGetResult[float64](resBytes)
	if err != nil {
		return 0, err
	}
	return int(n), err
}

func (mpv *MPV) SetVol(vol int) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "volume", vol}),
	)
}

func (mpv *MPV) GetPos() (int, error) {
	resBytes := mpv.GetCmd([]string{"get_property", "time-pos"})
	n, err := processMPVGetResult[float64](resBytes)
	if err != nil {
		return 0, err
	}
	return int(n), err
}

func (mpv *MPV) GetTimeRemaining() (int, error) {
	resBytes := mpv.GetCmd([]string{"get_property", "time-remaining"})
	n, err := processMPVGetResult[float64](resBytes)
	if err != nil {
		return 0, err
	}
	return int(n), err

}

func (mpv *MPV) SetPos(pos int) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "time-pos", pos}),
	)
}

func (mpv *MPV) GetTracks() ([]models.Track, error) {
	resBytes := mpv.GetCmd([]string{"get_property", "track-list"})
	return processMPVGetResult[[]models.Track](resBytes)
}

func (mpv *MPV) SetSubVisibility(visible bool) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "sub-visibility", visible}),
	)
}
