package atuhdsw52ed

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

//Volumes .
func (vs *AtlonaVideoSwitcher5x1) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	toReturn := make(map[string]int)

	vs.once.Do(vs.createPool)

	var roomInfo room
	var bytes []byte

	err := vs.pool.Do(ctx, func(ws *websocket.Conn) error {
		body := `{
			"jsonrpc": "2.0",
			"id": "<configuration_id>",
			"method": "config_get",
			"params": {
				"sections": [
					"AV Settings"
				]
			}
		}`

		vs.pool.Logger.Infof("writing message to Get Volume")

		err := ws.WriteMessage(websocket.TextMessage, []byte(body))
		if err != nil {
			return fmt.Errorf("failed to write message: %s", err.Error())
		}

		timeout := time.Now()
		timeout = timeout.Add(time.Second * 5)

		err = ws.SetReadDeadline(timeout)
		if err != nil {
			return fmt.Errorf("failed to set readDeadline: %s", err)
		}

		_, bytes, err = ws.ReadMessage()
		if err != nil {
			vs.Logger.Errorf("failed reading message from websocket: %s", err)
			return fmt.Errorf("failed to read message: %s", err)
		}

		vs.pool.Logger.Infof("read message from Get volume")

		return nil
	})
	if err != nil {
		return toReturn, fmt.Errorf("failed to read message from channel: %s", err.Error())
	}

	err = json.Unmarshal(bytes, &roomInfo)

	if err != nil {
		return toReturn, fmt.Errorf("failed to unmarshal response: %s", err.Error())
	}

	volumeLevel, err := strconv.Atoi(roomInfo.Result.AVSettings.Volume)
	if err != nil {
		return toReturn, fmt.Errorf("failed to convert volume to int: %s", err.Error())
	}

	if volumeLevel < -35 {
		toReturn[""] = 0
	} else {
		volume := ((volumeLevel + 35) * 2)
		if volume%2 != 0 {
			volume = volume + 1
		}
		toReturn[""] = volume
	}

	return toReturn, nil
}

//SetVolume .
func (vs *AtlonaVideoSwitcher5x1) SetVolume(ctx context.Context, output string, level int) error {
	vs.once.Do(vs.createPool)

	if level == 0 {
		level = -80
	} else {
		convertedVolume := -35 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}

	err := vs.pool.Do(ctx, func(ws *websocket.Conn) error {
		body := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"id": "<configuration_id>",
			"method": "config_set",
			"params": {
			  "AV Settings": {
				"Volume": "%v"
			  }
			}
		  }`, level)

		vs.pool.Logger.Infof("writing message to Set Volume")

		err := ws.WriteMessage(websocket.TextMessage, []byte(body))
		if err != nil {
			return fmt.Errorf("failed to write message: %s", err.Error())
		}

		vs.pool.Logger.Infof("successfully wrote to Set Volume")

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read message from channel: %s", err.Error())
	}

	if err != nil {
		return fmt.Errorf("failed to read message from channel: %s", err.Error())
	}

	return nil
}

//Mutes .
func (vs *AtlonaVideoSwitcher5x1) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	toReturn := make(map[string]bool)

	for _, block := range blocks {
		vs.once.Do(vs.createPool)

		var roomInfo room
		var bytes []byte

		err := vs.pool.Do(ctx, func(ws *websocket.Conn) error {
			body := `{
				"jsonrpc": "2.0",
				"id": "<configuration_id>",
				"method": "config_get",
				"params": {
					"sections": [
						"AV Settings"
					]
				}
			}`

			vs.pool.Logger.Infof("writing message to Get Muted")

			err := ws.WriteMessage(websocket.TextMessage, []byte(body))
			if err != nil {
				return fmt.Errorf("failed to write message: %s", err.Error())
			}

			timeout := time.Now()
			timeout = timeout.Add(time.Second * 5)

			err = ws.SetReadDeadline(timeout)
			if err != nil {
				return fmt.Errorf("failed to set readDeadline: %s", err)
			}

			_, bytes, err = ws.ReadMessage()
			if err != nil {
				vs.Logger.Errorf("failed reading message from websocket: %s", err)
				return fmt.Errorf("failed to read message: %s", err)
			}

			vs.pool.Logger.Infof("read message from Get Muted")

			return nil
		})
		if err != nil {
			return toReturn, fmt.Errorf("failed to read message from channel: %s", err.Error())
		}

		err = json.Unmarshal(bytes, &roomInfo)

		if err != nil {
			return toReturn, fmt.Errorf("failed to unmarshal response: %s", err.Error())
		}

		switch block {
		case "HDMI":
			isMuted, err := strconv.ParseBool(fmt.Sprintf("%v", roomInfo.Result.AVSettings.HDMIAudioMute))
			if err != nil {
				return toReturn, fmt.Errorf("failed to parse bool: %s", err.Error())
			}
			toReturn[block] = isMuted
		case "HDBT":
			isMuted, err := strconv.ParseBool(fmt.Sprintf("%v", roomInfo.Result.AVSettings.HDBTAudioMute))
			if err != nil {
				return toReturn, fmt.Errorf("failed to parse bool: %s", err.Error())
			}
			toReturn[block] = isMuted
		default:
			// Analog
			isMuted, err := strconv.ParseBool(fmt.Sprintf("%v", roomInfo.Result.AVSettings.AnalogAudioMute))
			if err != nil {
				return toReturn, fmt.Errorf("failed to parse bool: %s", err.Error())
			}
			toReturn[block] = isMuted
		}
	}

	return toReturn, nil
}

//SetMute .
func (vs *AtlonaVideoSwitcher5x1) SetMute(ctx context.Context, output string, muted bool) error {
	vs.once.Do(vs.createPool)

	var audioBlock string
	muteInt := 0

	if muted {
		muteInt = 1
	}

	switch output {
	case "HDMI":
		audioBlock = fmt.Sprintf(`"HDMI Audio Mute": %v`, muteInt)
	case "HDBT":
		audioBlock = fmt.Sprintf(`"HDBT Audio Mute": %v`, muteInt)
	default:
		// Analog
		audioBlock = fmt.Sprintf(`"Analog Audio Mute": %v`, muteInt)
	}

	err := vs.pool.Do(ctx, func(ws *websocket.Conn) error {
		body := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"id": "<configuration_id>",
			"method": "config_set",
			"params": {
			  "AV Settings": {
				%s
			  }
			}
		  }`, audioBlock)

		vs.pool.Logger.Infof("writing message to set Muted")

		err := ws.WriteMessage(websocket.TextMessage, []byte(body))
		if err != nil {
			return fmt.Errorf("failed to write message: %s", err.Error())
		}

		vs.pool.Logger.Infof("wrote message to set Muted")

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to read message from channel: %s", err.Error())
	}

	return nil

}
