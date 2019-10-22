package models

// WorklogList represents response from JIRA API of worklogs for issue
type IssueList struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}
