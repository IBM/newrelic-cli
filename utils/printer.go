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
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

type Printer interface {
	Print(obj interface{}, out io.Writer)
}

func GetArg(cmd *cobra.Command, arg string) (string, error) {
	flags := cmd.Flags()

	var val string
	var err error
	if flags.Lookup(arg) != nil {
		val, err = cmd.Flags().GetString(arg)
		if err != nil {
			return "", fmt.Errorf("error accessing flag %s for command %s: %v", arg, cmd.Name(), err)
		}
	}

	return val, err
}

func NewPriter(cmd *cobra.Command) (Printer, error) {
	flags := cmd.Flags()

	var output string
	var err error
	if flags.Lookup("output") != nil {
		output, err = cmd.Flags().GetString("output")
		if err != nil {
			return nil, fmt.Errorf("error accessing flag %s for command %s: %v", "output", cmd.Name(), err)
		}
	}

	var printer Printer
	switch output {
	case "yaml", "yml":
		printer = &YAMLPrinter{}
	case "json":
		printer = &JSONPrinter{}
	default:
		printer = &TablePrinter{}
	}

	return printer, nil
}

type JSONPrinter struct{}

type YAMLPrinter struct{}

type TablePrinter struct{}

func (p *JSONPrinter) Print(obj interface{}, out io.Writer) {
	if output, err := json.MarshalIndent(obj, "", "  "); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", string(output))
	}
}

func (p *YAMLPrinter) Print(obj interface{}, out io.Writer) {
	if output, err := yaml.Marshal(obj); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", string(output))
	}
}

func (p *TablePrinter) Print(obj interface{}, out io.Writer) {
	w := new(tabwriter.Writer)
	w.Init(out, 5, 1, 3, ' ', 0)

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// smartly detect if obj is a List or plain Object
	if val.NumField() == 1 {
		val = val.Field(0)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() == reflect.Slice {
			// iterate each Object
			heading := []string{}
			rows := []string{}
			for i := 0; i < val.Len(); i++ {
				e := val.Index(i)
				if e.Kind() == reflect.Ptr {
					e = e.Elem()
				}
				t := e.Type()
				row := []string{}
				for j := 0; j < e.NumField(); j++ {
					f := e.Field(j)
					if f.Kind() == reflect.Ptr {
						f = f.Elem()
					}
					if i == 0 {
						heading = append(heading, fmt.Sprint(t.Field(j).Name))
					}
					row = append(row, fmt.Sprint(f.Interface()))
					// fmt.Printf("%d: %s %s = %v\n", j, t.Field(j).Name, f.Type(), f.Elem().Interface())
				}
				rows = append(rows, strings.Join(row, "\t"))
			}
			if len(heading) == 0 {
				fmt.Fprintln(w, "No resources found.")
			} else {
				fmt.Fprintln(w, strings.Join(heading, "\t"))
				fmt.Fprintln(w, strings.Join(rows, "\n"))
			}
		} else if val.Kind() == reflect.Struct {
			t := val.Type()
			heading := []string{}
			row := []string{}
			for j := 0; j < val.NumField(); j++ {
				f := val.Field(j)
				heading = append(heading, fmt.Sprint(t.Field(j).Name))
				row = append(row, fmt.Sprint(f.Elem().Interface()))
				// fmt.Printf("%d: %s %s = %v\n", j, t.Field(j).Name, f.Type(), f.Elem().Interface())
			}
			if len(heading) == 0 {
				fmt.Fprintln(w, "No resources found.")
			} else {
				fmt.Fprintln(w, strings.Join(heading, "\t"))
				fmt.Fprintln(w, strings.Join(row, "\t"))
			}
		}
	} else {
		fmt.Println("[WARNING] unsupported data format.")
	}

	w.Flush()
}
