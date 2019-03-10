// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jiraApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi/models"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strings"
	"time"
)

// Initialize method is used to initialize API client
func Initialize(serverUrl string, username string, password string) {
	resty.SetHostURL(serverUrl)
	resty.SetTimeout(1 * time.Minute)
	resty.SetBasicAuth(username, password)
	// Headers for all request
	resty.SetHeader("Accept", "application/json")
	resty.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "jira-cli",
	})
}

func execute(method string, endpoint string, payload interface{}, response interface{}) (code int, error error) {
	r := resty.R()
	if payload != nil {
		r.SetBody(payload)
	}
	res, err := r.Execute(method, endpoint)
	if err != nil {
		logrus.Errorln(err)
		return 1, err
	}

	if res.StatusCode() >= 400 {
		logrus.Debugf("GET: %s\nStatus code: %d\nResponse: %s\n", endpoint, res.StatusCode(), string(res.Body()))
		return res.StatusCode(), errors.New(fmt.Sprintf("http error: %d", res.StatusCode()))
	}

	if res.StatusCode() == 204 {
		return 204, nil
	}

	jsonErr := json.Unmarshal(res.Body(), response)

	if jsonErr != nil {
		logrus.Errorf("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		return 1, errors.New("unmarshalling error")
	}

	return 0, nil
}

// GetVersions method returns JIRA versions of given project
func GetVersions(projectKey string) []models.Version {
	versions := make([]models.Version, 0)
	execute("GET", fmt.Sprintf("rest/api/2/project/%s/versions", projectKey), nil, &versions)
	return versions
}

func find(vs []models.Version, name string) (models.Version, error) {
	for _, v := range vs {
		if v.Name == name {
			return v, nil
		}
	}
	return models.Version{}, errors.New("version not found")
}

// GetVersion method returns JIRA version details
func GetVersion(projectKey string, version string) (models.Version, error) {
	versions := GetVersions(projectKey)
	return find(versions, version)
}

// CreateVersion creates new version in specified project
func CreateVersion(projectKey string, version string) models.Version {
	payload := models.Version{}
	payload.Name = version
	payload.Project = projectKey
	response := models.Version{}
	execute("POST", "rest/api/2/version", payload, &response)
	return response
}

func updateVersion(versionId string, payload models.Version) models.Version {
	response := models.Version{}
	execute("PUT", fmt.Sprintf("rest/api/2/version/%s", versionId), payload, &response)
	return response
}

// ReleaseVersion method changes project status to "released"
func ReleaseVersion(projectKey string, version string) {
	versionFromServer, _ := GetVersion(projectKey, version)
	payload := models.Version{}
	payload.Released = true
	updateVersion(versionFromServer.Id, payload)
}

// GetProject method returns project details
func GetProject(projectKey string) models.Project {
	project := models.Project{}
	execute("GET", fmt.Sprintf("rest/api/2/project/%s", projectKey), nil, &project)
	return project
}

// GetProjects method list all projects
func GetProjects() []models.Project {
	projects := make([]models.Project, 0)
	execute("GET", "rest/api/2/project", nil, &projects)
	return projects
}

func mapVersionName(versions []models.Version) []string {
	result := make([]string, 0)
	for _, v := range versions {
		result = append(result, v.Name)
	}
	return result
}

// SetFixVersion method sets fix version of issue. When version is already set it won't be modified
func SetFixVersion(issueKey string, version string) error {
	response, err := GetIssue(issueKey)
	if err != nil {
		return err
	}
	if len(response.Fields.FixVersions) > 0 {
		logrus.Warnf("Fix version is already set to: %#v\n", mapVersionName(response.Fields.FixVersions))
		return errors.New("fix version is already set")
	}
	project := strings.Split(issueKey, "-")[0]
	_, err = GetVersion(project, version)
	if err != nil {
		CreateVersion(project, version)
	}
	_, err = execute("PUT", fmt.Sprintf("rest/api/2/issue/%s", issueKey), fmt.Sprintf("{\"update\":{\"fixVersions\":[{\"set\":[{\"name\":\"%s\"}]}]}}", version), &response)
	if err != nil {
		return err
	}
	return nil
}

// GetIssue method returns issue details
func GetIssue(issueKey string) (i models.Issue, error error) {
	issue := models.Issue{}
	_, err := execute("GET", fmt.Sprintf("rest/api/2/issue/%s", issueKey), nil, &issue)
	if err != nil {
		return issue, err
	}
	return issue, nil
}

// Worklog method add worklog to issue
func Worklog(key string, min uint64, com string) error {
	sec := min * 60
	logrus.Infof("Attempting to add %d[sec] for issue %s.", sec, key)
	payload := models.AddWorklog{}
	payload.Comment = com
	payload.TimeSpentSeconds = sec
	wr := models.WorklogResp{}
	_, err := execute("POST", fmt.Sprintf("rest/api/2/issue/%s/worklog", key), payload, &wr)
	if err == nil && len(wr.Id) > 0 {
		logrus.Infof("Successfully added %d[sec] to issue %s.", sec, key)
		return nil
	} else {
		msg := fmt.Sprintf("There was an error adding your time to issue %s.", key)
		logrus.Error(msg)
		return errors.New(msg)
	}
}

// TransitionIssue method executes issue transition
func TransitionIssue(workflowPath string, issueKey string, targetStatus string) (status int, error error) {
	workflow, err := ReadWorkflow(workflowPath)
	if err != nil {
		return 1, err
	}
	for i := 0; i < 20; i++ {
		issue, err := GetIssue(issueKey)
		if err != nil {
			return 1, err
		}
		currentStatus := strings.ToLower(issue.Fields.Status.Name)
		logrus.Infof("%s: current status: '%s', target status: '%s'\n", issueKey, currentStatus, targetStatus)
		if currentStatus == targetStatus {
			break
		}
		transitionName, err := workflow.GetOrDefault(currentStatus, targetStatus)
		transition, err := GetTransitionByName(issueKey, transitionName)
		if err != nil {
			return 1, err
		}
		payload := models.Transitions{}
		payload.Transition = transition
		logrus.Infof("%s: executing transition: '%s'\n", issueKey, transition.Name)
		status, err := execute("POST", fmt.Sprintf("rest/api/2/issue/%s/transitions", issueKey), payload, nil)
		if err != nil {
			logrus.Errorf("%s: executing transition: '%s'\n", issueKey, transition.Name)
			return status, err
		}
	}
	return 0, nil
}

// GetTransitionByName method returns transition details from issue
func GetTransitionByName(issueKey string, transitionName string) (models.Transition, error) {
	transitions := GetTransitions(issueKey)
	for _, transition := range transitions {
		if strings.ToLower(transition.Name) == strings.ToLower(transitionName) {
			return transition, nil
		}
	}
	return models.Transition{}, errors.New(fmt.Sprintf("transition '%s' is not found in transitions of issue: %s\n", transitionName, issueKey))
}

// GetTransitions method returns available transitions for issue
func GetTransitions(issueKey string) []models.Transition {
	transitions := models.Transitions{}
	execute("GET", fmt.Sprintf("rest/api/2/issue/%s/transitions", issueKey), nil, &transitions)
	return transitions.Transitions
}

// TestTransitions method run through all transitions to test workflow definition
func TestTransitions(workflowPath string, issueKey string) error {
	_, err := ReadWorkflow(workflowPath)
	if err != nil {
		return err
	}
	workflow := viper.GetStringMap("workflow")
	for fromState := range workflow {
		logrus.Infof("\tTesting transitions from state: '%s'\n", fromState)
		TransitionIssue(workflowPath, issueKey, fromState)
		for toState := range workflow {
			logrus.Infof("\tto state: '%s'\n", toState)
			TransitionIssue(workflowPath, issueKey, toState)
		}
	}
	return nil
}
