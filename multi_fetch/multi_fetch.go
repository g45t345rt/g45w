package multi_fetch

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/g45t345rt/g45w/app_db"
)

func Fetch(url string) (*http.Response, error) {
	if strings.HasPrefix(url, "ipfs://") {
		cId := strings.Replace(url, "ipfs://", "", -1)
		return IPFSFetch(cId)
	} else if strings.HasPrefix(url, "http") {
		return HttpFetch(url, 5*time.Second)
	} else {
		return nil, fmt.Errorf("url scheme not supported")
	}
}

func HttpFetch(url string, timeout time.Duration) (*http.Response, error) {
	client := new(http.Client)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func IPFSFetch(cId string) (*http.Response, error) {
	gateways, err := app_db.GetIPFSGateways(app_db.GetIPFSGatewaysParams{
		Active: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	for _, gateway := range gateways {
		res, err := gateway.Fetch(cId, 3*time.Second)
		if err != nil {
			continue
		}

		if res.StatusCode != 200 {
			continue
		}

		return res, nil
	}

	return nil, fmt.Errorf("unavailable")
}
