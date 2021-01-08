package atgain60

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestGetInfo(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username:     "admin",
		Password:     "Atlona", // default username/password
		Address:      "TMCB-173-DSP1.byu.edu",
		RequestDelay: 1000 * time.Millisecond,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := amp.Info(ctx)
	is.NoErr(err)

	s, ok := info.(status)
	is.True(ok)
	is.True(s.Model != "")
	is.True(s.Firmware != "")
	is.True(s.MACAddress != "")
	is.True(s.SerialNumber != "")
	is.True(s.OperatingTime != "")
}
