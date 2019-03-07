package models

// AddWorklog represents request to JIRA API to add worklog
type AddWorklog struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
}
