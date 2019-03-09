package jiraApi

import (
	"github.com/spf13/viper"
	"testing"
)

func TestReadWorkflowFromVar(t *testing.T) {
	setWorkflow()
	workflow := ReadWorkflow("")

	tables := []struct {
		currentStatus      string
		targetStatus       string
		expectedTransition string
	}{
		{"code review", "code review", "ready to test"},
		{"code review", "in test", "ready to test"},
		{"code review", "to do", "ready to test"},
		{"code review", "in progress", "ready to test"},
		{"code review", "done", "ready to test"},
		{"code review", "rejected", "ready to test"},
		{"in test", "done", "done"},
		{"in test", "rejected", "bug found"},
	}

	for _, table := range tables {
		transition, err := workflow.GetOrDefault(table.currentStatus, table.targetStatus)
		if err != nil {
			t.Error(err)
		}
		if transition != table.expectedTransition {
			t.Errorf("Transition is invalid, got: '%s', want: '%s'", transition, table.expectedTransition)
		}
	}
}

func setWorkflow() {
	viper.Set(
		"JIRA_WORKFLOW_CONTENT",
		`
workflow:
 code review:
   default: ready to test
 in test:
   done: done
   default: bug found
 to do:
   rejected: reject
   default: start progress
 in progress:
   default: code review
 done:
   default: reopen
 rejected:
   default: reopen
`)
	viper.SetConfigType("yaml")
}
