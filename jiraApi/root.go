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
	"fmt"
	"github.com/sotomskir/jira-cli/logger"
	"gopkg.in/resty.v1"
	"log"
	"os"
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
	Id  string `json:"id"`
	Key string `json:"key"`
	Name string `json:"name"`
}

type Fields struct {
	FixVersions []Version `json:"fixVersions"`
}

type Issue struct {
	Id  string `json:"id"`
	Key string `json:"key"`
	Summary string `json:"summary"`
	Fields Fields `json:"fields"`
}

func Initialize(serverUrl string, authHeader string) {
	resty.SetHostURL(serverUrl)
	resty.SetTimeout(1 * time.Minute)

	// Headers for all request
	resty.SetHeader("Accept", "application/json")
	resty.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"User-Agent":    "jira-cli",
		"Authorization": authHeader,
	})
}

func get(endpoint string, response interface{}) {
	res, err := resty.R().Get(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}

	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nResponse: %s\n", res.StatusCode(), string(res.Body()))
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(res.Body(), response)

	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func post(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Post(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nRequest: %#v\nResponse: %s\n", res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func put(endpoint string, payload interface{}, response interface{}) {
	res, err := resty.R().SetBody(payload).Put(endpoint)
	if err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
	if res.StatusCode() >= 400 {
		logger.ErrorF("Status code: %d\nRequest: %#v\nResponse: %s\n", res.StatusCode(), payload, string(res.Body()))
		os.Exit(1)
	}
	if res.StatusCode() == 204 {
		return
	}
	jsonErr := json.Unmarshal(res.Body(), response)
	if jsonErr != nil {
		logger.ErrorF("StatusCode: %d\nServer responded with invalid JSON: %s\nResponse: %s\n", res.StatusCode(), jsonErr, string(res.Body()))
		os.Exit(1)
	}
}

func GetVersions(projectKey string) []Version {
	// TODO paginacja
	versions := make([]Version, 0)
	get(fmt.Sprintf("rest/api/2/project/%s/versions", projectKey), &versions)
	return versions
}

func find(vs []Version, name string) Version {
	for _, v := range vs {
		if v.Name == name {
			return v
		}
	}
	log.Fatalf("Version %s not found\n", name)
	return Version{}
}

func GetVersion(projectKey string, version string) Version {
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
	versionFromServer := GetVersion(projectKey, version)
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
		logger.WarnF("Fix version is already set to: %#v\n", mapVersionName(response.Fields.FixVersions))
		return false
	}
	put(fmt.Sprintf("rest/api/2/issue/%s", issueKey), fmt.Sprintf("{\"update\":{\"fixVersions\":[{\"set\":[{\"name\":\"%s\"}]}]}}", version), &response)
	return true
}

func GetIssue(issueKey string) Issue {
	issue := Issue{}
	get(fmt.Sprintf("rest/api/2/issue/%s", issueKey), &issue)
	return issue
}

func IssueTransition(issueKey string, targetState string)  {
	logger.ErrorLn("Issue transition not implemented")
	os.Exit(1)
}