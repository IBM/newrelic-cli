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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Delete one synthetics monitor by id.",
	Aliases: []string{"m"},
	Example: "nr delete monitor <id>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("length of [flags] should be 1 instead of %d", len(args))
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := utils.GetNewRelicClient("synthetics")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		id := string(args[0])
		resp, err := client.SyntheticsMonitors.DeleteByID(context.Background(), &id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		} else {
			fmt.Println(resp.StatusCode)
			if resp.StatusCode >= 400 {
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
	},
}

func DeleteMonitorByID(id string) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_MONITOR, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	resp, err := client.SyntheticsMonitors.DeleteByID(context.Background(), &id)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_MONITOR, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_DELETE_MONITOR, resp.StatusCode, "monitor id: "+id)
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Remove monitor id '%s'\n", statusCode, id)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_DELETE_MONITOR, nil, nil, "")
	return nil, ret
}

func init() {
	DeleteCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
