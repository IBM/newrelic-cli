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
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

var dashboardsCmd = &cobra.Command{
	Use:     "dashboards",
	Short:   "Backup dashboards to a directory.",
	Example: "nr backup dashboards -d <Directory of backup dashboards files>",
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

			fmt.Printf("Start to backup all dashboards to '%s' folder\n", backupFolder)
		} else {
			fmt.Println("Please give backup folder.")
			os.Exit(1)
			return
		}

		var resultFileName string = ""
		if flags.Lookup("result-file-name") != nil {
			resultFileName, err = cmd.Flags().GetString("result-file-name")
		}
		if resultFileName == "" {
			resultFileName = "backup-dashboards-file-list.log"
		}

		var backupDashboardMetaList tracker.BackupDashboardMetaList = tracker.BackupDashboardMetaList{}
		var allBackupDashboardMeta []tracker.BackupDashboardMeta

		resultStr, err, returnValue := get.GetAllDashboards()
		if err != nil {
			fmt.Println(err)
			exitBackupDashboardWithError(returnValue, resultFileName)
			os.Exit(1)
			return
		}

		if returnValue.IsContinue == false {
			fmt.Println(returnValue.OriginalError)
			exitBackupDashboardWithError(returnValue, resultFileName)
			os.Exit(1)
			return
		}

		dashboardArr := gjson.Parse(resultStr).Get("dashboards").Array()
		for _, dashboard := range dashboardArr {
			var backupDashboardMeta tracker.BackupDashboardMeta = tracker.BackupDashboardMeta{}

			id := gjson.Parse(dashboard.String()).Get("id")
			title := gjson.Parse(dashboard.String()).Get("title")
			name := title.String()
			var fileName = backupFolder + "/" + name + "-" + id.String() + ".dashboard.bak"

			backupDashboardMeta.FileName = fileName
			backupDashboardMeta.OperationStatus = "fail"

			strDashboard, err, ret := get.GetDashboardByID(id.Int())
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				if ret.IsContinue == false {
					fmt.Println(ret.OriginalError)
					continue
				}
				jsonDashboard := pretty.Pretty([]byte(strDashboard))
				// fmt.Printf("%s\n", string(jsonDashboard))

				err = ioutil.WriteFile(fileName, jsonDashboard, 0666)
				if err != nil {
					fmt.Println(err)
				} else {
					backupDashboardMeta.OperationStatus = "success"
				}
			}
			allBackupDashboardMeta = append(allBackupDashboardMeta, backupDashboardMeta)
		}

		backupDashboardMetaList.AllBackupDashboardMeta = allBackupDashboardMeta

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		tracker.PrintStatisticsInfo(backupDashboardMetaList)
		fmt.Println()
		writeFailDashboardConditionsFileList(resultFileName, allBackupDashboardMeta)

		os.Exit(0)
	},
}

func writeFailDashboardConditionsFileList(resultFileName string, backupDashboardMetaArray []tracker.BackupDashboardMeta) {
	var totalCount = len(backupDashboardMetaArray)
	var successCount int = 0
	var failCount int = 0
	var failInfoContent string = ""
	for _, meta := range backupDashboardMetaArray {
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
	fmt.Printf("Backup dashboards, total: " + strconv.Itoa(totalCount) + ", success: " + strconv.Itoa(successCount) + ", fail: " + strconv.Itoa(failCount))
}

func exitBackupDashboardWithError(returnValue tracker.ReturnValue, resultFileName string) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to backup dashboards, exit.")

	var fileContent = "Backup failed."
	var fileLogName = resultFileName
	err := ioutil.WriteFile(fileLogName, []byte(fileContent), 0666)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(1)
}

func init() {
	BackupCmd.AddCommand(dashboardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
