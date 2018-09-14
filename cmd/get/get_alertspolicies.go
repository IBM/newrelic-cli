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
package get

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

// alertspoliciesCmd represents the alertspolicies command
var alertspoliciesCmd = &cobra.Command{
	Use:     "alertspolicies",
	Short:   "Display all alerts_policies.",
	Aliases: []string{"ap", "alertpolicy", "alertspolicy"},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		var opt *newrelic.AlertsPolicyListOptions
		if filter, err := utils.GetArg(cmd, "policy"); err == nil {
			opt = &newrelic.AlertsPolicyListOptions{
				NameOptions: filter,
			}
		}
		alertsPolicyList, resp, err := client.AlertsPolicies.ListAll(context.Background(), opt)
		if err != nil || resp.StatusCode >= 400 {
			fmt.Printf("%v:%v\n", resp.Status, err)
			os.Exit(1)
			return
		}
		printer, err := utils.NewPriter(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		printer.Print(alertsPolicyList, os.Stdout)

		os.Exit(0)
	},
}

func GetAllAlertPolicies() (*newrelic.AlertsPolicyList, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_POLICIES, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	var allAlertList *newrelic.AlertsPolicyList = &newrelic.AlertsPolicyList{}

	var opt *newrelic.AlertsPolicyListOptions
	opt = &newrelic.AlertsPolicyListOptions{}

	var pageCount = 1
	for {
		opt.Page = pageCount
		alertsPolicyList, resp, err := client.AlertsPolicies.ListAll(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_POLICIES, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		tracker.AppendRESTCallResult(client.AlertsPolicies, tracker.OPERATION_NAME_GET_ALERT_POLICIES, resp.StatusCode, "pageCount:"+strconv.Itoa(pageCount))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get one page alert policies, pageCount '%d'\n", statusCode, pageCount)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_POLICIES, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
		var policiesListLen = len(alertsPolicyList.AlertsPolicies)

		if policiesListLen == 0 {
			break
		} else {
			allAlertList.AlertsPolicies = utils.MergeAlertPolicyList(allAlertList.AlertsPolicies, alertsPolicyList.AlertsPolicies)

			pageCount++
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_ALERT_POLICIES, nil, nil, "")
	return allAlertList, err, ret
}

func IsPolicyNameExists(policyName string) (bool, *newrelic.AlertsPolicy, error, tracker.ReturnValue) {
	allPolicyList, err, returnValue := GetAllAlertPolicies()
	if returnValue.IsContinue == false {
		return false, nil, err, returnValue
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CHECK_ALERT_POLICY_NAME_EXISTS, nil, nil, "")

	for _, policy := range allPolicyList.AlertsPolicies {
		if *policy.Name == policyName {
			return true, policy, nil, ret
		}
	}

	return false, nil, err, ret
}

func init() {
	GetCmd.AddCommand(alertspoliciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	alertspoliciesCmd.Flags().StringP("policy", "p", "", "policy name to filter returned result")
}
