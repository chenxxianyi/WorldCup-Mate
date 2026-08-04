package footballdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func NewClient(baseURL, apiToken string) *Client {
	if baseURL == "" {
		baseURL = "https://api.football-data.org/v4"
	}
	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiToken: apiToken,
		httpClient: &http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   15 * time.Second,
		},
	}
}

func (c *Client) CompetitionMatches(ctx context.Context, competitionCode string, season int) (*MatchesResponse, error) {
	u, err := url.Parse(c.baseURL + "/competitions/" + url.PathEscape(competitionCode) + "/matches")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if season > 0 {
		q.Set("season", fmt.Sprintf("%d", season))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("football-data returned status %d", resp.StatusCode)
	}

	var data MatchesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// CompetitionStandings fetches the standings (TOTAL / HOME / AWAY) of a
// competition season. New method for the league extension; the existing
// CompetitionMatches is untouched.
func (c *Client) CompetitionStandings(ctx context.Context, competitionCode string, season int) (*StandingsResponse, error) {
	u, err := url.Parse(c.baseURL + "/competitions/" + url.PathEscape(competitionCode) + "/standings")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if season > 0 {
		q.Set("season", fmt.Sprintf("%d", season))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("football-data returned status %d", resp.StatusCode)
	}

	var data StandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
