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
package insert

import (
	"github.com/spf13/cobra"
)

// InsertCmd represents the insert command
var InsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Inster NewRelic custom events using specified subcommand.",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// InsertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	InsertCmd.PersistentFlags().StringP("file", "f", "", "File name to store custom events data in JSON.")
	InsertCmd.MarkPersistentFlagRequired("file")

	InsertCmd.PersistentFlags().StringP("insert-key", "i", "", "New Relic insert key.")
	InsertCmd.MarkPersistentFlagRequired("insert-key")

	InsertCmd.PersistentFlags().StringP("account-id", "a", "", "New Relic account ID.")
	InsertCmd.MarkPersistentFlagRequired("account-id")
}
