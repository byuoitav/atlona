package atgain60

/*

// Info gets the current amp status
func (a *Amp60) Info(ctx context.Context) (interface{}, error) {
	resp, err := a.sendReq(ctx, "devicestatus_get")
	if err != nil {
		return nil, fmt.Errorf("unable to get info: %w", err)
	}
	var info AmpStatus
	err = json.Unmarshal(resp, &info)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal into AmpStatus: %w", err)
	}
	return info, nil
}

func (a *Amp60) Healthy(ctx context.Context) error {
	_, err := a.Volumes(ctx, []string{})
	if err != nil {
		return fmt.Errorf("unable to get volume (not healthy): %s", err)
	}

	return nil
}
*/
