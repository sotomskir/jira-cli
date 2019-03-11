// Copyright © 2019 Robert Sotomski <sotomski@gmail.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package issue

import (
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"os"
)

// Cmd represents the issue command
var inspectCmd = &cobra.Command{
	Use:     "inspect ISSUE_KEY [ISSUE_KEY...]",
	Aliases: []string{"i"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Fetch data for given issue.",
	Run: func(cmd *cobra.Command, args []string) {
		keys := args[0:]
		issues := jiraApi.GetIssues(keys)

		if len(issues) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "KEY", "STATUS", "SUMMARY"})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
			for _, i := range issues {
				table.Append([]string{i.Id, i.Key, i.Fields.Status.Name, i.Fields.Summary})
			}
			table.Render()
		} else {
			logrus.Errorf("None of the provided issues %v servers are resolvable.", keys)
			os.Exit(1)
		}
	},
}

func init() {
}
