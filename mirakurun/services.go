package mirakurun

import (
	"context"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to create new request")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dispatch request")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Errorf("non-success status code %d", resp.StatusCode)
	}

	var services ServicesResponse
	if err := decodeBody(resp, &services); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return &services, nil
}
