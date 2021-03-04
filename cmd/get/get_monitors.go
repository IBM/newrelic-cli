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

// monitorsCmd represents the monitors command
var monitorsCmd = &cobra.Command{
	Use:   "monitors",
	Short: "Display all synthetics monitors.",
	Example: `* nr get monitors
* nr get monitors -o json
* nr get monitors -o yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		monitorArray, err, ret := GetMonitors()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		if ret.IsContinue == false {
			os.Exit(1)
			return
		}

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
		printer.Print(monitorArray, os.Stdout)

		os.Exit(0)
	},
}

func GetMonitors() ([]*newrelic.Monitor, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	var opt *newrelic.MonitorListOptions
	opt = &newrelic.MonitorListOptions{}
	var pageSize int = 50
	var pageOffset int = 0

	var mList *newrelic.MonitorList = &newrelic.MonitorList{}

	for {
		opt.PageLimitOptions.Limit = pageSize
		opt.PageLimitOptions.Offset = pageOffset
		monitorList, resp, err := client.SyntheticsMonitors.ListAll(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_GET_MONITORS, resp.StatusCode, "pageSize:"+strconv.Itoa(pageSize)+",pageOffset:"+strconv.Itoa(pageOffset))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get one page monitors, pageSize '%d', pageOffset '%d'\n", statusCode, pageSize, pageOffset)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
		var monitorListLen = len(monitorList.Monitors)

		if monitorListLen == 0 {
			break
		} else {
			mList.Monitors = utils.MergeMonitorList(mList.Monitors, monitorList.Monitors)

			pageOffset = pageOffset + pageSize
		}
	}

	var mListLen = len(mList.Monitors)
	var monitorArray = make([]*newrelic.Monitor, mListLen)
	for i := 0; i < mListLen; i++ {
		monitor := mList.Monitors[i]
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
					ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_SCRIPT, err, tracker.ERR_REST_CALL, "")
					return nil, err, ret
				}
				if resp.StatusCode >= 400 {
					var statusCode = resp.StatusCode
					fmt.Printf("Response status code: %d. Get one monitor script, monitor id '%s', monitor name '%s'\n", statusCode, id, *monitor.Name)
					ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITOR_SCRIPT, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "monitor id: "+id+", monitor name: "+(*monitor.Name))
					return nil, err, ret
				}
			}
			monitor.Script = scriptText
			monitorArray[i] = monitor

		} else {
			monitorArray[i] = monitor
		}

	}

	for _, m := range monitorArray {
		m.Labels = make([]*string, 0)
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_MONITORS, nil, nil, "")
	return monitorArray, err, ret
}

func init() {
	GetCmd.AddCommand(monitorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:s
}
