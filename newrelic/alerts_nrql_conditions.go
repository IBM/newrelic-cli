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

type nrqlConditions service

type AlertsNRQLConditionList struct {
	AlertsNRQLConditions []*AlertsNRQLCondition `json:"nrql_conditions,omitempty"`
}

type AlertsNRQLCondition struct {
	Type                      *string                    `json:"type,omitempty"`
	ViolationTimeLimitSeconds *int64                     `json:"violation_time_limit_seconds"`
	ID                        *int64                     `json:"id,omitempty"`
	Name                      *string                    `json:"name,omitempty"`
	RunbookURL                *string                    `json:"runbook_url,omitempty"`
	Enabled                   *bool                      `json:"enabled,omitempty"`
	Terms                     []*AlertsConditionTerm     `json:"terms,omitempty"`
	ValueFunction             *string                    `json:"value_function,omitempty"`
	NRQL                      *AlertsNRQLConditionNRQL   `json:"nrql,omitempty"`
	Signal                    *AlertsNRQLConditionSignal `json:"signal,omitempty"`
}

type AlertsNRQLConditionSignal struct {
	AggregationWindow *string `json:"aggregation_window,omitempty"`
	EvaluationOffset  *string `json:"evaluation_offset,omitempty"`
	FillOption        *string `json:"fill_option,omitempty"`
	FillValue         *string `json:"fill_value,omitempty"`
}

type AlertsNRQLConditionNRQL struct {
	Query      *string `json:"query,omitempty"`
	SinceValue *string `json:"since_value,omitempty"`
}

type AlertsNRQLConditionEntity struct {
	AlertsNRQLCondition *AlertsNRQLCondition `json:"nrql_condition,omitempty"`
}

func (s *nrqlConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_nrql_conditions.json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsNRQLConditionList = new(AlertsNRQLConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsNRQLConditionList)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *nrqlConditions) deleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_nrql_conditions/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *nrqlConditions) create(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_nrql_conditions/policies/%v.json", policyID)
	if c.AlertsNRQLConditionEntity.AlertsNRQLCondition.ID != nil {
		c.AlertsNRQLConditionEntity.AlertsNRQLCondition.ID = nil
	}
	req, err := s.client.NewRequest("POST", u, c.AlertsNRQLConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsNRQLConditionEntity = new(AlertsNRQLConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsNRQLConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}

func (s *nrqlConditions) update(ctx context.Context, c *AlertsConditionEntity, id int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_nrql_conditions/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, c.AlertsNRQLConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsNRQLConditionEntity = new(AlertsNRQLConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsNRQLConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}
