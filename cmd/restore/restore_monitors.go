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

	"github.com/IBM/newrelic-cli/cmd/create"
	"github.com/IBM/newrelic-cli/cmd/delete"
	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/cmd/update"
	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/IBM/newrelic-cli/utils"
	"github.com/spf13/cobra"
)

var monitorsCmd = &cobra.Command{
	Use:     "monitors",
	Short:   "Restore monitors from directory.",
	Example: "nr restore monitors -d <Directory name where are files to restore>",
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

					isBak := strings.HasSuffix(restoreFileName, ".monitor.bak")
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
				isBak := strings.HasSuffix(fileName, ".monitor.bak")
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

		fmt.Println("Start to restore all monitors")
		fmt.Println()

		var resultFileName string = ""
		if flags.Lookup("result-file-name") != nil {
			resultFileName, err = cmd.Flags().GetString("result-file-name")
		}
		if resultFileName == "" {
			resultFileName = "fail-restore-monitors-file-list.log"
		}

		var updateMode string = "skip"
		if flags.Lookup("update-mode") != nil {
			updateMode, err = cmd.Flags().GetString("update-mode")
		}
		fmt.Printf("Using update mode - %s\n", updateMode)

		if updateMode == "clean" {
			var rmmArray []tracker.RestoreMonitorMeta
			for _, restoreFileName := range restoreFileNameList {
				var restoreMonitorMeta tracker.RestoreMonitorMeta = tracker.RestoreMonitorMeta{}
				restoreMonitorMeta.FileName = restoreFileName

				rmmArray = append(rmmArray, restoreMonitorMeta)
			}

			//delete all monitors
			monitors, err, returnValue := get.GetMonitors()
			if err != nil {
				fmt.Println(err)
				writeFailRestoreMonitorsFileList(resultFileName, rmmArray)
				exitRestoreMonitorWithError(returnValue)
				os.Exit(1)
				return
			}
			if returnValue.IsContinue == false {
				fmt.Println(err)
				writeFailRestoreMonitorsFileList(resultFileName, rmmArray)
				exitRestoreMonitorWithError(returnValue)
				os.Exit(1)
				return
			}
			fmt.Println()
			fmt.Println("Deleting all monitors...")
			for _, monitor := range monitors {
				err, returnValue := delete.DeleteMonitorByID(*monitor.ID)
				if err != nil {
					fmt.Println(err)
					writeFailRestoreMonitorsFileList(resultFileName, rmmArray)
					exitRestoreMonitorWithError(returnValue)
					os.Exit(1)
					return
				}
				if returnValue.IsContinue == false {
					fmt.Println(err)
					writeFailRestoreMonitorsFileList(resultFileName, rmmArray)
					exitRestoreMonitorWithError(returnValue)
					os.Exit(1)
					return
				}
			}
		}

		var restoreMonitorMetaArray []tracker.RestoreMonitorMeta

		for _, restoreFileName := range restoreFileNameList {
			var restoreMonitorMeta tracker.RestoreMonitorMeta = tracker.RestoreMonitorMeta{}
			restoreMonitorMeta.FileName = restoreFileName

			fmt.Println("start to restore monitor in file: " + restoreFileName)
			isBak := strings.HasSuffix(restoreFileName, ".monitor.bak")
			if isBak == true {

				restoreFile, err := os.Open(restoreFileName)
				defer restoreFile.Close()
				if err != nil {
					fmt.Printf("Unable to open file '%v': %v\n", restoreFileName, err)
					restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
					continue
				}
				// validation
				decorder := utils.NewYAMLOrJSONDecoder(restoreFile, 4096)
				var p = new(*newrelic.Monitor)
				err = decorder.Decode(p)
				if err != nil {
					fmt.Printf("Unable to decode %q: %v\n", restoreFileName, err)
					restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
					continue
				}
				if reflect.DeepEqual(new([]*newrelic.Monitor), p) {
					fmt.Printf("Error validating %q.\n", restoreFileName)
					restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
					continue
				}

				monitor := *p

				// restoreMonitorMeta.Name = *monitor.Name
				restoreMonitorMeta.Type = *monitor.Type
				if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
					restoreMonitorMeta.Script = true
				} else {
					restoreMonitorMeta.Script = false
				}
				var labelLen = len(monitor.Labels)
				restoreMonitorMeta.LabelCount = labelLen
				if labelLen > 0 {
					for _, label := range monitor.Labels {
						restoreMonitorMeta.Labels = append(restoreMonitorMeta.Labels, *label)
					}
				}
				restoreMonitorMeta.OperationStatus = "fail"

				if updateMode == "clean" {
					//create all monitors
					fmt.Println()
					fmt.Println("Creating one monitor...")

					var scriptTextEncoded *newrelic.Script = nil
					if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
						scriptTextEncoded = monitor.Script
					}
					_, err, returnValue := create.CreateMonitor(monitor, scriptTextEncoded)
					if err != nil {
						restoreMonitorMeta.OperationStatus = "fail"
						restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
						continue
					} else {
						if returnValue.IsContinue == false {
							restoreMonitorMeta.OperationStatus = "fail"

							restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
							continue
						} else {
							restoreMonitorMeta.OperationStatus = "success"
						}
					}
					fmt.Println()
					fmt.Println("Restore monitor done, file name: " + restoreFileName)
				} else {

					var scriptTextEncoded *newrelic.Script = nil
					if *monitor.Type == "SCRIPT_BROWSER" || *monitor.Type == "SCRIPT_API" {
						scriptTextEncoded = monitor.Script
					}
					backupMonitorId := monitor.ID
					//try to create, if response status code is 400, monitor exist, then update
					// _, err, returnValue := create.CreateMonitor(monitor, scriptTextEncoded)
					isExists, _, err, returnValue := get.IsMonitorNameExists(*monitor.Name)
					if err != nil {
						restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
						continue
					} else {
						if returnValue.IsContinue == false {
							restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
							continue
						}
						if isExists == false {
							_, err, returnValue := create.CreateMonitor(monitor, scriptTextEncoded)
							if err != nil {
								restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
								continue
							} else {
								if returnValue.IsContinue == false {
									restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
									continue
								}
								restoreMonitorMeta.OperationStatus = "success"
							}
						} else {
							if updateMode == "override" {
								//update monitor
								monitor.ID = backupMonitorId
								err, ret := update.UpdateMonitorByName(monitor, scriptTextEncoded)
								if err != nil {
									restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
									continue
								} else {
									if ret.IsContinue == false {
										restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
										continue
									}
									restoreMonitorMeta.OperationStatus = "success"
								}
							} else if updateMode == "skip" {
								//skip, do nothing
								fmt.Printf("Monitor '%s' already exists, skip.\n", *monitor.Name)
								restoreMonitorMeta.OperationStatus = "success"
								// restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
							}
							fmt.Println("Restore monitor done, file name: " + restoreFileName)
						}
					}
				}

				restoreMonitorMetaArray = append(restoreMonitorMetaArray, restoreMonitorMeta)
			}
		}

		var restoreMonitorMetaList tracker.RestoreMonitorMetaList = tracker.RestoreMonitorMetaList{}
		restoreMonitorMetaList.AllRestoreMonitorMeta = restoreMonitorMetaArray

		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		tracker.PrintStatisticsInfo(restoreMonitorMetaList)

		writeFailRestoreMonitorsFileList(resultFileName, restoreMonitorMetaArray)

		fmt.Println()
		fmt.Printf("Restore monitors done.\n")

		os.Exit(0)
	},
}

func writeFailRestoreMonitorsFileList(resultFileName string, restoreMonitorMetaArray []tracker.RestoreMonitorMeta) {
	var totalCount = len(restoreMonitorMetaArray)
	var successCount int = 0
	var failCount int = 0
	var failInfoContent string = ""
	for _, meta := range restoreMonitorMetaArray {
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
	fmt.Printf("Monitors to restore, total: " + strconv.Itoa(totalCount) + ", success: " + strconv.Itoa(successCount) + ", fail: " + strconv.Itoa(failCount))

	if failCount > 0 {
		os.Exit(1)
	}
}

func exitRestoreMonitorWithError(returnValue tracker.ReturnValue) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to restore monitors, exit.")
	os.Exit(1)
}

func init() {
	RestoreCmd.AddCommand(monitorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
