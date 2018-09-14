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
package patch

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/utils"
)

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Short:   "Patch monitor from a file.",
	Example: "nr patch monitor -f <example.yaml>",
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

		resp, err := client.SyntheticsMonitors.Patch(context.Background(), p, p.ID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		} else {
			var statusCode = resp.StatusCode
			fmt.Printf("Response status code: %d. Patch monitor '%s', monitor id: '%s'\n", statusCode, *p.Name, *p.ID)
			if resp.StatusCode >= 400 {
				os.Exit(1)
				return
			}
		}

		os.Exit(0)
	},
}

func init() {
	PatchCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alertspoliciesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertspoliciesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
