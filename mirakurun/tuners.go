package mirakurun

import (
	"context"
	"github.com/pkg/errors"
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
		StreamInfo *map[int]struct {
			Packet int `json:"packet"`
			Drop   int `json:"drop"`
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
		return nil, errors.Wrap(err, "failed to create new request")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dispatch request")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Errorf("non-success status code %d", resp.StatusCode)
	}

	var tuners TunersResponse
	if err := decodeBody(resp, &tuners); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return &tuners, nil
}
