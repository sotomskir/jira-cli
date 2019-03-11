package models

// WorklogResp type represents JIRA worklog resource response
type WorklogResp struct {
	Id        string `json:"id"`
	TimeSpent int    `json:"timeSpentSeconds"`
	Author    Author `json:"author"`
}
