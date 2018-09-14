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

import "context"

type AlertsIncidentService service

type AlertsIncidentList struct {
	Incidents []*AlertsIncident        `json:"incidents,omitempty"`
	Links     *AlertsIncidentListLinks `json:"links,omitempty"`
}

type AlertsIncident struct {
	ID                 *int64               `json:"id,omitempty"`
	OpenedAt           *int64               `json:"opened_at,omitempty"`
	IncidentPreference *string              `json:"incident_preference,omitempty"`
	Links              *AlertsIncidentLinks `json:"links,omitempty"`
}

type AlertsIncidentLinks struct {
	Violations []*int64 `json:"violations,omitempty"`
	PolicyID   *int64   `json:"policy_id,omitempty"`
}

type AlertsIncidentListLinks struct {
	IncidentPolicyID *string `json:"incident.policy_id,omitempty"`
}

type AlertsIncidentListOptions struct {
	PageOptions
	OnlyOpen bool `url:"only_open,omitempty"`
}

func (s *AlertsIncidentService) ListAll(ctx context.Context, opt *AlertsIncidentListOptions) (*AlertsIncidentList, *Response, error) {
	u, err := addOptions("alerts_incidents.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alertsIncidentList := new(AlertsIncidentList)
	resp, err := s.client.Do(ctx, req, alertsIncidentList)
	if err != nil {
		return nil, nil, err
	}

	return alertsIncidentList, resp, nil
}
