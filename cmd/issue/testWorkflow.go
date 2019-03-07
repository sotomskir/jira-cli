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

	"github.com/spf13/cobra"
)

// testWorkflowCmd represents the issueTransitionTest command
var testWorkflowCmd = &cobra.Command{
	Use:   "test ISSUE_KEY",
	Short: "Run through all transitions to test workflow definition yaml file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		workflow, err := cmd.Flags().GetString("workflow")
		if err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
		jiraApi.TestTransitions(workflow, args[0])
		logrus.Infoln("issueTransitionTest PASSED")
	},
}

func init() {
	Cmd.AddCommand(testWorkflowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//testWorkflowCmd.PersistentFlags().StringP("workflow", "w", "workflow.yaml", "Workflow definition file")
	testWorkflowCmd.Flags().StringP("workflow", "w", "workflow.yaml", "Workflow definition file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
