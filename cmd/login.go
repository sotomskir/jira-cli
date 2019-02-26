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

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"l"},
	Short:   "Login to Atlassian Jira server",
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("JIRA_SERVER_URL")
		if server == "" {
			server = getInput("JIRA server URL: ")
			viper.Set("JIRA_SERVER_URL",server)
		}
		user := viper.GetString("JIRA_USER")
		if user == "" {
			user = getInput("Username: ")
			viper.Set("JIRA_USER",user)
		}
		password := viper.GetString("JIRA_PASSWORD")
		if password == "" {
			password = getPasswd()
			//@FIXME password store in plain text
			viper.Set("JIRA_PASSWORD",password)
		}

		loggedIn := login(server, user, password)

		if loggedIn {
			viper.WriteConfig()
			logrus.Infof("Success, Logged in to: %s as: %s\n", server, user)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loginCmd.Flags().StringP("server", "s", "", "Jira server url. Also read from JIRA_SERVER_URL env variable")
	loginCmd.Flags().StringP("user", "u", "", "Jira username. Also read from JIRA_USER env variable")
	loginCmd.Flags().StringP("password", "p", "", "Jira password. Also read from JIRA_PASSWORD env variable")
	viper.BindPFlag("JIRA_SERVER_URL", loginCmd.Flags().Lookup("server"))
	viper.BindPFlag("JIRA_USER", loginCmd.Flags().Lookup("user"))
	viper.BindPFlag("JIRA_PASSWORD", loginCmd.Flags().Lookup("password"))
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.Trim(text, "\n")
}

func getPasswd() string {
	fmt.Print("Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(bytePassword)
}

type JiraUser struct {
	Name string `json:"name"`
}

func login(server string, user string, password string) bool {
	url := server + "/rest/auth/1/session"

	httpClient := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	req.SetBasicAuth(user, password)
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		logrus.Errorln(getErr)
		os.Exit(1)
	}

	if res.StatusCode == 401 {
		logrus.Errorln("Server responded with status: 401 Unauthorized")
		os.Exit(1)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logrus.Errorln("Cannot read response body: " + readErr.Error())
		os.Exit(1)
	}

	jiraUser := JiraUser{}
	jsonErr := json.Unmarshal(body, &jiraUser)
	if jsonErr != nil {
		logrus.Errorf("%s\n", body)
		logrus.Errorln("Server responded invalid JSON")
		os.Exit(1)
	}

	return true
}
