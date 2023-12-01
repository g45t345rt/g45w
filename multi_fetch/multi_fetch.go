package multi_fetch

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Fetch(url string) (*http.Response, error) {
	if strings.HasPrefix(url, "http") {
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
