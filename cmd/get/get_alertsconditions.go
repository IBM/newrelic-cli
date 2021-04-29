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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions",
	Short:   "Display alert conditions by alert policy id.",
	Aliases: []string{"ac", "alertcondition", "alertscondition"},
	Example: "nr get alertsconditions <policy_id> -o json",
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
		var conditionsOptions *newrelic.AlertsConditionsOptions
		conditionsOptions = new(newrelic.AlertsConditionsOptions)
		conditionsOptions.PolicyIDOptions = strconv.FormatInt(id, 10)

		var alertsConditionList *newrelic.AlertsConditionList

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
			if conditionType == "plugins" || conditionType == "synthetics" || conditionType == "ext" || conditionType == "nrql" || conditionType == "conditions" || conditionType == "infrastructure" {
				var cat newrelic.ConditionCategory
				switch conditionType {
				case "plugins":
					cat = newrelic.ConditionPlugins
				case "synthetics":
					cat = newrelic.ConditionSynthetics
				case "ext":
					cat = newrelic.ConditionExternalService
				case "nrql":
					cat = newrelic.ConditionNRQL
				case "infrastructure":
					cat = newrelic.ConditionInfrastructure
				default:
					cat = newrelic.ConditionDefault
				}
				list, resp, err := client.AlertsConditions.List(context.Background(), conditionsOptions, cat)
				if err != nil || resp.StatusCode >= 400 {
					fmt.Printf("%v\n", err)
					os.Exit(1)
					return
				}
				alertsConditionList = list
			} else {
				list, err := client.AlertsConditions.ListAll(context.Background(), conditionsOptions)
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
					return
				}
				alertsConditionList = list
				// location_failure_conditions uses different methods for pagination,
				// if a page requested doesn't exist, it returns the first page instead of empty entires.
				// use seperate logic to get location_failure_conditions, instead of adding it in ListAll
				list, resp, err := client.AlertsConditions.List(context.Background(), conditionsOptions, newrelic.ConditionLocation)
				if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
					alertsConditionList.AlertsLocationConditionList = list.AlertsLocationConditionList
				}
			}
		} else {
			list, err := client.AlertsConditions.ListAll(context.Background(), conditionsOptions)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
				return
			}
			alertsConditionList = list
			// location_failure_conditions uses different methods for pagination,
			// if a page requested doesn't exist, it returns the first page instead of empty entires.
			// use seperate logic to get location_failure_conditions, instead of adding it in ListAll
			list, resp, err := client.AlertsConditions.List(context.Background(), conditionsOptions, newrelic.ConditionLocation)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				alertsConditionList.AlertsLocationConditionList = list.AlertsLocationConditionList
			}
		}

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
			return
		}
		printer, err := utils.NewPriter(cmd)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		printer.Print(alertsConditionList, os.Stdout)

		os.Exit(0)
	},
}

func GetAllConditionsByAlertPolicyID(id int64) (*newrelic.AlertsConditionList, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	var allList *newrelic.AlertsConditionList = &newrelic.AlertsConditionList{}
	allList.AlertsDefaultConditionList = &newrelic.AlertsDefaultConditionList{}
	allList.AlertsExternalServiceConditionList = &newrelic.AlertsExternalServiceConditionList{}
	allList.AlertsNRQLConditionList = &newrelic.AlertsNRQLConditionList{}
	allList.AlertsPluginsConditionList = &newrelic.AlertsPluginsConditionList{}
	allList.AlertsSyntheticsConditionList = &newrelic.AlertsSyntheticsConditionList{}

	var conditionsOptions *newrelic.AlertsConditionsOptions
	conditionsOptions = new(newrelic.AlertsConditionsOptions)
	conditionsOptions.PolicyIDOptions = strconv.FormatInt(id, 10)

	var pageCount = 1

	for {
		conditionsOptions.Page = pageCount

		list, err := client.AlertsConditions.ListAll(context.Background(), conditionsOptions)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		var defaultConditionsLen = len(list.AlertsDefaultConditions)
		var externalServiceConditionsLen = len(list.AlertsExternalServiceConditions)
		var nrqlConditionsLen = len(list.AlertsNRQLConditions)
		var pluginsConditionsLen = len(list.AlertsPluginsConditions)
		var syntheticsConditionsLen = len(list.AlertsSyntheticsConditions)

		if defaultConditionsLen == 0 && externalServiceConditionsLen == 0 && nrqlConditionsLen == 0 && pluginsConditionsLen == 0 && syntheticsConditionsLen == 0 {
			break
		} else {
			//merge conditions list
			allList = utils.MergeAlertConditionList(allList, list)
			pageCount++
		}
	}

	// location_failure_conditions uses different methods for pagination,
	// if a page requested doesn't exist, it returns the first page instead of empty entires.
	// use seperate logic to get location_failure_conditions
	list, resp, err := client.AlertsConditions.List(context.Background(), conditionsOptions, newrelic.ConditionLocation)
	if err != nil || resp.StatusCode >= 400 {
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, fmt.Errorf("%v.Response: %v. Error: %v.", newrelic.ConditionLocation, resp, err), tracker.ERR_REST_CALL, "")
		return nil, err, ret
	}
	allList.AlertsLocationConditionList = &newrelic.AlertsLocationConditionList{}
	allList.AlertsLocationConditions = append(allList.AlertsLocationConditions, list.AlertsLocationConditions...)

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, nil, nil, "")

	return allList, err, ret
}

func GetConditionsByAlertPolicyIDAndConditionType(id int64, cat newrelic.ConditionCategory) (*newrelic.AlertsConditionList, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	var alertsConditionList *newrelic.AlertsConditionList = &newrelic.AlertsConditionList{}
	alertsConditionList.AlertsDefaultConditionList = &newrelic.AlertsDefaultConditionList{}
	alertsConditionList.AlertsExternalServiceConditionList = &newrelic.AlertsExternalServiceConditionList{}
	alertsConditionList.AlertsNRQLConditionList = &newrelic.AlertsNRQLConditionList{}
	alertsConditionList.AlertsPluginsConditionList = &newrelic.AlertsPluginsConditionList{}
	alertsConditionList.AlertsSyntheticsConditionList = &newrelic.AlertsSyntheticsConditionList{}
	alertsConditionList.AlertsInfrastructureConditionList = &newrelic.AlertsInfrastructureConditionList{}

	var conditionsOptions *newrelic.AlertsConditionsOptions
	conditionsOptions = new(newrelic.AlertsConditionsOptions)
	conditionsOptions.PolicyIDOptions = strconv.FormatInt(id, 10)

	var pageCount = 1

	for {
		conditionsOptions.Page = pageCount

		list, resp, err := client.AlertsConditions.List(context.Background(), conditionsOptions, cat)
		if err != nil {
			fmt.Printf("%v\n", err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		tracker.AppendRESTCallResult(client.AlertsConditions, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, resp.StatusCode, "pageCount:"+strconv.Itoa(pageCount))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get alert conditions, pageCount '%d'\n", statusCode, pageCount)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}

		var size = 0
		if cat == newrelic.ConditionDefault {
			size = len(list.AlertsDefaultConditions)
		} else if cat == newrelic.ConditionExternalService {
			size = len(list.AlertsExternalServiceConditions)
		} else if cat == newrelic.ConditionNRQL {
			size = len(list.AlertsNRQLConditions)
		} else if cat == newrelic.ConditionPlugins {
			size = len(list.AlertsPluginsConditions)
		} else if cat == newrelic.ConditionSynthetics {
			size = len(list.AlertsSyntheticsConditions)
		}

		if size == 0 {
			break
		} else {
			//merge conditions list
			alertsConditionList = utils.MergeAlertConditionList(alertsConditionList, list)
			pageCount++
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID, nil, nil, "")

	return alertsConditionList, err, ret
}

func IsConditionNameExists(alertPolicyID int64, condtionName string, cat newrelic.ConditionCategory) (bool, int64, error, tracker.ReturnValue) {
	list, err, ret := GetConditionsByAlertPolicyIDAndConditionType(alertPolicyID, cat)
	if ret.IsContinue == false {
		return false, -1, err, ret
	}

	if cat == newrelic.ConditionDefault {
		for _, condition := range list.AlertsDefaultConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	} else if cat == newrelic.ConditionExternalService {
		for _, condition := range list.AlertsExternalServiceConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	} else if cat == newrelic.ConditionNRQL {
		for _, condition := range list.AlertsNRQLConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	} else if cat == newrelic.ConditionPlugins {
		for _, condition := range list.AlertsPluginsConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	} else if cat == newrelic.ConditionSynthetics {
		for _, condition := range list.AlertsSyntheticsConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	} else if cat == newrelic.ConditionInfrastructure {
		for _, condition := range list.AlertsInfrastructureConditions {
			if *condition.Name == condtionName {
				return true, *condition.ID, err, ret
			}
		}
	}

	return false, -1, err, ret
}

func init() {
	GetCmd.AddCommand(alertsconditionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertsconditionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
