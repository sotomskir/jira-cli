package models

// Transition type represents JIRA workflow transition
type Transition struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
