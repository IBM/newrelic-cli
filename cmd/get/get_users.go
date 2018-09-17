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

	"github.com/spf13/cobra"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/utils"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Display all users.",
	Example: `* nr get users
* nr get users -o json
* nr get users -o yaml
* nr get users -i 2102902
* nr get users -i 2102902,+801314`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := utils.GetNewRelicClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		var opt *newrelic.UserListOptions
		if ids, err := utils.GetArg(cmd, "id"); err == nil {
			opt = &newrelic.UserListOptions{
				IDOptions: ids,
			}
		}
		if emails, err := utils.GetArg(cmd, "email"); err == nil {
			if opt == nil {
				opt = new(newrelic.UserListOptions)
			}
			opt.EmailOptions = emails
		}
		userList, resp, err := client.Users.ListAll(context.Background(), opt)
		if err != nil || resp.StatusCode >= 400 {
			fmt.Printf("%v:%v\n", resp.Status, err)
			os.Exit(1)
			return
		}
		printer, err := utils.NewPriter(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		printer.Print(userList, os.Stdout)

		os.Exit(0)
	},
}

func init() {
	GetCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	usersCmd.Flags().StringP("id", "i", "", "user id(s) to filter returned result. use ',+' to separate ids")
	usersCmd.Flags().StringP("email", "e", "", "email to filter returned result. can't specify emails")
}
