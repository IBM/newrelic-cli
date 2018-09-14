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
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

// dashboardsCmd represents the dashboards command
var dashboardsCmd = &cobra.Command{
	Use:     "dashboards",
	Short:   "Display all dashboards.",
	Example: `* nr get dashboards`,
	Run: func(cmd *cobra.Command, args []string) {

		resultStr, err, ret := GetAllDashboards()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		if ret.IsContinue == false {
			fmt.Println(ret.OriginalError)
			os.Exit(1)
			return
		}

		result := pretty.Pretty([]byte(resultStr))
		fmt.Printf("%s\n", string(result))

		os.Exit(0)
	},
}

func GetAllDashboards() (string, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARDS, err, tracker.ERR_CREATE_NR_CLINET, "")
		return "", err, ret
	}

	var dashboardJson = `{"dashboards":[]}`

	var pageCount = 1
	var opt *newrelic.DashboardListOptions = &newrelic.DashboardListOptions{}

	for {
		opt.Page = pageCount

		resp, bytes, err := client.Dashboards.ListAll(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARDS, err, tracker.ERR_REST_CALL, "")
			return "", err, ret
		}

		tracker.AppendRESTCallResult(client.Labels, tracker.OPERATION_NAME_GET_DASHBOARDS, resp.StatusCode, "pageCount:"+strconv.Itoa(pageCount))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get one page dashboards, pageCount '%d'\n", statusCode, pageCount)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARDS, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return "", err, ret
		}

		strContent := string(bytes)
		dashboardArr := gjson.Parse(strContent).Get("dashboards").Array()
		var lenDashboards = len(dashboardArr)

		if lenDashboards == 0 {
			break
		} else {
			for _, dashboard := range dashboardArr {
				if dashboard.String() != "" {
					dashboardJson, _ = sjson.Set(dashboardJson, "dashboards.-1", dashboard.String())
				}
			}
		}
		pageCount++
	}

	var resultStr = ""
	var dashboardListSize = 0
	dashboardArr := gjson.Parse(dashboardJson).Get("dashboards").Array()
	var lenArr = len(dashboardArr)

	resultStr = resultStr + `{"dashboards":[`
	for index, dashboard := range dashboardArr {
		if index < (lenArr - 1) {
			resultStr = resultStr + dashboard.String() + ","
		} else {
			resultStr = resultStr + dashboard.String()
		}

		// fmt.Println(dashboard.String())
		// fmt.Println()
		// id := gjson.Parse(dashboard.String()).Get("id")
		// title := gjson.Parse(dashboard.String()).Get("title")
		// fmt.Println(id)
		// fmt.Println(title)
		// fmt.Println()
		// fmt.Println()

		dashboardListSize++
	}
	resultStr = resultStr + `]}`

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_DASHBOARDS, nil, nil, "")

	return resultStr, nil, ret
}

func init() {
	GetCmd.AddCommand(dashboardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dashboardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
