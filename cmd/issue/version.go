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
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"sync"
)

// versionCmd represents the issueVersion command
var versionCmd = &cobra.Command{
	Use:   "version VERSION ISSUE_KEY [ISSUE_KEY...]",
	Short: "Set issue fix version",
	Long: `Set issue fix version. 
If version is already set it will not be overwritten. 
If version does not exist it will be created`,
	Aliases: []string{"v"},
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		issueKeys := args[1:]

		var wg sync.WaitGroup
		for _, issueKey := range issueKeys {
			wg.Add(1)
			go func(issueKey string, version string) {
				defer wg.Done()
				_, err := jiraApi.SetFixVersion(issueKey, version)
				if err == nil {
					logrus.Infof("Success version %s set for issue %s\n", version, issueKey)
				}
			}(issueKey, version)
		}
		wg.Wait()
	},
}

func init() {
	Cmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
