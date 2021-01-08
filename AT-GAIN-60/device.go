package atgain60

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Amp represents an Atlona 60 watt amplifier
type Amp struct {
	Username string
	Password string
	Address  string
	Log      *zap.Logger

	RequestDelay time.Duration

	loginMu sync.Mutex
	once    sync.Once
	limiter *rate.Limiter
}

func (a *Amp) init() {
	a.limiter = rate.NewLimiter(rate.Every(a.RequestDelay), 1)
}

func (a *Amp) r() string {
	return fmt.Sprintf("%v", rand.Float64())
}

func (a *Amp) login(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/action=compare&701=%s&702=%s&r=%s", a.Address, a.Username, a.Password, a.r())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	client := &http.Client{
		Transport: &transport{},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	var login struct {
		Login string `json:"Login"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&login); err != nil {
		return fmt.Errorf("unable to decode response: %w", err)
	}

	if login.Login != "True" {
		return fmt.Errorf("unexpected login result: %s", login.Login)
	}

	return nil
}

func (a *Amp) doReq(req *http.Request) ([]byte, error) {
	a.once.Do(a.init)

	if err := a.limiter.Wait(req.Context()); err != nil {
		return nil, fmt.Errorf("unable to wait for ratelimit: %w", err)
	}

	// probably not the best solution...
	a.loginMu.Lock()
	defer a.loginMu.Unlock()

	login := false

	for {
		if login {
			a.Log.Info("Logging in to amp")

			if err := a.login(req.Context()); err != nil {
				return nil, err
			}

			a.Log.Info("Successfully logged in; retrying previous request")
		}

		client := &http.Client{
			Transport: &transport{},
		}

		a.Log.Debug("Doing request", zap.String("url", req.URL.String()))

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("unable to do request: %w", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read body")
		}

		a.Log.Debug("Response", zap.ByteString("body", body))

		if login {
			return body, nil
		}

		// see if we are unauthorized
		var auth struct {
			LoggedIn string `json:"909"`
		}

		_ = json.Unmarshal(body, &auth)
		if auth.LoggedIn == "" {
			return body, nil
		}

		login = true
	}
}
