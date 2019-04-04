package models

// Fields type represents fields of JIRA issue
type Fields struct {
	FixVersions []Version  `json:"fixVersions,omitempty"`
	Status      *Status    `json:"status,omitempty"`
	Summary     string     `json:"summary,omitempty"`
	Project     *Project   `json:"project,omitempty"`
	IssueType   *IssueType `json:"issuetype,omitempty"`
	Description string     `json:"description,omitempty"`
}
