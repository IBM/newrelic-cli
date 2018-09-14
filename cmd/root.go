/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmd

import (
	"fmt"
	"os"

	addCmd "github.com/IBM/newrelic-cli/cmd/add"
	backupCmd "github.com/IBM/newrelic-cli/cmd/backup"
	createCmd "github.com/IBM/newrelic-cli/cmd/create"
	deleteCmd "github.com/IBM/newrelic-cli/cmd/delete"
	getCmd "github.com/IBM/newrelic-cli/cmd/get"
	insertCmd "github.com/IBM/newrelic-cli/cmd/insert"
	patchCmd "github.com/IBM/newrelic-cli/cmd/patch"
	restoreCmd "github.com/IBM/newrelic-cli/cmd/restore"
	updateCmd "github.com/IBM/newrelic-cli/cmd/update"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nr",
	Short: "nr is a command line tool for NewRelic",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(getCmd.GetCmd)
	rootCmd.AddCommand(deleteCmd.DeleteCmd)
	rootCmd.AddCommand(createCmd.CreateCmd)
	rootCmd.AddCommand(updateCmd.UpdateCmd)
	rootCmd.AddCommand(patchCmd.PatchCmd)
	rootCmd.AddCommand(backupCmd.BackupCmd)
	rootCmd.AddCommand(restoreCmd.RestoreCmd)
	rootCmd.AddCommand(addCmd.AddCmd)
	rootCmd.AddCommand(insertCmd.InsertCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}
