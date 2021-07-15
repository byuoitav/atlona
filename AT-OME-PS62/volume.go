package atomeps62

import (
	"context"
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

	// Get the digital audio out on zoneOut1 volume
	if config.Audio.AudOut.ZoneOut1.AudioVol < -50 {
		vols["zoneOut1-Digital"] = 0
	} else {
		vols["zoneOut1-Digital"] = 2 * (config.Audio.AudOut.ZoneOut1.AudioVol + 50)
	}

	// Get the analog audio out on zoneOut1 volume
	if config.Audio.AudOut.ZoneOut1.AnalogOut.AudioVol < -50 {
		vols["zoneOut1-Analog"] = 0
	} else {
		vols["zoneOut1-Analog"] = 2 * (config.Audio.AudOut.ZoneOut1.AnalogOut.AudioVol + 50)
	}

	// Get the digital audio out on zoneOut2 volume
	if config.Audio.AudOut.ZoneOut2.AudioVol < -50 {
		vols["zoneOut2-Digital"] = 0
	} else {
		vols["zoneOut2-Digital"] = 2 * (config.Audio.AudOut.ZoneOut2.AudioVol + 50)
	}

	// Get the analog audio out on zoneOut2 volume
	if config.Audio.AudOut.ZoneOut2.AnalogOut.AudioVol < -50 {
		vols["zoneOut1-Analog"] = 0
	} else {
		vols["zoneOut1-Analog"] = 2 * (config.Audio.AudOut.ZoneOut2.AnalogOut.AudioVol + 50)
	}

	return vols, nil
}

//SetVolume .
func (vs *AtlonaVideoSwitcher6x2) SetVolume(ctx context.Context, block string, level int) error {
	zblock := ""
	if block == "zoneOut1" || block == "zoneOut2" {

		// Atlona volume levels are from -90 to 10 and the number we receive is 0-100
		// If volume level is supposed to be zero set it -90 on atlona
		if level == 0 {
			level = -90
		} else {
			convertedVolume := -50 + math.Round(float64(level/2))
			level = int(convertedVolume)
		}

		// Set digital and analog audio together for the audio block
		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "audioVol": %d, "analogOut":{"audioVol": %d }}}}}}`, block, level, level)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil
	} else if block == "zoneOut1Analog" || block == "zoneOut2Analog" {
		if block == "zoneOut1Analog" {
			zblock = "zoneOut1"
		} else {
			zblock = "zoneOut2"
		}
		// Atlona volume levels are from -90 to 10 and the number we receive is 0-100
		// If volume level is supposed to be zero set it -90 on atlona
		if level == 0 {
			level = -90
		} else {
			convertedVolume := -50 + math.Round(float64(level/2))
			level = int(convertedVolume)
		}

		// Set digital and analog audio together for the audio block
		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "analogOut":{"audioVol": %d }}}}}}`, zblock, level)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil
	} else if block == "zoneOut1Digital" || block == "zoneOut2Digital" {
		if block == "zoneOut1Digital" {
			zblock = "zoneOut1"
		} else {
			zblock = "zoneOut2"
		}
		// Atlona volume levels are from -90 to 10 and the number we receive is 0-100
		// If volume level is supposed to be zero set it -90 on atlona
		if level == 0 {
			level = -90
		} else {
			convertedVolume := -50 + math.Round(float64(level/2))
			level = int(convertedVolume)
		}

		// Set digital and analog audio together for the audio block
		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "audioVol": %d }}}}}`, zblock, level)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil
	} else {
		return fmt.Errorf("Unable to set config: Block %v is not a valid block", block)
	}
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
	mutes["zoneOut1-Analog"] = config.Audio.AudOut.ZoneOut1.AnalogOut.AudioMute
	mutes["zoneOut1-Digital"] = config.Audio.AudOut.ZoneOut1.VideoOut.AudioMute
	mutes["zoneOut2-Analog"] = config.Audio.AudOut.ZoneOut2.AnalogOut.AudioMute
	mutes["zoneOut2-Digital"] = config.Audio.AudOut.ZoneOut2.VideoOut.AudioMute

	return mutes, nil
}

//SetMute for all of the audio objects within the block
func (vs *AtlonaVideoSwitcher6x2) SetMute(ctx context.Context, block string, muted bool) error {
	zblock := ""
	if block == "zoneOut1" || block == "zoneOut2" {
		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "videoOut": { "audioMute": %t }, "analogOut": { "audioMute": %t }}}}}}`, block, muted, muted)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil
	} else if block == "zoneOut1Analog" || block == "zoneOut2Analog" {
		if block == "zoneOut1Analog" {
			zblock = "zoneOut1"
		} else {
			zblock = "zoneOut2"
		}

		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "analogOut": { "audioMute": %t }}}}}}`, zblock, muted)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil

	} else if block == "zoneOut1Digital" || block == "zoneOut2Digital" {
		if block == "zoneOut1Digital" {
			zblock = "zoneOut1"
		} else {
			zblock = "zoneOut2"
		}

		body := fmt.Sprintf(`{ "setConfig": { "audio": { "audOut": { "%s": { "videoOut": { "audioMute": %t }}}}}}`, zblock, muted)
		if err := vs.setConfig(ctx, body); err != nil {
			return fmt.Errorf("unable to set config: %w", err)
		}

		return nil

	} else {
		return fmt.Errorf("Unable to set config: Block %v is not a valid block", block)
	}

}
