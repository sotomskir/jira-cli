package models

// Status type represents JIRA issue status
type Status struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
