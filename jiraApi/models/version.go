package models

type Version struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Released  bool   `json:"released,omitempty"`
	ProjectId int    `json:"projectId,omitempty"`
	Project   string `json:"project,omitempty"`
}
