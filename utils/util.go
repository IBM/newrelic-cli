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
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/IBM/newrelic-cli/newrelic"
)

// GetNewRelicClient returns a NewRelicClient, if env var NEW_RELIC_APIKEY is set
func GetNewRelicClient(ctype ...string) (*newrelic.Client, error) {
	var client *newrelic.Client

	var httpClient *http.Client

	proxyStr := os.Getenv("NEW_RELIC_PROXY")
	if proxyStr != "" {
		url, _ := url.Parse(proxyStr)
		proxyURL := http.ProxyURL(url)
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: proxyURL,
			},
		}
	}

	var needCheckAPIKey = true

	if len(ctype) == 0 {
		client = newrelic.NewClient(httpClient, "default")
	} else {
		if ctype[0] == "synthetics" {
			client = newrelic.NewClient(httpClient, "synthetics")
		} else if ctype[0] == "labelSynthetics" {
			client = newrelic.NewClient(httpClient, "labelSynthetics")
		} else if ctype[0] == "insights" {
			client = newrelic.NewClient(httpClient, "insights")
			needCheckAPIKey = false
		} else if ctype[0] == "infrastructure" {
			client = newrelic.NewClient(httpClient, "infrastructure")
		}
	}

	if needCheckAPIKey == true {
		apikey := os.Getenv("NEW_RELIC_APIKEY")
		if apikey == "" {
			return nil, fmt.Errorf("No NEW_RELIC_APIKEY detected.")
		}

		client.XApiKey = apikey
	}

	proxyAuth := os.Getenv("PROXY_AUTH")
	client.ProxyAuth = proxyAuth

	retries := os.Getenv("RETRIES")
	if retries != "" {
		var err error
		client.Retries, err = strconv.Atoi(retries)
		if err != nil {
			client.Retries = 3
		}
	}

	return client, nil
}

func MergeMonitorList(s1 []*newrelic.Monitor, s2 []*newrelic.Monitor) (slice []*newrelic.Monitor) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.Monitor, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return
}

func MergeLabelList(s1 []*newrelic.Label, s2 []*newrelic.Label) (slice []*newrelic.Label) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.Label, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return
}

func MergeMonitorReflList(s1 []*newrelic.MonitorRef, s2 []*newrelic.MonitorRef) (slice []*newrelic.MonitorRef) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.MonitorRef, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return
}

func MergeStringList(s1 []*string, s2 []*string) (slice []*string) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*string, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return
}

func MergeAlertPolicyList(s1 []*newrelic.AlertsPolicy, s2 []*newrelic.AlertsPolicy) (slice []*newrelic.AlertsPolicy) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsPolicy, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertDefaultConditionsList(s1 []*newrelic.AlertsDefaultCondition, s2 []*newrelic.AlertsDefaultCondition) (slice []*newrelic.AlertsDefaultCondition) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsDefaultCondition, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertExternalServiceConditionsList(s1 []*newrelic.AlertsExternalServiceCondition, s2 []*newrelic.AlertsExternalServiceCondition) (slice []*newrelic.AlertsExternalServiceCondition) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsExternalServiceCondition, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertNRQLConditionsList(s1 []*newrelic.AlertsNRQLCondition, s2 []*newrelic.AlertsNRQLCondition) (slice []*newrelic.AlertsNRQLCondition) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsNRQLCondition, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertPluginsConditionsList(s1 []*newrelic.AlertsPluginsCondition, s2 []*newrelic.AlertsPluginsCondition) (slice []*newrelic.AlertsPluginsCondition) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsPluginsCondition, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertSyntheticsConditionsList(s1 []*newrelic.AlertsSyntheticsCondition, s2 []*newrelic.AlertsSyntheticsCondition) (slice []*newrelic.AlertsSyntheticsCondition) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsSyntheticsCondition, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}

func MergeAlertConditionList(s1 *newrelic.AlertsConditionList, s2 *newrelic.AlertsConditionList) *newrelic.AlertsConditionList {
	var newList *newrelic.AlertsConditionList = &newrelic.AlertsConditionList{}
	newList.AlertsDefaultConditionList = &newrelic.AlertsDefaultConditionList{}
	newList.AlertsExternalServiceConditionList = &newrelic.AlertsExternalServiceConditionList{}
	newList.AlertsNRQLConditionList = &newrelic.AlertsNRQLConditionList{}
	newList.AlertsPluginsConditionList = &newrelic.AlertsPluginsConditionList{}
	newList.AlertsSyntheticsConditionList = &newrelic.AlertsSyntheticsConditionList{}

	//default
	if s1.AlertsDefaultConditionList != nil && s2.AlertsDefaultConditionList != nil {
		newList.AlertsDefaultConditions = MergeAlertDefaultConditionsList(s1.AlertsDefaultConditionList.AlertsDefaultConditions, s2.AlertsDefaultConditionList.AlertsDefaultConditions)
	}

	//ExternalService
	if s1.AlertsExternalServiceConditionList != nil && s2.AlertsExternalServiceConditionList != nil {
		newList.AlertsExternalServiceConditions = MergeAlertExternalServiceConditionsList(s1.AlertsExternalServiceConditionList.AlertsExternalServiceConditions, s2.AlertsExternalServiceConditionList.AlertsExternalServiceConditions)
	}
	//nrql
	if s1.AlertsNRQLConditionList != nil && s2.AlertsNRQLConditionList != nil {
		newList.AlertsNRQLConditions = MergeAlertNRQLConditionsList(s1.AlertsNRQLConditionList.AlertsNRQLConditions, s2.AlertsNRQLConditionList.AlertsNRQLConditions)
	}
	//plugin
	if s1.AlertsPluginsConditionList != nil && s2.AlertsPluginsConditionList != nil {
		newList.AlertsPluginsConditions = MergeAlertPluginsConditionsList(s1.AlertsPluginsConditionList.AlertsPluginsConditions, s2.AlertsPluginsConditionList.AlertsPluginsConditions)
	}
	//synthetics
	if s1.AlertsSyntheticsConditionList != nil && s2.AlertsSyntheticsConditionList != nil {
		newList.AlertsSyntheticsConditions = MergeAlertSyntheticsConditionsList(s1.AlertsSyntheticsConditionList.AlertsSyntheticsConditions, s2.AlertsSyntheticsConditionList.AlertsSyntheticsConditions)
	}
	return newList
}

func MergeAlertChannelList(s1 []*newrelic.AlertsChannel, s2 []*newrelic.AlertsChannel) (slice []*newrelic.AlertsChannel) {
	var len1 = len(s1)
	var len2 = len(s2)
	var len3 = len1 + len2
	slice = make([]*newrelic.AlertsChannel, len3)

	for i := 0; i < len1; i++ {
		slice[i] = s1[i]
	}
	for j := 0; j < len2; j++ {
		slice[j+len1] = s2[j]
	}

	return slice
}
