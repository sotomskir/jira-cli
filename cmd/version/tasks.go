// Copyright © 2019 Bartosz Wojtyła <bartosz.wojtyla@gmail.com>
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

package version

import (
	"github.com/olekukonko/tablewriter"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// versionCreateCmd represents the versionCreate command
var tasksCmd = &cobra.Command{
	Use:     "tasks PROJECT_KEY VERSION ISSUETYPES",
	Aliases: []string{"c"},
	Short:   "Get tasks in version",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]
		version := args[1]
		issueTypes := args[2]
		response, _ := jiraApi.GetIssuesInVersions(projectKey, version, issueTypes)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"KEY", "SUMMARY"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		for _, p := range response.Issues {
			table.Append([]string{p.Key, p.Fields.Summary})
		}
		table.SetFooter([]string{"Total", strconv.Itoa(response.Total)})
		table.Render()
	},
}

func init() {
	Cmd.AddCommand(tasksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
