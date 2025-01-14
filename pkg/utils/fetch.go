package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/goodwithtech/dockle/pkg/log"
)

var versionPattern = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

func fetchURL(ctx context.Context, url string, cookie *http.Cookie) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.AddCookie(cookie)
	resp, err := (&http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 3,
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return nil, fmt.Errorf("HTTP error code : %d, url : %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

func FetchLatestVersion(ctx context.Context) (version string, err error) {
	log.Logger.Debug("Fetch latest version from github")
	body, err := fetchURL(
		ctx,
		"https://github.com/goodwithtech/dockle/releases/latest",
		&http.Cookie{Name: "user_session", Value: "guard"},
	)
	if err != nil {
		return "", err
	}
	if versionMatched := versionPattern.FindString(string(body)); versionMatched != "" {
		return versionMatched, nil
	}
	return "", errors.New("not found version patterns")
}
