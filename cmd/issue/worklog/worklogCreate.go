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
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
)

//Cmd workload add command
var worklogCreateCmd = &cobra.Command{
	Use:     "add TIME_IN_MINUTES ISSUE_KEY [ISSUE_KEY...]",
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Add worklog for given tasks",
	Run: func(cmd *cobra.Command, args []string) {
		issueKeys := args[1:]
		min, _ := strconv.ParseUint(args[0], 0, 64)
		com, err := cmd.Flags().GetString("comment")
		date, _ := cmd.Flags().GetString("date")
		time, _ := cmd.Flags().GetString("time")

		if err != nil || len(com) == 0 {
			com = `Automatically added by jira-cli. 
			We'd love to accept your patches!
			Come to the Dark Side - we have cookies. 
			Project url: [https://github.com/sotomskir/jira-cli]
			`
		}

		var wg sync.WaitGroup
		for _, issueKey := range issueKeys {
			wg.Add(1)
			go func(issueKey string, min uint64, com string, date string, time string) {
				defer wg.Done()
				jiraApi.AddWorklog(issueKey, min, com, date, time)
			}(issueKey, min, com, date, time)
		}
		wg.Wait()
	},
}

func init() {
	worklogCreateCmd.Flags().StringP("comment", "c", "", "Comment for worklog entry.")
	worklogCreateCmd.Flags().StringP("date", "d", "", "Explicit date for worklog entry.\nMust adhere to format: YYYY-MM-DD (eg. 2019-04-01).\n[ Default: current date ]")
	worklogCreateCmd.Flags().StringP("time", "t", "", "Explicit time for worklog entry.\nMust adhere to format: HH:ss (eg. 12:30).\n[ Default: 08:00 ]")

}
