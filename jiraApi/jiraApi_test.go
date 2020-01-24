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
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/sotomskir/jira-cli/jiraApi/models"
	"github.com/spf13/viper"
	"gopkg.in/jarcoal/httpmock.v1"
	"gopkg.in/resty.v1"
	"gotest.tools/assert"
	"io/ioutil"
	"testing"
)

var fakeClock = clockwork.NewFakeClock()

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
		t.Error(err)
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

	version, _ := CreateVersion("TEST", "1.2.0")
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

	status, err := TransitionIssue("", "TEST-1", "code review", "")

	if err != nil {
		t.Error(err)
	}
	if status != 0 {
		t.Errorf("TestTransitionIssue: expected status: 0, got: %d", status)
	}
}

func TestTransitionIssueWorkflowError(t *testing.T) {
	viper.Reset()
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")

	httpmock.RegisterResponder("GET", "http://example.com",
		httpmock.NewErrorResponder(errors.New("not found")))

	_, err := TransitionIssue("http://example.com", "TEST-1", "code review", "")

	if err == nil {
		t.Error("TestTransitionIssueWorkflowError: should return error")
	}
	expectedError := " Get http://example.com: not found"
	if err.Error() != expectedError {
		t.Errorf("TestTransitionIssueWorkflowError: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestTransitionIssueGetIssueError(t *testing.T) {
	viper.Reset()
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	viper.Set("JIRA_WORKFLOW_CONTENT", getWorkflowString())

	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(400, ""))

	_, err := TransitionIssue("", "TEST-1", "code review", "")

	if err == nil {
		t.Error("TestTransitionIssueGetIssueError: should return error")
	}
	expectedError := "http error: 400"
	if err.Error() != expectedError {
		t.Errorf("TestTransitionIssueGetIssueError: expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestWorklog(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/worklog.json")
	httpmock.RegisterResponder("POST", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewStringResponder(200, response))

	AddWorklog("TEST-1", 60, "comment", "", "")

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	count := info["POST https://jira.example.com/rest/api/2/issue/TEST-1/worklog"]
	if count != 1 {
		t.Errorf("TestWorklog: expected api calls: 1, got: %d", count)
	}
}

func TestWorklogError(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.RegisterResponder("POST", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewErrorResponder(errors.New("ERROR")))

	AddWorklog("TEST-1", 60, "comment", "", "")

	httpmock.GetTotalCallCount()
	info := httpmock.GetCallCountInfo()
	count := info["POST https://jira.example.com/rest/api/2/issue/TEST-1/worklog"]
	if count != 1 {
		t.Errorf("TestWorklog: expected api calls: 1, got: %d", count)
	}
}

func TestWorklogPayloadError(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.RegisterResponder("POST", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewErrorResponder(errors.New("ERROR")))

	_, e := AddWorklog("TEST-1", 60, "comment", "wwf", "1")

	assert.Error(t, e, "If provided the date and time must adhere to formats: [YYYY-MM-DD] and [HH:ss]. You provided: date=[ wwf ] and time=[ 1 ]\n")
}

func TestDeleteWorklogForUser(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	//given
	httpmock.Activate()
	Initialize("https://jira.example.com", "jenkins_jira", "pass")
	//-------------------------
	//when - cant list worklogs
	//-------------------------
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewStringResponder(404, ""))

	sumOk, sumError, err := DeleteWorklogForUser("jenkins_jira", "TEST-1")
	//-------------------------
	//then - assert 404 error
	//-------------------------
	expectedError := "http error: 404"
	if err.Error() != expectedError {
		t.Errorf("TestDeleteWorklogForUser: expected error: %s, got: %s", expectedError, err.Error())
	}
	//-------------------------
	//when - cant delete worklog
	//-------------------------
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewStringResponder(200, readResponse("./responses/issue/TEST-1/worklogsList.json")))

	sumOk, sumError, err = DeleteWorklogForUser("jenkins_jira", "TEST-1")
	//-------------------------
	//then assert error deleting
	//-------------------------
	if sumError != 1 && sumOk == 0 {
		t.Errorf("TestDeleteWorklogForUser: expected bad request for DELETE worklog 367175, errors slice should be equal 1 element, got: %x", sumError)
	}
	//-------------------------
	//when - positive scenario we can delete 367175
	//-------------------------
	httpmock.RegisterResponder("DELETE", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog/367175",
		httpmock.NewStringResponder(204, ""))

	sumOk, sumError, err = DeleteWorklogForUser("jenkins_jira", "TEST-1")
	//-------------------------
	//then assert 1 deleted, no errors
	//-------------------------
	sumOk, sumError, err = DeleteWorklogForUser("jenkins_jira", "TEST-1")

	if err != nil && sumError != 0 && sumOk != 1 {
		t.Errorf("TestDeleteWorklogForUser: expected no errors, no bad DELETE requests, one successful request, but got: errors: %v, bad DELETE requests: %d, successful request: %d", err, sumError, sumOk)
	}
}

func TestGetTransitionByName(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/transitions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/transitions",
		httpmock.NewStringResponder(200, response))

	transition, err := GetTransitionByName("TEST-1", "Reviewed")
	if err != nil {
		t.Errorf("TestGetTransitionByName: %s", err)
	}
	if transition.Id != "81" {
		t.Errorf("TestGetTransitionByName: expected id: 81, got: %s", transition.Id)
	}
	if transition.Name != "Reviewed" {
		t.Errorf("TestGetTransitionByName: expected id: Reviewed, got: %s", transition.Name)
	}
}

func TestGetTransitionByNameError(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/transitions.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/transitions",
		httpmock.NewStringResponder(200, response))

	_, err := GetTransitionByName("TEST-1", "Non existent transition")
	if err == nil {
		t.Error("TestGetTransitionByNameError: should return error when transition not exist")
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

	if len(transitions) != 2 {
		t.Errorf("TestGetTransitions: expected length: 1, got: %d", len(transitions))
	}
	if transitions[0].Id != "81" {
		t.Errorf("TestGetTransitions: expected id: 81, got: %s", transitions[0].Id)
	}
	if transitions[0].Name != "Reviewed" {
		t.Errorf("TestGetTransitions: expected id: Reviewed, got: %s", transitions[0].Name)
	}
}

func TestListWorklog(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/issue/TEST-1/worklogsList.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog",
		httpmock.NewStringResponder(200, response))

	worklogs, _ := ListWorklog("TEST-1")

	if worklogs.Total != 2 {
		t.Errorf("TestListWorklog: expected length: 2, got: %d", worklogs.Total)
	}

	expFirstId := "359998"
	if worklogs.Worklogs[0].Id != expFirstId {
		t.Errorf("TestListWorklog: expected first worklog id: 359998, got: %s", worklogs.Worklogs[0].Id)
	}
}

func TestDeleteWorklog(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.RegisterResponder("DELETE", "https://jira.example.com/rest/api/2/issue/TEST-1/worklog/666",
		httpmock.NewStringResponder(204, ""))

	status, _ := DeleteWorklog("TEST-1", "666")

	if status != 204 {
		t.Errorf("TestDeleteWorklog: expected status 204 got: %d", status)
	}
}

func TestGetIssues(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response1 := readResponse("./responses/issue/TEST-1.json")
	response2 := readResponse("./responses/issue/TEST-2.json")
	response404 := readResponse("./responses/issue/404.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response1))
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-2",
		httpmock.NewStringResponder(200, response2))
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-3",
		httpmock.NewStringResponder(404, response404))

	issues := GetIssues([]string{"TEST-1", "TEST-2", "TEST-3"})

	if len(issues) != 2 {
		t.Errorf("TestGetIssues: expected length: 2, got: %d", len(issues))
	}

	if issues[0].Id != "10000" {
		t.Errorf("TestGetIssues: expected first element Id 10000, got: %s", issues[0].Id)
	}

	if issues[1].Id != "10001" {
		t.Errorf("TestGetIssues: expected first element Id 10001, got: %s", issues[1].Id)
	}
}

func TestCreateIssue(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := "{\"id\":\"10109\",\"key\":\"TEST-16\",\"self\":\"http://jira.example.com/rest/api/2/issue/10109\"}"
	httpmock.RegisterResponder(resty.MethodPost, "https://jira.example.com/rest/api/2/issue",
		httpmock.NewStringResponder(201, response))
	issue, err := CreateIssue("TEST", "test", "test", "Task")
	if err != nil {
		t.Errorf("TestCreateIssue: unexpected error %#v\n", err)
	}
	expected := models.Issue{Id: "10109", Key: "TEST-16", Self: "http://jira.example.com/rest/api/2/issue/10109"}
	assert.DeepEqual(t, issue, expected)
}

func TestGetIssuesInVersions(t *testing.T) {

	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	response := readResponse("./responses/version/list.json")
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/search",
		httpmock.NewStringResponder(200, response))

	issuesInVersions, err := GetIssuesInVersions("TEST", "1.0.0", "story")

	if err != nil {
		t.Errorf("TestGetIssuesInVersions: unexpected error %#v\n", err)
	}
	expected := models.IssueList{
		Total: 1,
		Issues: []models.Issue{
			{
				Id:     "1350616",
				Key:    "TEST-3192",
				Fields: models.Fields{Summary: "Test test tesT Test test tesT Test test tesT Test test tesT"},
				Self:   "https://jira:8080/rest/api/2/issue/1350616",
			},
		},
	}
	assert.DeepEqual(t, issuesInVersions, expected)

}

func TestGetIssueWorkflow(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	workflowResponse := readResponse("./responses/workflows/workflow.json")
	httpmock.RegisterResponder("GET", "https://user:pass@jira.example.com/browse/TEST-1",
		httpmock.NewStringResponder(200, "<html><a href=\"workflowName=test-workflow&test=test\" class=\"jira-workflow-designer-link\"></a></html>"))
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/workflowDesigner/latest/workflows?name=test-workflow",
		httpmock.NewStringResponder(200, workflowResponse))

	workflow, err := GetIssueWorkflow("TEST-1")

	if err != nil {
		t.Errorf("TestGetIssueWorkflow: unexpected error: %#v\n", err)
	}
	expected := &models.Workflow{
		Layout: models.WorkflowLayout{
			Statuses:    []models.Status{
				{
					Id:       "S<6>",
					Name:     "Code review",
					StepId:   6,
					StatusId: "10601",
				},
			},
			Transitions: []models.Transition{
				{
					Id:               "A<111:S<1>:S<7>>",
					Name:             "Need Info",
					SourceId:         "S<1>",
					TargetId:         "S<7>",
					ActionId:         111,
					Initial:          false,
					GlobalTransition: false,
					LoopedTransition: false,
				},
			},
		},
	}
	assert.DeepEqual(t, workflow, expected)
}

func TestBuildWorkflow(t *testing.T) {
	defer httpmock.DeactivateAndReset()
	httpmock.Activate()
	Initialize("https://jira.example.com", "user", "pass")
	workflowResponse := readResponse("./responses/workflows/workflow_simple.json")
	httpmock.RegisterResponder("GET", "https://user:pass@jira.example.com/browse/TEST-1",
		httpmock.NewStringResponder(200, "<html><a href=\"workflowName=test-workflow&test=test\" class=\"jira-workflow-designer-link\"></a></html>"))
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/workflowDesigner/latest/workflows?name=test-workflow",
		httpmock.NewStringResponder(200, workflowResponse))

	workflow, err := GetIssueWorkflow("TEST-1")
	if err != nil {
		t.Errorf("TestBuildWorkflow: unexpected error: %#v\n", err)
	}

	transitionsMap := BuildWorkflow(workflow, "In Progress", "Code Review")

	expected := &WorkflowTransitionsMap{workflow: map[string]interface{}{}}

	assert.DeepEqual(t, transitionsMap, expected)
}
