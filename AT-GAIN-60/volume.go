package atgain60

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type audio struct {
	Volume string `json:"608"`
	Muted  string `json:"609"`
}

func (a *Amp) Volumes(ctx context.Context, _ []string) (map[string]int, error) {
	url := fmt.Sprintf("http://%s/action=deviceaudio_get&r=%s", a.Address, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	body, err := a.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	var audio audio
	if err := json.Unmarshal(body, &audio); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	vol, err := strconv.Atoi(audio.Volume)
	if err != nil {
		return nil, fmt.Errorf("unable to convert volume: %w", err)
	}

	return map[string]int{"": vol}, nil
}

func (a *Amp) Mutes(ctx context.Context, _ []string) (map[string]bool, error) {
	url := fmt.Sprintf("http://%s/action=deviceaudio_get&r=%s", a.Address, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	body, err := a.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	var audio audio
	if err := json.Unmarshal(body, &audio); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return map[string]bool{"": audio.Muted == "1"}, nil
}

func (a *Amp) SetVolume(ctx context.Context, _ string, volume int) error {
	url := fmt.Sprintf("http://%s/action=deviceaudio_set&608=%d&r=%s", a.Address, volume, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	_, err = a.doReq(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}

	return nil
}

func (a *Amp) SetMute(ctx context.Context, _ string, muted bool) error {
	muteStr := "0"
	if muted {
		muteStr = "1"
	}

	url := fmt.Sprintf("http://%s/action=deviceaudio_set&609=%s&r=%s", a.Address, muteStr, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	_, err = a.doReq(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}

	return nil
}
