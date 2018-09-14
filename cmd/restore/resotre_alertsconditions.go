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
package restore

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/IBM/newrelic-cli/cmd/backup"
	"github.com/IBM/newrelic-cli/cmd/create"
	"github.com/IBM/newrelic-cli/cmd/delete"
	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/cmd/update"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions",
	Short:   "Restore alertsconditions from directory.",
	Example: "nr restore alertsconditions -d <Directory name where are files to restore>",
	Run: func(cmd *cobra.Command, args []string) {

		var restoreFileFolder string
		var err error
		flags := cmd.Flags()
		if flags.Lookup("dir") != nil {
			restoreFileFolder, err = cmd.Flags().GetString("dir")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}

		}

		var restoreFileNamesParam string
		var restoreFileNames []string
		flags = cmd.Flags()
		if flags.Lookup("files") != nil {
			restoreFileNamesParam, err = cmd.Flags().GetString("files")

			if restoreFileNamesParam != "" {
				restoreFileNames = strings.Split(restoreFileNamesParam, ",")
			}
		}

		var restoreLoggingFileName string
		flags = cmd.Flags()
		if flags.Lookup("File") != nil {
			restoreLoggingFileName, err = cmd.Flags().GetString("File")
		}

		if restoreFileFolder == "" && restoreFileNamesParam == "" && restoreLoggingFileName == "" {
			fmt.Println("Please give restore folder name, files path or logging file name.")
			os.Exit(1)
			return
		}

		var restoreFileNameList []string

		if restoreFileFolder != "" {
			f, err := os.Open(restoreFileFolder)
			defer f.Close()
			if err != nil {
				fmt.Printf("Unable to open file '%v': %v\n", restoreFileFolder, err)
				os.Exit(1)
				return
			}

			dir, err := ioutil.ReadDir(restoreFileFolder)

			for _, fileInfo := range dir {
				if fileInfo.IsDir() == false {
					restoreFileName := fileInfo.Name()

					isBak := strings.HasSuffix(restoreFileName, ".alert-conditions.bak")
					if isBak == true {
						var fileName = restoreFileFolder + "/" + restoreFileName
						restoreFileNameList = append(restoreFileNameList, fileName)
					}
				}
			}
		}

		if restoreLoggingFileName != "" {
			f, err := os.Open(restoreLoggingFileName)
			defer f.Close()
			if err != nil {
				fmt.Printf("Unable to open retore logging file '%v': %v\n", restoreLoggingFileName, err)
				os.Exit(1)
				return
			}

			buf := bufio.NewReader(f)
			for {
				line, err := buf.ReadString('\n')
				line = strings.TrimSpace(line)
				var fileName = line
				isBak := strings.HasSuffix(fileName, ".alert-conditions.bak")
				if isBak == true {
					restoreFileNameList = append(restoreFileNameList, fileName)
				}
				if err != nil {
					if err != io.EOF {
						os.Exit(1)
						return
					} else {
						//end of file
						break
					}
				}
			}
		}

		if len(restoreFileNames) > 0 {
			for _, fileName := range restoreFileNames {
				restoreFileNameList = append(restoreFileNameList, fileName)
			}
		}

		var filesLen = len(restoreFileNameList)

		if filesLen == 0 {
			fmt.Printf("No files to restore found!")
			os.Exit(1)
			return
		}

		var resultFileName string = ""
		if flags.Lookup("result-file-name") != nil {
			resultFileName, err = cmd.Flags().GetString("result-file-name")
		}
		if resultFileName == "" {
			resultFileName = "fail-restore-alert-conditions-file-list.log"
		}

		fmt.Printf("Start to restore all alertsconditions.")

		var updateMode string = "skip"
		if flags.Lookup("update-mode") != nil {
			updateMode, err = cmd.Flags().GetString("update-mode")
		}
		fmt.Printf("Using update mode - %s\n", updateMode)

		var restoreAlertPolicyMetaArray []tracker.RestoreAlertPolicyMeta

		if updateMode == "clean" {
			var rapmArray []tracker.RestoreAlertPolicyMeta
			for _, restoreFileName := range restoreFileNameList {
				var restoreAlertPolicyMeta tracker.RestoreAlertPolicyMeta = tracker.RestoreAlertPolicyMeta{}
				restoreAlertPolicyMeta.FileName = restoreFileName

				rapmArray = append(rapmArray, restoreAlertPolicyMeta)
			}

			//delete all alert policies

			allPolicyList, err, returnValue := get.GetAllAlertPolicies()
			if err != nil {
				exitRestoreAlertCondtionsWithError(returnValue)
				writeFailRestoreConditionsFileList(resultFileName, rapmArray)
				os.Exit(1)
				return
			}
			if returnValue.IsContinue == false {
				exitRestoreAlertCondtionsWithError(returnValue)
				writeFailRestoreConditionsFileList(resultFileName, rapmArray)
				os.Exit(1)
				return
			}

			for _, alertsPolicy := range allPolicyList.AlertsPolicies {
				err, returnValue := delete.DeletePolicyByName(*alertsPolicy.Name)
				if err != nil {
					exitRestoreAlertCondtionsWithError(returnValue)
					writeFailRestoreConditionsFileList(resultFileName, rapmArray)
					os.Exit(1)
					return
				}
				if returnValue.IsContinue == false {
					exitRestoreAlertCondtionsWithError(returnValue)
					writeFailRestoreConditionsFileList(resultFileName, rapmArray)
					os.Exit(1)
					return
				}
			}
		}

		for _, restoreFileName := range restoreFileNameList {

			var restoreAlertPolicyMeta tracker.RestoreAlertPolicyMeta = tracker.RestoreAlertPolicyMeta{}
			restoreAlertPolicyMeta.FileName = restoreFileName
			restoreAlertPolicyMeta.OperationStatus = "fail"

			restoreFile, err := os.Open(restoreFileName)
			defer restoreFile.Close()
			if err != nil {
				fmt.Printf("Unable to open file '%v': %v\n", restoreFileName, err)
				continue
			}

			// validation
			decorder := utils.NewYAMLOrJSONDecoder(restoreFile, 4096)
			var p = new(backup.OneAlertBackup)
			err = decorder.Decode(p)
			if err != nil {
				fmt.Printf("Unable to decode %q: %v\n", restoreFileName, err)
				return
			}
			if reflect.DeepEqual(new(backup.OneAlertBackup), p) {
				fmt.Printf("Error validating %q.\n", restoreFileName)
				return
			}

			alertPolicySet := p.AlertPolicySet
			alertPolicy := alertPolicySet.AlertsPolicy
			// fmt.Println(*alertPolicy.Name)

			newAlertPolicy, isPolicyCreated, err, ret := RestoreOnePolicy(alertPolicy, updateMode)
			if err != nil {
				fmt.Println(err)
				goto next
			}
			if ret.IsContinue == false {
				goto next
			} else {
				if newAlertPolicy == nil {
					//skip
					restoreAlertPolicyMeta.OperationStatus = "success"
					goto next
				}
			}

			if alertPolicySet.AlertsConditionList != nil {
				if alertPolicySet.AlertsConditionList.AlertsDefaultConditionList != nil {
					var isErr bool = false
					for _, defaultCondition := range alertPolicySet.AlertsConditionList.AlertsDefaultConditionList.AlertsDefaultConditions {
						var cat newrelic.ConditionCategory = newrelic.ConditionDefault
						var ac = new(newrelic.AlertsConditionEntity)
						ac.AlertsDefaultConditionEntity = &newrelic.AlertsDefaultConditionEntity{}
						ac.AlertsDefaultConditionEntity.AlertsDefaultCondition = defaultCondition
						err, ret := RestoreOneCondition(*newAlertPolicy.ID, cat, ac, updateMode, isPolicyCreated)
						if err != nil || ret.IsContinue == false {
							isErr = true
							break
						}
					}
					if isErr == true {
						goto next
					}
				}
			}

			if alertPolicySet.AlertsConditionList != nil {
				if alertPolicySet.AlertsConditionList.AlertsExternalServiceConditionList != nil {
					var isErr bool = false
					for _, externalServiceCondition := range alertPolicySet.AlertsConditionList.AlertsExternalServiceConditionList.AlertsExternalServiceConditions {
						var cat newrelic.ConditionCategory = newrelic.ConditionExternalService
						var ac = new(newrelic.AlertsConditionEntity)
						ac.AlertsExternalServiceConditionEntity = &newrelic.AlertsExternalServiceConditionEntity{}
						ac.AlertsExternalServiceConditionEntity.AlertsExternalServiceCondition = externalServiceCondition
						err, ret := RestoreOneCondition(*newAlertPolicy.ID, cat, ac, updateMode, isPolicyCreated)
						if err != nil || ret.IsContinue == false {
							isErr = true
							break
						}
					}
					if isErr == true {
						goto next
					}
				}
			}

			if alertPolicySet.AlertsConditionList != nil {
				if alertPolicySet.AlertsConditionList.AlertsNRQLConditionList != nil {
					var isErr bool = false
					for _, nrqlCondition := range alertPolicySet.AlertsConditionList.AlertsNRQLConditionList.AlertsNRQLConditions {
						var cat newrelic.ConditionCategory = newrelic.ConditionNRQL
						var ac = new(newrelic.AlertsConditionEntity)
						ac.AlertsNRQLConditionEntity = &newrelic.AlertsNRQLConditionEntity{}
						ac.AlertsNRQLConditionEntity.AlertsNRQLCondition = nrqlCondition
						err, ret := RestoreOneCondition(*newAlertPolicy.ID, cat, ac, updateMode, isPolicyCreated)
						if err != nil || ret.IsContinue == false {
							isErr = true
							break
						}
					}
					if isErr == true {
						goto next
					}
				}
			}

			if alertPolicySet.AlertsConditionList != nil {
				if alertPolicySet.AlertsConditionList.AlertsPluginsConditionList != nil {
					var isErr bool = false
					for _, pluginsCondition := range alertPolicySet.AlertsConditionList.AlertsPluginsConditionList.AlertsPluginsConditions {
						var cat newrelic.ConditionCategory = newrelic.ConditionPlugins
						var ac = new(newrelic.AlertsConditionEntity)
						ac.AlertsPluginsConditionEntity = &newrelic.AlertsPluginsConditionEntity{}
						ac.AlertsPluginsConditionEntity.AlertsPluginsCondition = pluginsCondition
						err, ret := RestoreOneCondition(*newAlertPolicy.ID, cat, ac, updateMode, isPolicyCreated)
						if err != nil || ret.IsContinue == false {
							isErr = true
							break
						}
					}
					if isErr == true {
						goto next
					}
				}
			}

			if alertPolicySet.AlertsConditionList != nil {
				if alertPolicySet.AlertsConditionList.AlertsSyntheticsConditionList != nil {
					var isErr bool = false
					for _, syntheticsCondition := range alertPolicySet.AlertsConditionList.AlertsSyntheticsConditionList.AlertsSyntheticsConditions {
						var cat newrelic.ConditionCategory = newrelic.ConditionSynthetics
						var ac = new(newrelic.AlertsConditionEntity)
						ac.AlertsSyntheticsConditionEntity = &newrelic.AlertsSyntheticsConditionEntity{}
						ac.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition = syntheticsCondition
						err, ret := RestoreOneConditionForSynthetics(*newAlertPolicy.ID, cat, ac, updateMode, p.AlertDependencies.MonitorMap, isPolicyCreated)
						if err != nil || ret.IsContinue == false {
							isErr = true
							break
						}
					}
					if isErr == true {
						goto next
					}
				}
			}

			//restore policy channels associations
			err, ret = RestorePolicyChannels(*newAlertPolicy.ID, alertPolicySet.AlertsChannels, updateMode, isPolicyCreated)
			if ret.IsContinue == false {
				goto next
			} else {
				restoreAlertPolicyMeta.OperationStatus = "success"
				goto next
			}

		next:
			restoreAlertPolicyMetaArray = append(restoreAlertPolicyMetaArray, restoreAlertPolicyMeta)
		}

		var restoreAlertPolicyMetaList tracker.RestoreAlertPolicyMetaList = tracker.RestoreAlertPolicyMetaList{}
		restoreAlertPolicyMetaList.AllRestoreAlertPolicyMeta = restoreAlertPolicyMetaArray

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		tracker.PrintStatisticsInfo(restoreAlertPolicyMetaList)

		writeFailRestoreConditionsFileList(resultFileName, restoreAlertPolicyMetaArray)
		fmt.Println()

		fmt.Printf("Restore alert conditions done.")

		os.Exit(0)
	},
}

func RestoreOnePolicy(alertsPolicy *newrelic.AlertsPolicy, mode string) (*newrelic.AlertsPolicy, bool, error, tracker.ReturnValue) {
	isExist, policy, err, ret := get.IsPolicyNameExists(*alertsPolicy.Name)
	if err != nil {
		fmt.Println(err)
		return nil, false, err, ret
	} else {
		if ret.IsContinue == false {
			return nil, false, err, ret
		}
	}
	var newAlertsPolicy *newrelic.AlertsPolicy
	if mode == "skip" {
		if isExist == true {
			ret.IsContinue = true
			return nil, false, err, ret
		} else {
			//create the policy
			newAlertsPolicy, err, ret := create.CreateAlertsPolicy(alertsPolicy)
			if err != nil {
				fmt.Println(err)
				ret.IsContinue = false
				return nil, false, err, ret
			} else {
				if ret.IsContinue == false {
					return nil, false, err, ret
				}
				return newAlertsPolicy, true, err, ret
			}
		}
	} else if mode == "override" {
		if isExist == true {
			//update all by name
			newAlertsPolicy, err, ret := update.UpdateByPolicyName(policy, *policy.Name)
			if err != nil {
				fmt.Println(err)
				return nil, false, err, ret
			}
			if ret.IsContinue == false {
				return nil, false, err, ret
			}
			return newAlertsPolicy, false, err, ret
		} else {
			//create the policy
			newAlertsPolicy, err, ret := create.CreateAlertsPolicy(alertsPolicy)
			if err != nil {
				fmt.Println(err)
				return nil, false, err, ret
			} else {
				if ret.IsContinue == false {
					return nil, false, err, ret
				}
				return newAlertsPolicy, true, err, ret
			}
		}
	} else if mode == "clean" {
		if isExist == true {
			//delete policy by name
			err, ret := delete.DeletePolicyByName(*alertsPolicy.Name)
			if err != nil {
				fmt.Println(err)
				return nil, false, err, ret
			}
			if ret.IsContinue == false {
				return nil, false, err, ret
			}
		}
		//create the policy
		newAlertsPolicy, err, ret = create.CreateAlertsPolicy(alertsPolicy)
		if err != nil {
			fmt.Println(err)
			return nil, false, err, ret
		} else {
			if ret.IsContinue == false {
				return nil, false, err, ret
			}
			return newAlertsPolicy, true, err, ret
		}
	}
	if err != nil {
		fmt.Println(err)
		return nil, false, err, ret
	}

	return nil, false, err, ret
}

func RestoreOneCondition(alertPolicyID int64, cat newrelic.ConditionCategory, c *newrelic.AlertsConditionEntity, mode string, isPolicyCreated bool) (error, tracker.ReturnValue) {
	var err error
	if isPolicyCreated == true {
		//create conditions directly, because the alert policy was new created.
		_, err, ret := create.CreateCondition(cat, c, alertPolicyID)
		if err != nil {
			fmt.Println(err)
			ret.IsContinue = false
			return err, ret
		} else {
			if ret.IsContinue == false {
				return err, ret
			}
		}
	} else {
		//first check if condition exists by name, because the alert policy was updated.
		//if so, update condtion, if not, create condition
		//
		var conditionName string = ""
		switch cat {
		case newrelic.ConditionDefault:
			conditionName = *c.AlertsDefaultCondition.Name
		case newrelic.ConditionExternalService:
			conditionName = *c.AlertsExternalServiceCondition.Name
		case newrelic.ConditionNRQL:
			conditionName = *c.AlertsNRQLCondition.Name
		case newrelic.ConditionPlugins:
			conditionName = *c.AlertsPluginsCondition.Name
		case newrelic.ConditionSynthetics:
			conditionName = *c.AlertsSyntheticsCondition.Name
		}
		isConditionExists, conditionId, err, ret := get.IsConditionNameExists(alertPolicyID, conditionName, cat)
		if err != nil {
			fmt.Println(err)
			ret.IsContinue = false
			return err, ret
		}
		if ret.IsContinue == false {
			return err, ret
		}
		if isConditionExists == true {
			if mode == "skip" {
				//do nothing
			} else if mode == "override" {
				//update condtion
				_, err, ret := update.UpdateCondition(cat, c, conditionId)
				if ret.IsContinue == false {
					return err, ret
				}
			} else if mode == "clean" {
				//delete current condition
				err, ret := delete.DeleteCondition(cat, conditionId)
				if err != nil {
					fmt.Println(err)
					ret.IsContinue = false
					return err, ret
				}
				if ret.IsContinue == false {
					return err, ret
				}
				//create condition
				create.CreateCondition(cat, c, alertPolicyID)
			}
		} else {
			//create condtion directly
			create.CreateCondition(cat, c, alertPolicyID)
		}

	}
	ret := tracker.ToReturnValue(true, "Restore one condtion", nil, nil, "")
	return err, ret
}

func RestoreOneConditionForSynthetics(alertPolicyID int64, cat newrelic.ConditionCategory, c *newrelic.AlertsConditionEntity, mode string, monitorMap map[string]*newrelic.Monitor, isPolicyCreated bool) (error, tracker.ReturnValue) {
	monitorId := c.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.MonitorID
	monitorName := monitorMap[*monitorId].Name
	//check if monitor exist by monitorId
	isExists, monitor, err, ret := get.IsMonitorNameExists(*monitorName)
	if err != nil {
		fmt.Println(err)
		ret.IsContinue = false
		return err, ret
	}
	if ret.IsContinue == false {
		return err, ret
	}
	if isExists == true {
		c.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.MonitorID = monitor.ID
		if mode == "skip" {
			//do nothing
		} else {
			//update monitor by monitor name
			err, returnValue := update.UpdateMonitorByID(monitor.ID, monitor, monitor.Script)
			if err != nil {
				fmt.Println(err)
				return err, returnValue
			}
			if returnValue.IsContinue == false {
				return err, returnValue
			}
		}
	} else {
		//create this synthetics monitor
		monitor := monitorMap[*monitorId]

		newMonitorId, err, returnValue := create.CreateMonitor(monitor, monitor.Script)
		if err != nil {
			fmt.Println(err)
			return err, returnValue
		}
		if returnValue.IsContinue == false {

			return returnValue.OriginalError, returnValue
		} else {
			c.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.MonitorID = &newMonitorId
		}
	}

	if isPolicyCreated == true {
		//create conditions directly, because the alert policy was new created.
		_, err, ret := create.CreateCondition(cat, c, alertPolicyID)
		if err != nil {
			fmt.Println(err)
			return err, ret
		}
		if ret.IsContinue == false {
			return ret.OriginalError, ret
		}
	} else {
		//first check if condition exists by name, because the alert policy was updated.
		//if so, update condtion, if not, create condition
		//
		isConditionExists, conditionId, err, ret := get.IsConditionNameExists(alertPolicyID, *c.AlertsSyntheticsCondition.Name, cat)
		if err != nil {
			fmt.Println(err)
			ret.IsContinue = false
			return err, ret
		}
		if ret.IsContinue == false {
			return ret.OriginalError, ret
		}
		if isConditionExists == true {
			if mode == "skip" {
				//do nothing
			} else if mode == "override" {
				//update condtion
				_, err, ret := update.UpdateCondition(cat, c, conditionId)
				if err != nil {
					fmt.Println(err)
					return err, ret
				}
				if ret.IsContinue == false {
					return ret.OriginalError, ret
				}
			} else if mode == "clean" {
				//delete current condition
				err, ret := delete.DeleteCondition(cat, conditionId)
				if err != nil {
					fmt.Println(err)
					return err, ret
				}
				if ret.IsContinue == false {
					return ret.OriginalError, ret
				}
				//create condition
				_, err, ret = create.CreateCondition(cat, c, alertPolicyID)
				if ret.IsContinue == false {
					return ret.OriginalError, ret
				}
			}
		} else {
			//create condtion directly
			_, err, ret = create.CreateCondition(cat, c, alertPolicyID)
			if ret.IsContinue == false {
				return ret.OriginalError, ret
			}
		}

	}

	ret = tracker.ToReturnValue(true, "Restore one synthetics condtion", nil, nil, "")
	return nil, ret
}

func RestorePolicyChannels(policyId int64, channels []*newrelic.AlertsChannel, mode string, isPolicyCreated bool) (error, tracker.ReturnValue) {
	var channelIds []*int64
	for _, channel := range channels {
		isChannelExists, newChannel, err, ret := get.IsChannelNameExists(*channel.Name)
		if err != nil {
			fmt.Println(err)
			return err, ret
		}
		if ret.IsContinue == false {
			return ret.OriginalError, ret
		}
		if isChannelExists == true {
			channelIds = append(channelIds, newChannel.ID)
		}

	}

	var size = len(channelIds)
	if size == 0 {
		ret := tracker.ToReturnValue(false, "Restore policy channels", tracker.ERR_REST_CHANNEL_NOT_EXIST, tracker.ERR_REST_CHANNEL_NOT_EXIST, "")
		return nil, ret
	}

	if isPolicyCreated == true {
		err, ret := update.UpdatePolicyChannels(policyId, channelIds)
		if err != nil {
			fmt.Println(err)
			return err, ret
		}
		if ret.IsContinue == false {
			return ret.OriginalError, ret
		}
	} else {
		if mode == "skip" {
			//do nothing
		} else if mode == "override" {
			err, ret := update.UpdatePolicyChannels(policyId, channelIds)
			if err != nil {
				fmt.Println(err)
				ret.IsContinue = false
				return err, ret
			}
			if ret.IsContinue == false {
				return ret.OriginalError, ret
			}
		} else if mode == "clean" {
			//create associations
			err, ret := update.UpdatePolicyChannels(policyId, channelIds)
			if err != nil {
				fmt.Println(err)
				ret.IsContinue = false
				return err, ret
			}
			if ret.IsContinue == false {
				return ret.OriginalError, ret
			}
		}
	}

	ret := tracker.ToReturnValue(true, "Restore policy channels", nil, nil, "")
	return nil, ret
}

func writeFailRestoreConditionsFileList(resultFileName string, restoreAlertPolicyMetaArray []tracker.RestoreAlertPolicyMeta) {
	var totalCount = len(restoreAlertPolicyMetaArray)
	var successCount int = 0
	var failCount int = 0
	var failInfoContent string = ""
	for _, meta := range restoreAlertPolicyMetaArray {
		if meta.OperationStatus == "fail" {
			failCount++
			failInfoContent = failInfoContent + meta.FileName + "\r\n"
		} else {
			successCount++
		}
	}
	if failCount == 0 {
		failInfoContent = "" //empty string to represent "no failed"
	}
	var fileLogName = resultFileName
	err := ioutil.WriteFile(fileLogName, []byte(failInfoContent), 0666)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fmt.Printf("Alert conditoins to restore, total: " + strconv.Itoa(totalCount) + ", success: " + strconv.Itoa(successCount) + ", fail: " + strconv.Itoa(failCount))

	if failCount > 0 {
		os.Exit(1)
	}
}

func exitRestoreAlertCondtionsWithError(returnValue tracker.ReturnValue) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to restore alert conditions, exit.")
	os.Exit(1)
}

func init() {
	RestoreCmd.AddCommand(alertsconditionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
