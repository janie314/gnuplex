package mpv

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type MPVQuery struct {
	Command []interface{} `json:"command"`
}

type ResponseDatum interface {
	float64 | ~string | bool
}

type EmptyResponse struct {
	RequestId int    `json:"request_id"`
	Err       string `json:"error"`
}

type Response[T ResponseDatum] struct {
	Data      T      `json:"data"`
	RequestId int    `json:"request_id"`
	Err       string `json:"error"`
}

/*
 * MPV command public fxns
 * Refer to the following MPV documentation
 *     - https://mpv.io/manual/stable/#list-of-input-commands
 *     - https://mpv.io/manual/stable/#properties
 *     - https://mpv.io/manual/stable/#json-ipc
 */

// Set mpv video's play/paused state (which is a boolean `paused`).
func (mpv *MPV) SetPause(paused bool) error {
	_, err := MPVCmd(mpv, "pause", false, false, true)
	return err
}

// Toggles play/pause.
// Returns the `paused` boolean status after the toggle operation is complete.
func (mpv *MPV) Toggle() (bool, error) {
	paused, err := MPVCmd[bool](mpv, "pause", false, true, true)
	if err != nil {
		return false, err
	} else {
		_, err = MPVCmd(mpv, "pause", !paused, false, true)
	}
	if err != nil {
		return false, errors.New("issue with play/pause cmd")
	} else {
		return !paused, nil
	}
}

// Returns the current `paused` boolean status: whether or not the mpv
// video is paused.
func (mpv *MPV) IsPaused() (bool, error) {
	return MPVCmd[bool](mpv, "pause", false, true, true)
}

func (mpv *MPV) GetMedia() (string, error) {
	return MPVCmd[string](mpv, "path", "", true, true)
}

func (mpv *MPV) SetMedia(filepath string) error {
	_, err := MPVCmd(mpv, "loadfile", filepath, false, false)
	return err
}

func (mpv *MPV) GetVolume() (float64, error) {
	return MPVCmd[float64](mpv, "volume", 0, true, true)
}

func (mpv *MPV) SetVolume(vol float64) error {
	_, err := MPVCmd(mpv, "volume", vol, false, true)
	return err
}

func (mpv *MPV) GetPos() (float64, error) {
	return MPVCmd[float64](mpv, "time-pos", 0, true, true)
}

func (mpv *MPV) SetPos(pos float64) error {
	_, err := MPVCmd(mpv, "time-pos", pos, false, true)
	return err
}

func (mpv *MPV) GetTimeRemaining() (float64, error) {
	return MPVCmd[float64](mpv, "time-remaining", 0, true, true)
}

func (mpv *MPV) IncPos(pos float64) error {
	_, err := MPVCmd(mpv, "seek", pos, false, false)
	return err
}

func (mpv *MPV) Screenshot() error {
	_, err := MPVCmd(mpv, "screenshot", "", false, false)
	return err
}

/*
 * Execute a read or write query to MPV.
 *
 * cmd_or_prop and argrepresent two arguments from the MPV command-line interface.
 * (cmd_or_prop, arg) could be:
 *    - ("seek", 30)
 *    - ("screenshot", "")
 *    - ("time-pos", 444)
 * etc.
 *
 * read_query represents whether the command is a read-query or write-query (does it get
 * the current media file, or set it?). We ignore the response value of write queries and
 * always return the zero value for the first return argument T.
 *
 * is_prop represents if the command is executed via MPV's get_property (read) or
 * set_property (write) API, e.g.
 *
 *    - set_property time-pos 444
 */
func MPVCmd[T ResponseDatum](mpv *MPV, cmd_or_prop string, arg T, read_query, is_prop bool) (T, error) {
	// set up query to mpv
	var zero T
	var query_part []interface{}
	var prop_str string
	if is_prop {
		if read_query {
			prop_str = "get_property"
		} else {
			prop_str = "set_property"
		}
	}
	if is_prop {
		if arg == zero && read_query {
			query_part = []interface{}{prop_str, cmd_or_prop}
		} else {
			query_part = []interface{}{prop_str, cmd_or_prop, arg}
		}
	} else {
		if arg == zero && read_query {
			query_part = []interface{}{cmd_or_prop}
		} else {
			query_part = []interface{}{cmd_or_prop, arg}
		}
	}
	query_struct := MPVQuery{Command: query_part}
	query, err := json.Marshal(query_struct)
	// make query and parse result
	res_bytes := mpv.unixCmd(query)
	if read_query {
		var response Response[T]
		err = json.Unmarshal(res_bytes, &response)
		if err != nil {
			return zero, err
		} else if response.Err != "success" {
			log.Println("err", response.Err)
			return zero, errors.New("failure from mpv query")
		} else {
			return response.Data, nil
		}
	} else {
		// ignore response if it's a write query
		var response EmptyResponse
		err = json.Unmarshal(res_bytes, &response)
		if err != nil {
			return zero, err
		} else if response.Err != "success" {
			log.Println("err", response.Err)
			return zero, errors.New("failure from mpv query")
		} else {
			return zero, nil
		}

	}
}

/*
 * Send a message to the MPV unix socket.
 *
 * This function handles its own locking.
 */
func (mpv *MPV) unixCmd(msg []byte) []byte {
	mpv.mu.Lock()
	defer mpv.mu.Unlock()
	log.Println("debug\tsending\t", string(msg[:]))
	_, err := mpv.conn.Write(append(msg, '\n'))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(mpv.conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "request_id") {
			log.Println("debug\treceiving\t", line)
			return []byte(line)
		}
	}
	return []byte{}
}
