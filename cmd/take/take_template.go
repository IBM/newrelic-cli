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
package take

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Take JSON templdate file for creating target by template type name.",
	Example: "nr take template [ monitor_simple | monitor_script_inline | alertspolicies | " +
		"dashboard | alertsconditions_infra | alertsconditions_nrql | alertsconditions_synthetics | " +
		"alertsconditions_plugin | alertsconditions_ext | alertsconditions_apm | " +
		"alertschannels_campfire | alertschannels_email | alertschannels_hipchat | " +
		"alertschannels_opsgenie | alertschannels_pagerduty | alertschannels_victorops | " +
		"alertschannels_webhook_json | alertschannels_webhook_form ]",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			var err = fmt.Errorf("length of [flags] should be 1 instead of %d", len(args))
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		var templateName = args[0]
		if templateName == "" {
			var err = fmt.Errorf("Please provide template type name.\n")
			fmt.Println(err)
			os.Exit(1)
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var templateName = args[0]
		templateContent := GetTemplateByName(templateName)
		fmt.Println(templateContent)
		os.Exit(0)
	},
}

func GetTemplateByName(name string) string {
	return name + " file content"
}

func init() {
	TakeCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	templateCmd.PersistentFlags().Arg(1)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
