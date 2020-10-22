package athdvs210u

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// GetAudioVideoInputs .
func (vs *AtlonaVideoSwitcher2x1) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)

	var resp wallPlateStruct
	url := fmt.Sprintf("http://%s/aj.html?a=avs", vs.Address)
	body, gerr := vs.make2x1request(ctx, url)
	if gerr != nil {
		return toReturn, fmt.Errorf("An error occured while making the call: %w", gerr)
	}
	err := json.Unmarshal([]byte(body), &resp) // here!
	if err != nil {
		return toReturn, fmt.Errorf("error when unmarshalling the response: %w", err)
	}

	in := strconv.Itoa(resp.Inp)

	toReturn[""] = in
	return toReturn, nil
}

// SetAudioVideoInput .
func (vs *AtlonaVideoSwitcher2x1) SetAudioVideoInput(ctx context.Context, output, input string) error {
	intInput, nerr := strconv.Atoi(input)
	if nerr != nil {
		return fmt.Errorf("failed to convert input from string to int: %w", nerr)
	}
	if intInput != 1 && intInput != 2 {
		return fmt.Errorf("Invalid Input, the input you sent was %v the valid inputs are 1 or 2", intInput)
	}
	url := fmt.Sprintf("http://%s/aj.html?a=command&cmd=x%sAVx1", vs.Address, input)
	_, gerr := vs.make2x1request(ctx, url)
	if gerr != nil {
		return fmt.Errorf("An error occured while making the call: %w", gerr)
	}
	return nil
}
