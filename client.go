package zenhub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const apiURL = "https://api.zenhub.io"

// Client represents a ZenHub API HTTP client.
type Client struct {
	url   string
	token string
}

// NewClient returns a new ZenHub client for the authentication token passed in.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("invalid token")
	}

	c := &Client{
		url:   apiURL,
		token: token,
	}

	return c, nil
}

// GetIssueEvents returns all events available for an issue.
func (c *Client) GetIssueEvents(ctx context.Context, repoID, issueNumber int) ([]Event, error) {
	url := fmt.Sprintf("%s/p1/repositories/%d/issues/%d/events", c.url, repoID, issueNumber)

	body, err := get(ctx, url, c.token)
	if err != nil {
		return nil, err
	}

	var payload []event

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, errors.Errorf("failed to unmarshal payload %q %s %v", url, body, err)
	}

	events := make([]Event, len(payload))

	for i, p := range payload {
		switch t := p.Type; t {
		case "estimateIssue":
			event := EstimateIssueEvent{
				UserID:       p.UserID,
				CreatedAt:    p.CreatedAt,
				FromEstimate: p.FromEstimate,
				ToEstimate:   p.ToEstimate,
			}

			events[i] = event
		case "transferIssue":
			event := TransferIssueEvent{
				UserID:       p.UserID,
				CreatedAt:    p.CreatedAt,
				FromPipeline: p.FromPipeline,
				ToPipeline:   p.ToPipeline,
			}

			events[i] = event
		default:
			return nil, ErrUnknownEventType{t: t}
		}

	}

	return events, nil
}

// GetBoard returns the board information for a repository.
func (c *Client) GetBoard(ctx context.Context, repoID int) (Board, error) {
	url := fmt.Sprintf("%s/p1/repositories/%d/board", c.url, repoID)

	board := Board{}

	body, err := get(ctx, url, c.token)
	if err != nil {
		return board, err
	}

	err = json.Unmarshal(body, &board)
	if err != nil {
		return board, errors.Errorf("failed to unmarshal payload %q %s", url, body)
	}

	return board, nil
}

func get(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request %q", url)
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authentication-Token", token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send request %q", url)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read response body %q", url)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to send request [%d] %q %s", res.StatusCode, url, body)
	}

	return body, nil
}
