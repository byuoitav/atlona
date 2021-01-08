package atomeps62

import (
	"context"
	"fmt"
	"strconv"
)

//AudioVideoInputs .
func (vs *AtlonaVideoSwitcher6x2) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	body := `{ "getConfig": { "video": { "vidOut": { "hdmiOut": {}}}}}`

	config, err := vs.getConfig(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("unable to get config: %w", err)
	}

	inputs := make(map[string]string)
	if config.Video.VidOut.HdmiOut.Mirror.Status {
		inputs["mirror"] = strconv.Itoa(config.Video.VidOut.HdmiOut.Mirror.VideoSrc)
	} else {
		inputs["hdmiOutA"] = strconv.Itoa(config.Video.VidOut.HdmiOut.HdmiOutA.VideoSrc)
		inputs["hdmiOutB"] = strconv.Itoa(config.Video.VidOut.HdmiOut.HdmiOutB.VideoSrc)
	}

	return inputs, nil
}

//SetAudioVideoInput .
func (vs *AtlonaVideoSwitcher6x2) SetAudioVideoInput(ctx context.Context, output, input string) error {
	in, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("input must be an int: %w", err)
	}

	body := fmt.Sprintf(`{ "setConfig": { "video": { "vidOut": { "hdmiOut": { "%s": { "videoSrc": %v }}}}}}`, output, in)
	if err := vs.setConfig(ctx, body); err != nil {
		return fmt.Errorf("unable to set config: %w", err)
	}

	return nil
}
