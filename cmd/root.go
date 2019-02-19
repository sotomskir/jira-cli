// Copyright © 2019 Robert Sotomski <sotomskie@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/sotomskir/jira-cli/jiraApi"
	"github.com/sotomskir/jira-cli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var cfgFile string
var noColor bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jira-cli",
	Short: "CLI client for Atlassian Jira REST API.",
	Long: `
jira-cli is a CLI for Atlassian Jira REST API.
It can be used with CI/CD pipelines to automate workflow.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorLn(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jira-cli.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable ANSI color output")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "admin", "Jira username")
	//rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "admin", "Jira password")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if noColor {
		color.NoColor = true
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logger.ErrorLn(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jira-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jira-cli")
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logger.ErrorLn(err)
			os.Exit(1)
		}
		viper.WriteConfigAs(path.Join(home, "/.jira-cli.yaml"))
	}
	jiraApi.Initialize(viper.GetString("JIRA_SERVER_URL"), viper.GetString("JIRA_USER"), viper.GetString("JIRA_PASSWORD"))
}
