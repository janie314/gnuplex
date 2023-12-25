package mpv

import (
	"encoding/json"
)

/*
 * MPV command public fxns
 */
func Play() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", false})
}

func Pause() []byte {
	return mpvSetCmd([]interface{}{"set_property", "pause", true})
}

func IsPaused() []byte {
	return mpvGetCmd([]string{"get_property", "pause"})
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

func mpvSetCmd(cmd []interface{}) []byte {
	query := IMPVQuery{Command: cmd}
	jsonData, err := json.Marshal(query)
	if err != nil {
		return []byte{}
	}
	return unixMsg(jsonData)
}
