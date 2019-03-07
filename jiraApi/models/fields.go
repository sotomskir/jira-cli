package models

// Fields type represents fields of JIRA issue
type Fields struct {
	FixVersions []Version `json:"fixVersions"`
	Status      Status    `json:"status"`
}
