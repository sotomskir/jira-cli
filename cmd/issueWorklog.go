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

package cmd

import (
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"strconv"
)

// worklogCmd represents the worklog command
var issueWorklogCmd = &cobra.Command{
	Use:     "worklog ISSUE_KEY TIME_IN_MINUTES",
	Aliases: []string{"w"},
	Args:    cobra.ExactArgs(2),
	Short:   "Manage worklogs for given task",
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		min, _ := strconv.ParseUint(args[1], 0, 64)
		com, err := cmd.Flags().GetString("comment")
		if err != nil || len(com) == 0 {
			com = `Automatically added by jira-cli. 
			We'd love to accept your patches!
			Come to the Dark Side - we have cookies. 
			Project url: [https://github.com/sotomskir/jira-cli]
			`
		}
		jiraApi.Worklog(key, min, com)
	},
}

func init() {
	issueCmd.AddCommand(issueWorklogCmd)
	issueWorklogCmd.Flags().StringP("comment", "c", "", "Comment for worklog entry.")
}
