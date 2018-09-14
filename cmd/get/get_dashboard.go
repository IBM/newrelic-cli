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

	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

// dashboardsCmd represents the dashboards command
var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Display one single dashboard by id.",
	Example: `* nr get dashboard <id>`,
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

		id, _ := strconv.ParseInt(args[0], 10, 64)

		result, err, ret := GetDashboardByID(id)
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

		jsonStr := pretty.Pretty([]byte(result))
		fmt.Printf("%s\n", string(jsonStr))

		os.Exit(0)
	},
}

func GetDashboardByID(id int64) (string, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARD_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return "", err, ret
	}

	resp, bytes, err := client.Dashboards.GetByID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARD_BY_ID, err, tracker.ERR_REST_CALL, "")
		return "", err, ret
	}
	tracker.AppendRESTCallResult(client.Labels, tracker.OPERATION_NAME_GET_DASHBOARD_BY_ID, resp.StatusCode, "id:"+strconv.FormatInt(id, 10))

	if resp.StatusCode >= 400 {
		var statusCode = resp.StatusCode
		fmt.Printf("Response status code: %d. Get one specific dashboard, id '%d'\n", statusCode, id)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_DASHBOARD_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
		return "", err, ret
	}

	strContent := string(bytes)

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_DASHBOARD_BY_ID, nil, nil, "")

	return strContent, nil, ret
}

func IsDashboardTitleExists(dashboardTitle string) (bool, string, error, tracker.ReturnValue) {
	resultStr, err, returnValue := GetAllDashboards()
	if returnValue.IsContinue == false {
		return false, "", err, returnValue
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CHECK_DASHBOARD_TITLE_EXISTS, nil, nil, "")

	dashboardArr := gjson.Parse(resultStr).Get("dashboards").Array()
	for _, dashboard := range dashboardArr {

		titleResult := gjson.Parse(dashboard.String()).Get("title")
		title := titleResult.String()
		if dashboardTitle == title {
			return true, dashboard.String(), err, ret
		}
	}

	return false, "", err, ret
}

func init() {
	GetCmd.AddCommand(dashboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dashboardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
