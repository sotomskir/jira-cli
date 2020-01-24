package models

// Transition type represents JIRA workflow transition
type Transition struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	SourceId         string `json:"sourceId"`
	TargetId         string `json:"targetId"`
	ActionId         uint   `json:"actionId"`
	Initial          bool   `json:"initial"`
	GlobalTransition bool   `json:"globalTransition"`
	LoopedTransition bool   `json:"loopedTransition"`
}
