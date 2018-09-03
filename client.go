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

type Client struct {
	url   string
	token string
}

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

func (c *Client) GetIssueEvents(ctx context.Context, repoID, issueNumber int) ([]Event, error) {

	url := fmt.Sprintf("%s/p1/repositories/%d/issues/%d/events", c.url, repoID, issueNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request %q", url)
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authentication-Token", c.token)

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

	var payload []event

	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, errors.Errorf("failed to unmarshal payload %q %s %v", url, body, err)
	}

	events := make([]Event, len(payload))

	for i, p := range payload {
		switch p.Type {
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
			return nil, errors.Errorf("unknown event type %q", p.Type)
		}

	}

	return events, nil
}
