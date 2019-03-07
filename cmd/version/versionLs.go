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

package version

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"os"
)

// versionLsCmd represents the versionLs command
var versionLsCmd = &cobra.Command{
	Use:   "ls PROJECY_KEY",
	Short: "List versions for project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]
		versions := jiraApi.GetVersions(projectKey)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "NAME", "ARCHIVED", "RELEASED", "PROJECT ID"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		for _, v := range versions {
			table.Append([]string{v.Id, v.Name, fmt.Sprintf("%t", v.Archived), fmt.Sprintf("%t", v.Released), fmt.Sprintf("%d", v.ProjectId)})
		}
		table.Render() // Send output
	},
}

func init() {
	Cmd.AddCommand(versionLsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionLsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionLsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
