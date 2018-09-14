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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/IBM/newrelic-cli/cmd/create"
	"github.com/IBM/newrelic-cli/cmd/delete"
	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/cmd/update"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var dashboardsCmd = &cobra.Command{
	Use:     "dashboards",
	Short:   "Restore dashboards from directory.",
	Example: "nr restore dashboards -d <Directory name where are files to restore>",
	Run: func(cmd *cobra.Command, args []string) {

		var restoreFileFolder string
		var err error
		flags := cmd.Flags()
		if flags.Lookup("dir") != nil {
			restoreFileFolder, err = cmd.Flags().GetString("dir")
		} else {
			fmt.Println("Please give restore folder name.")
			os.Exit(1)
			return
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

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
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

					isBak := strings.HasSuffix(restoreFileName, ".dashboard.bak")
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
				isBak := strings.HasSuffix(fileName, ".dashboard.bak")
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

		fmt.Println("Start to restore all dashboards")
		fmt.Println()

		var resultFileName string = ""
		if flags.Lookup("result-file-name") != nil {
			resultFileName, err = cmd.Flags().GetString("result-file-name")
		}
		if resultFileName == "" {
			resultFileName = "fail-restore-dashboards-file-list.log"
		}

		var updateMode string = "skip"
		if flags.Lookup("update-mode") != nil {
			updateMode, err = cmd.Flags().GetString("update-mode")
		}
		fmt.Printf("Using update mode - %s\n", updateMode)

		if updateMode == "clean" {
			var ramArray []tracker.RestoreDashboardMeta
			for _, restoreFileName := range restoreFileNameList {
				var restoreDashboardMeta tracker.RestoreDashboardMeta = tracker.RestoreDashboardMeta{}
				restoreDashboardMeta.FileName = restoreFileName

				ramArray = append(ramArray, restoreDashboardMeta)
			}

			//delete all dashboards
			resultStr, err, returnValue := get.GetAllDashboards()
			if err != nil {
				fmt.Println(err)
				exitRestoreDashboardsWithError(returnValue)
				writeFailRestoreDashboardsFileList(resultFileName, ramArray)
				os.Exit(1)
				return
			}
			if returnValue.IsContinue == false {
				fmt.Println(err)
				exitRestoreDashboardsWithError(returnValue)
				writeFailRestoreDashboardsFileList(resultFileName, ramArray)
				os.Exit(1)
				return
			}
			dashboardArr := gjson.Parse(resultStr).Get("dashboards").Array()

			for _, dashboard := range dashboardArr {
				if dashboard.String() != "" {
					id := gjson.Parse(dashboard.String()).Get("id")
					title := gjson.Parse(dashboard.String()).Get("title")
					fmt.Println(id.String())
					fmt.Println(title.String())

					dashboardId, _ := strconv.ParseInt(id.String(), 10, 64)
					err, returnValue := delete.DeleteDashboardByID(dashboardId)
					if err != nil {
						fmt.Println(err)
						exitRestoreDashboardsWithError(returnValue)
						writeFailRestoreDashboardsFileList(resultFileName, ramArray)
						os.Exit(1)
						return
					}
					if returnValue.IsContinue == false {
						fmt.Println(err)
						exitRestoreDashboardsWithError(returnValue)
						writeFailRestoreDashboardsFileList(resultFileName, ramArray)
						os.Exit(1)
						return
					}
				}
			}

			fmt.Println("Delete all dashboards completed.")
		}

		var restoreDashboardMetaArray []tracker.RestoreDashboardMeta

		for _, restoreFileName := range restoreFileNameList {

			var restoreDashboardMeta tracker.RestoreDashboardMeta = tracker.RestoreDashboardMeta{}
			restoreDashboardMeta.FileName = restoreFileName
			restoreDashboardMeta.OperationStatus = "fail"

			restoreFile, err := os.Open(restoreFileName)
			defer restoreFile.Close()
			if err != nil {
				fmt.Printf("Unable to open file '%v': %v\n", restoreFileName, err)
				continue
			}

			bytes, err := ioutil.ReadFile(restoreFileName)
			if err != nil {
				fmt.Print(err)
				continue
			}
			fileContent := string(bytes)
			if !gjson.Valid(fileContent) {
				fmt.Printf("Incorrect JSON format: %v.\n", errors.New("invalid json"))
				continue
			}

			isRestored, err, ret := RestoreOneDashboard(fileContent, updateMode)
			if err != nil {
				fmt.Println(err)
				goto next
			}
			if ret.IsContinue == false {
				goto next
			} else {
				if isRestored == true {
					restoreDashboardMeta.OperationStatus = "success"
				}
				goto next
			}

		next:
			restoreDashboardMetaArray = append(restoreDashboardMetaArray, restoreDashboardMeta)
		}

		var restoreDashboardMetaList tracker.RestoreDashboardMetaList = tracker.RestoreDashboardMetaList{}
		restoreDashboardMetaList.AllRestoreDashboardMeta = restoreDashboardMetaArray

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		tracker.PrintStatisticsInfo(restoreDashboardMetaList)

		writeFailRestoreDashboardsFileList(resultFileName, restoreDashboardMetaArray)
		fmt.Println()

		fmt.Printf("Restore alert dashboards done.")

		os.Exit(0)
	},
}

func RestoreOneDashboard(dashboardContent string, mode string) (bool, error, tracker.ReturnValue) {

	title := gjson.Parse(dashboardContent).Get("dashboard.title").String()
	isExist, _, err, ret := get.IsDashboardTitleExists(title)
	if err != nil {
		fmt.Println(err)
		return false, err, ret
	} else {
		if ret.IsContinue == false {
			return false, err, ret
		}
	}

	// var newDashboard string
	if mode == "skip" {
		if isExist == true {
			ret.IsContinue = true
			return true, err, ret
		} else {
			//create the dashboard
			_, err, ret := create.CreateDashboard(dashboardContent)
			if err != nil {
				fmt.Println(err)
				ret.IsContinue = false
				return false, err, ret
			} else {
				if ret.IsContinue == false {
					return false, err, ret
				}
				return true, err, ret
			}
		}
	} else if mode == "override" {
		if isExist == true {
			//update all by title
			title := gjson.Parse(dashboardContent).Get("dashboard.title").String()
			err, ret := update.UpdateByDashboardTitle(dashboardContent, title)
			if err != nil {
				fmt.Println(err)
				return false, err, ret
			}
			if ret.IsContinue == false {
				return false, err, ret
			}
			return true, err, ret
		} else {
			//create the dashboard
			_, err, ret := create.CreateDashboard(dashboardContent)
			if err != nil {
				fmt.Println(err)
				return false, err, ret
			} else {
				if ret.IsContinue == false {
					return false, err, ret
				}
				return true, err, ret
			}
		}
	} else if mode == "clean" {
		if isExist == true {
			//delete dashbaord by title
			title := gjson.Parse(dashboardContent).Get("dashboard.title").String()
			err, ret := delete.DeleteByDashboardTitle(title)
			if err != nil {
				fmt.Println(err)
				return false, err, ret
			}
			if ret.IsContinue == false {
				return false, err, ret
			}
		}
		//create the dashboard
		_, err, ret = create.CreateDashboard(dashboardContent)
		if err != nil {
			fmt.Println(err)
			return false, err, ret
		} else {
			if ret.IsContinue == false {
				return false, err, ret
			}
			return true, err, ret
		}
	}

	if err != nil {
		fmt.Println(err)
		return false, err, ret
	}

	return true, nil, ret
}

func exitRestoreDashboardsWithError(returnValue tracker.ReturnValue) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to restore dashboards, exit.")
	os.Exit(1)
}

func writeFailRestoreDashboardsFileList(resultFileName string, restoreDashboardMetaArray []tracker.RestoreDashboardMeta) {
	var totalCount = len(restoreDashboardMetaArray)
	var successCount int = 0
	var failCount int = 0
	var failInfoContent string = ""
	for _, meta := range restoreDashboardMetaArray {
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
	fmt.Printf("Dashboards to restore, total: " + strconv.Itoa(totalCount) + ", success: " + strconv.Itoa(successCount) + ", fail: " + strconv.Itoa(failCount))

	if failCount > 0 {
		os.Exit(1)
	}
}

func init() {
	RestoreCmd.AddCommand(dashboardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
