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
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var customeventsCmd = &cobra.Command{
	Use:     "customevents",
	Short:   "Insert custom events.",
	Aliases: []string{"m"},
	Example: `nr insert customevents`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		flags := cmd.Flags()

		var insertKey = ""
		if flags.Lookup("insert-key") != nil {
			insertKey, err = cmd.Flags().GetString("insert-key")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			if insertKey == "" {
				fmt.Println("Please give New Relic insert key.")
			}
		}

		var accountID = ""
		if flags.Lookup("account-id") != nil {
			accountID, err = cmd.Flags().GetString("account-id")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			if accountID == "" {
				fmt.Println("Please give New Relic account ID.")
			}
		}

		var filePath = ""

		if flags.Lookup("file") != nil {
			filePath, err = cmd.Flags().GetString("file")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
		}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Unable to open file '%v': %v\n", filePath, err)
			os.Exit(1)
			return
		}

		var fileContent = string(bytes)

		client, err := utils.GetNewRelicClient("insights")
		if err != nil {
			fmt.Println(err)
			tracker.ToReturnValue(false, tracker.OPERATION_NAME_INSERT_CUSTOM_EVENTS, err, tracker.ERR_CREATE_NR_CLINET, "")
			os.Exit(1)
			return
		}

		resp, bytes, err := client.CustomEvents.Insert(context.Background(), insertKey, accountID, fileContent)
		retMsg := string(bytes)

		if err != nil {
			tracker.ToReturnValue(false, tracker.OPERATION_NAME_INSERT_CUSTOM_EVENTS, err, tracker.ERR_REST_CALL, "")
			fmt.Printf("Failed to insert custom events.")
			fmt.Println(err)
			os.Exit(1)
			return
		}
		fmt.Printf("Response status code: %d.\n", resp.StatusCode)

		tracker.AppendRESTCallResult(client.CustomEvents, tracker.OPERATION_NAME_INSERT_CUSTOM_EVENTS, resp.StatusCode, retMsg)

		var ret tracker.ReturnValue
		if resp.StatusCode >= 400 {
			ret = tracker.ToReturnValue(false, tracker.OPERATION_NAME_INSERT_CUSTOM_EVENTS, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, retMsg)
			fmt.Printf("Failed to insert custom events. status code is not 200, Response status code: %d\n", resp.StatusCode)
			fmt.Println(retMsg)
		} else {
			ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_INSERT_CUSTOM_EVENTS, nil, nil, "")
		}

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()

		if ret.IsContinue == false {
			fmt.Printf("Failed to insert custom events.")
			os.Exit(1)
			return
		}

		fmt.Println("Insert custom events successful.")
		os.Exit(0)
	},
}

func init() {
	InsertCmd.AddCommand(customeventsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
