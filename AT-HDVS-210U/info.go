package athdvs210u

import (
	"context"
	"fmt"
)

//Info .
func (vs *AtlonaVideoSwitcher2x1) Info(ctx context.Context) (interface{}, error) {
	var info interface{}
	return info, fmt.Errorf("not currently implemented")
}

func (vs *AtlonaVideoSwitcher2x1) Healthy(ctx context.Context) error {
	_, err := vs.AudioVideoInputs(ctx)
	if err != nil {
		return fmt.Errorf("unable to get inputs (not healthy): %s", err)
	}

	return nil
}
