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
	"github.com/sotomskir/jira-cli/cmd/issue/transition"
	"github.com/sotomskir/jira-cli/cmd/issue/version"
	"github.com/sotomskir/jira-cli/cmd/issue/worklog"
	"github.com/spf13/cobra"
)

// Cmd represents the issue command
var Cmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"i"},
	Short:   "Manage Jira issues",
}

func init() {
	Cmd.AddCommand(inspectCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(worklog.Cmd)
	Cmd.AddCommand(version.VersionCmd)
	Cmd.AddCommand(transition.TransitionCmd)
}
