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
package add

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var labelsmonitorsCmd = &cobra.Command{
	Use:     "labelsmonitors",
	Short:   "Add one lable to a specific monitor by monitor id.",
	Aliases: []string{"m"},
	Example: `nr add labelsmonitors <id> <category:label>
	* nr add labelsmonitors xxx-xxxx-xxx Env:Staging`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("length of [flags] should be 2 instead of %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		monitorId := string(args[0])
		label := string(args[1])
		var monitorLabel *newrelic.MonitorLabel
		monitorLabel = &newrelic.MonitorLabel{}
		arr := strings.Split(label, ":")
		monitorLabel.Category = &arr[0]
		monitorLabel.Label = &arr[1]
		err, returnValue := AddLabelToMonitor(monitorId, monitorLabel)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		if returnValue.IsContinue == false {
			fmt.Println(returnValue.OriginalError)
			fmt.Println(returnValue.TypicalError)
			os.Exit(1)
			return
		}

		os.Exit(0)
	},
}

func AddLabelToMonitor(monitorId string, monitorLabel *newrelic.MonitorLabel) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("labelSynthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_ADD_LABEL_MONITOR, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	resp, err := client.LabelsSynthetics.AddLabelToMonitor(context.Background(), monitorId, monitorLabel)
	var label = *monitorLabel.Category + ":" + *monitorLabel.Label
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_ADD_LABEL_MONITOR, err, tracker.ERR_REST_CALL, "label: "+", monitor id: ")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.LabelsSynthetics, tracker.OPERATION_NAME_ADD_LABEL_MONITOR, resp.StatusCode, "label: "+", monitor id: ")
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Adding label '%s' to monitor '%s'\n", statusCode, label, monitorId)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_ADD_LABEL_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_ADD_LABEL_MONITOR, nil, nil, "")
	return nil, ret
}

func init() {
	AddCmd.AddCommand(labelsmonitorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
