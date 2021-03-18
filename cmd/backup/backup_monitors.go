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

	"github.com/IBM/newrelic-cli/cmd/get"
	"github.com/IBM/newrelic-cli/tracker"
	"github.com/spf13/cobra"
)

var monitorsCmd = &cobra.Command{
	Use:     "monitors",
	Short:   "Backup monitors to a directory.",
	Example: "nr backup monitors -d <Directory of backup monitors files>",
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

			fmt.Printf("Start to backup all monitors to '%s' folder\n", backupFolder)
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
			resultFileName = "backup-monitors-file-list.log"
		}

		//
		//get all monitors
		monitorArray, err, returnValue := get.GetMonitors()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// 	return
		// }
		if returnValue.IsContinue == false {
			exitBackupMonitorWithError(returnValue, resultFileName)
			return
		}

		bSingle, _ := cmd.Flags().GetBool("single-file")

		if bSingle == true {
			fileContentBundle, err := json.MarshalIndent(monitorArray, "", "  ")
			if err != nil {
				fmt.Println(err)
			}
			var fileName = backupFolder + "/all-in-one-bundle.monitor.bak"
			err = ioutil.WriteFile(fileName, fileContentBundle, 0666)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// var backupFileNameStr string = ""

			for _, monitor := range monitorArray {
				var name = *monitor.Name
				fileContent, err := json.MarshalIndent(monitor, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				var fileName = backupFolder + "/" + name + ".monitor.bak"
				err = ioutil.WriteFile(fileName, fileContent, 0666)
				if err != nil {
					fmt.Println(err)
				}
				// else {
				// 	backupFileNameStr = backupFileNameStr + fileName + "\r\n"
				// }

			}
		}

		var fileContent = "No failed."
		var fileLogName = resultFileName
		err = ioutil.WriteFile(fileLogName, []byte(fileContent), 0666)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println()
		//print REST call
		tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
		fmt.Println()
		//print statistics monitor list
		tracker.PrintBackupMonitorInfo(monitorArray)

		fmt.Println()

		fmt.Printf("Backup monitors done. folder: %v\n", backupFolder)
		fmt.Println()
		fmt.Printf("Backup monitor file names of failure status listed in this file: " + fileLogName)
		fmt.Println()

		os.Exit(0)
	},
}

func exitBackupMonitorWithError(returnValue tracker.ReturnValue, resultFileName string) {
	//print REST call
	tracker.PrintStatisticsInfo(tracker.GlobalRESTCallResultList)
	fmt.Println()

	fmt.Println(returnValue.OriginalError)
	fmt.Println(returnValue.TypicalError)
	fmt.Println(returnValue.Description)
	fmt.Println()

	fmt.Println("Failed to backup monitors, exit.")

	var fileContent = "Backup failed."
	var fileLogName = resultFileName
	err := ioutil.WriteFile(fileLogName, []byte(fileContent), 0666)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(1)
}

func init() {
	BackupCmd.AddCommand(monitorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

}
