package models

type Workflow struct {
	Layout WorkflowLayout `json:"layout"`
}

type WorkflowLayout struct {
	Statuses    []Status     `json:"statuses"`
	Transitions []Transition `json:"transitions"`
}
