package atomeps62

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

//AudioVideoInputs .
func (vs *AtlonaVideoSwitcher6x2) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)

	for i := 1; i < 3; i++ {
		var resp atlonaVideo
		url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)

		requestBody := fmt.Sprintf(`
		{
			"getConfig": {
				"video": {
					"vidOut": {
						"hdmiOut": {
						}
					}
				}
			}
		}`)

		body, gerr := vs.make6x2request(ctx, url, requestBody)
		if gerr != nil {
			return toReturn, fmt.Errorf("An error occured while making the call: %w", gerr)
		}

		err := json.Unmarshal([]byte(body), &resp)
		if err != nil {
			fmt.Printf("%s/n", body)
			return toReturn, fmt.Errorf("error when unmarshalling the response: %w", err)
		}

		//Get the inputsrc for the requested output
		input := ""
		if i == 1 {
			input = strconv.Itoa(resp.Video.VidOut.HdmiOut.HdmiOutA.VideoSrc)
		} else if i == 2 {
			input = strconv.Itoa(resp.Video.VidOut.HdmiOut.HdmiOutB.VideoSrc)
		} else {
			input = strconv.Itoa(resp.Video.VidOut.HdmiOut.Mirror.VideoSrc)
		}

		toReturn[strconv.Itoa(i)] = input
	}

	return toReturn, nil
}

//SetAudioVideoInput .
func (vs *AtlonaVideoSwitcher6x2) SetAudioVideoInput(ctx context.Context, output, input string) error {
	in, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("error when making call: %w", err)
	}
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)
	requestBody := ""
	if output == "1" {
		requestBody = fmt.Sprintf(`
		{
			"setConfig":{
				"video":{
					"vidOut":{
						"hdmiOut":{
							"hdmiOutA":{
								"videoSrc":%v
							}
						}
					}
				}
			}
		}`, in)
	} else if output == "2" {
		requestBody = fmt.Sprintf(`
		{
			"setConfig":{
				"video":{
					"vidOut":{
						"hdmiOut":{
							"hdmiOutB":{
								"videoSrc":%v
							}
						}
					}
				}
			}
		}`, in)
	} else {
		requestBody = fmt.Sprintf(`
		{"setConfig":{"video":{"vidOut":{"hdmiOut":{"mirror":{"videoSrc":%v}}}}}}
		`, in)
	}

	_, gerr := vs.make6x2request(ctx, url, requestBody)
	if gerr != nil {
		return fmt.Errorf("An error occured while making the call: %w", gerr)
	}
	return nil
}
