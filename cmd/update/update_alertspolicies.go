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
package update

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var alertspoliciesCmd = &cobra.Command{
	Use:     "alertspolicies",
	Short:   "Update alerts_policies from a file.",
	Aliases: []string{"ap", "alertpolicy", "alertspolicy"},
	Example: "nr update alertspolicies -f <example.yaml>",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := utils.GetArg(cmd, "file")
		if err != nil {
			fmt.Printf("Unable to get argument 'file': %v\n", err)
			os.Exit(1)
			return
		}
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			fmt.Printf("Unable to open file '%v': %v\n", file, err)
			os.Exit(1)
			return
		}
		// validation
		decorder := utils.NewYAMLOrJSONDecoder(f, 4096)
		var p = new(newrelic.AlertsPolicyEntity)
		err = decorder.Decode(p)
		if err != nil {
			fmt.Printf("Unable to decode %q: %v\n", file, err)
			os.Exit(1)
			return
		}
		if reflect.DeepEqual(new(newrelic.AlertsPolicy), p) {
			fmt.Printf("Error validating %q.\n", file)
			os.Exit(1)
			return
		}
		if p.AlertsPolicy.ID == nil {
			fmt.Printf("Can't find {.policy.id} in %q.\n", file)
			os.Exit(1)
			return
		}
		// start Update
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		_, resp, err := client.AlertsPolicies.Update(context.Background(), p, *p.AlertsPolicy.ID)
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

func UpdateByPolicyID(policy *newrelic.AlertsPolicyEntity, alertPolicyID int64) (*newrelic.AlertsPolicyEntity, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	policyEntity, resp, err := client.AlertsPolicies.Update(context.Background(), policy, alertPolicyID)
	if err != nil {
		fmt.Printf("Failed to update alert policy, %v\n", err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID, err, tracker.ERR_REST_CALL, "")
		return nil, err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsPolicies, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID, resp.StatusCode, "alert policy id:"+strconv.FormatInt(alertPolicyID, 10)+",alert policy name:"+(*policy.AlertsPolicy.Name))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Update alert policy, alert policy id: '%d'\n", statusCode, alertPolicyID)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID, nil, nil, "")
	return policyEntity, err, ret
}

func UpdateByPolicyName(policy *newrelic.AlertsPolicy, policyName string) (*newrelic.AlertsPolicy, error, tracker.ReturnValue) {
	list, err, ret := get.GetAllAlertPolicies()
	if err != nil {
		fmt.Println(err)
		return nil, err, ret
	} else {
		if ret.IsContinue == false {
			return nil, err, ret
		}
	}
	for _, p := range list.AlertsPolicies {
		if *p.Name == policyName {
			policy.ID = p.ID
			var pEntity *newrelic.AlertsPolicyEntity = &newrelic.AlertsPolicyEntity{}
			pEntity.AlertsPolicy = policy
			newPolicy, err, ret := UpdateByPolicyID(pEntity, *policy.ID)
			if err != nil {
				fmt.Println(err)
				return nil, err, ret
			} else {
				if ret.IsContinue == false {
					return nil, err, ret
				}
			}
			return newPolicy.AlertsPolicy, err, ret
		}
	}
	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_BY_NAME, nil, nil, "")
	return nil, err, ret
}

func init() {
	UpdateCmd.AddCommand(alertspoliciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
