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

// API: https://rpm.newrelic.com/api/explore/alerts_external_service_conditions/list

type externalServiceConditions service

type AlertsExternalServiceConditionList struct {
	AlertsExternalServiceConditions []*AlertsExternalServiceCondition `json:"external_service_conditions,omitempty"`
}

type AlertsExternalServiceCondition struct {
	ID                 *int64                 `json:"id,omitempty"`
	Type               *string                `json:"type,omitempty"`
	Name               *string                `json:"name,omitempty"`
	Enabled            *bool                  `json:"enabled,omitempty"`
	Entities           []*string              `json:"entities,omitempty"`
	ExternalServiceURL *string                `json:"external_service_url,omitempty"`
	Metric             *string                `json:"metric,omitempty"`
	RunbookURL         *string                `json:"runbook_url,omitempty"`
	Terms              []*AlertsConditionTerm `json:"terms,omitempty"`
}

type AlertsExternalServiceConditionEntity struct {
	AlertsExternalServiceCondition *AlertsExternalServiceCondition `json:"external_service_condition,omitempty"`
}

func (s *externalServiceConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_external_service_conditions.json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsExternalServiceConditionList = new(AlertsExternalServiceConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsExternalServiceConditionList)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *externalServiceConditions) deleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_external_service_conditions/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *externalServiceConditions) create(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_external_service_conditions/policies/%v.json", policyID)
	if c.AlertsExternalServiceConditionEntity.AlertsExternalServiceCondition.ID != nil {
		c.AlertsExternalServiceConditionEntity.AlertsExternalServiceCondition.ID = nil
	}
	req, err := s.client.NewRequest("POST", u, c.AlertsExternalServiceConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsExternalServiceConditionEntity = new(AlertsExternalServiceConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsExternalServiceConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}

func (s *externalServiceConditions) update(ctx context.Context, c *AlertsConditionEntity, id int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_external_service_conditions/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, c.AlertsExternalServiceConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsExternalServiceConditionEntity = new(AlertsExternalServiceConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsExternalServiceConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}
