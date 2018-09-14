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

type AlertsViolationService service

type AlertsViolationList struct {
	Violations []*AlertsViolation        `json:"violations,omitempty"`
	Links      *AlertsViolationListLinks `json:"links,omitempty"`
}

type AlertsViolation struct {
	ID            *int64                 `json:"id,omitempty"`
	Label         *string                `json:"label,omitempty"`
	Duration      *int64                 `json:"duration,omitempty"`
	PolicyName    *string                `json:"policy_name,omitempty"`
	ConditionName *string                `json:"condition_name,omitempty"`
	Priority      *string                `json:"priority,omitempty"`
	OpenedAt      *int64                 `json:"opened_at,omitempty"`
	Entity        *AlertsViolationEntity `json:"entity,omitempty"`
	Links         *AlertsViolationLinks  `json:"links,omitempty"`
	ClosedAt      *int64                 `json:"closed_at,omitempty"`
}

type AlertsViolationEntity struct {
	Product *string `json:"product,omitempty"`
	Type    *string `json:"type,omitempty"`
	GroupID *int64  `json:"group_id,omitempty"`
	ID      *int64  `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type AlertsViolationLinks struct {
	PolicyID    *int64 `json:"policy_id,omitempty"`
	ConditionID *int64 `json:"condition_id,omitempty"`
}

type AlertsViolationListLinks struct {
	ViolationPolicyID    *string `json:"violation.policy_id,omitempty"`
	ViolationConditionID *string `json:"violation.condition_id,omitempty"`
}

type AlertsViolationListOptions struct {
	PageOptions
	OnlyOpen  bool   `url:"only_open,omitempty"`
	StartDate string `url:"start_date,omitempty"`
	EndDate   string `url:"end_date,omitempty"`
}

func (s *AlertsViolationService) ListAll(ctx context.Context, opt *AlertsViolationListOptions) (*AlertsViolationList, *Response, error) {
	u, err := addOptions("alerts_violations.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alertsViolationList := new(AlertsViolationList)
	resp, err := s.client.Do(ctx, req, alertsViolationList)
	if err != nil {
		return nil, nil, err
	}

	return alertsViolationList, resp, nil
}
