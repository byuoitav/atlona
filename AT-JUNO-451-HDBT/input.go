package atjuno451hdbt

import (
	"context"
	"fmt"
	"strconv"
)

// AudioVideoInputs returns the current input
func (vs *AtlonaVideoSwitcher4x1) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)

	var settings AVSettings
	err := getPage(ctx, vs.Address, avSettingsPage, &settings)
	if err != nil {
		return toReturn, fmt.Errorf("unable to get input: %w", err)
	}

	toReturn[""] = fmt.Sprintf("%v", settings.Input-1)
	return toReturn, nil
}

// SetAudioVideoInput changes the input on the given output to input
func (vs *AtlonaVideoSwitcher4x1) SetAudioVideoInput(ctx context.Context, output, input string) error {
	// atlona switchers are 1-based

	in, gerr := strconv.Atoi(input)
	if gerr != nil {
		return fmt.Errorf("unable to switch input on %s:%w", vs.Address, gerr)
	}

	in++

	// validate that input/output are valid numbers
	var settings AVSettings
	err := getPage(ctx, vs.Address, avSettingsPage, &settings)
	if err != nil {
		return fmt.Errorf("unable to switch input: %w", err)
	}

	if in > len(settings.HDCPSettings) || in <= 0 {
		return fmt.Errorf("unable to switch input on %s - input %s is out of range", vs.Address, input)
	}

	err = sendCommand(ctx, vs.Address, fmt.Sprintf("x%vAVx%v", in, 1))
	if err != nil {
		return fmt.Errorf("unable to switch input: %w", err)
	}

	return nil
}
