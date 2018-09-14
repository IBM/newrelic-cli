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
package delete

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Delete one dashboard by id.",
	Aliases: []string{"m"},
	Example: "nr delete dashboard <id>",
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
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		id, _ := strconv.ParseInt(args[0], 10, 64)
		resp, _, err := client.Dashboards.DeleteByID(context.Background(), id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		} else {
			fmt.Println(resp.StatusCode)
			if resp.StatusCode >= 400 {
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
	},
}

func DeleteDashboardByID(id int64) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_DASHBOARD_BY_ID, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}
	resp, bytes, err := client.Dashboards.DeleteByID(context.Background(), id)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_DASHBOARD_BY_ID, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		retMsg := string(bytes)
		tracker.AppendRESTCallResult(client.Dashboards, tracker.OPERATION_NAME_DELETE_DASHBOARD_BY_ID, resp.StatusCode, "dashboard id: "+strconv.FormatInt(id, 10))
		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Remove dashboard id '%d'\n", statusCode, id)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_DELETE_DASHBOARD_BY_ID, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, retMsg)
			return err, ret
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_DELETE_DASHBOARD_BY_ID, nil, nil, "")
	return nil, ret
}

func DeleteByDashboardTitle(title string) (error, tracker.ReturnValue) {
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
			dashboardId, _ := strconv.ParseInt(id.String(), 10, 64)
			err, ret := DeleteDashboardByID(dashboardId)
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

func init() {
	DeleteCmd.AddCommand(dashboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
