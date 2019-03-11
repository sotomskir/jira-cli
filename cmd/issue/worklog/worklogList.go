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
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

//Cmd workload add command
var worklogListCmd = &cobra.Command{
	Use:     "list ISSUE_KEY",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Short:   "List worklog for given task",
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		resp, err := jiraApi.ListWorklog(key)
		if err != nil {
			logrus.Errorf("There was an error while listing worklogs for issue %s.", key)
			os.Exit(1)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "AUTHOR", "TIME [m]"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		sum := 0
		for _, p := range resp.Worklogs {
			sum += p.TimeSpent
			table.Append([]string{p.Id, p.Author.Name, strconv.Itoa(p.TimeSpent / 60)})
		}
		table.SetFooter([]string{"", "Total", strconv.Itoa(sum / 60)})
		table.Render()
	},
}

func init() {
}
