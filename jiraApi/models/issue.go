package models

// Issue type represents JIRA issue resource
type Issue struct {
	Id     string `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	Fields Fields `json:"fields,omitempty"`
	Self   string `json:"self,omitempty"`
}
