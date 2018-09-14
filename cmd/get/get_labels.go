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

// labelsCmd represents the monitors command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Display all labels.",
	Example: `* nr get labels
* nr get labels -o json
* nr get labels -o yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		labelArray, err, _ := GetLabels()
		if err != nil {
			fmt.Println(err)
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
		printer.Print(labelArray, os.Stdout)

		os.Exit(0)
	},
}

func GetLabels() (*newrelic.LabelList, error, tracker.ReturnValue) {
	client, err := utils.GetNewRelicClient()
	if err != nil {
		fmt.Println(err)
		ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_LABELS, err, tracker.ERR_CREATE_NR_CLINET, "")
		return nil, err, ret
	}

	var opt *newrelic.LabelListOptions
	opt = &newrelic.LabelListOptions{}

	var allLabelList *newrelic.LabelList = &newrelic.LabelList{}

	var pageCount = 1
	for {
		opt.Page = pageCount
		labelList, resp, err := client.Labels.ListAll(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_LABELS, err, tracker.ERR_REST_CALL, "")
			return nil, err, ret
		}

		tracker.AppendRESTCallResult(client.Labels, tracker.OPERATION_NAME_GET_LABELS, resp.StatusCode, "pageCount:"+strconv.Itoa(pageCount))

		if resp.StatusCode >= 400 {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Get one page labels, pageCount '%d'\n", statusCode, pageCount)
			ret := tracker.ToReturnValue(false, tracker.OPERATION_NAME_GET_LABELS, tracker.ERR_REST_CALL_NOT_2XX, tracker.ERR_REST_CALL_NOT_2XX, "")
			return nil, err, ret
		}
		var labelListLen = len(labelList.Labels)

		if labelListLen == 0 {
			break
		} else {
			allLabelList.Labels = utils.MergeLabelList(allLabelList.Labels, labelList.Labels)

			pageCount++
		}
	}

	ret := tracker.ToReturnValue(true, tracker.OPERATION_NAME_GET_LABELS, nil, nil, "")

	return allLabelList, err, ret
}

func init() {
	GetCmd.AddCommand(labelsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:s
}
