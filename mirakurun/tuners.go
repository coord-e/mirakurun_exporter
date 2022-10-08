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
	"fmt"
)

type TunersResponse []struct {
	Index   int      `json:"index"`
	Name    string   `json:"name"`
	Types   []string `json:"types"`
	Command string   `json:"command"`
	PID     int      `json:"pid"`
	Users   []struct {
		ID             string  `json:"id"`
		Priority       int     `json:"priority"`
		Agent          *string `json:"agent"`
		URL            *string `json:"url"`
		DisableDecoder *bool   `json:"disableDecoder"`
		StreamSetting  *struct {
			Channel struct {
				Type       string  `json:"type"`
				Channel    string  `json:"channel"`
				Name       string  `json:"name"`
				Satelite   *string `json:"satelite"`
				ServiceID  *int    `json:"serviceId"`
				Space      *int    `json:"space"`
				Freq       *int    `json:"freq"`
				Polarity   *string `json:"polarity"`
				TSMFRelTS  *int    `json:"tsmfRelTs"`
				IsDisabled *bool   `json:"isDisabled"`
			} `json:"channel"`
			NetworkID *int  `json:"networkId"`
			ServiceID *int  `json:"serviceId"`
			EventID   *int  `json:"eventId"`
			NoProvide *bool `json:"noProvide"`
			ParseEIT  *bool `json:"parseEIT"`
			ParseSDT  *bool `json:"parseSDT"`
			ParseNIT  *bool `json:"parseNIT"`
		} `json:"streamSetting"`
		StreamInfo *map[uint16]struct {
			Packet int64 `json:"packet"`
			Drop   int64 `json:"drop"`
		} `json:"streamInfo"`
	} `json:"users"`
	IsAvailable bool `json:"isAvailable"`
	IsRemote    bool `json:"isRemote"`
	IsFree      bool `json:"isFree"`
	IsUsing     bool `json:"isUsing"`
	IsFault     bool `json:"isFault"`
}

func (c *Client) GetTuners(ctx context.Context) (*TunersResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/api/tuners", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to dispatch request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-success status code %d", resp.StatusCode)
	}

	var tuners TunersResponse
	if err := decodeBody(resp, &tuners); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &tuners, nil
}
