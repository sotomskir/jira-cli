package jiraApi

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"os"
	"strings"
)

type WorkflowTransitionsMap struct {
	workflow map[string]interface{}
}

func (workflow WorkflowTransitionsMap) GetOrDefault(currentStatus string, targetStatus string) (string, error) {
	currentStatusTransitions := workflow.workflow[currentStatus]
	if currentStatusTransitions == nil {
		return "", errors.New(fmt.Sprintf("workflow does not define transitions for status: %s\n", currentStatus))
	}
	if val, ok := cast.ToStringMap(cast.ToStringMap(currentStatusTransitions))[targetStatus]; ok {
		return cast.ToString(val), nil
	}
	if val, ok := cast.ToStringMap(cast.ToStringMap(currentStatusTransitions))["default"]; ok {
		return cast.ToString(val), nil
	}
	return "", errors.New(fmt.Sprintf("transition '%s' is not defined in workflow\n", targetStatus))
}

// ReadWorkflow method loads workflow definition from env var, http url or file
func ReadWorkflow(workflowPath string) (WorkflowTransitionsMap, error) {
	workflowContent := viper.GetString("JIRA_WORKFLOW_CONTENT")
	if workflowContent != "" {
		err := viper.MergeConfig(bytes.NewBuffer([]byte(workflowContent)))
		if err != nil {
			return WorkflowTransitionsMap{}, err
		}
		return WorkflowTransitionsMap{viper.GetStringMap("workflow")}, nil
	}
	if strings.HasPrefix(workflowPath, "http://") || strings.HasPrefix(workflowPath, "https://") {
		response, err := resty.New().R().Get(workflowPath)
		logrus.Debugf("%#v\n", response)
		if err != nil {
			return WorkflowTransitionsMap{}, errors.New(fmt.Sprintf("%s %s", response.Body(), err))
		}
		viper.MergeConfig(bytes.NewBuffer(response.Body()))
		return WorkflowTransitionsMap{viper.GetStringMap("workflow")}, nil
	}
	if _, err := os.Stat(workflowPath); err != nil {
		if os.IsNotExist(err) {
			return WorkflowTransitionsMap{}, errors.New(fmt.Sprintf("WorkflowTransitionsMap file not found: %s\n", workflowPath))
		}
	}
	viper.SetConfigFile(workflowPath)
	viper.MergeInConfig()
	return WorkflowTransitionsMap{viper.GetStringMap("workflow")}, nil
}
