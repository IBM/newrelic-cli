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

type pluginsConditions service

type AlertsPluginsConditionList struct {
	AlertsPluginsConditions []*AlertsPluginsCondition `json:"plugins_conditions,omitempty"`
}

type AlertsPluginsCondition struct {
	ID                *int64                        `json:"id,omitempty"`
	Name              *string                       `json:"name,omitempty"`
	Enabled           *bool                         `json:"enabled,omitempty"`
	Entities          []*string                     `json:"entities,omitempty"`
	MetricDescription *string                       `json:"metric_description,omitempty"`
	Metric            *string                       `json:"metric,omitempty"`
	ValueFunction     *string                       `json:"value_function,omitempty"`
	RunbookURL        *string                       `json:"runbook_url,omitempty"`
	Terms             []*AlertsConditionTerm        `json:"terms,omitempty"`
	Plugin            *AlertsPluginsConditionPlugin `json:"plugin,omitempty"`
}

type AlertsPluginsConditionPlugin struct {
	ID   *string `json:"id,omitempty"`
	GUID *string `json:"guid,omitempty"`
}

type AlertsPluginsConditionEntity struct {
	AlertsPluginsCondition *AlertsPluginsCondition `json:"plugins_condition,omitempty"`
}

func (s *pluginsConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("alerts_plugins_conditions.json", opt)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	list.AlertsPluginsConditionList = new(AlertsPluginsConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsPluginsConditionList)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *pluginsConditions) deleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_plugins_conditions/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *pluginsConditions) create(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_plugins_conditions/policies/%v.json", policyID)
	if c.AlertsPluginsConditionEntity.AlertsPluginsCondition.ID != nil {
		c.AlertsPluginsConditionEntity.AlertsPluginsCondition.ID = nil
	}
	req, err := s.client.NewRequest("POST", u, c.AlertsPluginsConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsPluginsConditionEntity = new(AlertsPluginsConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsPluginsConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}

func (s *pluginsConditions) update(ctx context.Context, c *AlertsConditionEntity, id int64) (*AlertsConditionEntity, *Response, error) {
	u := fmt.Sprintf("alerts_plugins_conditions/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, c.AlertsPluginsConditionEntity)
	if err != nil {
		return nil, nil, err
	}

	condition := new(AlertsConditionEntity)
	condition.AlertsPluginsConditionEntity = new(AlertsPluginsConditionEntity)
	resp, err := s.client.Do(ctx, req, condition.AlertsPluginsConditionEntity)
	if err != nil {
		return nil, resp, err
	}

	return condition, resp, nil
}
