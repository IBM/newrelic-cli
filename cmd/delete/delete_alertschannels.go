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
package delete

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var alertschannelsCmd = &cobra.Command{
	Use:     "alertschannels",
	Short:   "Delete alerts_channel by id.",
	Aliases: []string{"ac", "alertchannel", "alertschannel"},
	Example: "nr delete alertschannels <id>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("length of [flags] should be 1 instead of %d", len(args))
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			var err = fmt.Errorf("%q looks like a non-number.\n", args[0])
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		id, _ := strconv.ParseInt(args[0], 10, 64)
		resp, err := client.AlertsChannels.DeleteByID(context.Background(), id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		} else {
			fmt.Println(resp.Status)
			if resp.StatusCode >= 400 {
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
	},
}

func init() {
	DeleteCmd.AddCommand(alertschannelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertschannelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertschannelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
