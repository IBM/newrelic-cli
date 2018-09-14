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
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
)

var alertschannelsCmd = &cobra.Command{
	Use:     "alertschannels",
	Short:   "Update policy/channel associations from a file.",
	Example: "nr update alertschannels -f <example.json>",
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
		var p = new(newrelic.PolicyChannelsAssociation)
		err = decorder.Decode(p)
		if err != nil {
			fmt.Printf("Unable to decode %q: %v\n", file, err)
			os.Exit(1)
			return
		}
		if reflect.DeepEqual(new(newrelic.PolicyChannelsAssociation), p) {
			fmt.Printf("Error validating %q.\n", file)
			os.Exit(1)
			return
		}

		// start to udpate
		err, _ = UpdatePolicyChannels(*p.PolicyID, p.ChannelIDList)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}

		os.Exit(0)
	},
}

func UpdatePolicyChannels(policyId int64, channelIds []*int64) (error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL, err, tracker.ERR_CREATE_NR_CLINET, "")
		return err, ret
	}

	resp, err := client.AlertsChannels.UpdatePolicyChannels(context.Background(), policyId, channelIds)
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL, err, tracker.ERR_REST_CALL, "")
		return err, ret
	} else {
		tracker.AppendRESTCallResult(client.AlertsChannels, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL, resp.StatusCode, "")

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Update policy and channels associations.'\n", statusCode)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return err, ret
		}
	}
	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL, nil, nil, "")
	return nil, ret
}

func init() {
	UpdateCmd.AddCommand(alertschannelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
