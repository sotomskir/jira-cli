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
	"testing"
)

func TestGetIssue(t *testing.T) {
	Initialize("https://jira.example.com", "user", "pass")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	response := `{"id": "1", "key": "TEST-1", "summary": "Some test issue"}`
	httpmock.RegisterResponder("GET", "https://jira.example.com/rest/api/2/issue/TEST-1",
		httpmock.NewStringResponder(200, response))

	// do stuff that makes a request to articles.json}
	issue, _ := GetIssue("TEST-1")
	if issue.Id != "1" || issue.Key != "TEST-1" || issue.Summary != "Some test issue" {
		t.Errorf("Error: expected: %s, got: %#v\n", response, issue)
	}
}
