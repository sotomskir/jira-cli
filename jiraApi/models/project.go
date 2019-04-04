package models

// Project type represents JIRA project resource
type Project struct {
	Id   string `json:"id,omitempty"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
