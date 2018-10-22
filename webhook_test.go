package zenhub

import (
	"testing"
)

func TestParseWebhookEvent(t *testing.T) {

	tcs := []struct {
		scenario     string
		payload      string
		success      bool
		issueID      int
		organization string
		repository   string
		event        Event
		err          error
	}{
		{
			scenario:     "issue transfer",
			payload:      `type=issue_transfer&github_url=https%3A%2F%2Fgithub.com%2Fmanifoldco%2Fengineering%2Fissues%2F5675&organization=manifoldco&repo=engineering&user_name=luizbranco&issue_number=5675&issue_title=Test%20event%20listener&to_pipeline_name=In%20Progress&from_pipeline_name=Backlog`,
			success:      true,
			issueID:      5675,
			organization: "manifoldco",
			repository:   "engineering",
			event: IssueTransferWebhookEvent{
				GitHubURL:        "https://github.com/manifoldco/engineering/issues/5675",
				Organization:     "manifoldco",
				Repo:             "engineering",
				UserName:         "luizbranco",
				IssueTitle:       "Test event listener",
				IssueNumber:      "5675",
				FromPipelineName: "Backlog",
				ToPipelineName:   "In Progress",
			},
		},
		{
			scenario: "another event type",
			payload:  `type=estimate_set`,
			err:      ErrUnknownEventType{t: "estimate_set"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			event, err := ParseWebhookEvent([]byte(tc.payload))

			if tc.success {
				if err != nil {
					t.Error(err)
				}

				if event != tc.event {
					t.Errorf("expected event to eq %v, got %v", tc.event, event)
				}

			} else {
				if err.Error() != tc.err.Error() {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}

			}
		})
	}
}
