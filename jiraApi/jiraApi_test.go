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
	"github.com/spf13/viper"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"testing"
)

func readResponse(path string) string {
	json, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(json)
}

func TestGetIssue(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response))

	issue, _ := GetIssue("TEST-1")
	expectedId := "10000"
	expectedKey := "TEST-1"
	expectedSummary := "ax"
	if issue.Id != expectedId {
		t.Errorf("TestGetIssue: expected: %s, got: %s", expectedId, issue.Id)
	}
	if issue.Key != expectedKey {
		t.Errorf("TestGetIssue: expected: %s, got: %s", expectedKey, issue.Key)
	}
	if issue.Fields.Summary != expectedSummary {
		t.Errorf("TestGetIssue: expected: %s, got: %s", expectedSummary, issue.Fields.Summary)
	}
}

func TestGetIssueWithError400(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(400, "Bad request"))

	_, err := GetIssue("TEST-1")
	expectedError := "http error: 400"
	if err == nil {
		t.Errorf("TestGetIssueWithError400: should return error on http 400")
	}
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestGetIssueWithError400: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetIssueWithIncorrectJson(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, "some incorrect json"))

	_, err := GetIssue("TEST-1")
	expectedError := "unmarshalling error"
	if err == nil {
		t.Errorf("TestGetIssueWithError400: should return error when response is incorrect json")
	}
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestGetIssueWithError400: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetProjects(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/project.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/project",
		httpmock.NewStringResponder(200, response))

	projects := GetProjects()
	if len(projects) != 2 {
		t.Errorf("TestGetProjects: expected length: 2, got: %d", len(projects))
	}
	if projects[0].Id != "10001" {
		t.Errorf("TestGetProjects: expected id: 10001, got: %s", projects[0].Id)
	}
	if projects[1].Id != "10000" {
		t.Errorf("TestGetProjects: expected id: 10000, got: %s", projects[1].Id)
	}
}

func TestSetFixVersion(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response))

	httpmock.RegisterResponder("PUT", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(204, response))

	err := SetFixVersion("TEST-1", "1")
	if err != nil {
		t.Error(err)
	}
}

func TestSetFixVersion404(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(404, response))

	httpmock.RegisterResponder("PUT", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(204, response))

	err := SetFixVersion("TEST-1", "1")
	if err == nil {
		t.Error("TestSetFixVersion404: SetFixVersion should return error")
	}
	expectedError := "http error: 404"
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestGetIssueWithError404: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestSetFixVersion400(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response))

	httpmock.RegisterResponder("PUT", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(400, response))

	err := SetFixVersion("TEST-1", "1")
	if err == nil {
		t.Error("TestSetFixVersion400: SetFixVersion should return error")
	}
	expectedError := "http error: 400"
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestGetIssueWithError400: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestSetFixVersionAlreadySet(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-2.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-2",
		httpmock.NewStringResponder(200, response))

	err := SetFixVersion("TEST-2", "1")
	if err == nil {
		t.Error("TestSetFixVersionAlreadySet: SetFixVersion should return error")
	}
	expectedError := "fix version is already set"
	if err != nil && err.Error() != expectedError {
		t.Errorf("TestSetFixVersionAlreadySet: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetProject(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/project/TEST.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/project/TEST",
		httpmock.NewStringResponder(200, response))

	project := GetProject("TEST")
	if project.Id != "10001" {
		t.Errorf("TestGetProject: expected id: 10001, got: %s", project.Id)
	}
	if project.Key != "TEST" {
		t.Errorf("TestGetProject: expected key: TEST, got: %s", project.Key)
	}
	if project.Name != "Test2" {
		t.Errorf("TestGetProject: expected name: Test2, got: %s", project.Name)
	}
}

func TestGetVersions(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/project/TEST/versions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/project/TEST/versions",
		httpmock.NewStringResponder(200, response))

	versions := GetVersions("TEST")
	if len(versions) != 4 {
		t.Errorf("Error: expected len: 4, got: %d\n", len(versions))
	}
	if versions[0].Id != "10001" {
		t.Errorf("TestGetVersions: expected id: 10001, got: %s", versions[0].Id)
	}
	if versions[0].Name != "2.0.0" {
		t.Errorf("TestGetVersions: expected version: 2.0.0, got: %s", versions[0].Name)
	}
	if versions[0].Released != true {
		t.Errorf("TestGetVersions: expected released: true, got: %t", versions[0].Released)
	}
}

func TestGetVersion(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/project/TEST/versions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/project/TEST/versions",
		httpmock.NewStringResponder(200, response))

	version, err := GetVersion("TEST", "1.2.0")
	if err != nil {
		panic(err)
	}
	if version.Id != "10003" {
		t.Errorf("TestGetVersion: expected id: 10003, got: %s", version.Id)
	}
	if version.Name != "1.2.0" {
		t.Errorf("TestGetVersion: expected version: 1.2.0, got: %s", version.Name)
	}
	if version.Released != false {
		t.Errorf("TestGetVersion: expected released: false, got: %t", version.Released)
	}
}

func TestCreateVersion(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/version/10001.json")
	httpmock.RegisterResponder("POST", "https://jira.example.com/rest/api/2/version",
		httpmock.NewStringResponder(200, response))

	version := CreateVersion("TEST", "1.2.0")
	if version.Id != "10001" {
		t.Errorf("TestCreateVersion: expected id: 10001, got: %s", version.Id)
	}
	if version.Name != "2.0.0" {
		t.Errorf("TestCreateVersion: expected version: 2.0.0, got: %s", version.Name)
	}
	if version.Released != true {
		t.Errorf("TestCreateVersion: expected released: true, got: %t", version.Released)
	}
}

func TestReleaseVersion(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/project/TEST/versions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/project/TEST/versions",
		httpmock.NewStringResponder(200, response))
	response = ""
	httpmock.RegisterResponder("PUT", "https://jira.example.com/rest/api/2/version/10003",
		httpmock.NewStringResponder(204, response))

	ReleaseVersion("TEST", "1.2.0")

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	count1 := info["GET https://jira.example.com/rest/api/2/project/TEST/versions"]
	count2 := info["PUT https://jira.example.com/rest/api/2/version/10003"]
	if count1 != 1 {
		t.Errorf("TestReleaseVersion: expected api calls: 1, got: %d", count1)
	}
	if count2 != 1 {
		t.Errorf("TestReleaseVersion: expected api calls: 1, got: %d", count2)
	}
}

func TestTransitionIssue(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	viper.Set("JIRA_WORKFLOW_CONTENT", getWorkflowString())

	response := readResponse("./responses/issue/TEST-1.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response))

	response = readResponse("./responses/issue/TEST-1/transitions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/transitions",
		httpmock.NewStringResponder(200, response))

	status, err := TransitionIssue("", "TEST-1", "code review")

	if err != nil {
		panic(err)
	}
	if status != 0 {
		t.Errorf("TestTransitionIssue: expected status: 0, got: %d", status)
	}
}

func TestWorklog(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/worklog.json")
	httpmock.RegisterResponder("POST", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewStringResponder(200, response))

	Worklog("TEST-1", 60, "comment")

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	count := info["POST https://jira.example.com/rest/api/2/issue/TEST-1/worklog"]
	if count != 1 {
		t.Errorf("TestWorklog: expected api calls: 1, got: %d", count)
	}
}

func TestGetTransitionByName(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/transitions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/transitions",
		httpmock.NewStringResponder(200, response))

	transition := GetTransitionByName("TEST-1", "Reviewed")

	if transition.Id != "81" {
		t.Errorf("TestGetTransitions: expected id: 81, got: %s", transition.Id)
	}
	if transition.Name != "Reviewed" {
		t.Errorf("TestGetTransitions: expected id: Reviewed, got: %s", transition.Name)
	}
}

func TestGetTransitions(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/transitions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/transitions",
		httpmock.NewStringResponder(200, response))

	transitions := GetTransitions("TEST-1")

	if len(transitions) != 1 {
		t.Errorf("TestGetTransitions: expected length: 1, got: %d", len(transitions))
	}
	if transitions[0].Id != "81" {
		t.Errorf("TestGetTransitions: expected id: 81, got: %s", transitions[0].Id)
	}
	if transitions[0].Name != "Reviewed" {
		t.Errorf("TestGetTransitions: expected id: Reviewed, got: %s", transitions[0].Name)
	}
}
