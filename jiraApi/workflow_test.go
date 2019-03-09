package jiraApi

import (
	"github.com/spf13/viper"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

func TestReadWorkflowFromVar(t *testing.T) {
	viper.Set("JIRA_WORKFLOW_CONTENT", getWorkflowString())
	workflow := ReadWorkflow("")
	checkWorkflow(workflow, t)
}

func TestReadWorkflowFromHttp(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://example.com/workflow",
		httpmock.NewStringResponder(200, getWorkflowString()))
	workflow := ReadWorkflow("https://example.com/workflow")
	checkWorkflow(workflow, t)
}

func TestReadWorkflowFromFile(t *testing.T) {
	workflow := ReadWorkflow("./responses/workflow.yaml")
	checkWorkflow(workflow, t)
}

func checkWorkflow(workflow Workflow, t *testing.T) {
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

func getWorkflowString() string {
	viper.SetConfigType("yaml")
	return `
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
`
}
