package atgain60

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestVolume(t *testing.T) {
	t.SkipNow()
	is := is.New(t)

	amp := &Amp{
		Username: "admin",
		Password: "Atlona", // default username/password
		Address:  "TMCB-173-DSP1.byu.edu",
		Log:      zap.NewNop(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	vols, err := amp.Volumes(ctx, []string{})
	fmt.Printf("vols: %+v\n", vols)
	is.NoErr(err)
	is.True(len(vols) == 1)
}

func TestRepeatVolume(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username:     "admin",
		Password:     "Atlona", // default username/password
		Address:      "TMCB-173-DSP1.byu.edu",
		RequestDelay: 1 * time.Second,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			vols, err := amp.Volumes(ctx, []string{})
			is.NoErr(err)
			is.True(len(vols) == 1)
		}
	}
}
