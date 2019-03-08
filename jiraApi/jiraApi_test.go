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
		t.Errorf("TestGetIssue: expected: %s, got: %s\n", expectedId, issue.Id)
	}
	if issue.Key != expectedKey {
		t.Errorf("TestGetIssue: expected: %s, got: %s\n", expectedKey, issue.Key)
	}
	if issue.Fields.Summary != expectedSummary {
		t.Errorf("TestGetIssue: expected: %s, got: %s\n", expectedSummary, issue.Fields.Summary)
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
		t.Errorf("TestGetProjects: expected length: 2, got: %d\n", len(projects))
	}
	if projects[0].Id != "10001" {
		t.Errorf("TestGetProjects: expected id: 10001, got: %s\n", projects[0].Id)
	}
	if projects[1].Id != "10000" {
		t.Errorf("TestGetProjects: expected id: 10000, got: %s\n", projects[1].Id)
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

	status, err := SetFixVersion("TEST-1", "1")
	if err != nil {
		panic(err)
	}
	if status != 204 {
		t.Errorf("TestSetFixVersion: expected status: 204, got: %d\n", status)
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
		t.Errorf("TestGetProject: expected id: 10001, got: %s\n", project.Id)
	}
	if project.Key != "TEST" {
		t.Errorf("TestGetProject: expected key: TEST, got: %s\n", project.Key)
	}
	if project.Name != "Test2" {
		t.Errorf("TestGetProject: expected name: Test2, got: %s\n", project.Name)
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
		t.Errorf("TestGetVersions: expected id: 10001, got: %s\n", versions[0].Id)
	}
	if versions[0].Name != "2.0.0" {
		t.Errorf("TestGetVersions: expected version: 2.0.0, got: %s\n", versions[0].Name)
	}
	if versions[0].Released != true {
		t.Errorf("TestGetVersions: expected released: true, got: %t\n", versions[0].Released)
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
		t.Errorf("TestGetVersion: expected id: 10003, got: %s\n", version.Id)
	}
	if version.Name != "1.2.0" {
		t.Errorf("TestGetVersion: expected version: 1.2.0, got: %s\n", version.Name)
	}
	if version.Released != false {
		t.Errorf("TestGetVersion: expected released: false, got: %t\n", version.Released)
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
		t.Errorf("TestCreateVersion: expected id: 10001, got: %s\n", version.Id)
	}
	if version.Name != "2.0.0" {
		t.Errorf("TestCreateVersion: expected version: 2.0.0, got: %s\n", version.Name)
	}
	if version.Released != true {
		t.Errorf("TestCreateVersion: expected released: true, got: %t\n", version.Released)
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
		t.Errorf("TestReleaseVersion: expected api calls: 1, got: %d\n", count1)
	}
	if count2 != 1 {
		t.Errorf("TestReleaseVersion: expected api calls: 1, got: %d\n", count2)
	}
}
