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
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// transitionCmd represents the issueTransition command
var transitionCmd = &cobra.Command{
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
				jiraApi.TransitionIssue(workflow, issueKey, targetState)
			}(workflow, issueKey, targetState)
		}
		wg.Wait()
	},
}

func init() {
	Cmd.AddCommand(transitionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transitionCmd.PersistentFlags().String("foo", "", "A help for foo")
	transitionCmd.Flags().StringP("workflow", "w", "workflow.yaml", "Workflow definition local file or http URL")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transitionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
