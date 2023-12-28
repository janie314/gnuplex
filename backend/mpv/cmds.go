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
func Play() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", false})
}

func Pause() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", true})
}

type ResponseData interface {
	int | ~string | bool
}

type MPVResponse[T ResponseData] struct {
	data         T
	requested_id int
	error        string
}

func Toggle() error {
	paused, err := mpvGetProperty[bool]("pause")
	if err != nil {
		if paused {
			Play()
		} else {
			Pause()
		}
	}
	return err
}

func IsPaused() (bool, error) {
	return mpvGetProperty[bool]("pause")
}

func GetMedia() []byte {
	return mpvGetCmd([]string{"get_property", "path"})
}

func SetMedia(filepath string) []byte {
	return mpvSetCmd([]interface{}{"loadfile", filepath})
}

func GetVolume() []byte {
	return mpvGetCmd([]string{"get_property", "volume"})
}

func SetVolume(vol int) []byte {
	return mpvSetCmd([]interface{}{"set_property", "volume", vol})
}

func GetPos() []byte {
	return mpvGetCmd([]string{"get_property", "time-pos"})
}

func SetPos(pos int) []byte {
	return mpvSetCmd([]interface{}{"set_property", "time-pos", pos})
}

func IncPos(pos int) []byte {
	return mpvSetCmd([]interface{}{"seek", pos})
}

func Screenshot() []byte {
	return mpvSetCmd([]interface{}{"screenshot"})
}

/*
 * key wrapper functions
 */
func mpvGetCmd(cmd []string) []byte {
	query := IMPVQueryString{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(jsonData)
}

func mpvGetProperty[T ResponseData](prop string) (T, error) {
	// set up query to mpv
	query_part := []string{"get_property", prop}
	query_struct := IMPVQueryString{Command: query_part}
	query, err := json.Marshal(query_struct)
	// return values
	var response MPVResponse[T]
	var defaultT T
	if err != nil {
		return defaultT, err
	}
	// make query and parse result
	res_bytes := unixMsg(query)
	err = json.Unmarshal(res_bytes, &defaultT)
	if err != nil {
		return defaultT, err
	} else if response.error != "success" {
		return defaultT, errors.New("failure from mpv query")
	} else {
		return response.data, nil
	}
}

func mpvSetCmd(cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(jsonData)
}
