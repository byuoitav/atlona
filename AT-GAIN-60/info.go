package atgain60

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type status struct {
	Model         string `json:"101"`
	Firmware      string `json:"102"`
	MACAddress    string `json:"103"`
	SerialNumber  string `json:"104"`
	OperatingTime string `json:"106"`
}

func (a *Amp) Info(ctx context.Context) (interface{}, error) {
	url := fmt.Sprintf("http://%s/action=devicestatus_get&r=%s", a.Address, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	body, err := a.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("unable to do request: %w", err)
	}

	var status status
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("unable to decode response: %w", err)
	}

	return status, nil
}

func (a *Amp) Healthy(ctx context.Context) error {
	_, err := a.Info(ctx)
	return err
}
