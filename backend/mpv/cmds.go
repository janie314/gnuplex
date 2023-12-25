package mpv

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
