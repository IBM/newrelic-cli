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

type AlertsEventService service

type AlertsEventList struct {
	RecentEvents []*RecentEvent `json:"recent_events,omitempty"`
}

type RecentEvent struct {
	ID            *int64  `json:"id,omitempty"`
	EventType     *string `json:"event_type,omitempty"`
	Description   *string `json:"description,omitempty"`
	Timestamp     *int64  `json:"timestamp,omitempty"`
	IncidentID    *int64  `json:"incident_id,omitempty"`
	Product       *string `json:"product,omitempty"`
	EntityType    *string `json:"entity_type,omitempty"`
	EntityGroupID *int64  `json:"entity_group_id,omitempty"`
	EntityID      *int64  `json:"entity_id,omitempty"`
	Priority      *string `json:"priority,omitempty"`
}

type AlertsEventListOptions struct {
	PageOptions
	Product       string `url:"filter[product],omitempty"`
	EntityType    string `url:"filter[entity_type],omitempty"`
	EntityGroupID int64  `url:"filter[entity_group_id],omitempty"`
	EntityID      int64  `url:"filter[entity_id],omitempty"`
	EventType     string `url:"filter[event_type],omitempty"`
	IncidentID    int64  `url:"filter[incident_id],omitempty"`
}

func (s *AlertsEventService) ListAll(ctx context.Context, opt *AlertsEventListOptions) (*AlertsEventList, *Response, error) {
	u, err := addOptions("alerts_events.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alertsEventList := new(AlertsEventList)
	resp, err := s.client.Do(ctx, req, alertsEventList)
	if err != nil {
		return nil, nil, err
	}

	return alertsEventList, resp, nil
}
