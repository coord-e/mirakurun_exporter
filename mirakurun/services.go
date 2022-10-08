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

type ServicesResponse []struct {
	ID                 int64  `json:"id"`
	ServiceID          int    `json:"serviceId"`
	NetworkID          int    `json:"networkId"`
	Name               string `json:"name"`
	Type               int    `json:"type"`
	LogoID             *int   `json:"logoId"`
	RemoteControlKeyID *int   `json:"remoteControlKeyId"`
	EpgReady           *bool  `json:"epgReady"`
	EpgUpdatedAt       *int64 `json:"epgUpdatedAt"`
	Channel            *struct {
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
	HasLogoData *bool `json:"hasLogoData"`
}

func (c *Client) GetServices(ctx context.Context) (*ServicesResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/api/services", nil)
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

	var services ServicesResponse
	if err := decodeBody(resp, &services); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &services, nil
}
