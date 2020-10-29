package atjuno451hdbt

import (
	"context"
	"fmt"
)

// Info is the response from the switcher for the info page
type SysInfo struct {
	SystemInfo []string      `json:"info_val1"`
	VideoInfo  []interface{} `json:"info_val2"`
	LoggedIn   int           `json:"login_ur"`
}

type Info struct {
	ModelName       string
	FirmwareVersion string
}

// GetHardwareInfo returns a hardware info struct
func (vs *AtlonaVideoSwitcher4x1) Info(ctx context.Context) (interface{}, error) {
	var toReturn Info

	var info SysInfo
	err := getPage(ctx, vs.Address, infoPage, &info)
	if err != nil {
		return toReturn, fmt.Errorf("unable to get hardware info: %w", err)
	}

	// fill in the hwinfo
	if len(info.SystemInfo) >= 1 {
		toReturn.ModelName = info.SystemInfo[0]
	}

	if len(info.SystemInfo) >= 2 {
		toReturn.FirmwareVersion = info.SystemInfo[1]
	}

	return toReturn, nil
}

func (vs *AtlonaVideoSwitcher4x1) Healthy(ctx context.Context) error {
	_, err := vs.AudioVideoInputs(ctx)
	if err != nil {
		return fmt.Errorf("unable to get inputs (not healthy): %s", err)
	}

	return nil
}
