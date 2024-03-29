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

type StatusResponse struct {
	Time    int64  `json:"time"`
	Version string `json:"version"`
	Process struct {
		Arch        string            `json:"arch"`
		Platform    string            `json:"platform"`
		Versions    map[string]string `json:"versions"`
		Env         map[string]string `json:"env"`
		PID         int               `json:"pid"`
		MemoryUsage struct {
			RSS          int64 `json:"rss"`
			HeapTotal    int64 `json:"heapTotal"`
			HeapUsed     int64 `json:"heapUsed"`
			External     int64 `json:"external"`
			ArrayBuffers int64 `json:"arrayBuffers"`
		} `json:"memoryUsage"`
	} `json:"process"`
	EPG struct {
		GatheringNetworks []int64 `json:"gatheringNetworks"`
		StoredEvents      int64   `json:"storedEvents"`
	} `json:"EPG"`
	RPCCount    *int `json:"rpcCount"` // available since 3.9.0-beta.0
	StreamCount struct {
		TunerDevice int `json:"tunerDevice"`
		TSFilter    int `json:"tsFilter"`
		Decoder     int `json:"decoder"`
	} `json:"streamCount"`
	ErrorCount struct {
		UncaughtException  int `json:"uncaughtException"`
		UnhandledRejection int `json:"unhandledRejection"`
		BufferOverflow     int `json:"bufferOverflow"`
		TunerDeviceRespawn int `json:"tunerDeviceRespawn"`
		DecoderRespawn     int `json:"decoderRespawn"`
	} `json:"errorCount"`
	TimerAccuracy struct {
		Last float64 `json:"last"`
		M1   struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"m1"`
		M5 struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"m5"`
		M15 struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"m15"`
	} `json:"timerAccuracy"`
}

func (c *Client) GetStatus(ctx context.Context) (*StatusResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/api/status", nil)
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

	var status StatusResponse
	if err := decodeBody(resp, &status); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &status, nil
}
