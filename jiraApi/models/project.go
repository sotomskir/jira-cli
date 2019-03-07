package models

// Project type represents JIRA project resource
type Project struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}
