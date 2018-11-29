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
	"strconv"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions",
	Short:   "Create alerts_conditions from a file.",
	Aliases: []string{"ac", "alertcondition", "alertscondition"},
	Example: "nr create alertsconditions -f <example.yaml>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("length of [flags] should be 1 instead of %d", len(args))
			fmt.Printf("%v\n", err)
			os.Exit(1)
			return err
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			var err = fmt.Errorf("%q looks like a non-number.\n", args[0])
			fmt.Printf("%v\n", err)
			os.Exit(1)
			return err
		}
		return nil
	},
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

		alertPolicyID, _ := strconv.ParseInt(args[0], 10, 64)

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
			var ac = new(newrelic.AlertsConditionEntity)

			switch conditionType {
			case "plugins":
				cat = newrelic.ConditionPlugins
				var ace = new(newrelic.AlertsPluginsConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for plugins type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsPluginsConditionEntity), ace) {
					fmt.Printf("Error validating for plugins type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsPluginsConditionEntity = ace

				cat = newrelic.ConditionPlugins
			case "synthetics":
				var ace = new(newrelic.AlertsSyntheticsConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for synthetics type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsSyntheticsConditionEntity), ace) {
					fmt.Printf("Error validating for synthetics type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsSyntheticsConditionEntity = ace

				cat = newrelic.ConditionSynthetics
			case "ext":
				var ace = new(newrelic.AlertsExternalServiceConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for ext type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsExternalServiceConditionEntity), ace) {
					fmt.Printf("Error validating for ext type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsExternalServiceConditionEntity = ace

				cat = newrelic.ConditionExternalService
			case "nrql":
				var ace = new(newrelic.AlertsNRQLConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for nrql type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsNRQLConditionEntity), ace) {
					fmt.Printf("Error validating for nrql type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsNRQLConditionEntity = ace

				cat = newrelic.ConditionNRQL
			case "infrastructure":
				var ace = new(newrelic.AlertsInfrastructureConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for infrastructure type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsInfrastructureConditionEntity), ace) {
					fmt.Printf("Error validating for infrastructure type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsInfrastructureConditionEntity = ace

				cat = newrelic.ConditionInfrastructure
			default:
				var ace = new(newrelic.AlertsDefaultConditionEntity)
				err = decorder.Decode(ace)
				if err != nil {
					fmt.Printf("Unable to decode for default type condition %q: %v\n", file, err)
					os.Exit(1)
					return
				}
				if reflect.DeepEqual(new(newrelic.AlertsDefaultConditionEntity), ace) {
					fmt.Printf("Error validating for default type condition %q.\n", file)
					os.Exit(1)
					return
				}
				ac.AlertsDefaultConditionEntity = ace

				cat = newrelic.ConditionDefault
			}
			// start to create
			client, err := utils.GetNewRelicClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			_, resp, err := client.AlertsConditions.Create(context.Background(), cat, ac, alertPolicyID)
			if err != nil {
				fmt.Printf("Failed to create condition, %v\n", err)
			} else {
				var statusCode = resp.StatusCode
				fmt.Printf("Response status code: %d. Create alert conditions, alert policy id: '%d'\n", statusCode, alertPolicyID)
				if resp.StatusCode >= 400 {
					os.Exit(1)
				}
			}

		}

		os.Exit(0)
	},
}

func CreateCondition(cat newrelic.ConditionCategory, ac *newrelic.AlertsConditionEntity, alertPolicyID int64) (*newrelic.AlertsConditionEntity, error, tracker.ReturnValue) {
	// start to create
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_ALERT_CONDITIION, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	alertsConditionEntity, resp, err := client.AlertsConditions.Create(context.Background(), cat, ac, alertPolicyID)
	if err != nil {
		fmt.Printf("Failed to create alert condition, %v\n", err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_ALERT_CONDITIION, err, tracker.ERR_REST_CALL, "")
		return nil, err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsConditions, tracker.OPERATION_NAME_CREATE_ALERT_CONDITIION, resp.StatusCode, "")

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Create alert conditions, alert policy id: '%d'\n", statusCode, alertPolicyID)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_ALERT_CONDITIION, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_ALERT_CONDITIION, nil, nil, "")
	return alertsConditionEntity, err, ret
}

func init() {
	CreateCmd.AddCommand(alertsconditionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertsconditionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertsconditionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
