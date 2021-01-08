package atomeps62

import (
	"context"
	"errors"
	"fmt"
	"math"
)

//Volumes .
func (vs *AtlonaVideoSwitcher6x2) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	body := `{ "getConfig": { "audio": { "audOut": {}}}}`

	config, err := vs.getConfig(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("unable to get config: %w", err)
	}

	// always return all of the blocks, regardless of `blocks`
	// (since we don't have to do any extra work)
	vols := make(map[string]int)

	// zoneOut1 volume
	if config.Audio.AudOut.ZoneOut1.AudioVol < -50 {
		vols["zoneOut1"] = 0
	} else {
		vols["zoneOut1"] = 2 * (config.Audio.AudOut.ZoneOut1.AudioVol + 50)
	}

	// zoneOut2 volume
	if config.Audio.AudOut.ZoneOut2.AudioVol < -50 {
		vols["zoneOut2"] = 0
	} else {
		vols["zoneOut2"] = 2 * (config.Audio.AudOut.ZoneOut2.AudioVol + 50)
	}

	return vols, nil
}

//SetVolume .
func (vs *AtlonaVideoSwitcher6x2) SetVolume(ctx context.Context, block string, level int) error {
	if block != "zoneOut1" && block != "zoneOut2" {
		return errors.New("invalid block")
	}

	// Atlona volume levels are from -90 to 10 and the number we receive is 0-100
	// If volume level is supposed to be zero set it -90 on atlona
	if level == 0 {
		level = -90
	} else {
		convertedVolume := -50 + math.Round(float64(level/2))
		level = int(convertedVolume)
	}

	body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "audioVol": %d }}}}}`, block, level)
	if err := vs.setConfig(ctx, body); err != nil {
		return fmt.Errorf("unable to set config: %w", err)
	}

	return nil
}

//Mutes .
func (vs *AtlonaVideoSwitcher6x2) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	body := `{ "getConfig": { "audio": { "audOut": {}}}}`

	config, err := vs.getConfig(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("unable to get config: %w", err)
	}

	// always return all of the blocks, regardless of `blocks`
	// (since we don't have to do any extra work)
	mutes := make(map[string]bool)
	mutes["zoneOut1"] = config.Audio.AudOut.ZoneOut1.AnalogOut.AudioMute
	mutes["zoneOut2"] = config.Audio.AudOut.ZoneOut2.AnalogOut.AudioMute

	return mutes, nil
}

//SetMute .
func (vs *AtlonaVideoSwitcher6x2) SetMute(ctx context.Context, block string, muted bool) error {
	if block != "zoneOut1" && block != "zoneOut2" {
		return errors.New("invalid block")
	}

	body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "analogOut": { "audioMute": %t }}}}}}`, block, muted)
	if err := vs.setConfig(ctx, body); err != nil {
		return fmt.Errorf("unable to set config: %w", err)
	}

	return nil
}
