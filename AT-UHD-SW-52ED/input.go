package atuhdsw52ed

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

//AudioVideoInputs .
func (vs *AtlonaVideoSwitcher5x1) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)
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

		vs.pool.Logger.Infof("writing message to Get Input")

		err := ws.WriteMessage(websocket.TextMessage, []byte(body))
		if err != nil {
			return fmt.Errorf("failed to write message: %s", err.Error())
		}

		if vs.Logger != nil {
			vs.Logger.Infof("reading message from websocket")
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

		vs.pool.Logger.Infof("read message from Get Input")

		return nil
	})

	if err != nil {
		return toReturn, fmt.Errorf("failed to read message from channel: %s", err.Error())
	}

	err = json.Unmarshal(bytes, &roomInfo)

	if err != nil {
		return toReturn, fmt.Errorf("failed to unmarshal message: %s", err.Error())
	}

	toReturn[""] = roomInfo.Result.AVSettings.Source[6:]
	return toReturn, nil
}

//SetAudioVideoInput .
func (vs *AtlonaVideoSwitcher5x1) SetAudioVideoInput(ctx context.Context, output, input string) error {
	vs.once.Do(vs.createPool)

	intInput, nerr := strconv.Atoi(input)

	if nerr != nil {
		return fmt.Errorf("error occured when converting input to int: %w", nerr)
	}

	if intInput == 0 || intInput > 5 {
		return fmt.Errorf("Invalid Input. The input requested must be between 1-5. The input you requested was %v", intInput)
	}

	err := vs.pool.Do(ctx, func(ws *websocket.Conn) error {
		body := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"id": "<configuration_id>",
			"method": "config_set",
			"params": {
			  "AV Settings": {
				"source": "input %s"
			  }
			}
		  }`, input)

		if vs.Logger != nil {
			vs.Logger.Infof("writing message")
		}

		vs.pool.Logger.Infof("writing message to Set Input")

		err := ws.WriteMessage(websocket.TextMessage, []byte(body))
		if err != nil {
			return fmt.Errorf("failed to write message: %s", err.Error())
		}

		vs.pool.Logger.Infof("successful wrote message to set Input")

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

func (vs *AtlonaVideoSwitcher5x1) Healthy(ctx context.Context) error {
	_, err := vs.AudioVideoInputs(ctx)
	if err != nil {
		return fmt.Errorf("unable to get inputs (not healthy): %s", err)
	}

	return nil
}
