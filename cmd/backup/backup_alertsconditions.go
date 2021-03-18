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
package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/spf13/cobra"
)

const (
	MaxConcurrentTask int = 10
)

type AlertPolicySet struct {
	AlertsPolicy        *newrelic.AlertsPolicy        `json:"policy,omitempty"`
	AlertsConditionList *newrelic.AlertsConditionList `json:"alerts_conditions,omitempty"`
	AlertsChannels      []*newrelic.AlertsChannel     `json:"alerts_channels,omitempty"`
}

type AlertDependencies struct {
	MonitorMap map[string]*newrelic.Monitor `json:"dependent_monitors,omitempty"`
}

type AlertBackup struct {
	AlertPolicySetList []AlertPolicySet   `json:"policies,omitempty"`
	AlertDependencies  *AlertDependencies `json:"dependencies,omitempty"`
}

type OneAlertBackup struct {
	AlertPolicySet    AlertPolicySet     `json:"policy,omitempty"`
	AlertDependencies *AlertDependencies `json:"dependencies,omitempty"`
}

var alertsconditionsCmd = &cobra.Command{
	Use:     "alertsconditions",
	Short:   "Backup alertsconditions to a directory.",
	Example: "nr backup alertsconditions -d <<Directory of backup alert conditions files>",
	Run: func(cmd *cobra.Command, args []string) {

		var backupFolder string
		var err error
		flags := cmd.Flags()
		if flags.Lookup("dir") != nil {
			backupFolder, err = cmd.Flags().GetString("dir")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
				return
			}
			if backupFolder == "" {
				fmt.Println("Please give backup folder.")
				os.Exit(1)
				return
			}
			fileInfo, err := os.Stat(backupFolder)
			if err != nil {
				fmt.Println("The folder did not exist.")
				os.Exit(1)
				return
			}
			if fileInfo.IsDir() == false {
				fmt.Println(backupFolder + "is not folder.")
				os.Exit(1)
				return
			}
			fmt.Printf("Start to backup all alertsconditions to '%s' folder\n", backupFolder)
		} else {
			os.Exit(1)
			return
		}

		var resultFileName string = ""
		if flags.Lookup("result-file-name") != nil {
			resultFileName, err = cmd.Flags().GetString("result-file-name")
		}
		if resultFileName == "" {
			resultFileName = "fail-backup-alert-conditions-file-list.log"
		}

		bSingle, _ := cmd.Flags().GetBool("single-file")

		var backupPolicyMetaList tracker.BackupPolicyMetaList = tracker.BackupPolicyMetaList{}
		var allBackupPolicyMeta []tracker.BackupPolicyMeta

		var alertBackup AlertBackup = AlertBackup{}
		var allAlertPolicySet []AlertPolicySet

		alertBackup.AlertDependencies = &AlertDependencies{}
		alertBackup.AlertDependencies.MonitorMap = map[string]*newrelic.Monitor{}

		allChannelList, err, returnValue := get.GetAllAlertsChannels()
		if returnValue.IsContinue == false {
			exitBackupAlertConditionsWithError(returnValue, resultFileName)
			return
		}

		allPolicyList, err, returnValue := get.GetAllAlertPolicies()
		if returnValue.IsContinue == false {
			exitBackupAlertConditionsWithError(returnValue, resultFileName)
			return
		}
		conditionChMap := make(map[int64]chan *newrelic.AlertsConditionList)
		chTaskCtrl := make(chan struct{}, MaxConcurrentTask)
		defer close(chTaskCtrl)

		for _, alertsPolicy := range allPolicyList.AlertsPolicies {
			name := *alertsPolicy.Name
			alertPolicyID := *alertsPolicy.ID
			r := make(chan *newrelic.AlertsConditionList)
			go func() {
				defer close(r)
				chTaskCtrl <- struct{}{}
				fmt.Printf("Fetching alert conditions for Policy: %s\n", name)
				conditionList, _, returnValue := get.GetAllConditionsByAlertPolicyID(alertPolicyID)
				<-chTaskCtrl
				if returnValue.IsContinue == false {
					r <- nil
					return
				}
				r <- conditionList
				return
			}()
			conditionChMap[alertPolicyID] = r
		}

		for _, alertsPolicy := range allPolicyList.AlertsPolicies {
			var backupPolicyMeta tracker.BackupPolicyMeta = tracker.BackupPolicyMeta{}

			var alertPolicySet AlertPolicySet = AlertPolicySet{}

			alertPolicySet.AlertsPolicy = alertsPolicy

			var policyName = *alertsPolicy.Name
			var ID = *alertsPolicy.ID
			var fileNamePrefix = policyName + "-" + strconv.FormatInt(ID, 10)
			if bSingle == true {
				backupPolicyMeta.Policy = fileNamePrefix
				backupPolicyMeta.FileName = backupFolder + "/all-in-one-bundle.alert-conditions.bak"
			} else {
				backupPolicyMeta.Policy = strconv.FormatInt(ID, 10)
				backupPolicyMeta.FileName = backupFolder + "/" + fileNamePrefix + ".alert-conditions.bak"
			}
			backupPolicyMeta.OperationStatus = "fail"
			// backupPolicyMeta.PolicyName = policyName

			var alertPolicyID = alertsPolicy.ID
			conditionList := <-conditionChMap[*alertPolicyID]
			if conditionList == nil {
				backupPolicyMeta.OperationStatus = "fail"
				allBackupPolicyMeta = append(allBackupPolicyMeta, backupPolicyMeta)
				continue
			}

			bNodeps, _ := cmd.Flags().GetBool("no-deps")

			if bNodeps == false {
				syntheticsArray := conditionList.AlertsSyntheticsConditions

				// alertPolicySet.MonitorList := []*newrelic.Monitor
				for _, monitor := range syntheticsArray {
					if monitor.MonitorID != nil {
						fmt.Printf("Calling  GetMonitorByID() func, monitor id: %s\n", *monitor.MonitorID)
						m, err, ret := get.GetMonitorByID(*monitor.MonitorID)
						if err != nil {
							fmt.Println(err)
						}
						if ret.IsContinue == false {
							//ignore err
						}
						alertBackup.AlertDependencies.MonitorMap[*monitor.MonitorID] = m
					}
				}
			}
			alertPolicySet.AlertsConditionList = conditionList

			//process channels
			for _, channel := range allChannelList.AlertsChannels {
				for _, policyID := range channel.Links.PolicyIDs {
					if *policyID == *alertPolicyID {
						alertPolicySet.AlertsChannels = append(alertPolicySet.AlertsChannels, channel)
					}
				}
			}
			////

			allAlertPolicySet = append(allAlertPolicySet, alertPolicySet)

			backupPolicyMeta.OperationStatus = "success"
			allBackupPolicyMeta = append(allBackupPolicyMeta, backupPolicyMeta)

		}

		alertBackup.AlertPolicySetList = allAlertPolicySet

		backupPolicyMetaList.AllBackupPolicyMeta = allBackupPolicyMeta

		fmt.Println()

		if bSingle == true {
			fileContentBundle, err := json.MarshalIndent(alertBackup.AlertPolicySetList, "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			var fileName = backupFolder + "/all-in-one-bundle.alert-conditions.bak"
			err = ioutil.WriteFile(fileName, fileContentBundle, 0666)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			for _, policy := range alertBackup.AlertPolicySetList {
				var onePolicy OneAlertBackup = OneAlertBackup{}

				var name = *policy.AlertsPolicy.Name
				var ID = *policy.AlertsPolicy.ID
				var fileNamePrefix = name + "-" + strconv.FormatInt(ID, 10)

				onePolicy.AlertPolicySet = policy
				onePolicy.AlertDependencies = alertBackup.AlertDependencies

				fileContent, err := json.MarshalIndent(onePolicy, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				var fileName = backupFolder + "/" + fileNamePrefix + ".alert-conditions.bak"
				err = ioutil.WriteFile(fileName, fileContent, 0666)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		tracker.PrintStatisticsInfo(backupPolicyMetaList)
		fmt.Println()
		writeFailBackupConditionsFileList(resultFileName, allBackupPolicyMeta)
		fmt.Println()

		os.Exit(0)
	},
}

func writeFailBackupConditionsFileList(resultFileName string, backupPolicyMetaArray []tracker.BackupPolicyMeta) {
	var totalCount = len(backupPolicyMetaArray)
	var successCount int = 0
	var failCount int = 0
	var failInfoContent string = ""
	for _, meta := range backupPolicyMetaArray {
		if meta.OperationStatus == "fail" {
			failCount++
			failInfoContent = failInfoContent + meta.FileName + "\r\n"
		} else {
			successCount++
		}
	}
	if failCount == 0 {
		failInfoContent = "No failed"
	}
	var fileLogName = resultFileName
	err := ioutil.WriteFile(fileLogName, []byte(failInfoContent), 0666)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	fmt.Printf("Backup alert conditions, total: " + strconv.Itoa(totalCount) + ", success: " + strconv.Itoa(successCount) + ", fail: " + strconv.Itoa(failCount))
}

func exitBackupAlertConditionsWithError(returnValue tracker.ReturnValue, resultFileName string) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to backup alert conditions, exit.")

	var fileContent = "Backup failed."
	var fileLogName = resultFileName
	err := ioutil.WriteFile(fileLogName, []byte(fileContent), 0666)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(1)
}

func init() {
	BackupCmd.AddCommand(alertsconditionsCmd)
	alertsconditionsCmd.PersistentFlags().BoolP("no-deps", "n", false, "Don't get associated monitor confiugration")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
