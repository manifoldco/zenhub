package zenhub

import (
	"fmt"
	"time"
)

// EventType represents what type an event is.
type EventType string

const (
	// EventTypeEstimateIssue is the type for the estimate issue event
	EventTypeEstimateIssue EventType = "estimateIssue"

	// EventTypeTransferIssue is the type for the transfer issue event
	EventTypeTransferIssue EventType = "transferIssue"

	// EventTypeIssueEstiTransfer is the type for issue transfer webhook event
	EventTypeIssueTransfer EventType = "issue_transfer"
)

// ErrUnknownEventType is an error returned when an event cannot be parsed.
type ErrUnknownEventType struct {
	t string
}

// Error returns a custom error message.
func (e ErrUnknownEventType) Error() string {
	return fmt.Sprintf("unknown event type %q", e.t)
}

// Event is an interface implemented by all event types. EventType() can be used
// to find its type and cast the event to that type.
type Event interface {
	EventType() EventType
}

// EstimateIssueEvent is an event when an issue has its estimate value set or
// unset.
type EstimateIssueEvent struct {
	UserID       int
	CreatedAt    time.Time
	FromEstimate *Estimate
	ToEstimate   *Estimate
}

// EventType returns the type of event.
func (e EstimateIssueEvent) EventType() EventType {
	return EventTypeEstimateIssue
}

// TransferIssueEvent is an event when an issue moves from one pipeline to
// another.
type TransferIssueEvent struct {
	UserID       int
	CreatedAt    time.Time
	FromPipeline *Pipeline
	ToPipeline   *Pipeline
}

// EventType returns the type of event.
func (e TransferIssueEvent) EventType() EventType {
	return EventTypeTransferIssue
}

// event combines all possible event fields for json deserialization.
type event struct {
	UserID       int       `json:"user_id"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	FromEstimate *Estimate `json:"from_estimate"`
	ToEstimate   *Estimate `json:"to_estimate"`
	FromPipeline *Pipeline `json:"from_pipeline"`
	ToPipeline   *Pipeline `json:"to_pipeline"`
}

// Estimate represents the estimate value for an issue.
type Estimate struct {
	Value int `json:"value"`
}
