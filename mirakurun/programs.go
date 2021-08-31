package mirakurun

import (
	"context"

	"github.com/pkg/errors"
)

type ProgramsResponse []struct {
	ID          int64   `json:"id"`
	EventID     int     `json:"eventId"`
	ServiceID   int     `json:"serviceId"`
	NetworkID   int     `json:"networkId"`
	StartAt     int64   `json:"startAt"`
	Duration    int64   `json:"duration"`
	IsFree      bool    `json:"isFree"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Genres      *[]struct {
		Lv1 int `json:"lv1"`
		Lv2 int `json:"lv2"`
		Un1 int `json:"un1"`
		Un2 int `json:"un2"`
	} `json:"genres"`
	Video *struct {
		Type          string `json:"type"`
		Resolution    string `json:"resolution"`
		StreamContent int    `json:"streamContent"`
		ComponentType int    `json:"componentType"`
	} `json:"video"`
	Audio *struct {
		ComponentType int       `json:"componentType"`
		ComponentTag  *int      `json:"componentTag"`
		IsMain        *bool     `json:"isMain"`
		SamplingRate  int       `json:"samplingRate"`
		Langs         *[]string `json:"langs"`
	} `json:"audio"`
	Series *struct {
		ID          int    `json:"id"`
		Repeat      int    `json:"repeat"`
		Pattern     int    `json:"pattern"`
		ExpiresAt   int64  `json:"expiresAt"`
		Episode     int    `json:"episode"`
		LastEpisode int    `json:"lastEpisode"`
		Name        string `json:"name"`
	} `json:"series"`
	Extended     *map[string]string `json:"extended"`
	RelatedItems *[]struct {
		Type      *string `json:"type"`
		NetworkID *int    `json:"networkId"`
		ServiceID int     `json:"serviceId"`
		EventID   int     `json:"eventId"`
	} `json:"relatedItems"`
}

func (c *Client) GetPrograms(ctx context.Context) (*ProgramsResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/api/programs", nil)
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

	var programs ProgramsResponse
	if err := decodeBody(resp, &programs); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return &programs, nil
}
