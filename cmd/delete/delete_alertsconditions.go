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

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions",
	Short:   "Delete alerts_conditions by id.",
	Aliases: []string{"ac", "alertcondition", "alertscondition"},
	Example: "nr delete alertsconditions <id>",
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

		conditionPolicyID, _ := strconv.ParseInt(args[0], 10, 64)

		flags := cmd.Flags()
		var conditionType string
		var errConditionType error
		if flags.Lookup("type-condition") != nil {
			conditionType, errConditionType = cmd.Flags().GetString("type-condition")
			if errConditionType != nil {
				fmt.Printf("error accessing flag %s for command %s: %v\n", "type-condition", cmd.Name(), errConditionType)
				os.Exit(1)
				return
			}

			var cat newrelic.ConditionCategory
			if conditionType == "plugins" {
				cat = newrelic.ConditionPlugins
			} else if conditionType == "synthetics" {
				cat = newrelic.ConditionSynthetics
			} else if conditionType == "ext" {
				cat = newrelic.ConditionExternalService
			} else if conditionType == "nrql" {
				cat = newrelic.ConditionNRQL
			} else {
				cat = newrelic.ConditionDefault
			}
			// start to delete
			client, err := utils.GetNewRelicClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			resp, err := client.AlertsConditions.DeleteByID(context.Background(), cat, conditionPolicyID)
			if err != nil {
				fmt.Printf("Failed to delete condition, %v\n", err)
				os.Exit(1)
				return
			} else {
				fmt.Println(resp.Status)
				if resp.StatusCode >= 400 {
					os.Exit(1)
					return
				}
			}
		} else {
			fmt.Println("Can not find type-condition argument.")
			os.Exit(1)
			return
		}

		os.Exit(0)
	},
}

func DeleteCondition(cat newrelic.ConditionCategory, conditionPolicyID int64) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_CONDITION, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	resp, err := client.AlertsConditions.DeleteByID(context.Background(), cat, conditionPolicyID)
	if err != nil {
		fmt.Printf("Failed to delete condition, %v\n", err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_CONDITION, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsConditions, tracker.OPERATION_NAME_DELETE_ALERT_CONDITION, resp.StatusCode, "monitor id: "+strconv.FormatInt(conditionPolicyID, 10))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Remove alert condition id '%d'\n", statusCode, conditionPolicyID)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_ALERT_CONDITION, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_DELETE_ALERT_CONDITION, nil, nil, "")
	return nil, ret
}

func init() {
	DeleteCmd.AddCommand(alertsconditionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertsconditionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertsconditionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
