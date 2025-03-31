package mpv

import (
	"bufio"
	"encoding/json"
	"errors"
	"gnuplex/models"
	"log"
	"strings"
)

/*
 * Types
 */
type MPVQuery struct {
	Command []interface{} `json:"command"`
}

type MPVQueryString struct {
	Command []string `json:"command"`
}

type MPVResponseBool struct {
	Data bool `json:"data"`
}

type MPVResponseString struct {
	Data string `json:"data"`
}

type MPVResponseInt struct {
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

func processMPVGetResult[T string | int | float64 | []models.Track](resBytes []byte, err error) (T, error) {
	var defaultVal T
	if err != nil {
		return defaultVal, err
	}
	var res MPVGetResult[T]
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("7b289386-c838-457d-8545-2f56bddb3746 MPV reported error in API call:", err)
		return defaultVal, err
	} else if res.Error != "success" {
		log.Println("b8d14416-bbec-4297-ada8-3d13bf9c6e79 MPV reported error in API Call:", err)
		return defaultVal, errors.New(res.Error)
	}
	return res.Data, nil

}

func processMPVSetResult(resBytes []byte, err error) error {
	if err != nil {
		return err
	}
	var res MPVSetResult
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		log.Println("3989f92e-5230-4403-a60c-1ee8429c0088 MPV reported error in API call:", err)
		return err
	} else if res.Error != "success" {
		log.Println("d1a6f614-5eb6-4fa5-aedb-31b1940cb58e MPV reported error in API call:", err)
		return errors.New(res.Error)
	}
	return nil

}

func (mpv *MPV) unixMsg(msg []byte) ([]byte, error) {
	mpv.Mu.Lock()
	defer mpv.Mu.Unlock()
	_, err := mpv.Conn.Write(append(msg, '\n'))
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(mpv.Conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "request_id") {
			return []byte(line), nil
		}
	}
	return []byte{}, nil
}

func (mpv *MPV) GetCmd(cmd []string) ([]byte, error) {
	query := MPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}, nil
	}
	return mpv.unixMsg(jsonData)
}

func (mpv *MPV) SetCmd(cmd []interface{}) ([]byte, error) {
	query := MPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}, nil
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

func (mpv *MPV) GetCurrentFilepath() (string, error) {
	return processMPVGetResult[string](mpv.GetCmd([]string{"get_property", "path"}))
}

func (mpv *MPV) SetNowPlaying(filepath string) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"loadfile", filepath}),
	)
}

func (mpv *MPV) ReplaceQueueAndPlay(filepath string) error {
	return processMPVSetResult(mpv.SetCmd([]interface{}{"loadfile", filepath}))
}

func (mpv *MPV) QueueMedia(filepath string) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"loadfile", filepath, "append-play"}),
	)
}

func (mpv *MPV) GetVol() (int, error) {
	n, err := processMPVGetResult[float64](mpv.GetCmd([]string{"get_property", "volume"}))
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
	n, err := processMPVGetResult[float64](mpv.GetCmd([]string{"get_property", "time-pos"}))
	if err != nil {
		return 0, err
	}
	return int(n), err
}

func (mpv *MPV) GetTimeRemaining() (int, error) {
	n, err := processMPVGetResult[float64](mpv.GetCmd([]string{"get_property", "time-remaining"}))
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
	return processMPVGetResult[[]models.Track](mpv.GetCmd([]string{"get_property", "track-list"}))
}

func (mpv *MPV) SetSubVisibility(visible bool) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "sub-visibility", visible}),
	)
}

func (mpv *MPV) SetSubTrack(trackID int64) error {
	return processMPVSetResult(
		mpv.SetCmd([]interface{}{"set_property", "sid", trackID}),
	)
}

func (mpv *MPV) Ping() (int, error) {
	return processMPVGetResult[int](mpv.GetCmd([]string{"get_property", "pid"}))
}
