package mpv

import (
	"encoding/json"
	"errors"
	"log"
)

/*
 * MPV command public fxns
 * Refer to the following MPV documentation
 *     - https://mpv.io/manual/stable/#list-of-input-commands
 *     - https://mpv.io/manual/stable/#json-ipc
 *     - https://mpv.io/manual/stable/#properties
 */

// Set mpv video's play/paused state (which is a boolean `paused`).
func (mpv *MPV) SetPause(paused bool) error {
	return SetMPVProperty(mpv, "pause", false)
}

type ResponseData interface {
	float64 | ~string | bool
}

type EmptyResponse struct {
	RequestId int    `json:"request_id"`
	Err       string `json:"error"`
}

type Response[T ResponseData] struct {
	Data      T      `json:"data"`
	RequestId int    `json:"request_id"`
	Err       string `json:"error"`
}

// Toggles play/pause.
// Returns the `paused` boolean status after the toggle operation is complete.
func (mpv *MPV) Toggle() (bool, error) {
	paused, err := GetMPVProperty[bool](mpv, "pause")
	if err != nil {
		err = SetMPVProperty(mpv, "pause", !paused)
	} else {
		return false, err
	}
	if err != nil {
		return !paused, nil
	} else {
		return false, errors.New("issue with play/pause cmd")
	}
}

// Returns the current `paused` boolean status: whether or not the mpv
// video is paused.
func (mpv *MPV) IsPaused() (bool, error) {
	return GetMPVProperty[bool](mpv, "pause")
}

func (mpv *MPV) GetMedia() (string, error) {
	return GetMPVProperty[string](mpv, "path")
}

func (mpv *MPV) SetMedia(filepath string) error {
	// TODO addhist here
	return mpv.SetCmd([]string{"loadfile", filepath})
}

func (mpv *MPV) GetVolume() (float64, error) {
	return GetMPVProperty[float64](mpv, "volume")
}

func (mpv *MPV) SetVolume(vol float64) error {
	return SetMPVProperty(mpv, "volume", vol)
}

func (mpv *MPV) GetPos() (float64, error) {
	return GetMPVProperty[float64](mpv, "time-pos")
}

func (mpv *MPV) SetPos(pos float64) error {
	return SetMPVProperty(mpv, "time-pos", pos)
}

func (mpv *MPV) IncPos(pos float64) error {
	return SetMPVProperty(mpv, "seek", pos)
}

func (mpv *MPV) Screenshot() error {
	return mpv.SetCmd([]string{"screenshot"})
}

/*
 * key wrapper functions
 */
func GetMPVProperty[T ResponseData](mpv *MPV, prop string) (T, error) {
	// set up query to mpv
	query_part := []string{"get_property", prop}
	query_struct := IMPVQueryString{Command: query_part}
	query, err := json.Marshal(query_struct)
	// return values
	var response Response[T]
	var defaultT T
	if err != nil {
		return defaultT, err
	}
	// make query and parse result
	res_bytes := mpv.UnixMsg(query)
	err = json.Unmarshal(res_bytes, &response)
	log.Println("debug", string(res_bytes[:]))
	log.Println("debug", response.Data, response.Err, response.RequestId)
	if err != nil {
		return defaultT, err
	} else if response.Err != "success" {
		log.Println("err", response.Err)
		return defaultT, errors.New("failure from mpv query")
	} else {
		return response.Data, nil
	}
}

func SetMPVProperty[T ResponseData](mpv *MPV, prop string, val T) error {
	// set up query to mpv
	query_part := []interface{}{"set_property", prop, val}
	query_struct := IMPVQuery{Command: query_part}
	query, err := json.Marshal(query_struct)
	// make query and parse result
	res_bytes := mpv.UnixMsg(query)
	var response Response[T]
	err = json.Unmarshal(res_bytes, &response)
	log.Println("debug", string(res_bytes[:]))
	log.Println("debug", response.Data, response.Err, response.RequestId)
	if err != nil {
		return err
	} else if response.Err != "success" {
		log.Println("err", response.Err)
		return errors.New("failure from mpv query")
	} else {
		return nil
	}
}

func (mpv *MPV) SetCmd(args []string) error {
	// set up query to mpv
	query_struct := IMPVQueryString{Command: args}
	query, err := json.Marshal(query_struct)
	// make query and parse result
	res_bytes := mpv.UnixMsg(query)
	var response Response[bool]
	err = json.Unmarshal(res_bytes, &response)
	log.Println("debug", string(res_bytes[:]))
	log.Println("debug", response.Data, response.Err, response.RequestId)
	if err != nil {
		return err
	} else if response.Err != "success" {
		log.Println("err", response.Err)
		return errors.New("failure from mpv query")
	} else {
		return nil
	}
}
