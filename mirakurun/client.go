// Copyright 2021 coord_e
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mirakurun

import (
	"context"
	"encoding/json"
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
		return nil, fmt.Errorf("missing URL")
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
