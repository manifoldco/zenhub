package zenhub

// Board represents a ZenHub board for a repository.
type Board struct {
	Pipelines []Pipeline `json:"pipelines"`
}

// Pipeline represents a ZenHub pipeline. On events, a pipeline doesn't return
// its id or issues list.
type Pipeline struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Issues []Issue `json:"issues"`
}

// Pipeline represents a ZenHub issue.
type Issue struct {
	IssueNumber int       `json:"issue_number"`
	Estimate    *Estimate `json:"estimate"`
	Position    int       `json:"position"`
	IsEpic      bool      `json:"is_epic"`
}
