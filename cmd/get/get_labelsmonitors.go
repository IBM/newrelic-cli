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

// labelsmonitorsCmd represents the monitors command
var labelsmonitorsCmd = &cobra.Command{
	Use:   "labelsmonitors",
	Short: "Display monitors by label.",
	Example: `* nr get labelsmonitors <category:label>
* nr get labelsmonitors Category:Label -o json
* nr get labelsmonitors Category:Label -o yaml
* nr get labelsmonitors Env:Staging -o json`,
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
		label := args[0]

		monitorRefList, err, ret := GetMonitorsByLabel(label)
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
		printer.Print(monitorRefList, os.Stdout)

		os.Exit(0)
	},
}

func GetMonitorsByLabel(label string) (*newrelic.MonitorRefList, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("labelSynthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS_BY_LABEL, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	var opt *newrelic.PageLimitOptions
	opt = &newrelic.PageLimitOptions{}
	var pageSize int = 20
	var pageOffset int = 0
	var allMonitorRefList *newrelic.MonitorRefList = &newrelic.MonitorRefList{}
	for {
		opt.Limit = pageSize
		opt.Offset = pageOffset
		labelSynthetics, resp, err := client.LabelsSynthetics.GetMonitorsByLabel(context.Background(), opt, label)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS_BY_LABEL, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}
		tracker.AppendRESTCallResult(client.LabelsSynthetics, tracker.OPERATION_NAME_GET_MONITORS_BY_LABEL, resp.StatusCode, "pageSize:"+strconv.Itoa(pageSize)+",pageOffset:"+strconv.Itoa(pageOffset))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get monitors by specific label, pageSize '%d', pageOffset '%d'\n", statusCode, pageSize, pageOffset)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_MONITORS_BY_LABEL, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}

		var refsLen = len(labelSynthetics.PagedData.MonitorRefs)
		if refsLen == 0 {
			break
		} else {
			allMonitorRefList.MonitorRefs = utils.MergeMonitorReflList(allMonitorRefList.MonitorRefs, labelSynthetics.PagedData.MonitorRefs)

			pageOffset = pageOffset + pageSize
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_MONITORS_BY_LABEL, nil, nil, "")
	return allMonitorRefList, err, ret
}

func GetLabelsByMonitorID(monitorId string) ([]*string, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_LABELS_BY_MONITOR_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}
	monitor, _, err := client.SyntheticsMonitors.GetByID(context.Background(), monitorId)

	var labels []*string

	//get all labels
	lablesArray, err, returnValue := GetLabels()
	if returnValue.IsContinue == false {
		return nil, err, returnValue
	}

	for index, _ := range lablesArray.Labels {
		l := lablesArray.Labels[index]
		key := fmt.Sprintf("%v:%v", *l.Category, *l.Name)

		labelSynthetics, err, returnValue := GetMonitorsByLabel(key)
		if returnValue.IsContinue == false {
			return nil, err, returnValue
		}

		// if err != nil {
		// 	fmt.Println(err)
		// 	return nil, err, returnValue
		// }

		monitorRefList := labelSynthetics.MonitorRefs
		var refListLen = len(monitorRefList)
		// fmt.Println(refListLen)

		if refListLen > 0 {
			//the label was added to monitors
			for _, ref := range monitorRefList {
				//get monitor id
				mId := *ref.ID
				if mId == (*monitor.ID) {
					labLen := len(monitor.Labels)
					var newLabelList = make([]*string, (labLen + 1))
					for i := 0; i < labLen; i++ {
						newLabelList[i] = monitor.Labels[i]
					}
					newLabelList[labLen] = &key
					labels = utils.MergeStringList(labels, newLabelList)
				}
			}
		} else {
			//the label was not used by any monitors, skip
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_LABELS_BY_MONITOR_ID, nil, nil, "")

	return labels, err, ret
}

func init() {
	GetCmd.AddCommand(labelsmonitorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:s
}
