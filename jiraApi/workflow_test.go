package jiraApi

import (
	"errors"
	"github.com/spf13/viper"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
)

func TestReadWorkflowFromHttp(t *testing.T) {
	viper.Reset()
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://example.com/workflow",
		httpmock.NewStringResponder(200, getWorkflowString()))
	workflow, _ := ReadWorkflow("https://example.com/workflow")
	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	count := info["GET https://example.com/workflow"]
	if count != 1 {
		t.Errorf("TestReadWorkflowFromHttp: expected api calls: 1, got: %d", count)
	}
	checkWorkflow(workflow, t)
}

func TestReadWorkflowFromVar(t *testing.T) {
	viper.Reset()
	viper.Set("JIRA_WORKFLOW_CONTENT", getWorkflowString())
	workflow, _ := ReadWorkflow("")
	checkWorkflow(workflow, t)
}

func TestReadWorkflowErrorsFileNotFound(t *testing.T) {
	viper.Reset()
	_, err := ReadWorkflow("")
	if err == nil {
		t.Errorf("TestReadWorkflowErrorsFileNotFound: ReadWorkflow should return error when file not exist")
	}
	if err.Error() != "WorkflowTransitionsMap file not found: \n" {
		t.Errorf("TestReadWorkflowErrorsFileNotFound: expected error: WorkflowTransitionsMap file not found: \n, got: %s", err.Error())
	}
}

func TestReadWorkflowErrorsHttpError(t *testing.T) {
	viper.Reset()
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "http://example.com",
		httpmock.NewErrorResponder(errors.New("bad request")))
	_, err := ReadWorkflow("http://example.com")
	if err == nil {
		t.Errorf("TestReadWorkflowErrorsHttpError: ReadWorkflow should return http error")
	}
	expectedError := " Get http://example.com: bad request"
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestReadWorkflowErrorsHttpError: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestReadWorkflowFromFile(t *testing.T) {
	viper.Reset()
	workflow, _ := ReadWorkflow("./responses/workflow.yaml")
	checkWorkflow(workflow, t)
}

func TestWorkflow_GetOrDefault(t *testing.T) {
	viper.Reset()
	viper.Set("JIRA_WORKFLOW_CONTENT", getWorkflowString())
	workflow, _ := ReadWorkflow("")

	// Should throw error when currentStatus not exist in workflow
	status, err := workflow.GetOrDefault("non existent status", "")
	if err == nil {
		t.Errorf("GetOrDefault should return error when currentStatus not exist")
	}
	if status != "" {
		t.Errorf("GetOrDefault: status shoult be empty, got: %s", status)
	}

	// Should throw error when target status, or default not exist in workflow
	status, err = workflow.GetOrDefault("to do", "non existent status")
	if err == nil {
		t.Errorf("GetOrDefault should return error when targetStatus not exist")
	}
	if status != "" {
		t.Errorf("GetOrDefault: status shoult be empty, got: %s", status)
	}
}

func checkWorkflow(workflow WorkflowTransitionsMap, t *testing.T) {
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
 in progress:
   default: code review
 done:
   default: reopen
 rejected:
   default: reopen
`
}
