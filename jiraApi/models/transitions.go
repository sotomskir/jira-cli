package models

// Transitions type represents JIRA workflow transitions
type Transitions struct {
	Transitions []Transition `json:"transitions,omitempty"`
	Transition  Transition   `json:"transition,omitempty"`
}
