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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sotomskir/jira-cli/logger"
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

var user string
var password string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login SERVER",
	Aliases: []string{"l"},
	Short:   "Login to Atlassian Jira server",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		logger.InfoF("Login to %s\n", server)
		if user == "" {
			user = getInput("Username: ")
		}
		if password == "" {
			password = getPasswd()
		}

		loggedIn := login(server, user, password)

		if loggedIn {
			saveConfig(server, user, password)
			logger.SuccessF("Success, Logged in to %s\n", server)
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
	loginCmd.Flags().StringVarP(&user, "user", "u", "", "Jira username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Jira password")

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
		logger.ErrorLn(err)
		os.Exit(1)
	}

	req.SetBasicAuth(user, password)
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		logger.ErrorLn(getErr)
		os.Exit(1)
	}

	if res.StatusCode == 401 {
		logger.ErrorLn("Server responded with status: 401 Unauthorized")
		os.Exit(1)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logger.ErrorLn("Cannot read response body: " + readErr.Error())
		os.Exit(1)
	}

	jiraUser := JiraUser{}
	jsonErr := json.Unmarshal(body, &jiraUser)
	if jsonErr != nil || jiraUser.Name != user {
		logger.ErrorF("%s\n", body)
		logger.ErrorLn("Server responded invalid JSON")
		os.Exit(1)
	}

	return true
}

func saveConfig(server string, user string, password string) {
	viper.Set("server", server)
	viper.Set("auth", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))))
	viper.WriteConfig()
}
