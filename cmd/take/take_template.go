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
	"encoding/base64"
	"fmt"
	"io/ioutil"
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

		fileName, templateContent, err := getTemplateByName(templateName)
		if fileName == "unknown" {
			var err = fmt.Errorf("Please provide correct template type name.\n")
			fmt.Println(err)
			os.Exit(1)
			return
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writeTemplateToDisk(fileName, templateContent)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Failed to write to disk and generate template file.")
			os.Exit(1)
		} else {
			fmt.Println()
			fmt.Printf("%s, template file generated.", fileName)
			fmt.Println()
			fmt.Println()
		}

		fmt.Println(">>>>template>>>>")
		fmt.Println(templateContent)
		fmt.Println("<<<<template<<<<")

		os.Exit(0)
	},
}

func getTemplateByName(name string) (string, string, error) {
	fileName := ""
	templateContentEncoded := ""
	switch name {
	case "monitor_simple":
		fileName, templateContentEncoded = getTemplateFileContent_monitor_simple()
	case "monitor_script_inline":
		fileName, templateContentEncoded = getTemplateFileContent_monitor_script_inline()
	case "alertspolicies":
		fileName, templateContentEncoded = getTemplateFileContent_alertspolicies()
	case "dashboard":
		fileName, templateContentEncoded = getTemplateFileContent_dashboard()
	case "alertsconditions_infra":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_infra()
	case "alertsconditions_nrql":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_nrql()
	case "alertsconditions_synthetics":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_synthetics()
	case "alertsconditions_plugin":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_plugin()
	case "alertsconditions_ext":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_ext()
	case "alertsconditions_apm":
		fileName, templateContentEncoded = getTemplateFileContent_alertsconditions_apm()
	case "alertschannels_campfire":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_campfire()
	case "alertschannels_email":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_email()
	case "alertschannels_hipchat":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_hipchat()
	case "alertschannels_opsgenie":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_opsgenie()
	case "alertschannels_pagerduty":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_pagerduty()
	case "alertschannels_victorops":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_victorops()
	case "alertschannels_webhook_json":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_webhook_json()
	case "alertschannels_webhook_form":
		fileName, templateContentEncoded = getTemplateFileContent_alertschannels_webhook_form()
	default:
		return "unknown", "unknown", nil
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(templateContentEncoded)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return "", "", err
	}

	templateContent := string(decodeBytes)

	return fileName, templateContent, nil
}

func writeTemplateToDisk(fileName string, templateContent string) error {
	err := ioutil.WriteFile(fileName, []byte(templateContent), 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getTemplateFileContent_monitor_simple() (string, string) {
	fileName := "monitor_simple.json"
	templateContent := "ewogICAgIm5hbWUiOiAiJHttb25pdG9yX25hbWV9IiwKICAgICJ0eXBlIjogIlNJTVBMRSIsCiAgICAiZnJlcXVlbmN5IjogMTAsCiAgICAidXJpIjogIiR7dXJsfSIsCiAgICAibG9jYXRpb25zIjogWwogICAgICAiQVdTX1VTX1dFU1RfMSIKICAgIF0sCiAgICAic3RhdHVzIjogIkVOQUJMRUQiLAogICAgInNsYVRocmVzaG9sZCI6IDcsCiAgICAidXNlcklkIjogMCwKICAgICJvcHRpb25zIjoge30KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_monitor_script_inline() (string, string) {
	fileName := "monitor_script_inline.json"
	templateContent := "ewogICJuYW1lIjogIiR7bW9uaXRvcl9uYW1lfSIsCiAgInR5cGUiOiAiU0NSSVBUX0JST1dTRVIiLAogICJmcmVxdWVuY3kiOiAxNDQwLAogICJ1cmkiOiAiIiwKICAibG9jYXRpb25zIjogWwogICAgIkFXU19VU19XRVNUXzEiCiAgXSwKICAic3RhdHVzIjogIkVOQUJMRUQiLAogICJzbGFUaHJlc2hvbGQiOiA3LAogICJ1c2VySWQiOiAwLAogICJvcHRpb25zIjoge30sCiAgInNjcmlwdCI6IHsKICAgICJzY3JpcHRUZXh0IjogIiR7YmFzZTY0X2VuY29kZWRfY29udGVudH0iCiAgfQp9Cg=="
	return fileName, templateContent
}

func getTemplateFileContent_alertspolicies() (string, string) {
	fileName := "alertspolicies.json"
	templateContent := "ewogICAgInBvbGljeSI6ewogICAgICAgICJpbmNpZGVudF9wcmVmZXJlbmNlIjogIlBFUl9QT0xJQ1kiLAogICAgICAgICJuYW1lIjogIiR7cG9saWN5X25hbWV9IgogICAgfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_dashboard() (string, string) {
	fileName := "dashboard.json"
	templateContent := "ewogICJkYXNoYm9hcmQiOiB7CiAgICAidGl0bGUiOiAiJHtkYXNoYm9hcmRfbmFtZX0iLCAKICAgICJkZXNjcmlwdGlvbiI6IG51bGwsIAogICAgImljb24iOiAiYmFyLWNoYXJ0IiwKICAgICJ2aXNpYmlsaXR5IjogImFsbCIsIAogICAgImVkaXRhYmxlIjogImVkaXRhYmxlX2J5X2FsbCIsCiAgICAib3duZXJfZW1haWwiOiAiJHtvd25lcl9lbWFpbH0iLCAKICAgICJtZXRhZGF0YSI6IHsKICAgICAgInZlcnNpb24iOiAxCiAgICB9LCAKICAgICJ3aWRnZXRzIjogWwogICAgICB7CiAgICAgICAgInZpc3VhbGl6YXRpb24iOiAiY29tcGFyaXNvbl9saW5lX2NoYXJ0IiwgCiAgICAgICAgImxheW91dCI6IHsKICAgICAgICAgICJ3aWR0aCI6IDEsIAogICAgICAgICAgImhlaWdodCI6IDEsIAogICAgICAgICAgInJvdyI6IDEsIAogICAgICAgICAgImNvbHVtbiI6IDEKICAgICAgICB9LAogICAgICAgICJkYXRhIjogWwogICAgICAgICAgewogICAgICAgICAgICAibnJxbCI6ICIke25ycWx9IgogICAgICAgICAgfQogICAgICAgIF0sIAogICAgICAgICJwcmVzZW50YXRpb24iOiB7CiAgICAgICAgICAidGl0bGUiOiAiJHt3aWRnZXRfdGl0bGV9IiwgCiAgICAgICAgICAibm90ZXMiOiBudWxsCiAgICAgICAgfQogICAgICB9CiAgICBdLCAKICAgICJmaWx0ZXIiOiBudWxsCiAgfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_infra() (string, string) {
	fileName := "alertsconditions_infra.json"
	templateContent := "ewoJImRhdGEiOiB7CgkJImNvbXBhcmlzb24iOiAiYWJvdmUiLAoJCSJjcml0aWNhbF90aHJlc2hvbGQiOiB7CgkJCSJkdXJhdGlvbl9taW51dGVzIjogMywKCQkJInRpbWVfZnVuY3Rpb24iOiAiYWxsIiwKCQkJInZhbHVlIjogOTAKCQl9LAoJCSJlbmFibGVkIjogdHJ1ZSwKCQkiZXZlbnRfdHlwZSI6ICIke2V2ZW50X3R5cGV9IiwKCQkibmFtZSI6ICIke25hbWV9IiwKCQkic2VsZWN0X3ZhbHVlIjogImNwdVBlcmNlbnQiLAoJCSJ0eXBlIjogIiR7dHlwZX0iCgl9Cn0="
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_nrql() (string, string) {
	fileName := "alertsconditions_nrql.json"
	templateContent := "ewogICJucnFsX2NvbmRpdGlvbiI6IHsKICAgICJuYW1lIjogInN0cmluZyIsCiAgICAicnVuYm9va191cmwiOiAic3RyaW5nIiwKICAgICJlbmFibGVkIjogImJvb2xlYW4iLAogICAgImV4cGVjdGVkX2dyb3VwcyI6ICJpbnRlZ2VyIiwKICAgICJpZ25vcmVfb3ZlcmxhcCI6ICJib29sZWFuIiwKICAgICJ2YWx1ZV9mdW5jdGlvbiI6ICJzdHJpbmciLAogICAgInRlcm1zIjogWwogICAgICB7CiAgICAgICAgImR1cmF0aW9uIjogInN0cmluZyIsCiAgICAgICAgIm9wZXJhdG9yIjogInN0cmluZyIsCiAgICAgICAgInByaW9yaXR5IjogInN0cmluZyIsCiAgICAgICAgInRocmVzaG9sZCI6ICJzdHJpbmciLAogICAgICAgICJ0aW1lX2Z1bmN0aW9uIjogInN0cmluZyIKICAgICAgfQogICAgXSwKICAgICJucnFsIjogewogICAgICAicXVlcnkiOiAic3RyaW5nIiwKICAgICAgInNpbmNlX3ZhbHVlIjogInN0cmluZyIKICAgIH0KICB9Cn0="
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_synthetics() (string, string) {
	fileName := "alertsconditions_synthetics.json"
	templateContent := "ewoJInN5bnRoZXRpY3NfY29uZGl0aW9uIjogewoJICAibmFtZSI6ICJzdHJpbmciLAoJICAibW9uaXRvcl9pZCI6ICJzdHJpbmciLAoJICAicnVuYm9va191cmwiOiAic3RyaW5nIiwKCSAgImVuYWJsZWQiOiAiYm9vbGVhbiIKCX0KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_plugin() (string, string) {
	fileName := "alertsconditions_plugin.json"
	templateContent := "ewogICJwbHVnaW5zX2NvbmRpdGlvbiI6IHsKICAgICJuYW1lIjogInN0cmluZyIsCiAgICAiZW5hYmxlZCI6ICJib29sZWFuIiwKICAgICJlbnRpdGllcyI6IFsKICAgICAgImludGVnZXIiCiAgICBdLAogICAgIm1ldHJpY19kZXNjcmlwdGlvbiI6ICJzdHJpbmciLAogICAgIm1ldHJpYyI6ICJzdHJpbmciLAogICAgInZhbHVlX2Z1bmN0aW9uIjogInN0cmluZyIsCiAgICAicnVuYm9va191cmwiOiAic3RyaW5nIiwKICAgICJ0ZXJtcyI6IFsKICAgICAgewogICAgICAgICJkdXJhdGlvbiI6ICJzdHJpbmciLAogICAgICAgICJvcGVyYXRvciI6ICJzdHJpbmciLAogICAgICAgICJwcmlvcml0eSI6ICJzdHJpbmciLAogICAgICAgICJ0aHJlc2hvbGQiOiAic3RyaW5nIiwKICAgICAgICAidGltZV9mdW5jdGlvbiI6ICJzdHJpbmciCiAgICAgIH0KICAgIF0sCiAgICAicGx1Z2luIjogewogICAgICAiaWQiOiAic3RyaW5nIiwKICAgICAgImd1aWQiOiAic3RyaW5nIgogICAgfQogIH0KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_ext() (string, string) {
	fileName := "alertsconditions_ext.json"
	templateContent := "ewogICJleHRlcm5hbF9zZXJ2aWNlX2NvbmRpdGlvbiI6IHsKICAgICJ0eXBlIjogInN0cmluZyIsCiAgICAibmFtZSI6ICJzdHJpbmciLAogICAgImVuYWJsZWQiOiAiYm9vbGVhbiIsCiAgICAiZW50aXRpZXMiOiBbCiAgICAgICJpbnRlZ2VyIgogICAgXSwKICAgICJleHRlcm5hbF9zZXJ2aWNlX3VybCI6ICJzdHJpbmciLAogICAgIm1ldHJpYyI6ICJzdHJpbmciLAogICAgInJ1bmJvb2tfdXJsIjogInN0cmluZyIsCiAgICAidGVybXMiOiBbCiAgICAgIHsKICAgICAgICAiZHVyYXRpb24iOiAic3RyaW5nIiwKICAgICAgICAib3BlcmF0b3IiOiAic3RyaW5nIiwKICAgICAgICAicHJpb3JpdHkiOiAic3RyaW5nIiwKICAgICAgICAidGhyZXNob2xkIjogInN0cmluZyIsCiAgICAgICAgInRpbWVfZnVuY3Rpb24iOiAic3RyaW5nIgogICAgICB9CiAgICBdCiAgfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertsconditions_apm() (string, string) {
	fileName := "alertsconditions_apm.json"
	templateContent := "ewogICJjb25kaXRpb24iOiB7CiAgICAidHlwZSI6ICJzdHJpbmciLAogICAgIm5hbWUiOiAic3RyaW5nIiwKICAgICJlbmFibGVkIjogImJvb2xlYW4iLAogICAgImVudGl0aWVzIjogWwogICAgICAiaW50ZWdlciIKICAgIF0sCiAgICAibWV0cmljIjogInN0cmluZyIsCiAgICAiZ2NfbWV0cmljIjogInN0cmluZyIsCiAgICAiY29uZGl0aW9uX3Njb3BlIjogInN0cmluZyIsCiAgICAidmlvbGF0aW9uX2Nsb3NlX3RpbWVyIjogImludGVnZXIiLAogICAgInRlcm1zIjogWwogICAgICB7CiAgICAgICAgImR1cmF0aW9uIjogInN0cmluZyIsCiAgICAgICAgIm9wZXJhdG9yIjogInN0cmluZyIsCiAgICAgICAgInByaW9yaXR5IjogInN0cmluZyIsCiAgICAgICAgInRocmVzaG9sZCI6ICJzdHJpbmciLAogICAgICAgICJ0aW1lX2Z1bmN0aW9uIjogInN0cmluZyIKICAgICAgfQogICAgXSwKICAgICJ1c2VyX2RlZmluZWQiOiB7CiAgICAgICJtZXRyaWMiOiAic3RyaW5nIiwKICAgICAgInZhbHVlX2Z1bmN0aW9uIjogInN0cmluZyIKICAgIH0KICB9Cn0="
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_campfire() (string, string) {
	fileName := "alertschannels_campfire.json"
	templateContent := "ewoJImNoYW5uZWwiOiB7CgkJImNvbmZpZ3VyYXRpb24iOiB7CgkJCSJyb29tIjogIiR7cm9vbX0iLAogICAgICAgICAgICAic3ViZG9tYWluIjogIiR7c3ViZG9tYWlufSIsCiAgICAgICAgICAgICJ0b2tlbiI6ICIke3Rva2VufSIKCQl9LAoJCSJuYW1lIjogIiR7bmFtZX0iLAoJCSJ0eXBlIjogImNhbXBmaXJlIgoJfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_email() (string, string) {
	fileName := "alertschannels_email.json"
	templateContent := "ewogICAgImNoYW5uZWwiOnsKICAgICAgICAiY29uZmlndXJhdGlvbiI6IHsKICAgICAgICAgICAgInJlY2lwaWVudHMiIDogIiR7ZW1haWx9IiwKICAgICAgICAgICAgImluY2x1ZGVfanNvbl9hdHRhY2htZW50IiA6IHRydWUKICAgICAgICB9LAogICAgICAgICJuYW1lIjogIiR7bmFtZX0iLAogICAgICAgICJ0eXBlIjogImVtYWlsIgogICAgfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_hipchat() (string, string) {
	fileName := "alertschannels_hipchat.json"
	templateContent := "ewogICAgImNoYW5uZWwiOnsKICAgICAgICAiY29uZmlndXJhdGlvbiI6IHsKICAgICAgICAgICAgImF1dGhfdG9rZW4iOiAiJHt0b2tlbn0iLAogICAgICAgICAgICAicm9vbV9pZCI6ICIke3Jvb21faWR9IgogICAgICAgIH0sCiAgICAgICAgIm5hbWUiOiAiJHtuYW1lfSIsCiAgICAgICAgInR5cGUiOiAiaGlwY2hhdCIKICAgIH0KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_opsgenie() (string, string) {
	fileName := "alertschannels_opsgenie.json"
	templateContent := "ewogICAgImNoYW5uZWwiOnsKICAgICAgICAiY29uZmlndXJhdGlvbiI6IHsKICAgICAgICAgICAgImFwaV9rZXkiOiAiJHthcGlfa2V5fSIsCiAgICAgICAgICAgICJ0ZWFtcyI6ICIke3RlYW1zfSIsCiAgICAgICAgICAgICJ0YWdzIjogIiR7dGFnfSIsCiAgICAgICAgICAgICJyZWNpcGllbnRzIjogIiR7cmVjaXBpZW50c30iCiAgICAgICAgfSwKICAgICAgICAibmFtZSI6ICIke25hbWV9IiwKICAgICAgICAidHlwZSI6ICJvcHNnZW5pZSIKICAgIH0KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_pagerduty() (string, string) {
	fileName := "alertschannels_pagerduty.json"
	templateContent := "ewogICAgImNoYW5uZWwiOnsKICAgICAgICAiY29uZmlndXJhdGlvbiI6IHsKICAgICAgICAgICAgInNlcnZpY2Vfa2V5IjogIiR7c2VydmljZV9rZXl9IgogICAgICAgIH0sCiAgICAgICAgIm5hbWUiOiAiJHtuYW1lfSIsCiAgICAgICAgInR5cGUiOiAicGFnZXJkdXR5IgogICAgfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_victorops() (string, string) {
	fileName := "alertschannels_victorops.json"
	templateContent := "ewogICAgImNoYW5uZWwiOnsKICAgICAgICAiY29uZmlndXJhdGlvbiI6IHsKICAgICAgICAgICAgImtleSI6ICIke2tleX0iLAogICAgICAgICAgICAicm91dGVfa2V5IjogIiR7cm91dGVfa2V5fSIKICAgICAgICB9LAogICAgICAgICJuYW1lIjogIiR7bmFtZX0iLAogICAgICAgICJ0eXBlIjogInZpY3Rvcm9wcyIKICAgIH0KfQ=="
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_webhook_json() (string, string) {
	fileName := "alertschannels_webhook_json.json"
	templateContent := "ewoJImNoYW5uZWwiOiB7CgkJImNvbmZpZ3VyYXRpb24iOiB7CgkJCSJhdXRoX3Bhc3N3b3JkIjogIiR7YXV0aF9wYXNzd29yZH0iLAoJCQkiYXV0aF91c2VybmFtZSI6ICIke2F1dGhfdXNlcm5hbWV9IiwKCQkJImJhc2VfdXJsIjogIiR7YmFzZV91cmx9IiwKCQkJImhlYWRlcnMiOiB7CgkJCQkiaGVhZGVyMSI6ICIke2hlYWRlcjF9IgoJCQl9LAoJCQkicGF5bG9hZCI6IHsKCQkJCSJhY2NvdW50X2lkIjogIiRBQ0NPVU5UX0lEIiwKCQkJCSJhY2NvdW50X25hbWUiOiAiJEFDQ09VTlRfTkFNRSIsCgkJCQkiY29uZGl0aW9uX2lkIjogIiRDT05ESVRJT05fSUQiLAoJCQkJImNvbmRpdGlvbl9uYW1lIjogIiRDT05ESVRJT05fTkFNRSIsCgkJCQkiY3VycmVudF9zdGF0ZSI6ICIkRVZFTlRfU1RBVEUiLAoJCQkJImRldGFpbHMiOiAiJEVWRU5UX0RFVEFJTFMiLAoJCQkJImV2ZW50X3R5cGUiOiAiJEVWRU5UX1RZUEUiLAoJCQkJImluY2lkZW50X2Fja25vd2xlZGdlX3VybCI6ICIkSU5DSURFTlRfQUNLTk9XTEVER0VfVVJMIiwKCQkJCSJpbmNpZGVudF9pZCI6ICIkSU5DSURFTlRfSUQiLAoJCQkJImluY2lkZW50X3VybCI6ICIkSU5DSURFTlRfVVJMIiwKCQkJCSJvd25lciI6ICIkRVZFTlRfT1dORVIiLAoJCQkJInBvbGljeV9uYW1lIjogIiRQT0xJQ1lfTkFNRSIsCgkJCQkicG9saWN5X3VybCI6ICIkUE9MSUNZX1VSTCIsCgkJCQkicnVuYm9va191cmwiOiAiJFJVTkJPT0tfVVJMIiwKCQkJCSJzZXZlcml0eSI6ICIkU0VWRVJJVFkiLAoJCQkJInRhcmdldHMiOiAiJFRBUkdFVFMiLAoJCQkJInRpbWVzdGFtcCI6ICIkVElNRVNUQU1QIiwKCQkJCSJ2aW9sYXRpb25fY2hhcnRfdXJsIjogIiRWSU9MQVRJT05fQ0hBUlRfVVJMIgoJCQl9LAoJCQkicGF5bG9hZF90eXBlIjogImFwcGxpY2F0aW9uL2pzb24iCgkJfSwKCQkibmFtZSI6ICIke25hbWV9IiwKCQkidHlwZSI6ICJ3ZWJob29rIgoJfQp9"
	return fileName, templateContent
}

func getTemplateFileContent_alertschannels_webhook_form() (string, string) {
	fileName := "alertschannels_webhook_form.json"
	templateContent := "ewoJImNoYW5uZWwiOiB7CgkJImNvbmZpZ3VyYXRpb24iOiB7CgkJCSJhdXRoX3Bhc3N3b3JkIjogIiR7YXV0aF9wYXNzd29yZH0iLAoJCQkiYXV0aF91c2VybmFtZSI6ICIke2F1dGhfdXNlcm5hbWV9IiwKCQkJImJhc2VfdXJsIjogIiR7YmFzZV91cmx9IiwKCQkJImhlYWRlcnMiOiB7CgkJCQkiaGVhZGVyMSI6ICIke2hlYWRlcjF9IgoJCQl9LAoJCQkicGF5bG9hZCI6IHsKCQkJCSJhY2NvdW50X2lkIjogIiRBQ0NPVU5UX0lEIiwKCQkJCSJhY2NvdW50X25hbWUiOiAiJEFDQ09VTlRfTkFNRSIsCgkJCQkiY29uZGl0aW9uX2lkIjogIiRDT05ESVRJT05fSUQiLAoJCQkJImNvbmRpdGlvbl9uYW1lIjogIiRDT05ESVRJT05fTkFNRSIsCgkJCQkiY3VycmVudF9zdGF0ZSI6ICIkRVZFTlRfU1RBVEUiLAoJCQkJImRldGFpbHMiOiAiJEVWRU5UX0RFVEFJTFMiLAoJCQkJImV2ZW50X3R5cGUiOiAiJEVWRU5UX1RZUEUiLAoJCQkJImluY2lkZW50X2Fja25vd2xlZGdlX3VybCI6ICIkSU5DSURFTlRfQUNLTk9XTEVER0VfVVJMIiwKCQkJCSJpbmNpZGVudF9pZCI6ICIkSU5DSURFTlRfSUQiLAoJCQkJImluY2lkZW50X3VybCI6ICIkSU5DSURFTlRfVVJMIiwKCQkJCSJvd25lciI6ICIkRVZFTlRfT1dORVIiLAoJCQkJInBvbGljeV9uYW1lIjogIiRQT0xJQ1lfTkFNRSIsCgkJCQkicG9saWN5X3VybCI6ICIkUE9MSUNZX1VSTCIsCgkJCQkicnVuYm9va191cmwiOiAiJFJVTkJPT0tfVVJMIiwKCQkJCSJzZXZlcml0eSI6ICIkU0VWRVJJVFkiLAoJCQkJInRhcmdldHMiOiAiJFRBUkdFVFMiLAoJCQkJInRpbWVzdGFtcCI6ICIkVElNRVNUQU1QIiwKCQkJCSJ2aW9sYXRpb25fY2hhcnRfdXJsIjogIiRWSU9MQVRJT05fQ0hBUlRfVVJMIgoJCQl9LAoJCQkicGF5bG9hZF90eXBlIjogImFwcGxpY2F0aW9uL3gtd3d3LWZvcm0tdXJsZW5jb2RlZCIKCQl9LAoJCSJuYW1lIjogIiR7bmFtZX0iLAoJCSJ0eXBlIjogIndlYmhvb2siCgl9Cn0="
	return fileName, templateContent
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
