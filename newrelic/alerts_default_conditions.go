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

// API: https://rpm.newrelic.com/api/explore/alerts_conditions/list

type defaultConditions service

type AlertsDefaultConditionList struct {
	AlertsDefaultConditions []*AlertsDefaultCondition `json:"conditions,omitempty"`
}

// AlertsDefaultCondition manages Conditions of
// * APM Application
// * Key Transaction
// * Browser application metric
// * Mobile application metric
// for your alert policies
//
// NOTICE: API documents says Entities are integers, but it's string actually
type AlertsDefaultCondition struct {
	ID                  *int64                             `json:"id,omitempty"`
	Type                *string                            `json:"type,omitempty"`
	Name                *string                            `json:"name,omitempty"`
	Enabled             *bool                              `json:"enabled,omitempty"`
	Entities            []*string                          `json:"entities,omitempty"`
	Metric              *string                            `json:"metric,omitempty"`
	GCMetric            *string                            `json:"gc_metric,omitempty"`
	RunbookURL          *string                            `json:"runbook_url,omitempty"`
	ConditionScope      *string                            `json:"condition_scope,omitempty"`
	ViolationCloseTimer *int64                             `json:"violation_close_timer,omitempty"`
	Terms               []*AlertsConditionTerm             `json:"terms,omitempty"`
	UserDefined         *AlertsDefaultConditionUserDefined `json:"user_defined,omitempty"`
}

type AlertsDefaultConditionEntity struct {
	AlertsDefaultCondition *AlertsDefaultCondition `json:"condition,omitempty"`
}

type AlertsDefaultConditionUserDefined struct {
	Metric        *string `json:"metric,omitempty"`
	ValueFunction *string `json:"value_function,omitempty"`
}

func (s *defaultConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_conditions.json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsDefaultConditionList = new(AlertsDefaultConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsDefaultConditionList)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *defaultConditions) deleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_conditions/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *defaultConditions) create(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_conditions/policies/%v.json", policyID)
	if c.AlertsDefaultConditionEntity.AlertsDefaultCondition.ID != nil {
		c.AlertsDefaultConditionEntity.AlertsDefaultCondition.ID = nil
	}
	req, err := s.client.NewRequest("POST", u, c.AlertsDefaultConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsDefaultConditionEntity = new(AlertsDefaultConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsDefaultConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}

func (s *defaultConditions) update(ctx context.Context, c *AlertsConditionEntity, id int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_conditions/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, c.AlertsDefaultConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsDefaultConditionEntity = new(AlertsDefaultConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsDefaultConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}
