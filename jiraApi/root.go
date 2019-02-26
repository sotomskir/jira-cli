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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"os"
	"strings"
	"time"
)

type Version struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Released  bool   `json:"released,omitempty"`
	ProjectId int    `json:"projectId,omitempty"`
	Project   string `json:"project,omitempty"`
}

type Project struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Fields struct {
	FixVersions []Version `json:"fixVersions"`
	Status      Status    `json:"status"`
}

type Issue struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Summary string `json:"summary"`
	Fields  Fields `json:"fields"`
}

type Status struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Transition struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Transitions struct {
	Transitions []Transition `json:"transitions,omitempty"`
	Transition  Transition   `json:"transition,omitempty"`
}

type AddWorklog struct {
	Comment          string `json:"comment"`
	TimeSpentSeconds uint64 `json:"timeSpentSeconds"`
}

type WorklogResp struct {
	Id string `json:"id"`
}

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

func get(endpoint string, response interface{}) {
	res, err := resty.R().Get(endpoint)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	if res.StatusCode() >= 400 {
		logrus.Errorf("GET: %s\nStatus code: %d\nResponse: %s\n", endpoint, res.StatusCode(), string(res.Body()))
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(res.Body(), response)

	if jsonErr != nil {
		logrus.Errorf("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func post(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Post(endpoint)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logrus.Errorf("POST: %s\nStatus code: %d\nRequest: %#v\nResponse: %s\n", endpoint, res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}
	if res.StatusCode() != 204 {
		jsonErr := json.Unmarshal(res.Body(), response)
		if jsonErr != nil {
			logrus.Errorf("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
			os.Exit(1)
		}
	}
}

func put(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Put(endpoint)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logrus.Errorf("PUT: %s\nStatus code: %d\nRequest: %#v\nResponse: %s\n", endpoint, res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}
	if res.StatusCode() == 204 {
		return
	}
	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logrus.Errorf("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func GetVersions(projectKey string) []Version {
	// TODO paginacja
	versions := make([]Version, 0)
	get(fmt.Sprintf("rest/api/2/project/%s/versions", projectKey), &versions)
	return versions
}

func find(vs []Version, name string) (Version, error) {
	for _, v := range vs {
		if v.Name == name {
			return v, nil
		}
	}
	return Version{}, errors.New("version not found")
}

func GetVersion(projectKey string, version string) (Version, error) {
	versions := GetVersions(projectKey)
	return find(versions, version)
}

func CreateVersion(projectKey string, version string) Version {
	payload := Version{}
	payload.Name = version
	payload.Project = projectKey
	response := Version{}
	post("rest/api/2/version", payload, &response)
	return response
}

func updateVersion(versionId string, payload Version) Version {
	response := Version{}
	put(fmt.Sprintf("rest/api/2/version/%s", versionId), payload, &response)
	return response
}

func ReleaseVersion(projectKey string, version string) {
	versionFromServer, _ := GetVersion(projectKey, version)
	payload := Version{}
	payload.Released = true
	updateVersion(versionFromServer.Id, payload)
}

func GetProject(projectKey string) Project {
	project := Project{}
	get(fmt.Sprintf("rest/api/2/project/%s", projectKey), &project)
	return project
}

func GetProjects() []Project {
	projects := make([]Project, 0)
	get("rest/api/2/project", &projects)
	return projects
}

func mapVersionName(versions []Version) []string {
	result := make([]string, 0)
	for _, v := range versions {
		result = append(result, v.Name)
	}
	return result
}

func SetFixVersion(issueKey string, version string) bool {
	response := GetIssue(issueKey)
	if len(response.Fields.FixVersions) > 0 {
		logrus.Warnf("Fix version is already set to: %#v\n", mapVersionName(response.Fields.FixVersions))
		return false
	}
	project := strings.Split(issueKey, "-")[0]
	_, err := GetVersion(project, version)
	if err != nil {
		CreateVersion(project, version)
	}
	put(fmt.Sprintf("rest/api/2/issue/%s", issueKey), fmt.Sprintf("{\"update\":{\"fixVersions\":[{\"set\":[{\"name\":\"%s\"}]}]}}", version), &response)
	return true
}

func GetIssue(issueKey string) Issue {
	issue := Issue{}
	get(fmt.Sprintf("rest/api/2/issue/%s", issueKey), &issue)
	return issue
}

func Worklog(key string, min uint64, com string) {
	sec := min * 60
	logrus.Infof("Attempting to add %d[sec] for issue %s.", sec, key)
	payload := AddWorklog{}
	payload.Comment = com
	payload.TimeSpentSeconds = sec
	wr := WorklogResp{}
	post(fmt.Sprintf("rest/api/2/issue/%s/worklog", key), payload, &wr)

	if len(wr.Id) > 0 {
		logrus.Infof("Successfully added %d[sec] to issue %s.", sec, key)
	} else {
		logrus.Errorf("There was an error adding your time to issue %s.", sec, key)
	}

}

func TransitionIssue(workflowPath string, issueKey string, targetStatus string) {
	readWorkflow(workflowPath)
	lowerTargetStatus := strings.ToLower(targetStatus)
	workflow := viper.GetStringMap("workflow")
	if workflow == nil {
		logrus.Errorln("workflow not present in config file")
		os.Exit(1)
	}
	for i := 0; i < 20; i++ {
		currentStatus := strings.ToLower(GetIssue(issueKey).Fields.Status.Name)
		logrus.Infof("current status: '%s', target status: '%s'\n", currentStatus, lowerTargetStatus)
		currentStatusTransitions := workflow[currentStatus]
		if currentStatusTransitions == nil {
			logrus.Errorf("workflow does not define transitions for status: %s\n", currentStatus)
			os.Exit(1)
		}
		if currentStatus == targetStatus {
			break
		}
		transition := GetTransitionByName(issueKey, getByNameOrDefault(cast.ToStringMap(currentStatusTransitions), targetStatus))
		payload := Transitions{}
		payload.Transition = transition
		logrus.Infof("executing transition: '%s'\n", transition.Name)
		post(fmt.Sprintf("rest/api/2/issue/%s/transitions", issueKey), payload, nil)
	}
}

func getByNameOrDefault(transitions map[string]interface{}, name string) string {
	if val, ok := cast.ToStringMap(transitions)[name]; ok {
		return cast.ToString(val)
	}
	if val, ok := cast.ToStringMap(transitions)["default"]; ok {
		return cast.ToString(val)
	}
	logrus.Errorf("transition '%s' is not defined in workflow\n", name)
	os.Exit(1)
	return ""
}

func GetTransitionByName(issueKey string, transitionName string) Transition {
	transitions := GetTransitions(issueKey)
	for _, t := range transitions {
		if strings.ToLower(t.Name) == transitionName {
			return t
		}
	}
	logrus.Errorf("transition '%s' is not found in transitions of issue: %s\n", transitionName, issueKey)
	os.Exit(1)
	return Transition{}
}

func GetTransitions(issueKey string) []Transition {
	transitions := Transitions{}
	get(fmt.Sprintf("rest/api/2/issue/%s/transitions", issueKey), &transitions)
	return transitions.Transitions
}

func TestTransitions(workflowPath string, issueKey string) {
	readWorkflow(workflowPath)
	workflow := viper.GetStringMap("workflow")
	if workflow == nil {
		logrus.Errorln("workflow not present in config file")
		os.Exit(1)
	}
	for fromState := range workflow {
		logrus.Infof("\tTesting transitions from state: '%s'\n", fromState)
		TransitionIssue(workflowPath, issueKey, fromState)
		for toState := range workflow {
			logrus.Infof("\tto state: '%s'\n", toState)
			TransitionIssue(workflowPath, issueKey, toState)
		}
	}
}

func readWorkflow(workflowPath string) {
	workflowContent := viper.GetString("JIRA_WORKFLOW_CONTENT")
	if workflowContent != "" {
		viper.MergeConfig(bytes.NewBuffer([]byte(workflowContent)))
		return
	}
	if strings.HasPrefix(workflowPath, "http://") || strings.HasPrefix(workflowPath, "https://") {
		response, err := resty.New().R().Get(workflowPath)
		logrus.Debugln(response)
		if err != nil {
			logrus.Fatalln(response.Body(), err)
		}
		viper.MergeConfig(bytes.NewBuffer(response.Body()))
		return
	}
	if _, err := os.Stat(workflowPath); err != nil {
		if os.IsNotExist(err) {
			logrus.Errorf("Workflow file not found: %s\n", workflowPath)
			os.Exit(1)
		}
	}
	viper.SetConfigFile(workflowPath)
	viper.MergeInConfig()
}
