package models

type AddWorklog struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
}
