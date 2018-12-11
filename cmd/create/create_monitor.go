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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/cmd/add"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Create monitor from a file.",
	Example: "nr create monitor -f <example.yaml>",
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
		var p = new(newrelic.Monitor)
		err = decorder.Decode(p)
		if err != nil {
			fmt.Printf("Unable to decode %q: %v\n", file, err)
			os.Exit(1)
			return
		}
		if reflect.DeepEqual(new(newrelic.Monitor), p) {
			fmt.Printf("Error validating %q.\n", file)
			os.Exit(1)
			return
		}
		// start to create

		var scriptTextEncoded *newrelic.Script
		scriptTextEncoded = &newrelic.Script{}

		if *p.Type == "SCRIPT_BROWSER" || *p.Type == "SCRIPT_API" {

			flags := cmd.Flags()
			if flags.Lookup("script-file") != nil {
				scriptFileName, err := cmd.Flags().GetString("script-file")
				if err != nil {
					fmt.Printf("error accessing flag %s for command %s: %v\n", "script-file", cmd.Name(), err)
					os.Exit(1)
					return
				}

				if scriptFileName != "" {
					sf, err := os.Open(scriptFileName)
					defer sf.Close()
					if err != nil {
						fmt.Printf("Unable to open monitor script file '%v': %v\n", file, err)
						os.Exit(1)
						return
					}
					byteArr, err := ioutil.ReadAll(sf)
					sfContentEncoded := base64.StdEncoding.EncodeToString(byteArr)
					scriptTextEncoded.ScriptText = &sfContentEncoded
				} else {
					scriptTextEncoded = p.Script
				}
			} else {
				scriptTextEncoded = p.Script
			}
		}

		_, err, returnValue := CreateMonitor(p, scriptTextEncoded)
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

		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()

		os.Exit(0)
	},
}

func CreateMonitor(p *newrelic.Monitor, scriptTextEncoded *newrelic.Script) (string, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_MONITOR, err, tracker.ERR_CREATE_NR_CLINET, "")
		return "", err, ret
	}

	p.ID = nil

	createdMonitor, resp, err := client.SyntheticsMonitors.Create(context.Background(), p)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_MONITOR, err, tracker.ERR_REST_CALL, "")
		return "", err, ret
	} else {
		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_CREATE_MONITOR, resp.StatusCode, "monitor name :"+(*p.Name))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Creating monitor '%s'\n", statusCode, *createdMonitor.Name)
			var ret tracker.ReturnValue
			if resp.StatusCode == 400 {
				ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_400, "")
			} else {
				ret = tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			}

			return "", err, ret
		}

	}

	if *p.Type == "SCRIPT_BROWSER" || *p.Type == "SCRIPT_API" {
		if scriptTextEncoded != nil && scriptTextEncoded.ScriptText != nil {
			id := *createdMonitor.ID
			resp, err := client.SyntheticsScript.UpdateByID(context.Background(), scriptTextEncoded, id)
			if err != nil {
				fmt.Println(err)
				ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR_SCRIPT, err, tracker.ERR_REST_CALL, "")
				return id, err, ret
			} else {
				tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_UPDATE_MONITOR_SCRIPT, resp.StatusCode, "monitor name :"+(*p.Name))
				if resp.StatusCode >= 400 {
					var statusCode = resp.StatusCode
					fmt.Printf("Response status code: %d. Update script to monitor '%s', monitor id: '%s'\n", statusCode, *createdMonitor.Name, id)
					ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR_SCRIPT, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
					return "", err, ret
				}
			}
		}
	}

	//add labels
	var labelsLen = len(p.Labels)
	if labelsLen > 0 {
		for _, label := range p.Labels {
			id := createdMonitor.ID
			var monitorLabel *newrelic.MonitorLabel
			monitorLabel = &newrelic.MonitorLabel{}
			arr := strings.Split(*label, ":")
			monitorLabel.Category = &arr[0]
			monitorLabel.Label = &arr[1]
			err, returnValue := add.AddLabelToMonitor(*id, monitorLabel)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return *id, "failed to add labels to monitor", err
			// }
			if returnValue.IsContinue == false {
				return "", err, returnValue
			}
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_MONITOR, nil, nil, "")

	return *createdMonitor.ID, err, ret
}

func init() {
	CreateCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
