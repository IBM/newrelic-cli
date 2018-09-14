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

// NewRelic API docs: https://docs.newrelic.com/docs/apis/rest-api-v2/labels-examples-v2/list-labels-v2

// LabelsService handles communication with label
// related method of the NewRelic API

type LabelsService service

type LabelEntity struct {
	Label *Label `json:"label,omitempty"`
}

type LabelList struct {
	Labels         []*Label        `json:"labels,omitempty"`
	LabelListLinks *LabelListLinks `json:"links,omitempty"`
}

type Label struct {
	Key                           *string             `json:"key,omitempty"`
	Category                      *string             `json:"category,omitempty"`
	Name                          *string             `json:"name,omitempty"`
	LabelsApplicationHealthStatus *LabelsHealthStatus `json:"application_health_status,omitempty"`
	LabelsServerHealthStatus      *LabelsHealthStatus `json:"server_health_status,omitempty"`
	LabelLinks                    *LabelLinks         `json:"links,omitempty"`
}

type LabelsHealthStatus struct {
	Green  []*int64 `json:"green,omitempty"`
	Orange []*int64 `json:"orange,omitempty"`
	Red    []*int64 `json:"red,omitempty"`
	Gray   []*int64 `json:"gray,omitempty"`
}

type LabelLinks struct {
	Applications []*int64 `json:"applications,omitempty"`
	Servers      []*int64 `json:"servers,omitempty"`
}

type LabelListLinks struct {
	LabelApplications *string `json:"label.applications,omitempty"`
	LabelServers      *string `json:"label.servers,omitempty"`
	LabelServer       *string `json:"label.server,omitempty"`
}

type LabelListOptions struct {
	PageOptions
}

func (s *LabelsService) ListAll(ctx context.Context, opt *LabelListOptions) (*LabelList, *Response, error) {
	u, err := addOptions("labels.json", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	labelList := new(LabelList)
	resp, err := s.client.Do(ctx, req, labelList)
	if err != nil {
		return nil, nil, err
	}

	return labelList, resp, nil
}

func (s *LabelsService) Create(ctx context.Context, l *LabelEntity) (*LabelEntity, *Response, error) {
	u := "labels.json"
	req, err := s.client.NewRequest("PUT", u, l)
	if err != nil {
		return nil, nil, err
	}

	labelEntity := new(LabelEntity)
	resp, err := s.client.Do(ctx, req, labelEntity)
	if err != nil {
		return nil, resp, err
	}

	return labelEntity, resp, nil
}

func (s *LabelsService) DeleteByKey(ctx context.Context, key string) (*Response, error) {
	u := fmt.Sprintf("labels/%v.json", key)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
