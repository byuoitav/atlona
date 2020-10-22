package atuhdsw52ed

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/byuoitav/wspool"
	"github.com/gorilla/websocket"
)

type AtlonaVideoSwitcher5x1 struct {
	Username string
	Password string
	Address  string
	once     sync.Once
	pool     wspool.Pool
	Logger   wspool.Logger
}

type room struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		AVSettings struct {
			Source          string `json:"source"`
			Autoswitch      int    `json:"Autoswitch"`
			Volume          string `json:"Volume"`
			HDMIAudioMute   int    `json:"HDMI Audio Mute"`
			HDBTAudioMute   int    `json:"HDBT Audio Mute"`
			AnalogAudioMute int    `json:"Analog Audio Mute"`
		} `json:"AV Settings"`
	} `json:"result"`
}

func (vs *AtlonaVideoSwitcher5x1) createPool() {
	if vs.Logger != nil {
		vs.Logger.Infof("creating pool")
	}

	vs.pool = wspool.Pool{
		NewConnection: createConnectionFunc(vs.Address),
		TTL:           10 * time.Second,
		Delay:         75 * time.Millisecond,
		Logger:        vs.Logger,
	}

}

func createConnectionFunc(address string) wspool.NewConnectionFunc {
	return func(ctx context.Context) (*websocket.Conn, error) {
		dialer := &websocket.Dialer{}
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		ws, _, err := dialer.DialContext(ctx, fmt.Sprintf("ws://%s:543", address), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to open websocket: %s", err.Error())
		}
		return ws, nil
	}
}
