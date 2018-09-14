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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertspoliciesCmd = &cobra.Command{
	Use:     "alertspolicies",
	Short:   "Delete alerts_policy by id.",
	Aliases: []string{"ap", "alertpolicy", "alertspolicy"},
	Example: "nr delete alertspolicies <id>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("length of [flags] should be 1 instead of %d", len(args))
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			return fmt.Errorf("%q looks like a non-number.\n", args[0])
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
		resp, err := client.AlertsPolicies.DeleteByID(context.Background(), id)
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

func DeletePolicyByID(alertPolicyID int64) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	resp, err := client.AlertsPolicies.DeleteByID(context.Background(), alertPolicyID)

	if err != nil {
		fmt.Printf("Failed to delete alert policy, %v\n", err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID, err, tracker.ERR_REST_CALL, "")
		return err, ret
	}

	tracker.AppendRESTCallResult(client.AlertsConditions, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID, resp.StatusCode, "alert policy id:"+strconv.FormatInt(alertPolicyID, 10))

	if resp.StatusCode >= 400 {
		var statusCode = resp.StatusCode
		fmt.Printf("response status code: %d. Remove alert condition id '%d'\n", statusCode, alertPolicyID)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
		return err, ret
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID, nil, nil, "")
	return nil, ret
}

func DeletePolicyByName(alertPolicyName string) (error, tracker.ReturnValue) {
	list, err, ret := get.GetAllAlertPolicies()
	if err != nil {
		fmt.Println(err)
		ret.IsContinue = false
		return err, ret
	}
	if ret.IsContinue == false {
		return ret.OriginalError, ret
	}
	for _, policy := range list.AlertsPolicies {
		if *policy.Name == alertPolicyName {
			err, ret := DeletePolicyByID(*policy.ID)
			if err != nil {
				ret.IsContinue = false
				return err, ret
			}
			if ret.IsContinue == false {
				return ret.OriginalError, ret
			}
		}
	}
	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_DELETE_ALERT_POLICY_BY_NAME, nil, nil, "")
	return nil, ret
}

func init() {
	DeleteCmd.AddCommand(alertspoliciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
