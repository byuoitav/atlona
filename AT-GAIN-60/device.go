package atgain60

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"go.uber.org/zap"
)

type Amp struct {
	Username string
	Password string
	Address  string
	Log      *zap.Logger
}

/*
// Amp60 represents an Atlona 60 watt amplifier
type Amp60 struct {
	Username string
	Password string
	Address  string
	Log      *zap.Logger
}
*/

/*
type ampResponse struct {
	Model         string `json:"101"`
	Firmware      string `json:"102"`
	MACAddress    string `json:"103"`
	SerialNumber  string `json:"104"`
	OperatingTime string `json:"105"`
}

// AmpStatus represents the current amp status
type AmpStatus struct {
	Model         string `json:"101"`
	Firmware      string `json:"102"`
	MACAddress    string `json:"103"`
	SerialNumber  string `json:"104"`
	OperatingTime string `json:"105"`
}

// AmpAudio represents an audio response from an Atlona 60 watt amp
type AmpAudio struct {
	Volume string `json:"608,omitempty"`
	Muted  string `json:"609,omitempty"`
}
*/

func (a *Amp) r() string {
	return fmt.Sprintf("%v", rand.Float64())
}

type loginResult struct {
	Login bool
}

func (a *Amp) login(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/action=compare&701=%s&702=%s&r=%s", a.Address, a.Username, a.Password, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	client := &http.Client{
		Transport: &transport{},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read body: %w", err)
	}

	fmt.Printf("body:\n%s\n", body)

	/*
		if err := json.NewDecoder(resp.Body).Decode(&login); err != nil {
		}
	*/

	return nil
}

/*
func getURL(address, endpoint string) string {
	return "http://" + address + "/action=" + endpoint + "&r=" + getR()
}

func (a *Amp60) getLoginUrl() string {
	return "http://" + a.Address + "/action=compare&701=" + a.Username + "&702=" + a.Password + "&r=" + getR()
}

func (a *Amp60) sendReq(ctx context.Context, endpoint string) ([]byte, error) {
	// checking to validate that it is logged in
	err := a.login(ctx)
	if err != nil {
		return nil, fmt.Errorf("Login failed to device: %v", err)
	}

	var toReturn []byte
	ampUrl := getURL(a.Address, endpoint)
	Client := http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequestWithContext(ctx, "GET", ampUrl, nil)
	req.Header.Set("Context-type", "application/json")
	//req, err := http.NewRequest("GET", ampUrl, nil)
	a.Log.Debug("Request Output", zap.Any("request", req))
	if err != nil {
		return toReturn, fmt.Errorf("unable to make new http request: %w", err)
	}
	resp, err := Client.Do(req)
	a.Log.Debug("RESP Output", zap.Any("response", resp))
	if err != nil {
		if nerr, ok := err.(*url.Error); ok {
			fmt.Printf("%v\n", nerr.Err)
			if !strings.Contains(nerr.Err.Error(), "malformed") {
				return toReturn, fmt.Errorf("unable to perform request: %w", err)
			}
		} else {
			return toReturn, fmt.Errorf("unable to perform request: %w", err)
		}
		return toReturn, nil
	}
	defer resp.Body.Close()
	toReturn, err = ioutil.ReadAll(resp.Body)
	s := string(toReturn)
	a.Log.Info("Response", zap.String("response", s))

	if err != nil {
		return toReturn, fmt.Errorf("unable to read resp body: %w", err)
	}
	return toReturn, nil
}

// login for device
func (a *Amp60) login(ctx context.Context) error {
	// Check if we are currently logged in
	resp, err := http.Get(a.getLoginUrl())
	if err != nil {
		return fmt.Errorf("Unable to log in: %v", err)
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	s := string(out)
	if err != nil {
		return fmt.Errorf("Cannot read body of test: %v", err)
	}

	if strings.Contains(s, "404") == true {
		var toReturn []byte
		loginURL := a.getLoginUrl()
		Client := http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequestWithContext(ctx, "GET", loginURL, nil)
		if err != nil {
			return fmt.Errorf("Unable to create request: %v", err)
		}
		resp, err := Client.Do(req)
		if err != nil {
			return fmt.Errorf("Unable to connect to device: %v", err)
		}
		defer resp.Body.Close()
		toReturn, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Cannot read the body of the response")
		}
		data := loginResult{}
		json.Unmarshal(toReturn, &data)
		if data.Login != true {
			return fmt.Errorf("Not able to login: %v", err)
		}
		return nil
	}

	return nil

}
*/
