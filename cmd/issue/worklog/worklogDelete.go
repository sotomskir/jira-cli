// Copyright © 2019 Maciej 'Cichy' Świderski  <mmswiderski@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package worklog

import (
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sync"
)

//Cmd workload add command
var worklogDeleteCmd = &cobra.Command{
	Use:     "remove ISSUE_KEY",
	Aliases: []string{"r"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Delete all worklogs for logged user from provided ISSUE_KEY",
	Run: func(cmd *cobra.Command, args []string) {
		user := viper.GetString("JIRA_USER")
		issueKeys := args[1:]
		var wg sync.WaitGroup
		for _, issueKey := range issueKeys {
			wg.Add(1)
			go func(user string, issueKey string) {
				defer wg.Done()
				sumOk, sumError, err := jiraApi.DeleteWorklogForUser(user, issueKey)
				if err != nil {
					logrus.Errorln(err)
				}

				if sumError == 0 && sumOk == 0 {
					logrus.Infof("There was no worklogs for user %s in issue %s.", user, issueKey)
				} else {
					logrus.Infof("%s Success: %d | Failed: %d", issueKey, sumOk, sumError)
				}
			} (user, issueKey)
		}
		wg.Wait()
	},
}

func init() {
}
