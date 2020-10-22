package athdvs210u

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AtlonaVideoSwitcher2x1 struct {
	Username string
	Password string
	Address  string
}

type wallPlateStruct struct {
	LoginUr   int    `json:"login_ur"`
	LoginUser string `json:"login_user"`
	Inp       int    `json:"inp"`
	Asw       int    `json:"asw"`
	Preport   int    `json:"preport"`
	Aswtime   int    `json:"aswtime"`
	HDMIAud   int    `json:"HDMIAud"`
	HDCPSet   []int  `json:"HDCPSet"`
}

func (vs *AtlonaVideoSwitcher2x1) make2x1request(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error when creting the request: %w", err)
	}
	req = req.WithContext(ctx)
	res, gerr := http.DefaultClient.Do(req)
	if gerr != nil {
		return nil, fmt.Errorf("error when making call: %w", gerr)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error when making call: %w", gerr)
	}
	return body, nil
}
