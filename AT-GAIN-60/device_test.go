package atgain60

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestLogin(t *testing.T) {
	is := is.New(t)

	amp := &Amp{
		Username: "admin",
		Password: "Atlona", // default username/password
		Address:  "TMCB-173-DSP1.byu.edu",
		Log:      zap.NewNop(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	is.NoErr(amp.login(ctx))
}
