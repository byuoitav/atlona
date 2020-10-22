package atomeps62

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AtlonaVideoSwitcher6x2 struct {
	Username string
	Password string
	Address  string
}

type atlonaVideo struct {
	Video struct {
		VidOut struct {
			HdmiOut struct {
				HdmiOutA struct {
					VideoSrc int `json:"videoSrc"`
				} `json:"hdmiOutA"`
				HdmiOutB struct {
					VideoSrc int `json:"videoSrc"`
				} `json:"hdmiOutB"`
				Mirror struct {
					VideoSrc int `json:"videoSrc"`
				}
			} `json:"hdmiOut"`
		} `json:"vidOut"`
	} `json:"video"`
}

type atlonaAudio struct {
	Audio struct {
		AudOut struct {
			ZoneOut1 struct {
				AnalogOut struct {
					AudioMute  bool `json:"audioMute"`
					AudioDelay int  `json:"audioDelay"`
				} `json:"analogOut"`
				AudioVol int `json:"audioVol"`
			} `json:"zoneOut1"`
			ZoneOut2 struct {
				AnalogOut struct {
					AudioMute  bool `json:"audioMute"`
					AudioDelay int  `json:"audioDelay"`
				} `json:"analogOut"`
				AudioVol int `json:"audioVol"`
			} `json:"zoneOut2"`
		} `json:"audOut"`
	} `json:"audio"`
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
