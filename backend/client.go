package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const openF1BaseURL = "https://api.openf1.org/v1"

type OpenF1Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewOpenF1Client() *OpenF1Client {
	return &OpenF1Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: openF1BaseURL,
	}
}

func (c *OpenF1Client) getJSON(ctx context.Context, path string, q url.Values, v any) error {
	u := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("openf1: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *OpenF1Client) GetCarData(ctx context.Context, f CarDataFilter) ([]CarData, error) {
	q := url.Values{}

	if f.DriverNumber != nil {
		q.Set("driver_number", strconv.Itoa(*f.DriverNumber))
	}
	if f.SessionKey != nil {
		q.Set("session_key", strconv.Itoa(*f.SessionKey))
	}
	if f.MinSpeed != nil {
		// OpenF1 uses comparison operators via query params.
		// The documented curl is:
		//   ...car_data?driver_number=55&session_key=9159&speed%3E=315
		// which is key "speed>" with value "315".
		q.Set("speed>", fmt.Sprintf("%.0f", *f.MinSpeed))
	}

	var out []CarData
	if err := c.getJSON(ctx, "/car_data", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *OpenF1Client) GetDrivers(ctx context.Context, f DriverFilter) ([]Driver, error) {
	q := url.Values{}

	if f.SessionKey != nil {
		q.Set("session_key", strconv.Itoa(*f.SessionKey))
	}
	if f.MeetingKey != nil {
		q.Set("meeting_key", strconv.Itoa(*f.MeetingKey))
	}
	if f.DriverNumber != nil {
		q.Set("driver_number", strconv.Itoa(*f.DriverNumber))
	}

	var out []Driver
	if err := c.getJSON(ctx, "/drivers", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}
