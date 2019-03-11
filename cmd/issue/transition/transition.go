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

package transition

import (
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/jiraApi"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// TransitionCmd represents the issueTransition command
var TransitionCmd = &cobra.Command{
	Use:     "transition STATE ISSUE_KEY [ISSUE_KEY...]",
	Aliases: []string{"t"},
	Short:   "Transition issue status to given state",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetState := args[0]
		issueKeys := args[1:]
		workflow, err := cmd.Flags().GetString("workflow")
		if err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
		var wg sync.WaitGroup
		for _, issueKey := range issueKeys {
			wg.Add(1)
			go func(workflow string, issueKey string, targetState string) {
				defer wg.Done()
				if _, err := jiraApi.TransitionIssue(workflow, issueKey, targetState); err != nil {
					logrus.Errorln(err)
				}
			}(workflow, issueKey, targetState)
		}
		wg.Wait()
	},
}

func init() {
	TransitionCmd.AddCommand(testWorkflowCmd)
	TransitionCmd.Flags().StringP("workflow", "w", "workflow.yaml", "Workflow definition local file or http URL")
}
