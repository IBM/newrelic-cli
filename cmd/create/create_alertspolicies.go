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
package create

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertspoliciesCmd = &cobra.Command{
	Use:     "alertspolicies",
	Short:   "Create alerts_policies from a file.",
	Aliases: []string{"ap", "alertpolicy", "alertspolicy"},
	Example: "nr create alertspolicies -f <example.yaml>",
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
		if reflect.DeepEqual(new(newrelic.AlertsPolicyEntity), p) {
			fmt.Printf("Error validating %q.\n", file)
			os.Exit(1)
			return
		}
		// start to create
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		alertsPolicy, resp, err := client.AlertsPolicies.Create(context.Background(), p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println(resp.Status)
		}

		printer, err := utils.NewPriter(cmd)

		if err != nil {
			fmt.Println(err)
			return
		}
		printer.Print(alertsPolicy, os.Stdout)

		if resp.StatusCode >= 400 {
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func CreateAlertsPolicyEntity(alertsPolcyEntity *newrelic.AlertsPolicyEntity) (*newrelic.AlertsPolicy, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_ALERT_POLICY, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	alertsPolicy, resp, err := client.AlertsPolicies.Create(context.Background(), alertsPolcyEntity)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS, err, tracker.ERR_REST_CALL, "")
		return nil, err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsPolicies, tracker.OPERATION_NAME_CREATE_ALERT_POLICY, resp.StatusCode, "alert policy name:"+(*alertsPolcyEntity.AlertsPolicy.Name))
		var statusCode = resp.StatusCode
		fmt.Printf("Response status code: %d. Create alert policy '%s', alert policy id: '%d'\n", statusCode, *alertsPolicy.AlertsPolicy.Name, *alertsPolicy.AlertsPolicy.ID)
		if resp.StatusCode >= 400 {
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_ALERT_POLICY, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_ALERT_POLICY, nil, nil, "")
	return alertsPolicy.AlertsPolicy, err, ret
}

func CreateAlertsPolicy(alertsPolcy *newrelic.AlertsPolicy) (*newrelic.AlertsPolicy, error, tracker.ReturnValue) {
	var alertsPolcyEntity *newrelic.AlertsPolicyEntity = &newrelic.AlertsPolicyEntity{}
	alertsPolcyEntity.AlertsPolicy = alertsPolcy
	return CreateAlertsPolicyEntity(alertsPolcyEntity)
}

func init() {
	CreateCmd.AddCommand(alertspoliciesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
