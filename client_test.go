package zenhub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestNewClient(t *testing.T) {
	tcs := []struct {
		scenario string
		token    string
		success  bool
		err      error
	}{
		{
			scenario: "with a valid token",
			token:    "ABC123",
			success:  true,
		},
		{
			scenario: "with an empty token",
			token:    "",
			success:  false,
			err:      errors.New("invalid token"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			c, err := NewClient(tc.token)

			if tc.success {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if c.url != apiURL {
					t.Fatalf("invalid API URL %q", c.url)
				}
			} else {
				if err.Error() != tc.err.Error() {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
			}

		})
	}

}

func TestGetIssueEvents(t *testing.T) {
	c, err := NewClient("ABC")
	if err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		scenario string
		status   int
		response string
		success  bool
		events   []Event
		err      error
	}{
		{
			scenario: "successful request",
			status:   200,
			response: `
		[
		{
			"user_id": 16717,
			"type": "estimateIssue",
			"created_at": "2015-12-11T19:43:22.296Z",
			"from_estimate": {
				"value": 8
			}
		},
		{
			"user_id": 16717,
			"type": "estimateIssue",
			"created_at": "2015-12-11T18:43:22.296Z",
			"from_estimate": {
				"value": 4
			},
			"to_estimate": {
				"value": 8
			}
		},
		{
			"user_id": 16717,
			"type": "estimateIssue",
			"created_at": "2015-12-11T13:43:22.296Z",
			"to_estimate": {
				"value": 4
			}
		},
		{
			"user_id": 16717,
			"type": "transferIssue",
			"created_at": "2015-12-11T12:43:22.296Z",
			"from_pipeline": {
				"name": "Backlog"
			},
			"to_pipeline": {
				"name": "In progress"
			}
		},
		{
			"user_id": 16717,
			"type": "transferIssue",
			"created_at": "2015-12-11T11:43:22.296Z",
			"to_pipeline": {
				"name": "Backlog"
			}
		}
		]
			`,
			success: true,
			events: []Event{
				EstimateIssueEvent{
					UserID:    16717,
					CreatedAt: time.Date(2015, time.December, 11, 19, 43, 22, 296000000, time.UTC),
					FromEstimate: &Estimate{
						Value: 8,
					},
				},
				EstimateIssueEvent{
					UserID:    16717,
					CreatedAt: time.Date(2015, time.December, 11, 18, 43, 22, 296000000, time.UTC),
					FromEstimate: &Estimate{
						Value: 4,
					},
					ToEstimate: &Estimate{
						Value: 8,
					},
				},
				EstimateIssueEvent{
					UserID:    16717,
					CreatedAt: time.Date(2015, time.December, 11, 13, 43, 22, 296000000, time.UTC),
					ToEstimate: &Estimate{
						Value: 4,
					},
				},
				TransferIssueEvent{
					UserID:    16717,
					CreatedAt: time.Date(2015, time.December, 11, 12, 43, 22, 296000000, time.UTC),
					FromPipeline: &Pipeline{
						Name: "Backlog",
					},
					ToPipeline: &Pipeline{
						Name: "In progress",
					},
				},
				TransferIssueEvent{
					UserID:    16717,
					CreatedAt: time.Date(2015, time.December, 11, 11, 43, 22, 296000000, time.UTC),
					ToPipeline: &Pipeline{
						Name: "Backlog",
					},
				},
			},
		},
		{
			scenario: "server error",
			status:   500,
			response: "Internal Server Error",
			success:  false,
			err:      errors.New("failed to send request [500]"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Header().Set("Content-type", "application/json")
				w.Write([]byte(tc.response))
			}))
			defer srv.Close()

			c.url = srv.URL

			ctx := context.Background()

			events, err := c.GetIssueEvents(ctx, 123, 456)

			if tc.success {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if !reflect.DeepEqual(tc.events, events) {
					t.Fatalf("expected events to eq %v, got %v", tc.events, events)
				}
			} else {
				if !strings.HasPrefix(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
			}

		})
	}
}
