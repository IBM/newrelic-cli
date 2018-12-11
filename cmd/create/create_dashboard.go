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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Create dashboard from a json file.",
	Example: "nr create dashboard -f <example.yaml>",
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := utils.GetArg(cmd, "file")
		if err != nil {
			fmt.Printf("Unable to get argument 'file': %v\n", err)
			os.Exit(1)
			return
		}
		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
			return
		}
		fileContent := string(bytes)
		if !gjson.Valid(fileContent) {
			fmt.Printf("Incorrect JSON format: %v.\n", errors.New("invalid json"))
			os.Exit(1)
			return
		}
		// start to create
		retMsg, err, ret := CreateDashboard(fileContent)
		fmt.Printf("response message: %s\n", retMsg)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
			return
		}
		if ret.IsContinue == false {
			fmt.Print(ret.OriginalError)
			os.Exit(1)
			return
		}

		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()

		os.Exit(0)
	},
}

func CreateDashboard(dashboard string) (string, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_DASHBOARD, err, tracker.ERR_CREATE_NR_CLINET, "")
		return "", err, ret
	}
	title := gjson.Parse(dashboard).Get("dashboard.title").String()

	resp, bytes, err := client.Dashboards.Create(context.Background(), dashboard)

	var retMsg string
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_DASHBOARD, err, tracker.ERR_REST_CALL, "")
		return "", err, ret
	} else {
		retMsg = string(bytes)
		tracker.AppendRESTCallResult(client.Dashboards, tracker.OPERATION_NAME_CREATE_DASHBOARD, resp.StatusCode, "dashboard title:"+title)
		var statusCode = resp.StatusCode
		fmt.Printf("Response status code: %d. Create dashboard '%s''\n", statusCode, title)
		if resp.StatusCode >= 400 {
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_CREATE_DASHBOARD, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, retMsg)
			return retMsg, err, ret
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_DASHBOARD, nil, nil, "")
	return retMsg, err, ret
}

func init() {
	CreateCmd.AddCommand(dashboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertschannelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertschannelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
