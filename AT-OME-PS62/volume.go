package atomeps62

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

//Volumes .
func (vs *AtlonaVideoSwitcher6x2) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	toReturn := make(map[string]int)

	for _, block := range blocks {
		var resp atlonaAudio
		url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)
		requestBody := fmt.Sprintf(`
		{
			"getConfig": {
				"audio": {
					"audOut": {
						}
					}
				}
		}`)
		body, gerr := vs.make6x2request(ctx, url, requestBody)
		if gerr != nil {
			return toReturn, fmt.Errorf("An error occured while making the call: %w", gerr)
		}

		err := json.Unmarshal([]byte(body), &resp) // here!
		if err != nil {
			return toReturn, fmt.Errorf("error when unmarshalling the response: %w", err)
		}
		if block == "1" {
			if resp.Audio.AudOut.ZoneOut1.AudioVol < -40 {
				toReturn[block] = 0
			} else {
				volume := ((resp.Audio.AudOut.ZoneOut1.AudioVol + 40) * 2)
				toReturn[block] = volume
			}
		} else if block == "2" {
			toReturn[block] = resp.Audio.AudOut.ZoneOut2.AudioVol + 90
		} else {
			return toReturn, fmt.Errorf("invalid Output. Valid Output names are 1 and 2 you gave us %s", block)
		}
	}

	return toReturn, nil
}

//SetVolume .
func (vs *AtlonaVideoSwitcher6x2) SetVolume(ctx context.Context, output string, level int) error {
	//Atlona volume levels are from -90 to 10 and the number we recieve is 0-100
	//if volume level is supposed to be zero set it to zero (which is -90) on atlona

	if level == 0 {
		level = -90
	} else {
		convertedVolume := -40 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)
	if output == "1" || output == "2" {
		requestBody := fmt.Sprintf(`
		{
			"setConfig": {
				"audio": {
					"audOut": {
						"zoneOut%s": {
							"audioVol": %d
						}
					}
				}
			}
		}`, output, level)
		_, gerr := vs.make6x2request(ctx, url, requestBody)
		if gerr != nil {
			return fmt.Errorf("An error occured while making the call: %w", gerr)
		}
	} else {
		return fmt.Errorf("Invalid Output. Valid Audio Output names are Audio1 and Audio2: you gave us %s", output)
	}
	return nil
}

//Mutes .
func (vs *AtlonaVideoSwitcher6x2) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	toReturn := make(map[string]bool)

	for _, block := range blocks {
		var resp atlonaAudio
		if block == "1" || block == "2" {
			url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)
			requestBody := fmt.Sprintf(`
			{
				"getConfig": {
					"audio":{
						"audOut":{
							"zoneOut%s":{
								"analogOut": {				
								}
							}
						}
					}	
				}	
			}`, block)
			body, gerr := vs.make6x2request(ctx, url, requestBody)
			if gerr != nil {
				return toReturn, fmt.Errorf("An error occured while making the call: %w", gerr)
			}
			err := json.Unmarshal([]byte(body), &resp)
			if err != nil {
				return toReturn, fmt.Errorf("error when unmarshalling the response: %w", err)
			}
		} else {
			return toReturn, fmt.Errorf("Invalid Output. Valid Output names are 1 and 2 you gave us %s", block)
		}
		if block == "1" {
			toReturn[block] = resp.Audio.AudOut.ZoneOut1.AnalogOut.AudioMute
		} else if block == "2" {
			toReturn[block] = resp.Audio.AudOut.ZoneOut2.AnalogOut.AudioMute
		} else {
			return toReturn, fmt.Errorf("Invalid Output. Valid Output names are 1 and 2 you gave us %s", block)
		}
	}

	return toReturn, nil
}

//SetMute .
func (vs *AtlonaVideoSwitcher6x2) SetMute(ctx context.Context, output string, muted bool) error {
	url := fmt.Sprintf("http://%s/cgi-bin/config.cgi", vs.Address)
	if output == "1" || output == "2" {
		requestBody := fmt.Sprintf(`
		{
			"setConfig": {
				"audio": {
					"audOut": {
						"zoneOut%s": {
							"analogOut": {
								"audioMute": %v
							}
						}
					}
				}
			}
		}`, output, muted)
		_, gerr := vs.make6x2request(ctx, url, requestBody)
		if gerr != nil {
			return fmt.Errorf("An error occured while making the call: %w", gerr)
		}
	} else {
		return fmt.Errorf("Invalid Output. Valid Output names are Audio1 and Audio2 you gave us %s", output)
	}
	return nil
}
