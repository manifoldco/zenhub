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
			scenario: "unknown event type",
			status:   200,
			response: `
		[
		{
			"user_id": 16717,
			"type": "unknown event",
			"created_at": "2015-12-11T19:43:22.296Z",
			"from_estimate": {
				"value": 8
			}
		}
		]
		`,
			success: false,
			err:     errors.New("unknown event type"),
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

func TestGetBoard(t *testing.T) {
	c, err := NewClient("ABC")
	if err != nil {
		t.Fatal(err)
	}

	tcs := []struct {
		scenario string
		status   int
		response string
		success  bool
		board    Board
		err      error
	}{
		{
			scenario: "successful request",
			status:   200,
			response: `
						{
				"pipelines": [
				{
					"id": "595d430add03f01d32460080",
					"name": "New Issues",
					"issues": [
					{
						"issue_number": 279,
						"estimate": { "value": 40 },
						"position": 0,
						"is_epic": true
					},
					{
						"issue_number": 142,
						"is_epic": false
					}
					]
				},
				{
					"id": "595d430add03f01d32460081",
					"name": "Backlog",
					"issues": [
					{
						"issue_number": 303,
						"estimate": { "value": 40 },
						"position": 3,
						"is_epic": false
					}
					]
				},
				{
					"id": "595d430add03f01d32460082",
					"name": "To Do",
					"issues": [
					{
						"issue_number": 380,
						"estimate": { "value": 1 },
						"position": 0,
						"is_epic": true
					},
					{
						"issue_number": 284,
						"position": 2,
						"is_epic": false
					},
					{
						"issue_number": 329,
						"estimate": { "value": 8 },
						"position": 7,
						"is_epic": false
					}
					]
				}
				]
			}
			`,
			success: true,
			board: Board{
				Pipelines: []Pipeline{
					{
						ID:   "595d430add03f01d32460080",
						Name: "New Issues",
						Issues: []Issue{
							{
								IssueNumber: 279,
								Estimate:    &Estimate{Value: 40},
								Position:    0,
								IsEpic:      true,
							},
							{
								IssueNumber: 142,
								IsEpic:      false,
							},
						},
					},
					{
						ID:   "595d430add03f01d32460081",
						Name: "Backlog",
						Issues: []Issue{
							{
								IssueNumber: 303,
								Estimate:    &Estimate{Value: 40},
								Position:    3,
								IsEpic:      false,
							},
						},
					},
					{
						ID:   "595d430add03f01d32460082",
						Name: "To Do",
						Issues: []Issue{
							{
								IssueNumber: 380,
								Estimate:    &Estimate{Value: 1},
								IsEpic:      true,
							},
							{
								IssueNumber: 284,
								Position:    2,
								IsEpic:      false,
							},
							{
								IssueNumber: 329,
								Estimate:    &Estimate{Value: 8},
								Position:    7,
								IsEpic:      false,
							},
						},
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

			board, err := c.GetBoard(ctx, 123)

			if tc.success {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}

				if !reflect.DeepEqual(tc.board, board) {
					t.Fatalf("expected board to eq %v, got %v", tc.board, board)
				}
			} else {
				if !strings.HasPrefix(err.Error(), tc.err.Error()) {
					t.Fatalf("expected error %v, got %v", tc.err, err)
				}
			}

		})
	}
}

func TestGet(t *testing.T) {
	token := "ABC"

	t.Run("invalid url", func(t *testing.T) {
		ctx := context.Background()

		_, err := get(ctx, ":]/", token)
		if err == nil {
			t.Error("expected error, got none")
		}
	})

	t.Run("invalid server", func(t *testing.T) {
		ctx := context.Background()

		_, err := get(ctx, "/", token)
		if err == nil {
			t.Error("expected error, got none")
		}
	})

	t.Run("status not ok", func(t *testing.T) {
		ctx := context.Background()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
			w.Write(nil)
		}))
		defer srv.Close()

		_, err := get(ctx, srv.URL, token)
		if err == nil {
			t.Error("expected error, got none")
		}
	})
}
