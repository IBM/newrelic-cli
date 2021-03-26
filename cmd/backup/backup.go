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
package backup

import (
	"github.com/spf13/cobra"
)

// BackupCmd represents the update command
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a NewRelic resource using specified subcommand.",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// BackupCmd.PersistentFlags().String("foo", "", "A help for foo")
	BackupCmd.PersistentFlags().StringP("dir", "d", "", "Folder name to backup.")
	BackupCmd.MarkPersistentFlagRequired("dir")

	BackupCmd.PersistentFlags().BoolP("single-file", "s", false, "Save the configuration to a single file")

	BackupCmd.PersistentFlags().StringP("result-file-name", "r", "", "Result file name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// UpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
