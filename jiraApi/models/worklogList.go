package models

// WorklogList represents response from JIRA API of worklogs for issue
type WorklogList struct {
	Total    int           `json:"total"`
	Worklogs []WorklogResp `json:"worklogs"`
}
