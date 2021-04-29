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
package newrelic

import (
	"context"
)

type locationConditions service

type LocationTerm struct {
	Priority  *string `json:"priority,omitempty"`
	Threshold *int64  `json:"threshold,omitempty"`
}

type AlertsLocationConditionList struct {
	AlertsLocationConditions []*AlertsLocationCondition `json:"location_failure_conditions,omitempty"`
}

type AlertsLocationCondition struct {
	ID         *int64          `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	RunbookURL *string         `json:"runbook_url,omitempty"`
	Enabled    *bool           `json:"enabled,omitempty"`
	Entities   []*string       `json:"entities,omitempty"`
	Terms      []*LocationTerm `json:"terms,omitempty"`
	TimeLimit  *int64          `json:"violation_time_limit_seconds,omitempty"`
}

func (s *locationConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_location_failure_conditions/policies/"+opt.PolicyIDOptions+".json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsLocationConditionList = new(AlertsLocationConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsLocationConditionList)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
