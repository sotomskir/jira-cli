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
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/spf13/cobra"
	"sync"
)

var (
	summary string
	description string
	issueType string
	create bool
	deployment bool
)

// VersionCmd represents the issueVersion command
var VersionCmd = &cobra.Command{
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
				err := jiraApi.AssignVersion(issueKey, version, create, deployment, summary, description, issueType)
				if err == nil {
					logrus.Infof("Success version %s set for issue %s\n", version, issueKey)
				}
			}(issueKey, version)
		}
		wg.Wait()
	},
}

func init() {
	VersionCmd.Flags().StringVarP(&summary, "summary", "s", "", "Deployment issue summary.")
	VersionCmd.Flags().StringVarP(&description, "description", "d", "", "Deployment issue description.")
	VersionCmd.Flags().StringVarP(&issueType, "issue-type", "t", "", "Deployment issue type.")
	VersionCmd.Flags().BoolVarP(&create, "create", "c", true, "Create version if not exists.")
	VersionCmd.Flags().BoolVarP(&deployment, "create-deployment-issue", "i", true, "Create deployment issue for version if not exists.")
}
