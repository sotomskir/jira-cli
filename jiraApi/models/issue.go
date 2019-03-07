package models

// Issue type represents JIRA issue resource
type Issue struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Summary string `json:"summary"`
	Fields  Fields `json:"fields"`
}
