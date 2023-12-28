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
	int | ~string | bool
}

type EmptyResponse struct {
	request_id int
	error      string
}

type Response[T ResponseData] struct {
	data      T
	requested int
	error     string
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
	return SetMPVProperty(mpv, "loadfile", filepath)
}

func (mpv *MPV) GetVolume() (int, error) {
	return GetMPVProperty[int](mpv, "volume")
}

func (mpv *MPV) SetVolume(vol int) error {
	return SetMPVProperty(mpv, "volume", vol)
}

func (mpv *MPV) GetPos() (int, error) {
	return GetMPVProperty[int](mpv, "time-pos")
}

func (mpv *MPV) SetPos(pos int) error {
	return SetMPVProperty(mpv, "time-pos", pos)
}

func (mpv *MPV) IncPos(pos int) error {
	return SetMPVProperty(mpv, "seek", pos)
}

func (mpv *MPV) Screenshot() error {
	return mpv.SetCmd("screenshot")
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
	log.Println("debug", string(res_bytes[:]))
	err = json.Unmarshal(res_bytes, &response)
	if err != nil {
		return defaultT, err
	} else if response.error != "success" {
		return defaultT, errors.New("failure from mpv query")
	} else {
		return response.data, nil
	}
}

func SetMPVProperty[T ResponseData](mpv *MPV, prop string, val T) error {
	// set up query to mpv
	query_part := []interface{}{"set_property", prop, val}
	query_struct := IMPVQuery{Command: query_part}
	query, err := json.Marshal(query_struct)
	// make query and parse result
	res_bytes := mpv.UnixMsg(query)
	log.Println("debug", string(res_bytes[:]))
	var response Response[T]
	err = json.Unmarshal(res_bytes, &response)
	if err != nil {
		return err
	} else if response.error != "success" {
		return errors.New("failure from mpv query")
	} else {
		return nil
	}
}

func (mpv *MPV) SetCmd(cmd string) error {
	// set up query to mpv
	query_part := []string{cmd}
	query_struct := IMPVQueryString{Command: query_part}
	query, err := json.Marshal(query_struct)
	// make query and parse result
	res_bytes := mpv.UnixMsg(query)
	log.Println("debug", string(res_bytes[:]))
	var response Response[bool]
	err = json.Unmarshal(res_bytes, &response)
	if err != nil {
		return err
	} else if response.error != "success" {
		return errors.New("failure from mpv query")
	} else {
		return nil
	}
}
