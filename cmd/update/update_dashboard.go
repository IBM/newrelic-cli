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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Update dashboard from a json file.",
	Example: "nr update dashboard -f <example.json>",
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

		//start to update
		id := gjson.Parse(fileContent).Get("id").String()
		dashboardId, _ := strconv.ParseInt(id, 10, 64)
		UpdateDashboardByID(fileContent, dashboardId)
	},
}

func UpdateDashboardByID(dashboard string, id int64) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_DASHBOARD_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	title := gjson.Parse(dashboard).Get("title").String()

	resp, bytes, err := client.Dashboards.Update(context.Background(), dashboard, id)

	var retMsg string
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_DASHBOARD_BY_ID, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		retMsg = string(bytes)
		tracker.AppendRESTCallResult(client.Dashboards, tracker.OPERATION_NAME_UPDATE_DASHBOARD_BY_ID, resp.StatusCode, "dashboard title:"+title+", dashboard id:"+strconv.FormatInt(id, 10))
		var statusCode = resp.StatusCode
		fmt.Printf("Response status code: %d. Update dashboard '%s'', dashboard id: '%d'\n", statusCode, title, id)
		if resp.StatusCode >= 400 {
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_DASHBOARD_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, retMsg)
			return err, ret
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CREATE_DASHBOARD, nil, nil, "")
	return err, ret
}

func UpdateByDashboardTitle(dashboardContent string, title string) (error, tracker.ReturnValue) {
	resultStr, err, ret := get.GetAllDashboards()
	if err != nil {
		fmt.Println(err)
		return err, ret
	} else {
		if ret.IsContinue == false {
			return err, ret
		}
	}
	dashboardArr := gjson.Parse(resultStr).Get("dashboards").Array()
	for _, existDashboard := range dashboardArr {
		t := gjson.Parse(existDashboard.String()).Get("title").String()
		if t == title {
			id := gjson.Parse(existDashboard.String()).Get("id")
			dashboardContent2, _ := sjson.Set(dashboardContent, "dashboard.id", id.Num)
			fmt.Println(dashboardContent2)
			dashboardId, _ := strconv.ParseInt(id.String(), 10, 64)
			err, ret := UpdateDashboardByID(dashboardContent2, dashboardId)
			if err != nil {
				fmt.Println(err)
				return err, ret
			} else {
				if ret.IsContinue == false {
					return err, ret
				}
			}
			return err, ret
		}
	}
	ret = tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_DASHBOARD_BY_NAME, nil, nil, "")
	return err, ret
}
