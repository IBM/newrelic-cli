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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions <condition-id>",
	Short:   "Update alerts_conditions from a file.",
	Aliases: []string{"ac", "alertcondition", "alertscondition"},
	Example: "nr update alertsconditions <condition-id> -f <example.yaml>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("Please give condition id, (nr get alertsconditions <condition-id> [flags])")
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			var err = fmt.Errorf("%q looks like a non-number", args[0])
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		return nil
	}, Run: func(cmd *cobra.Command, args []string) {
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

		flags := cmd.Flags()
		var conditionType string
		var errConditionType error
		if flags.Lookup("type-condition") == nil {
			fmt.Println("Can not find type-condition argument.")
			os.Exit(1)
			return
		}
		alertConditionID, _ := strconv.ParseInt(args[0], 10, 64)
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
			cat = newrelic.ConditionPlugins
			ac.AlertsPluginsConditionEntity = ace
			// alertConditionID = *ac.AlertsPluginsConditionEntity.AlertsPluginsCondition.ID
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
			cat = newrelic.ConditionSynthetics
			ac.AlertsSyntheticsConditionEntity = ace
			// alertConditionID = *ac.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.ID
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
			cat = newrelic.ConditionExternalService
			ac.AlertsExternalServiceConditionEntity = ace
			// alertConditionID = *ac.AlertsExternalServiceConditionEntity.AlertsExternalServiceCondition.ID
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
			cat = newrelic.ConditionNRQL
			ac.AlertsNRQLConditionEntity = ace
			// alertConditionID = *ac.AlertsNRQLConditionEntity.AlertsNRQLCondition.ID
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
			cat = newrelic.ConditionInfrastructure
			ac.AlertsInfrastructureConditionEntity = ace
			// alertConditionID = *ac.AlertsInfrastructureConditionEntity.AlertsInfrastructureCondition.ID
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
			cat = newrelic.ConditionDefault
			ac.AlertsDefaultConditionEntity = ace
			// alertConditionID = *ac.AlertsDefaultConditionEntity.AlertsDefaultCondition.ID
		}
		// start to update
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		_, resp, err := client.AlertsConditions.Update(context.Background(), cat, ac, alertConditionID)
		if err != nil {
			fmt.Printf("Failed to update condition, %v\n", err)
			os.Exit(1)
			return
		}
		statusCode := resp.StatusCode
		fmt.Printf("Response status code: %d. Update alert conditions, condition id: '%d'\n", statusCode, alertConditionID)
		if statusCode >= 400 {
			os.Exit(1)
			return
		}
		os.Exit(0)
	},
}

func UpdateCondition(cat newrelic.ConditionCategory, ac *newrelic.AlertsConditionEntity, alertConditionID int64) (*newrelic.AlertsConditionEntity, error, tracker.ReturnValue) {
	// start to update
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	alertsConditionEntity, resp, err := client.AlertsConditions.Update(context.Background(), cat, ac, alertConditionID)
	if err != nil {
		fmt.Printf("Failed to update condition, %v\n", err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID, err, tracker.ERR_REST_CALL, "")
		return nil, err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsConditions, tracker.OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID, resp.StatusCode, "alert condtion id:"+strconv.FormatInt(alertConditionID, 10))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Update alert conditions, condition id: '%d'\n", statusCode, alertConditionID)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID, nil, nil, "")
	return alertsConditionEntity, err, ret
}

func init() {
	UpdateCmd.AddCommand(alertsconditionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertsconditionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertsconditionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
