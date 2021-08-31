package mirakurun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-cleanhttp"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "devel"
	}
	return strings.Trim(info.Main.Version, "v()")
}

var userAgent = fmt.Sprintf("MirakurunGoClient/%s (%s)", version(), runtime.Version())

type Client struct {
	URL           *url.URL
	HTTPClient    *http.Client
	DefaultHeader http.Header
	Logger        log.Logger
}

func NewClient(urlString string) (*Client, error) {
	if len(urlString) == 0 {
		return nil, errors.New("missing URL")
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	client := &Client{
		URL:           parsedURL,
		HTTPClient:    cleanhttp.DefaultClient(),
		DefaultHeader: make(http.Header),
		Logger:        log.NewNopLogger(),
	}
	client.DefaultHeader.Set("User-Agent", userAgent)
	client.DefaultHeader.Set("Accept", "application/json")

	return client, nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	for k, v := range c.DefaultHeader {
		req.Header[k] = v
	}

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
