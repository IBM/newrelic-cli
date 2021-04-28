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
	"fmt"
)

type syntheticsConditions service

type AlertsSyntheticsConditionList struct {
	AlertsSyntheticsConditions []*AlertsSyntheticsCondition `json:"synthetics_conditions,omitempty"`
}

type AlertsSyntheticsCondition struct {
	ID         *int64  `json:"id,omitempty"`
	Name       *string `json:"name,omitempty"`
	MonitorID  *string `json:"monitor_id,omitempty"`
	RunbookURL *string `json:"runbook_url,omitempty"`
	Enabled    *bool   `json:"enabled,omitempty"`
}

type AlertsSyntheticsConditionEntity struct {
	AlertsSyntheticsCondition *AlertsSyntheticsCondition `json:"synthetics_condition,omitempty"`
}

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

func (s *syntheticsConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_synthetics_conditions.json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsSyntheticsConditionList = new(AlertsSyntheticsConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsSyntheticsConditionList)
	if err != nil {
		return resp, err
	}

	u = "alerts_location_failure_conditions/policies/" + opt.PolicyIDOptions + ".json"
	req, err = s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsLocationConditionList = new(AlertsLocationConditionList)
	resp, err = s.client.Do(ctx, req, list.AlertsLocationConditionList)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *syntheticsConditions) deleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_synthetics_conditions/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *syntheticsConditions) create(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_synthetics_conditions/policies/%v.json", policyID)
	if c.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.ID != nil {
		c.AlertsSyntheticsConditionEntity.AlertsSyntheticsCondition.ID = nil
	}
	req, err := s.client.NewRequest("POST", u, c.AlertsSyntheticsConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsSyntheticsConditionEntity = new(AlertsSyntheticsConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsSyntheticsConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}

func (s *syntheticsConditions) update(ctx context.Context, c *AlertsConditionEntity, id int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_synthetics_conditions/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, c.AlertsSyntheticsConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsSyntheticsConditionEntity = new(AlertsSyntheticsConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsSyntheticsConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}
