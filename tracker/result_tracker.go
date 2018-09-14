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
package tracker

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/IBM/newrelic-cli/newrelic"
	"github.com/IBM/newrelic-cli/utils"
)

type RESTCallResultList struct {
	AllRESTCallResult []RESTCallResult
}

type BackupMonitorMetaList struct {
	AllBackupMonitorMeta []BackupMonitorMeta
}

type BackupMonitorMeta struct {
	Name       string
	Type       string
	Script     bool
	Labels     []string
	LabelCount int
	ID         string
}

type BackupPolicyMetaList struct {
	AllBackupPolicyMeta []BackupPolicyMeta
}

type BackupPolicyMeta struct {
	FileName        string
	OperationStatus string
}

type BackupDashboardMetaList struct {
	AllBackupDashboardMeta []BackupDashboardMeta
}

type BackupDashboardMeta struct {
	FileName        string
	OperationStatus string
}

type RestoreMonitorMetaList struct {
	AllRestoreMonitorMeta []RestoreMonitorMeta
}

type RestoreMonitorMeta struct {
	FileName string
	// Name             string
	Type       string
	Script     bool
	Labels     []string
	LabelCount int
	// PropertiesStatus string
	// ScriptStatus     string
	// LabelStatus      string
	OperationStatus string
}

type RestoreAlertPolicyMetaList struct {
	AllRestoreAlertPolicyMeta []RestoreAlertPolicyMeta
}

type RestoreAlertPolicyMeta struct {
	FileName        string
	OperationStatus string
}

type RestoreDashboardMetaList struct {
	AllRestoreDashboardMeta []RestoreDashboardMeta
}

type RestoreDashboardMeta struct {
	FileName        string
	OperationStatus string
}

type RESTCallResult struct {
	OperationName string
	StatusCode    int
	Description   string
	// HTTPMethod    string
	Message string
}

// type MonitorResult struct {
// 	Name            string
// 	Type            string
// 	Labels          []string
// 	OperationStatus string
// }

type AlertResult struct {
	AlertPolicyName string
	OperationStatus string
}

type AlertConditionResult struct {
	AlertConditionName string
	OperationStatus    string
}

type ReturnValue struct {
	OriginalError error
	IsContinue    bool
	TypicalError  error
	Description   string
	OperationName string
}

var OPERATION_NAME_GET_MONITORS = "Get Monitors"
var OPERATION_NAME_GET_MONITOR_BY_ID = "Get Monitor By ID"
var OPERATION_NAME_GET_MONITOR_BY_NAME = "Get Monitor By Name"
var OPERATION_NAME_GET_MONITOR_SCRIPT = "Get Monitor Script"
var OPERATION_NAME_GET_LABELS = "Get Labels"
var OPERATION_NAME_GET_LABELS_BY_MONITOR_ID = "Get Labels By Monitor ID"
var OPERATION_NAME_GET_MONITORS_BY_LABEL = "Get Monitors By Label"
var OPERATION_NAME_GET_ALERT_CHANNELS = "Get Alert Channels"
var OPERATION_NAME_GET_ALERT_POLICIES = "Get Alert Policies"
var OPERATION_NAME_GET_CONDITIONS_BY_POLICY_ID = "Get Conditions By Alert Policy ID"
var OPERATION_NAME_GET_DASHBOARDS = "Get Dashboards"
var OPERATION_NAME_GET_DASHBOARD_BY_ID = "Get Dashboard By ID"
var OPERATION_NAME_GET_DASHBOARD_BY_NAME = "Get Dashboard By Name"

var OPERATION_NAME_CHECK_ALERT_CHANNEL_NAME_EXISTS = "Check Alert Channel Name Exists"
var OPERATION_NAME_CHECK_ALERT_POLICY_NAME_EXISTS = "Check Alert Policy Name Exists"
var OPERATION_NAME_CHECK_MONITOR_NAME_EXISTS = "Check Monitor Name Exists"
var OPERATION_NAME_CHECK_DASHBOARD_TITLE_EXISTS = "Check Dashboard Title Exists"

var OPERATION_NAME_CREATE_MONITOR = "Create Monitor"
var OPERATION_NAME_CREATE_ALERT_POLICY = "Create Alert Policy"
var OPERATION_NAME_CREATE_ALERT_CONDITIION = "Create Alert Condition"
var OPERATION_NAME_CREATE_DASHBOARD = "Create Dashboard"

var OPERATION_NAME_UPDATE_MONITOR = "Update Monitor"
var OPERATION_NAME_UPDATE_MONITOR_SCRIPT = "Create Monitor"
var OPERATION_NAME_UPDATE_ALERT_POLICY_BY_NAME = "Update Alert Policy By Name"
var OPERATION_NAME_UPDATE_ALERT_POLICY_BY_ID = "Update Alert Policy By ID"
var OPERATION_NAME_UPDATE_ALERT_CONDITION_BY_ID = "Update Alert Condition By ID"
var OPERATION_NAME_UPDATE_ALERT_POLICY_CHANNEL = "Update Alert Policy Channel"
var OPERATION_NAME_UPDATE_DASHBOARD_BY_ID = "Update Dashboard By ID"
var OPERATION_NAME_UPDATE_DASHBOARD_BY_NAME = "Update Dashboard By Name"

var OPERATION_NAME_PATCH_MONITOR = "Patch Monitor"

var OPERATION_NAME_DELETE_MONITOR = "Delete Monitor"
var OPERATION_NAME_DELETE_LABEL_FROM_MONITOR = "Delete Label From Monitor"
var OPERATION_NAME_DELETE_ALERT_POLICY_BY_ID = "Delete Alert Policy By ID"
var OPERATION_NAME_DELETE_ALERT_POLICY_BY_NAME = "Delete Alert Policy By Name"
var OPERATION_NAME_DELETE_ALERT_CONDITION = "Delete Alert Condition"
var OPERATION_NAME_DELETE_DASHBOARD_BY_ID = "Delete Dashboard By ID"

var OPERATION_NAME_ADD_LABEL_MONITOR = "Add Label Monitor"

var OPERATION_NAME_INSERT_CUSTOM_EVENTS = "Insert Custom Events"

var ERR_CREATE_NR_CLINET = errors.New("Call NewRelic REST error")
var ERR_REST_CALL = errors.New("Call NewRelic REST error")
var ERR_REST_CALL_NOT_2XX = errors.New("Status code is not 2XX calling NewRelic REST")
var ERR_REST_CALL_400 = errors.New("Status code is 400 calling NewRelic REST")
var ERR_REST_CHANNEL_NOT_EXIST = errors.New("No any notification channels")

// var ERR_GET_MONITOR = errors.New("Get Monitor error")
// var ERR_CREATE_MONITOR = errors.New("Create Monitor error")
// var ERR_UPDATE_MONITOR = errors.New("Create Monitor error")
// var ERR_PATCH_MONITOR = errors.New("Create Monitor error")
// var ERR_DELETE_MONITOR = errors.New("Create Monitor error")

var STATUS_CODE_MAPPING_CONDITIONS = make(map[int]string)
var STATUS_CODE_MAPPING_MONITORS = make(map[int]string)
var STATUS_CODE_MAPPING_MONITORS_SCRIPT = make(map[int]string)
var STATUS_CODE_MAPPING_LABELS = make(map[int]string)
var STATUS_CODE_MAPPING_LABELS_SYNTHETICS = make(map[int]string)
var STATUS_CODE_MAPPING_ALERT_CHANNELS = make(map[int]string)
var STATUS_CODE_MAPPING_ALERT_POLICIES = make(map[int]string)
var STATUS_CODE_MAPPING_ALERT_CONDITIONS = make(map[int]string)
var STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS = make(map[int]string)

var GlobalRESTCallResultList RESTCallResultList

func init() {
	/*
		alert conditions:

		401 	Invalid API key
		401 	Invalid request, API key required
		403 	New Relic API access has not been enabled
		500 	A server error occurred, please contact New Relic support
		404 	No Alerts condition was found with the given ID
		406 	Bad entity type
		422 	Validation error occurred while trying to update the alert condition
	*/
	STATUS_CODE_MAPPING_CONDITIONS[200] = "Success"
	STATUS_CODE_MAPPING_CONDITIONS[400] = "Bad request"
	STATUS_CODE_MAPPING_CONDITIONS[401] = "Invalid request, API key required"
	STATUS_CODE_MAPPING_CONDITIONS[403] = "New Relic API access has not been enabled"
	STATUS_CODE_MAPPING_CONDITIONS[500] = "A server error occurred, please contact New Relic support"
	STATUS_CODE_MAPPING_CONDITIONS[404] = "No Alerts condition was found with the given ID"
	STATUS_CODE_MAPPING_CONDITIONS[406] = "Bad entity type"
	STATUS_CODE_MAPPING_CONDITIONS[422] = "Validation error occurred while trying to update the alert condition"

	/*
		synthetics monitors:
	*/
	STATUS_CODE_MAPPING_MONITORS[200] = "Success"
	STATUS_CODE_MAPPING_MONITORS[201] = "Success"
	STATUS_CODE_MAPPING_MONITORS[204] = "Success"

	STATUS_CODE_MAPPING_MONITORS[400] = "The monitor values is invalid, or the format of the request is invalid."
	STATUS_CODE_MAPPING_MONITORS[404] = "The specified monitor doesn not exist"
	STATUS_CODE_MAPPING_MONITORS[500] = "A server error occurred, please contact New Relic support"

	/*
		synthetics monitor script:
	*/
	STATUS_CODE_MAPPING_MONITORS_SCRIPT[200] = "Success"
	STATUS_CODE_MAPPING_MONITORS_SCRIPT[204] = "Success"

	STATUS_CODE_MAPPING_MONITORS[400] = "The monitor values is invalid, or the format of the request is invalid."
	STATUS_CODE_MAPPING_MONITORS[404] = "The specified monitor doesn not exist"
	STATUS_CODE_MAPPING_MONITORS[500] = "A server error occurred, please contact New Relic support"

	/*
		labels

		401 	Invalid API key
		401 	Invalid request, API key required
		403 	New Relic API access has not been enabled
		500 	A server error occurred, please contact New Relic support
	*/
	STATUS_CODE_MAPPING_LABELS[200] = "Success"
	STATUS_CODE_MAPPING_LABELS[400] = "Bad request"
	STATUS_CODE_MAPPING_LABELS[401] = "Invalid request, API key required"
	STATUS_CODE_MAPPING_LABELS[403] = "New Relic API access has not been enabled"
	STATUS_CODE_MAPPING_LABELS[500] = "A server error occurred, please contact New Relic support"

	/*
		label synthetics
	*/
	STATUS_CODE_MAPPING_LABELS_SYNTHETICS[200] = "Success"
	STATUS_CODE_MAPPING_LABELS_SYNTHETICS[204] = "Success"

	STATUS_CODE_MAPPING_LABELS_SYNTHETICS[400] = "Bad request, or the format of the request is invalid."
	STATUS_CODE_MAPPING_LABELS_SYNTHETICS[404] = "The specified label doesn not exist"
	STATUS_CODE_MAPPING_LABELS_SYNTHETICS[500] = "A server error occurred, please contact New Relic support"

	/*
		alert channels

		401 	Invalid API key
		401 	Invalid request, API key required
		403 	New Relic API access has not been enabled
		500 	A server error occurred, please contact New Relic support
		422 	Validation or internal error occurred
	*/

	STATUS_CODE_MAPPING_ALERT_CHANNELS[200] = "Success"
	STATUS_CODE_MAPPING_ALERT_CHANNELS[400] = "Bad request"
	STATUS_CODE_MAPPING_ALERT_CHANNELS[401] = "Invalid request, API key required"
	STATUS_CODE_MAPPING_ALERT_CHANNELS[403] = "New Relic API access has not been enabled"
	STATUS_CODE_MAPPING_ALERT_CHANNELS[500] = "A server error occurred, please contact New Relic support"
	STATUS_CODE_MAPPING_ALERT_CHANNELS[422] = "Validation or internal error occurred"

	/*
		alert policies

		401 	Invalid API key
		401 	Invalid request, API key required
		403 	New Relic API access has not been enabled
		500 	A server error occurred, please contact New Relic support
		422 	Validation or internal error occurred
	*/

	STATUS_CODE_MAPPING_ALERT_POLICIES[200] = "Success"
	STATUS_CODE_MAPPING_ALERT_POLICIES[400] = "Bad request"
	STATUS_CODE_MAPPING_ALERT_POLICIES[401] = "Invalid request, API key required"
	STATUS_CODE_MAPPING_ALERT_POLICIES[403] = "New Relic API access has not been enabled"
	STATUS_CODE_MAPPING_ALERT_POLICIES[500] = "A server error occurred, please contact New Relic support"
	STATUS_CODE_MAPPING_ALERT_POLICIES[422] = "Validation or internal error occurred"

	/*
		alert conditions

		401 	Invalid API key
		401 	Invalid request, API key required
		403 	New Relic API access has not been enabled
		500 	A server error occurred, please contact New Relic support
		404 	No Alerts policy was found for the given ID
		422 	Validation error occurred while trying to create the alert condition
	*/
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[200] = "Success"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[400] = "Bad request"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[401] = "Invalid request, API key required"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[403] = "New Relic API access has not been enabled"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[404] = "No Alerts policy was found for the given ID"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[500] = "A server error occurred, please contact New Relic support"
	STATUS_CODE_MAPPING_ALERT_CONDITIONS[422] = "Validation error occurred while trying to create the alert condition"

	/*
		Insights custom events
	*/
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[200] = "Success"
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[400] = "Bad request"
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[403] = "Invalid insert key"
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[408] = "Request timed out"
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[413] = "Content too large"
	STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[429] = "Too many requests"

	GlobalRESTCallResultList = RESTCallResultList{}
}

func AppendRESTCallResult(serviceInstance interface{}, operationName string, statusCode int, message string) {
	ret := ToRESTCallResult(serviceInstance, operationName, statusCode, message)
	GlobalRESTCallResultList.AllRESTCallResult = append(GlobalRESTCallResultList.AllRESTCallResult, ret)
}

func ToRESTCallResult(serviceInstance interface{}, operationName string, statusCode int, message string) RESTCallResult {
	var ret RESTCallResult = RESTCallResult{}
	ret.OperationName = operationName
	ret.StatusCode = statusCode
	ret.Message = message
	var desc string = ""
	desc = GetDescByStatusCode(serviceInstance, statusCode)
	ret.Description = desc
	return ret
}

func ToReturnValue(isContinue bool, operationName string, originalErr error, typicalErr error, description string) ReturnValue {
	var ret ReturnValue = ReturnValue{}
	ret.OriginalError = originalErr
	ret.IsContinue = isContinue
	ret.TypicalError = typicalErr
	ret.Description = description
	ret.OperationName = operationName
	return ret
}

func GetDescByStatusCode(serviceInstance interface{}, statusCode int) string {
	var desc string = ""

	if statusCode == 200 {
		desc = "Success"
		return desc
	}

	var typeName string = reflect.TypeOf(serviceInstance).String()

	if typeName == "" {

	} else if typeName == "*newrelic.SyntheticsService" {
		desc = STATUS_CODE_MAPPING_MONITORS[statusCode]
	} else if typeName == "*newrelic.ScriptService" {
		desc = STATUS_CODE_MAPPING_MONITORS_SCRIPT[statusCode]
	} else if typeName == "*newrelic.LabelsService" {
		desc = STATUS_CODE_MAPPING_LABELS[statusCode]
	} else if typeName == "*newrelic.LabelsSyntheticsService" {
		desc = STATUS_CODE_MAPPING_LABELS_SYNTHETICS[statusCode]
	} else if typeName == "*newrelic.AlertsChannelsService" {
		desc = STATUS_CODE_MAPPING_ALERT_CHANNELS[statusCode]
	} else if typeName == "*newrelic.AlertsPoliciesService" {
		desc = STATUS_CODE_MAPPING_ALERT_POLICIES[statusCode]
	} else if typeName == "*newrelic.AlertsConditionsService" {
		desc = STATUS_CODE_MAPPING_ALERT_CONDITIONS[statusCode]
	} else if typeName == "*newrelic.CustomEventService" {
		desc = STATUS_CODE_MAPPING_ALERT_CUSTOM_EVENTS[statusCode]
	}

	return desc
}

func GenerateBackupMonitorMeta(monitorList []*newrelic.Monitor) BackupMonitorMetaList {
	var backupMonitorMetaList []BackupMonitorMeta
	for _, monitor := range monitorList {
		var m BackupMonitorMeta = BackupMonitorMeta{}
		m.ID = *monitor.ID
		m.Name = *monitor.Name
		m.Type = *monitor.Type
		if monitor.Script != nil && monitor.Script.ScriptText != nil {
			m.Script = true
		} else {
			m.Script = false
		}
		if monitor.Labels != nil {
			var labels []string
			for _, label := range monitor.Labels {
				labels = append(labels, *label)
			}
			m.Labels = labels
			var labelLen = len(m.Labels)
			m.LabelCount = labelLen
		}
		backupMonitorMetaList = append(backupMonitorMetaList, m)
	}

	var allList BackupMonitorMetaList = BackupMonitorMetaList{}
	allList.AllBackupMonitorMeta = backupMonitorMetaList

	return allList
}

func PrintStatisticsInfo(obj interface{}) {
	//print statistics
	var printer utils.Printer = &utils.TablePrinter{}
	printer.Print(obj, os.Stdout)
}

func PrintBackupMonitorInfo(monitorList []*newrelic.Monitor) {
	allList := GenerateBackupMonitorMeta(monitorList)
	PrintStatisticsInfo(allList)
	var monitorLen = len(monitorList)
	var msg = strconv.Itoa(monitorLen) + " monitors backuped."
	fmt.Println()
	fmt.Println(msg)
}
