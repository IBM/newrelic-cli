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
	"path"
)

const (
	monitorTagsQueryString = "query ($CURSOR: String) {\n  actor {\n" +
		"    entitySearch(queryBuilder: {type: MONITOR}) {\n" +
		"      count\n" +
		"      results(cursor: $CURSOR) {\n" +
		"        nextCursor\n" +
		"        entities {\n" +
		"          tags {\n" +
		"            key\n" +
		"            values\n" +
		"          }\n" +
		"          ... on SyntheticMonitorEntityOutline {\n" +
		"            guid\n" +
		"            name\n" +
		"            monitorId\n" +
		"          }\n" +
		"        }\n" +
		"      }\n" +
		"    }\n" +
		"  }\n" +
		"}\n"
	// monitorTagsQueryString = "{\n  actor {\n    entitySearch(query: \"type IN (\u0027MONITOR\u0027)\") {\n      count\n    }\n  }\n}\n"
)

type Monitor struct {
	ID           *string        `json:"id,omitempty"`
	Name         *string        `json:"name,omitempty"`
	Type         *string        `json:"type,omitempty"`
	Frequency    *int64         `json:"frequency,omitempty"`
	URI          *string        `json:"uri,omitempty"`
	Locations    []*string      `json:"locations,omitempty"`
	Status       *string        `json:"status,omitempty"`
	SLAThreshold *float64       `json:"slaThreshold,omitempty"`
	UserID       *int64         `json:"userId,omitempty"`
	ApiVersion   *string        `json:"apiVersion,omitempty"`
	CreatedAt    *string        `json:"createdAt,omitempty"`
	UpdatedAt    *string        `json:"modifiedAt,omitempty"`
	Options      MonitorOptions `json:"options,omitempty"`
	Script       *Script        `json:"script,omitempty"`
	Labels       []*string      `json:"labels,omitempty"`
	Tags         []*Tag         `json:"tags,omitempty"`
}

type MonitorTagsGraphqlVar struct {
	Cursor *string `json:"CURSOR,omitempty"`
}

type Tag struct {
	Key    *string   `json:"key"`
	Values []*string `json:"values"`
}

type MonitorTagsQueryBody struct {
	Query     *string                `json:"query"`
	Variables *MonitorTagsGraphqlVar `json:"variables,omitempty"`
}

type MonitorOptions struct {
	ValidationString       *string `json:"validationString,omitempty"`
	VerifySSL              bool    `json:"verifySSL,omitempty"`
	BypassHEADRequest      bool    `json:"bypassHEADRequest,omitempty"`
	TreatRedirectAsFailure bool    `json:"treatRedirectAsFailure,omitempty"`
}

type MonitorList struct {
	Monitors []*Monitor `json:"monitors,omitempty"`
}

type PageLimitOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

type MonitorListOptions struct {
	PageOptions
	PageLimitOptions
}

type MonitorTagsResp struct {
	Data *Data `json:"data"`
}

type Data struct {
	Actor *Actor `json:"actor"`
}

type Actor struct {
	EntitySearch *EntitySearch `json:"entitySearch"`
}

type EntitySearch struct {
	Count   *int                 `json:"count"`
	Results *EntitySearchResults `json:"results"`
}

type EntitySearchResults struct {
	Entities   []*EntitySearchResultsMonitor `json:"entities"`
	NextCursor *string                       `json:"nextCursor"`
}

type EntitySearchResultsMonitor struct {
	Guid      *string `json:"guid"`
	MonitorId *string `json:"monitorId"`
	Name      *string `json:"name"`
	Tags      []*Tag  `json:"tags"`
}

type SyntheticsService service

func (s *SyntheticsService) ListAll(ctx context.Context, opt *MonitorListOptions) (*MonitorList, *Response, error) {
	u, err := addOptions("", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	monitorList := new(MonitorList)

	resp, err := s.client.Do(ctx, req, monitorList)
	if err != nil {
		return nil, resp, err
	}
	return monitorList, resp, nil
}

func (s *SyntheticsService) ListTags(ctx context.Context, cursor *string) (*MonitorTagsResp, *Response, error) {
	var q string = monitorTagsQueryString
	var body *MonitorTagsQueryBody = &MonitorTagsQueryBody{
		Query: &q}

	if cursor != nil {
		body.Variables = &MonitorTagsGraphqlVar{
			Cursor: cursor}
	}

	req, err := s.client.NewRequest("POST", "", body)

	if err != nil {
		return nil, nil, err
	}

	monitorTags := new(MonitorTagsResp)

	resp, err := s.client.Do(ctx, req, monitorTags)
	if err != nil {
		return nil, resp, err
	}
	return monitorTags, resp, nil
}

func (s *SyntheticsService) GetByID(ctx context.Context, id string) (*Monitor, *Response, error) {
	req, err := s.client.NewRequest("GET", id, nil)
	if err != nil {
		return nil, nil, err
	}

	monitor := new(Monitor)
	resp, err := s.client.Do(ctx, req, monitor)
	if err != nil {
		return nil, resp, err
	}

	return monitor, resp, nil
}

func (s *SyntheticsService) Create(ctx context.Context, monitor *Monitor) (*Monitor, *Response, error) {
	req, err := s.client.NewRequest("POST", "", monitor)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(ctx, req, monitor)
	if err != nil {
		return nil, resp, err
	}
	monitorUUID := path.Base(resp.Header.Get("Location"))
	monitor.ID = &monitorUUID
	return monitor, resp, nil
}

func (s *SyntheticsService) Update(ctx context.Context, monitor *Monitor, id *string) (*Response, error) {
	req, err := s.client.NewRequest("PUT", *id, monitor)
	if err != nil {
		return nil, err
	}

	m := new(Monitor)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SyntheticsService) DeleteByID(ctx context.Context, id *string) (*Response, error) {
	req, err := s.client.NewRequest("DELETE", *id, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *SyntheticsService) Patch(ctx context.Context, monitor *Monitor, id *string) (*Response, error) {
	req, err := s.client.NewRequest("Patch", *id, monitor)
	if err != nil {
		return nil, err
	}

	m := new(Monitor)
	resp, err := s.client.Do(ctx, req, m)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
