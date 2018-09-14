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

const (
	// https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#policies-create
	IncidentPerPolicy            IncidentPreferenceOption = "PER_POLICY"
	IncidentPerCondition         IncidentPreferenceOption = "PER_CONDITION"
	IncidentPerConditionAndTaret IncidentPreferenceOption = "PER_CONDITION_AND_TARGET"
)

type IncidentPreferenceOption string

type AlertsPoliciesService service

type AlertsPolicy struct {
	ID                 *int64                   `json:"id,omitempty"`
	IncidentPreference IncidentPreferenceOption `json:"incident_preference,omitempty"`
	Name               *string                  `json:"name,omitempty"`
	CreatedAt          *int64                   `json:"created_at,omitempty"`
	UpdatedAt          *int64                   `json:"updated_at,omitempty"`
}

type AlertsPolicyEntity struct {
	AlertsPolicy *AlertsPolicy `json:"policy,omitempty"`
}

type AlertsPolicyList struct {
	AlertsPolicies []*AlertsPolicy `json:"policies,omitempty"`
}

type AlertsPolicyListOptions struct {
	NameOptions string `url:"filter[name],omitempty"`

	PageOptions
}

func (s *AlertsPoliciesService) ListAll(ctx context.Context, opt *AlertsPolicyListOptions) (*AlertsPolicyList, *Response, error) {
	u, err := addOptions("alerts_policies.json", opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alertsPolicyList := new(AlertsPolicyList)
	resp, err := s.client.Do(ctx, req, alertsPolicyList)
	if err != nil {
		return nil, resp, err
	}

	return alertsPolicyList, resp, nil
}

// Create tries to create an alerts policy
// CAVEAT: it's good practice to check if a `p` existed
// If a `p` is "Created" twice, two alerts policies will be created with same Name but different ID
func (s *AlertsPoliciesService) Create(ctx context.Context, p *AlertsPolicyEntity) (*AlertsPolicyEntity, *Response, error) {
	u := "alerts_policies.json"
	req, err := s.client.NewRequest("POST", u, p)
	if err != nil {
		return nil, nil, err
	}

	policy := new(AlertsPolicyEntity)
	resp, err := s.client.Do(ctx, req, policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (s *AlertsPoliciesService) Update(ctx context.Context, p *AlertsPolicyEntity, id int64) (*AlertsPolicyEntity, *Response, error) {
	u := fmt.Sprintf("alerts_policies/%v.json", id)
	req, err := s.client.NewRequest("PUT", u, p)
	if err != nil {
		return nil, nil, err
	}

	policy := new(AlertsPolicyEntity)
	resp, err := s.client.Do(ctx, req, policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (s *AlertsPoliciesService) DeleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_policies/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
