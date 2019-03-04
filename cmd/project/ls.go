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

package project

import (
	"github.com/olekukonko/tablewriter"
	"github.com/sotomskir/jira-cli/jiraApi"
	"os"

	"github.com/spf13/cobra"
)

// projectLsCmd represents the projectLs command
var projectLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		projects := jiraApi.GetProjects()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "KEY", "NAME"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		for _, p := range projects {
			table.Append([]string{p.Id, p.Key, p.Name})
		}
		table.Render() // Send output
	},
}

func init() {
	Cmd.AddCommand(projectLsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectLsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectLsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
