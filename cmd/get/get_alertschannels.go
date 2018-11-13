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

var alertschannelsCmd = &cobra.Command{
	Use:     "alertschannels",
	Short:   "Display all alerts_channels.",
	Aliases: []string{"ac", "alertchannel", "alertschannel"},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		alertsChannelList, resp, err := client.AlertsChannels.ListAll(context.Background(), nil)
		if err != nil || resp.StatusCode >= 400 {
			fmt.Printf("%v:%v", resp.Status, err)
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
		printer.Print(alertsChannelList, os.Stdout)

		os.Exit(0)
	},
}

func GetAllAlertsChannels() (*newrelic.AlertsChannelList, error, tracker.ReturnValue) {
	var opt *newrelic.AlertsChannelListOptions = &newrelic.AlertsChannelListOptions{}

	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_CHANNELS, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	var allChannelList *newrelic.AlertsChannelList = &newrelic.AlertsChannelList{}

	var pageCount = 1
	for {
		opt.Page = pageCount
		alertsChannelList, resp, err := client.AlertsChannels.ListAll(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_CHANNELS, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		tracker.AppendRESTCallResult(client.AlertsChannels, tracker.OPERATION_NAME_GET_ALERT_CHANNELS, resp.StatusCode, "pageCount:"+strconv.Itoa(pageCount))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get alert channels, pageCount '%d'\n", statusCode, pageCount)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_ALERT_CHANNELS, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}

		var channelListLen = len(alertsChannelList.AlertsChannels)
		if channelListLen == 0 {
			break
		} else {
			allChannelList.AlertsChannels = utils.MergeAlertChannelList(allChannelList.AlertsChannels, alertsChannelList.AlertsChannels)
			pageCount++
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_ALERT_CHANNELS, nil, nil, "")

	return allChannelList, err, ret
}

func IsChannelNameExists(channelName string) (bool, *newrelic.AlertsChannel, error, tracker.ReturnValue) {
	list, err, returnValue := GetAllAlertsChannels()
	if returnValue.IsContinue == false {
		return false, nil, err, returnValue
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_CHECK_ALERT_CHANNEL_NAME_EXISTS, nil, nil, "")
	for _, channel := range list.AlertsChannels {
		if *channel.Name == channelName {
			return true, channel, nil, ret
		}
	}
	return false, nil, err, ret
}

func init() {
	GetCmd.AddCommand(alertschannelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertschannelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertschannelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
