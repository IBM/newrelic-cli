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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Display a single monitor by id.",
	Example: "nr get monitor <id>",
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

		id := args[0]

		monitor, err, _ := GetMonitorByID(id)

		var printer utils.Printer

		printer = &utils.JSONPrinter{}

		var output string
		flags := cmd.Flags()
		if flags.Lookup("output") != nil {
			output, err = cmd.Flags().GetString("output")
			if output == "yaml" {
				printer = &utils.YAMLPrinter{}
			}
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		printer.Print(monitor, os.Stdout)

		os.Exit(0)
	},
}

func GetMonitorByID(id string) (*newrelic.Monitor, error, tracker.ReturnValue) {
	fmt.Printf("Enter GetMonitorByID() func, monitor id: %s\n", id)
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	monitor, resp, err := client.SyntheticsMonitors.GetByID(context.Background(), id)

	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, err, tracker.ERR_REST_CALL, "")
		return nil, err, ret
	}

	tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, resp.StatusCode, "monitor id: "+id+", monitor name: "+(*monitor.Name))

	// if err != nil || resp.StatusCode >= 400 {
	// 	fmt.Printf("%v:%v\n", resp.Status, err)
	// 	return nil, err
	// }

	if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
		monitorID := monitor.ID
		var id string = ""
		id = *monitorID
		scriptText, resp, err := client.SyntheticsScript.GetByID(context.Background(), id)

		tracker.AppendRESTCallResult(client.SyntheticsScript, tracker.OPERATION_NAME_GET_MONITOR_SCRIPT, resp.StatusCode, "monitor id: "+id+", monitor name: "+(*monitor.Name))

		if resp.StatusCode == 404 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get one monitor script, monitor id '%s', monitor name '%s'\n", statusCode, id, *monitor.Name)
			s := new(newrelic.Script)
			s.ScriptText = new(string)
			scriptText = s
		} else {
			if err != nil {
				fmt.Println(err)
				// var st *newrelic.Script
				// st = &newrelic.Script{}
				// scriptText = st
				ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, err, tracker.ERR_REST_CALL, "")
				return nil, err, ret
			}
			if resp.StatusCode >= 400 {
				var statusCode = resp.StatusCode
				fmt.Printf("Response status code: %d. Get one monitor script, monitor id '%s', monitor name '%s'\n", statusCode, id, *monitor.Name)
				ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "monitor id: "+id+", monitor name: "+(*monitor.Name))
				return nil, err, ret
			}
		}
		monitor.Script = scriptText
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_MONITOR_BY_ID, nil, nil, "")
	return monitor, err, ret
}

var isAllMonitorsFetched bool = false
var allMonitors []*newrelic.Monitor

func IsMonitorNameExists(monitorName string) (bool, *newrelic.Monitor, error, tracker.ReturnValue) {
	// var allMonitorList []*newrelic.Monitor
	if isAllMonitorsFetched == false {
		allMonitorList, err, returnValue := GetMonitors()
		if returnValue.IsContinue == false {
			return false, nil, err, returnValue
		}
		allMonitors = allMonitorList
		isAllMonitorsFetched = true
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CHECK_MONITOR_NAME_EXISTS, nil, nil, "")
	for _, monitor := range allMonitors {
		if *monitor.Name == monitorName {
			return true, monitor, nil, ret
		}
	}

	return false, nil, nil, ret
}

func GetMonitorByName(monitorName string) (*newrelic.Monitor, error, tracker.ReturnValue) {

	isExist, monitor, err, ret := IsMonitorNameExists(monitorName)
	if ret.IsContinue == false {
		return nil, err, ret
	}

	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_MONITOR_BY_NAME, nil, nil, "")
	if isExist == true {
		return monitor, err, ret
	} else {
		return nil, err, ret
	}
}

func init() {
	GetCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
