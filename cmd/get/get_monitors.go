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

const (
	MaxConcurrentTask int = 10
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

	scriptChMap := make(map[string]chan *newrelic.Script)
	chTaskCtrl := make(chan struct{}, MaxConcurrentTask)
	defer close(chTaskCtrl)

	for i := 0; i < mListLen; i++ {
		monitor := mList.Monitors[i]
		if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
			r := make(chan *newrelic.Script)
			id := *(monitor.ID)
			name := *(monitor.Name)
			go func() {
				defer close(r)
				chTaskCtrl <- struct{}{}
				fmt.Printf("Fetching script for Monitor: %s\n", name)
				scriptText, resp, err := client.SyntheticsScript.GetByID(context.Background(), id)
				<-chTaskCtrl
				statusCode := -1
				if resp != nil {
					statusCode = resp.StatusCode
				}
				tracker.AppendRESTCallResult(client.SyntheticsScript, tracker.OPERATION_NAME_GET_MONITOR_SCRIPT, statusCode, "monitor id: "+id+", monitor name: "+name)
				if err != nil {
					fmt.Println(err)
					r <- nil
					return
				}
				r <- scriptText
				return
			}()
			scriptChMap[id] = r
		}
	}

	for i := 0; i < mListLen; i++ {
		monitor := mList.Monitors[i]
		if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
			monitorID := monitor.ID
			var id string = ""
			id = *monitorID
			scriptText := <-scriptChMap[id]
			monitor.Script = scriptText
			monitorArray[i] = monitor

		} else {
			monitorArray[i] = monitor
		}

	}
	tags, err := GetMonitorTags()
	if err == nil {
		for _, m := range monitorArray {
			if tags[*m.ID] != nil {
				m.Tags = tags[*m.ID].Tags
			}
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_MONITORS, err, nil, "")
	return monitorArray, err, ret
}

func GetMonitorTags() (map[string]*newrelic.EntitySearchResultsMonitor, error) {
	client, err := utils.GetNewRelicClient("graphql")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	m := make(map[string]*newrelic.EntitySearchResultsMonitor)
	var cursor *string = nil
	var c string = ""
	for {
		monitorTags, resp, err := client.SyntheticsMonitors.ListTags(context.Background(), cursor)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if cursor != nil {
			c = *cursor
		}
		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_GET_MONITORTAGS, resp.StatusCode, "cursor:"+c)

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get monitor tags, query cursor: '%s'\n", statusCode, c)
			return nil, tracker.ERR_REST_CALL_NOT_2XX
		}
		entities := monitorTags.Data.Actor.EntitySearch.Results.Entities

		for _, e := range entities {
			m[*e.MonitorId] = e
		}

		cursor = monitorTags.Data.Actor.EntitySearch.Results.NextCursor
		if cursor == nil {
			break
		}

	}
	return m, err
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
