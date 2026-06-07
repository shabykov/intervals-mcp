package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config struct {
	ApiURL    string
	ApiKey    string
	AthleteID string
}

type Client struct {
	apiURL    string
	apiKey    string
	athleteID string

	client *http.Client
}

func NewClient(config Config) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &Client{
		apiURL:    config.ApiURL,
		apiKey:    config.ApiKey,
		athleteID: config.AthleteID,
		client:    client,
	}
}

func (c *Client) ListActivities(ctx context.Context, _ *mcp.CallToolRequest, in RangeIn) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	q.Set("oldest", in.Oldest)
	q.Set("newest", in.Newest)
	q.Set("fields", "id,name,start_date_local,type,distance,moving_time,"+
		"icu_training_load,icu_intensity,average_watts,icu_weighted_avg_watts,"+
		"average_heartrate,max_heartrate,icu_ftp,icu_atl,icu_ctl")
	body, err := c.get(ctx, "/athlete/"+c.athleteID+"/activities", q)
	if err != nil {
		return nil, nil, err
	}
	return text(body), nil, nil
}

func (c *Client) GetActivity(ctx context.Context, _ *mcp.CallToolRequest, in ActivityIn) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if in.Intervals {
		q.Set("intervals", "true")
	}
	body, err := c.get(ctx, "/activity/"+in.ID, q)
	if err != nil {
		return nil, nil, err
	}
	return text(body), nil, nil
}

func (c *Client) GetStreams(ctx context.Context, _ *mcp.CallToolRequest, in StreamsIn) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if in.Types != "" {
		q.Set("types", in.Types)
	}
	body, err := c.get(ctx, "/activity/"+in.ID+"/streams", q)
	if err != nil {
		return nil, nil, err
	}
	return text(body), nil, nil
}

func (c *Client) GetWellness(ctx context.Context, _ *mcp.CallToolRequest, in RangeIn) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	q.Set("oldest", in.Oldest)
	q.Set("newest", in.Newest)
	body, err := c.get(ctx, "/athlete/"+c.athleteID+"/wellness", q)
	if err != nil {
		return nil, nil, err
	}
	return text(body), nil, nil
}

func (c *Client) GetAthlete(ctx context.Context, _ *mcp.CallToolRequest, _ Empty) (*mcp.CallToolResult, any, error) {
	body, err := c.get(ctx, "/athlete/"+c.athleteID, nil)
	if err != nil {
		return nil, nil, err
	}
	return text(body), nil, nil
}

func (c *Client) get(ctx context.Context, path string, q url.Values) (string, error) {
	u := c.apiURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("API_KEY", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("intervals API %d: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

func text(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: s}}}
}

type RangeIn struct {
	Oldest string `json:"oldest" jsonschema:"начало периода, YYYY-MM-DD"`
	Newest string `json:"newest" jsonschema:"конец периода, YYYY-MM-DD"`
}
type ActivityIn struct {
	ID        string `json:"id" jsonschema:"ID активности"`
	Intervals bool   `json:"intervals" jsonschema:"добавить разбивку по интервалам/лапам"`
}
type StreamsIn struct {
	ID    string `json:"id" jsonschema:"ID активности"`
	Types string `json:"types" jsonschema:"список рядов через запятую, напр. watts,heartrate,cadence,velocity_smooth,altitude"`
}
type Empty struct{}
