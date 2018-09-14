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
package create

import (
	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a NewRelic resource using specified subcommand.",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	CreateCmd.PersistentFlags().StringP("file", "f", "", "Filename to create resource with, yaml/json format is supported.")
	CreateCmd.MarkPersistentFlagRequired("file")

	CreateCmd.PersistentFlags().StringP("type-condition", "t", "default", "Alert condition type. Only used for 'alertsconditions' command. default|sythentics|ext|plugin|nrql are supported")

	CreateCmd.PersistentFlags().StringP("output", "o", "json", "Output format. json/yaml are supported")

	CreateCmd.PersistentFlags().StringP("script-file", "s", "", "Synthetics monitor javascript file name. Only used for 'monitor' command.")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// CreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
