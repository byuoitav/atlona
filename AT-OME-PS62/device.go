package atomeps62

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	_omePs62Endpoint = "/cgi-bin/config.cgi"
)

type AtlonaVideoSwitcher6x2 struct {
	Username string
	Password string
	Address  string

	// RequestDelay is the time to wait after sending one request to the videoswitcher
	// before sending another one. Once you have called a function on the struct,
	// changing it's value will not affect the delay.
	RequestDelay time.Duration

	once    sync.Once
	limiter *rate.Limiter
}

type config struct {
	Video videoConfig `json:"video"`
	Audio audioConfig `json:"audio"`
}

type videoConfig struct {
	VidOut struct {
		HdmiOut struct {
			Mirror struct {
				Status   bool `json:"status"`
				VideoSrc int  `json:"videoSrc"`
			} `json:"mirror"`

			HdmiOutA struct {
				VideoSrc int `json:"videoSrc"`
			} `json:"hdmiOutA"`

			HdmiOutB struct {
				VideoSrc int `json:"videoSrc"`
			} `json:"hdmiOutB"`
		} `json:"hdmiOut"`
	} `json:"vidOut"`
}

type audioConfig struct {
	AudOut struct {
		ZoneOut1 struct {
			AudioSource string `json:"audioSource"`
			AudioVol    int    `json:"audioVol"`
			VideoOut    struct {
				AudioMute bool `json:"audioMute"`
			} `json:"videoOut"`
			AnalogOut struct {
				AudioMute bool `json:"audioMute"`
				AudioVol  int  `json:"audioVol"`
			} `json:"analogOut"`
		} `json:"zoneOut1"`

		ZoneOut2 struct {
			AudioSource string `json:"audioSource"`
			AudioVol    int    `json:"audioVol"`
			VideoOut    struct {
				AudioMute bool `json:"audioMute"`
			} `json: "videoOut"`
			AnalogOut struct {
				AudioMute bool `json:"audioMute"`
				AudioVol  int  `json:"audioVol"`
			} `json:"analogOut"`
		} `json:"zoneOut2"`
	} `json:"audOut"`
}

type atlonaNetwork struct {
	Network struct {
		Eth0 struct {
			MacAddr    string `json:"macAddr"`
			DomainName string `json:"domainName"`
			DNSServer1 string `json:"dnsServer1"`
			DNSServer2 string `json:"dnsServer2"`
			IPSettings struct {
				TelnetPort int    `json:"telnetPort"`
				Ipaddr     string `json:"ipaddr"`
				Netmask    string `json:"netmask"`
				Gateway    string `json:"gateway"`
			} `json:"ipSettings"`
			LastIpaddr string `json:"lastIpaddr"`
			BootProto  string `json:"bootProto"`
		} `json:"eth0"`
	} `json:"network"`
}

//Atlona6x2HardwareInfo .
type atlonaHardwareInfo struct {
	System struct {
		PowerStatus     string `json:"powerStatus"`
		VendorID        string `json:"vendorID"`
		Model           string `json:"model"`
		SerialNumber    string `json:"serialNumber"`
		FirmwareVersion struct {
			Package          string `json:"package"`
			MasterMCU        string `json:"masterMCU"`
			TransceiverChipB string `json:"transceiverChip_B"`
			TransceiverChipC string `json:"transceiverChip_C"`
			TransceiverChipE string `json:"transceiverChip_E"`
			TransceiverChipF string `json:"transceiverChip_F"`
			Audio            string `json:"audio"`
			Fpga             string `json:"fpga"`
			Usb              string `json:"usb"`
			ScalerChip       string `json:"scalerChip"`
			ValensA          string `json:"valens_A"`
			ValensB          string `json:"valens_B"`
			ValensC          string `json:"valens_C"`
			SlaveMCU         string `json:"slaveMCU"`
			TransceiverChipA string `json:"transceiverChip_A"`
		} `json:"firmwareVersion"`
	} `json:"system"`
}

func (vs *AtlonaVideoSwitcher6x2) make6x2request(ctx context.Context, url, requestBody string) ([]byte, error) {
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("error when creting the request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	//This needs to be replaced with an environmental variable
	req.Header.Add("Authorization", "Basic YWRtaW46QXRsb25h")
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nil, fmt.Errorf("error when making call: %w", gerr)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error when making call: %w", gerr)
	}
	return body, nil
}

func (vs *AtlonaVideoSwitcher6x2) init() {
	vs.limiter = rate.NewLimiter(rate.Every(vs.RequestDelay), 1)
}

func (vs *AtlonaVideoSwitcher6x2) getConfig(ctx context.Context, body string) (config, error) {
	vs.once.Do(vs.init)

	var config config

	url := "http://" + vs.Address + _omePs62Endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return config, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	if len(vs.Username) > 0 {
		req.SetBasicAuth(vs.Username, vs.Password)
	}

	if err := vs.limiter.Wait(ctx); err != nil {
		return config, fmt.Errorf("unable to wait for ratelimit: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return config, fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return config, fmt.Errorf("unable to decode response: %w", err)
	}

	return config, nil
}

func (vs *AtlonaVideoSwitcher6x2) setConfig(ctx context.Context, body string) error {
	vs.once.Do(vs.init)

	url := "http://" + vs.Address + _omePs62Endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	if len(vs.Username) > 0 {
		req.SetBasicAuth(vs.Username, vs.Password)
	}

	if err := vs.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for ratelimit: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	var res struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}

	if !strings.EqualFold(res.Message, "OK") {
		return fmt.Errorf("bad response (%d): %s", res.Status, res.Message)
	}

	return nil
}
