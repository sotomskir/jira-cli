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
		t.Errorf("Error: expected: %s, got: %s\n", expectedId, issue.Id)
	}
	if issue.Key != expectedKey {
		t.Errorf("Error: expected: %s, got: %s\n", expectedKey, issue.Key)
	}
	if issue.Fields.Summary != expectedSummary {
		t.Errorf("Error: expected: %s, got: %s\n", expectedSummary, issue.Fields.Summary)
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
		t.Errorf("Error: expected length: 2, got: %d\n", len(projects))
	}
	if projects[0].Id != "10001" {
		t.Errorf("Error: expected id: 10001, got: %s\n", projects[0].Id)
	}
	if projects[1].Id != "10000" {
		t.Errorf("Error: expected id: 10000, got: %s\n", projects[1].Id)
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
		t.Errorf("Error: expected status: 204, got: %d\n", status)
	}
}
