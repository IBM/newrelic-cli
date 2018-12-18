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
package take

import (
	"github.com/spf13/cobra"
)

var TakeCmd = &cobra.Command{
	Use:   "take",
	Short: "Take template.",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// TakeCmd.PersistentFlags().String("foo", "", "A help for foo")
	TakeCmd.PersistentFlags().StringP("where", "w", "console", "Output where. console/file are supported")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// TakeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
