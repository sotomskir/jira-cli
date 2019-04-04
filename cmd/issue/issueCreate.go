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
	"os"
)

var summary string
var description string
var issueType string
var projectKey string

// Cmd represents the issue command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create new issue",
	Run: func(cmd *cobra.Command, args []string) {
		issue, err := jiraApi.CreateIssue(projectKey, summary, description, issueType)
		if err != nil {
			logrus.Errorln(err)
			os.Exit(1)
		}
		logrus.Infof("Created key: %s %s\n", issue.Key, issue.Self)
	},
}

func init() {
	createCmd.Flags().StringVarP(&summary, "summary", "s", "", "Issue summary")
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Issue description")
	createCmd.Flags().StringVarP(&issueType, "type", "t", "", "Issue type")
	createCmd.Flags().StringVarP(&projectKey, "project", "p", "", "Project key")
	createCmd.MarkFlagRequired("summary")
	createCmd.MarkFlagRequired("description")
	createCmd.MarkFlagRequired("type")
	createCmd.MarkFlagRequired("project")
}
