package zenhub

import (
	"net/url"

	"github.com/pkg/errors"
)

// IssueTransferWebhookEvent is a webhook event when an issue moves from one
// pipeline to another.
type IssueTransferWebhookEvent struct {
	GitHubURL        string
	Organization     string
	Repo             string
	UserName         string
	IssueNumber      string
	IssueTitle       string
	ToPipelineName   string
	FromPipelineName string
}

// EventType returns the type of event.
func (e IssueTransferWebhookEvent) EventType() EventType {
	return EventTypeEstimateIssue
}

// ParseWebhookEvent tries to create an event from ZenHub's webhook post.
func ParseWebhookEvent(b []byte) (Event, error) {
	q, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse data")
	}

	t := q.Get("type")

	switch EventType(t) {
	case EventTypeIssueTransfer:
		event := IssueTransferWebhookEvent{
			GitHubURL:        q.Get("github_url"),
			Organization:     q.Get("organization"),
			Repo:             q.Get("repo"),
			UserName:         q.Get("user_name"),
			IssueNumber:      q.Get("issue_number"),
			IssueTitle:       q.Get("issue_title"),
			ToPipelineName:   q.Get("to_pipeline_name"),
			FromPipelineName: q.Get("from_pipeline_name"),
		}

		return event, nil
	default:
		return nil, ErrUnknownEventType{t: t}
	}
}
