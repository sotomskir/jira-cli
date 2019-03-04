package models

type Transitions struct {
	Transitions []Transition `json:"transitions,omitempty"`
	Transition  Transition   `json:"transition,omitempty"`
}
