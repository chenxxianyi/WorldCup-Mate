package apifootball

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	if baseURL == "" {
		baseURL = "https://v3.football.api-sports.io"
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Transport: &http.Transport{Proxy: nil},
			Timeout:   15 * time.Second,
		},
	}
}

func (c *Client) TeamSquad(ctx context.Context, externalTeamID string) (*SquadResponse, error) {
	u, err := url.Parse(c.baseURL + "/players/squads")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("team", externalTeamID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apisports-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data SquadResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if hasAPIErrors(data.Errors) {
		return nil, fmt.Errorf("api-football returned errors: %v", data.Errors)
	}
	return &data, nil
}

func (c *Client) SearchTeams(ctx context.Context, keyword string) (*TeamSearchResponse, error) {
	u, err := url.Parse(c.baseURL + "/teams")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("search", keyword)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apisports-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data TeamSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if hasAPIErrors(data.Errors) {
		return nil, fmt.Errorf("api-football returned errors: %v", data.Errors)
	}
	return &data, nil
}

func (c *Client) FixtureLineups(ctx context.Context, fixtureID string) (*FixtureLineupsResponse, error) {
	u, err := url.Parse(c.baseURL + "/fixtures/lineups")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("fixture", fixtureID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apisports-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data FixtureLineupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if hasAPIErrors(data.Errors) {
		return nil, fmt.Errorf("api-football returned errors: %v", data.Errors)
	}
	return &data, nil
}

func (c *Client) SearchFixtures(ctx context.Context, date string, season int, teamID string) (*FixtureSearchResponse, error) {
	u, err := url.Parse(c.baseURL + "/fixtures")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("date", date)
	if season > 0 {
		q.Set("season", strconv.Itoa(season))
	}
	if strings.TrimSpace(teamID) != "" {
		q.Set("team", teamID)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apisports-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data FixtureSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if hasAPIErrors(data.Errors) {
		return nil, fmt.Errorf("api-football returned errors: %v", data.Errors)
	}
	return &data, nil
}

func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	const maxAttempts = 4
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := c.httpClient.Do(req.Clone(ctx))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}
		if resp.StatusCode != http.StatusTooManyRequests || attempt == maxAttempts {
			defer resp.Body.Close()
			return nil, fmt.Errorf("api-football returned status %d", resp.StatusCode)
		}

		wait := retryAfter(resp.Header.Get("Retry-After"))
		resp.Body.Close()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}
	return nil, fmt.Errorf("api-football request retry exhausted")
}

func retryAfter(value string) time.Duration {
	if seconds, err := strconv.Atoi(strings.TrimSpace(value)); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return 65 * time.Second
}

func hasAPIErrors(value interface{}) bool {
	switch errors := value.(type) {
	case nil:
		return false
	case []interface{}:
		return len(errors) > 0
	case map[string]interface{}:
		return len(errors) > 0
	case string:
		return strings.TrimSpace(errors) != ""
	default:
		return strings.TrimSpace(fmt.Sprint(errors)) != "" && strings.TrimSpace(fmt.Sprint(errors)) != "[]"
	}
}
