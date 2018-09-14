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
package restore

import (
	"github.com/spf13/cobra"
)

// RestoreCmd represents the update command
var RestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a NewRelic resource using specified subcommand.",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RestoreCmd.PersistentFlags().String("foo", "", "A help for foo")
	RestoreCmd.PersistentFlags().StringP("dir", "d", "", "Folder name where are files in to restore.")
	// RestoreCmd.MarkPersistentFlagRequired("dir")

	RestoreCmd.PersistentFlags().StringP("update-mode", "m", "skip", "Update mode. skip|override|clean are supported")

	RestoreCmd.PersistentFlags().StringP("files", "f", "", "File names to restore. Multiple files separated by commma.")
	RestoreCmd.PersistentFlags().StringP("File", "F", "", "Logging file stored failed resources to restore.")

	RestoreCmd.PersistentFlags().StringP("result-file-name", "r", "", "Result file name")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// UpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
