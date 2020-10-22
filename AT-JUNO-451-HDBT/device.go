package atjuno451hdbt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	avSettingsPage = "avs"
	infoPage       = "info"
)

type AtlonaVideoSwitcher4x1 struct {
	Username string
	Password string
	Address  string
}

// AVSettings is the response from the switcher for the av settings page
type AVSettings struct {
	HDMIInputAudioBreakout int   `json:"ARC"`
	HDCPSettings           []int `json:"HDCPSet"`
	AudioOutput            int   `json:"HDMIAud"`
	Toslink                int   `json:"Toslink"`
	AutoSwitch             int   `json:"asw"`
	Input                  int   `json:"inp"`
	LoggedIn               int   `json:"login_ur"`
}

func getPage(ctx context.Context, address, page string, structToFill interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/aj.html?a=%s", address, page), nil)
	if err != nil {
		return fmt.Errorf("unable to get page %s on %s", page, address)
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to get page %s on %s", page, address)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to get page %s on %s", page, address)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("unable to get page %s on %s - %v response recevied. body: %s", page, address, resp.StatusCode, b)
	}

	err = json.Unmarshal(b, structToFill)
	if err != nil {
		return fmt.Errorf("unable to get page %s on %s", page, address)
	}

	return nil
}

func sendCommand(ctx context.Context, address, command string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%v/aj.html?a=command&cmd=%s", address, command), nil)
	if err != nil {
		return fmt.Errorf("unable to send command '%s' to %s", command, address)
	}

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send command '%s' to %s", command, address)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to send command '%s' to %s", command, address)
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("unable to send command '%s' to %s - %v response received. body: %s", command, address, resp.StatusCode, b)
	}

	return nil
}
