package atgain60

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestGetVolumeParallel(t *testing.T) {
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

	wg := &sync.WaitGroup{}
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()

			vols, err := amp.Volumes(ctx, []string{})
			is.NoErr(err)
			is.True(len(vols) == 1)
		}()

		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
}

func TestGetVolumeRepeat(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username:     "admin",
		Password:     "Atlona", // default username/password
		Address:      "TMCB-173-DSP1.byu.edu",
		RequestDelay: 500 * time.Millisecond,
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
			if err != nil && strings.Contains(err.Error(), "would exceed context deadline") {
				return
			}

			is.NoErr(err)
			is.True(len(vols) == 1)
		}
	}
}

func TestSetVolumeRepeat(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username:     "admin",
		Password:     "Atlona", // default username/password
		Address:      "TMCB-173-DSP1.byu.edu",
		RequestDelay: 500 * time.Millisecond,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			vol := rand.Intn(101)

			err := amp.SetVolume(ctx, "", vol)
			if err != nil && strings.Contains(err.Error(), "would exceed context deadline") {
				return
			}

			is.NoErr(err)

			vols, err := amp.Volumes(ctx, []string{})
			if err != nil && strings.Contains(err.Error(), "would exceed context deadline") {
				return
			}

			is.NoErr(err)
			is.True(len(vols) == 1)
			is.True(vols[""] == vol)

		}
	}
}

func TestSetMuteRepeat(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username:     "admin",
		Password:     "Atlona", // default username/password
		Address:      "TMCB-173-DSP1.byu.edu",
		RequestDelay: 500 * time.Millisecond,
		Log:          zap.NewExample(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m := rand.Intn(2)

			err := amp.SetMute(ctx, "", m == 0)
			if err != nil && strings.Contains(err.Error(), "would exceed context deadline") {
				return
			}

			is.NoErr(err)

			mutes, err := amp.Mutes(ctx, []string{})
			if err != nil && strings.Contains(err.Error(), "would exceed context deadline") {
				return
			}

			is.NoErr(err)
			is.True(len(mutes) == 1)
			is.True(mutes[""] == (m == 0))

		}
	}
}
