package atomeps62

import (
	"context"
	"encoding/json"
	"fmt"
)

type Info struct {
	Hostname    string
	ModelName   string
	MACAddress  string
	IPAddress   string
	Gateway     string
	PowerStatus string
}

//Info .
func (vs *AtlonaVideoSwitcher6x2) Info(ctx context.Context) (interface{}, error) {
	var network atlonaNetwork
	var hardware atlonaHardwareInfo
	var resp Info
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)

	//Get network info
	requestBody := fmt.Sprintf(`
	{
		"getConfig": {
			"network": {
				"eth0":{
				}
			}
		}
	}`)
	body, gerr := vs.make6x2request(ctx, url, requestBody)
	if gerr != nil {
		return resp, fmt.Errorf("An error occured while making the call: %w", gerr)
	}

	err := json.Unmarshal([]byte(body), &network)

	if err != nil {
		return resp, fmt.Errorf("error when unmarshalling the response: %w", err)
	}

	//Get other hardware info
	requestBody = fmt.Sprintf(`
	{
		"getConfig": {
			"system": {}
		}
	}`)
	body, gerr = vs.make6x2request(ctx, url, requestBody)
	if gerr != nil {
		return resp, fmt.Errorf("An error occured while making the call: %w", gerr)
	}
	err = json.Unmarshal([]byte(body), &hardware)
	if err != nil {
		return resp, fmt.Errorf("error when unmarshalling the response: %w", err)
	}

	//Load up the hardware struct
	resp.Hostname = hardware.System.Model
	resp.ModelName = hardware.System.Model
	resp.MACAddress = network.Network.Eth0.MacAddr
	resp.IPAddress = network.Network.Eth0.IPSettings.Ipaddr
	resp.Gateway = network.Network.Eth0.IPSettings.Gateway
	resp.PowerStatus = hardware.System.PowerStatus
	return resp, nil
}

func (vs *AtlonaVideoSwitcher6x2) Healthy(ctx context.Context) error {
	_, err := vs.AudioVideoInputs(ctx)
	if err != nil {
		return fmt.Errorf("unable to get inputs (not healthy): %s", err)
	}

	return nil
}
