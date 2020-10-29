package atgain60

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/byuoitav/common/log"
)

// Volumes gets the current volume
func (a *Amp60) Volumes(ctx context.Context, blocks []string) (map[string]int, error) {
	resp, err := a.sendReq(ctx, "deviceaudio_get")
	if err != nil {
		return map[string]int{"": -1}, fmt.Errorf("unable to get volume: %w", err)
	}
	var info AmpAudio
	var test map[string]interface{}
	json.Unmarshal(resp, &test)
	for key, value := range test {
		log.L.Debug(key, value.(string))
	}
	log.L.Debug("Testing our json: %v", test)

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return map[string]int{"": -1}, fmt.Errorf("unable to unmarshal into AmpVolume in GetVolume: %w", err)
	}
	toReturn, err := strconv.Atoi(info.Volume)
	if err != nil {
		return map[string]int{"": -1}, fmt.Errorf("Volume is empty")
	}
	return map[string]int{"": toReturn}, nil
}

// Mutes gets the current mute status
func (a *Amp60) Mutes(ctx context.Context, blocks []string) (map[string]bool, error) {
	resp, err := a.sendReq(ctx, "deviceaudio_get")
	if err != nil {

		return map[string]bool{"": false}, fmt.Errorf("unable to get muted: %w", err)
	}
	var info AmpAudio
	err = json.Unmarshal(resp, &info)
	if err != nil {
		return map[string]bool{"": false}, fmt.Errorf("unable to unmarshal into AmpVolume in GetMuted: %w", err)
	}
	if info.Muted == "1" {
		return map[string]bool{"": true}, nil
	}
	return map[string]bool{"": false}, nil
}

// SetVolume sets the volume on the amp
func (a *Amp60) SetVolume(ctx context.Context, block string, volume int) error {
	_, err := a.sendReq(ctx, fmt.Sprintf("deviceaudio_set&608=%v", volume))
	if err != nil {
		return fmt.Errorf("unable to set volume: %w", err)
	}
	return nil
}

// SetMuted sets the current mute status on the amp
func (a *Amp60) SetMute(ctx context.Context, block string, muted bool) error {
	// open a connection with the dsp, set the muted status on block...
	mutedString := "0"
	if muted {
		mutedString = "1"
	}
	_, err := a.sendReq(ctx, fmt.Sprintf("deviceaudio_set&609=%v", mutedString))
	if err != nil {
		return fmt.Errorf("unable to set muted: %w", err)
	}
	return nil
}
