package models

// WorklogAdd represents request to JIRA API to add worklog
type WorklogAdd struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
}
