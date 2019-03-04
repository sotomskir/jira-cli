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
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/sotomskir/jira-cli/cmd/issue"
	"github.com/sotomskir/jira-cli/cmd/project"
	"github.com/sotomskir/jira-cli/cmd/version"
	"github.com/sotomskir/jira-cli/jiraApi"
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
	Version: "0.7.0-SNAPSHOT",
	Short: "CLI client for Atlassian Jira REST API.",
	Long: `
jira-cli is a CLI for Atlassian Jira REST API.
It can be used with CI/CD pipelines to automate workflow.`,
}
var debug bool
var trace bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorln(err)
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
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug output")
	rootCmd.PersistentFlags().BoolVar(&trace, "trace", false, "Trace output")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(issue.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(project.Cmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: noColor,
		ForceColors: true,
		DisableTimestamp: true,
		DisableLevelTruncation: true,
	})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			logrus.Errorln(err)
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
			logrus.Errorln(err)
			os.Exit(1)
		}
		viper.WriteConfigAs(path.Join(home, "/.jira-cli.yaml"))
	}
	jiraApi.Initialize(viper.GetString("JIRA_SERVER_URL"), viper.GetString("JIRA_USER"), viper.GetString("JIRA_PASSWORD"))
}
