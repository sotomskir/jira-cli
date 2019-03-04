package models

type Fields struct {
	FixVersions []Version `json:"fixVersions"`
	Status      Status    `json:"status"`
}

