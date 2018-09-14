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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/cmd/add"
	"github.com/IBM/newrelic-cli/cmd/delete"
	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Update monitor from a file.",
	Example: "nr update monitor -f <example.yaml>",
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
		// start to udpate
		client, err := utils.GetNewRelicClient("synthetics")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		if p.ID == nil {
			fmt.Printf("Can't find {.id} in %q.\n", file)
			os.Exit(1)
			return
		}

		resp, err := client.SyntheticsMonitors.Update(context.Background(), p, p.ID)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp.Status)
			fmt.Println(resp.StatusCode)
		}

		if *p.Type == "SCRIPT_BROWSER" || *p.Type == "SCRIPT_API" {
			var scriptTextEncoded *newrelic.Script
			scriptTextEncoded = &newrelic.Script{}

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

			if scriptTextEncoded != nil && scriptTextEncoded.ScriptText != nil {
				id := *p.ID
				resp, err := client.SyntheticsScript.UpdateByID(context.Background(), scriptTextEncoded, id)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
					return
				} else {
					fmt.Println(resp.Status)
					if resp.StatusCode >= 400 {
						os.Exit(1)
						return
					}
				}
			}
		}

		os.Exit(0)
	},
}

func UpdateMonitorByID(monitorId *string, p *newrelic.Monitor, scriptTextEncoded *newrelic.Script) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	//update monitor itself
	resp, err := client.SyntheticsMonitors.Update(context.Background(), p, monitorId)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_UPDATE_MONITOR, resp.StatusCode, "monitor id: "+(*monitorId)+",monitor name: "+(*p.Name))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Update monitor '%s', monitor id: '%s'\n", statusCode, *p.Name, *monitorId)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}

	}
	//update script if needed
	if scriptTextEncoded != nil && scriptTextEncoded.ScriptText != nil {
		id := *p.ID
		resp, err := client.SyntheticsScript.UpdateByID(context.Background(), scriptTextEncoded, id)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR_SCRIPT, err, tracker.ERR_REST_CALL, "")
			return err, ret
		} else {
			tracker.AppendRESTCallResult(client.SyntheticsScript, tracker.OPERATION_NAME_UPDATE_MONITOR, resp.StatusCode, "monitor id: "+(*monitorId)+",monitor name: "+(*p.Name))

			if resp.StatusCode >= 400 {
				var statusCode = resp.StatusCode
				fmt.Printf("Response status code: %d. Update script to monitor '%s', monitor id: '%s'\n", statusCode, *p.Name, id)
				ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
				return err, ret
			}
		}
	}
	//update labels if needed
	labelList, err, ret := get.GetLabelsByMonitorID(*p.ID)
	if ret.IsContinue == false {
		return err, ret
	}
	labelListLen := len(labelList)
	if labelListLen > 0 {
		//delete all labels on this monitor first
		for _, label := range labelList {
			err, ret := delete.DeleteLabelFromMonitor(*p.ID, *label)
			if ret.IsContinue == false {
				return err, ret
			}
		}
	}
	//and then, add new lables to this monitor
	newLabelList := p.Labels
	for _, label := range newLabelList {
		var monitorLabel *newrelic.MonitorLabel
		monitorLabel = &newrelic.MonitorLabel{}
		arr := strings.Split(*label, ":")
		monitorLabel.Category = &arr[0]
		monitorLabel.Label = &arr[1]
		err, ret := add.AddLabelToMonitor(*p.ID, monitorLabel)
		if ret.IsContinue == false {
			return err, ret
		}
	}

	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_MONITOR, nil, nil, "")

	return err, ret
}

func UpdateMonitorByName(p *newrelic.Monitor, scriptTextEncoded *newrelic.Script) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient("synthetics")
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}

	curMonitor, err, ret := get.GetMonitorByName(*p.Name)
	if ret.IsContinue == false {
		return err, ret
	}

	monitorId := curMonitor.ID

	//update monitor itself
	resp, err := client.SyntheticsMonitors.Update(context.Background(), p, monitorId)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.SyntheticsMonitors, tracker.OPERATION_NAME_UPDATE_MONITOR, resp.StatusCode, "monitor id: "+(*monitorId)+",monitor name: "+(*p.Name))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Update monitor '%s', monitor id: '%s'\n", statusCode, *p.Name, *monitorId)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}

	}
	//update script if needed
	if scriptTextEncoded != nil && scriptTextEncoded.ScriptText != nil {
		resp, err := client.SyntheticsScript.UpdateByID(context.Background(), scriptTextEncoded, *monitorId)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR_SCRIPT, err, tracker.ERR_REST_CALL, "")
			return err, ret
		} else {
			tracker.AppendRESTCallResult(client.SyntheticsScript, tracker.OPERATION_NAME_UPDATE_MONITOR, resp.StatusCode, "monitor id: "+(*monitorId)+",monitor name: "+(*p.Name))
			if resp.StatusCode >= 400 {
				var statusCode = resp.StatusCode
				fmt.Printf("Response status code: %d. Update script to monitor '%s', monitor id: '%s'\n", statusCode, *p.Name, *monitorId)
				ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_MONITOR, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
				return err, ret
			}
		}
	}
	//update labels if needed
	labelList, err, ret := get.GetLabelsByMonitorID(*monitorId)
	if ret.IsContinue == false {
		return err, ret
	}
	labelListLen := len(labelList)
	if labelListLen > 0 {
		//delete all labels on this monitor first
		for _, label := range labelList {
			err, ret := delete.DeleteLabelFromMonitor(*monitorId, *label)
			if ret.IsContinue == false {
				return err, ret
			}
		}
	}
	//and then, add new lables to this monitor
	newLabelList := p.Labels
	for _, label := range newLabelList {
		var monitorLabel *newrelic.MonitorLabel
		monitorLabel = &newrelic.MonitorLabel{}
		arr := strings.Split(*label, ":")
		monitorLabel.Category = &arr[0]
		monitorLabel.Label = &arr[1]
		err, ret := add.AddLabelToMonitor(*monitorId, monitorLabel)
		if ret.IsContinue == false {
			return err, ret
		}
	}

	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_MONITOR, nil, nil, "")

	return err, ret
}

func init() {
	UpdateCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
