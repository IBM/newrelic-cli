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

type LabelsSyntheticsService service

type MonitorRef struct {
	ID   *string `json:"id,omitempty"`
	HREF *string `json:"href,omitempty"`
}

type MonitorRefList struct {
	MonitorRefs []*MonitorRef `json:"monitorRefs,omitempty"`
}

// type PagedDataEntity struct {
// 	PagedData *MonitorRefList `json:"pagedData,omitempty"`
// }

type Metadata struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

type LabelSynthetics struct {
	PagedData *MonitorRefList `json:"pagedData,omitempty"`
	Metadata  *Metadata       `json:"metadata,omitempty"`
}

type MonitorLabel struct {
	Category *string `json:"id,omitempty"`
	Label    *string `json:"href,omitempty"`
}

func (s *LabelsSyntheticsService) GetMonitorsByLabel(ctx context.Context, opt *PageLimitOptions, label string) (*LabelSynthetics, *Response, error) {
	u := fmt.Sprintf("labels/%v", label)
	u, err := addOptions(u, opt)

	// if err != nil {
	// 	return nil, nil, err
	// }
	//DEBUG
	// var u = "labels/Monitor:PagerDuty"
	//DEBUG

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	labelSynthetics := new(LabelSynthetics)
	resp, err := s.client.Do(ctx, req, labelSynthetics)
	if err != nil {
		return nil, nil, err
	}

	return labelSynthetics, resp, nil
}

func (s *LabelsSyntheticsService) AddLabelToMonitor(ctx context.Context, monitorId string, monitorLabel *MonitorLabel) (*Response, error) {
	u := fmt.Sprintf("%v/labels", monitorId)

	var label string
	label = *monitorLabel.Category + ":" + *monitorLabel.Label
	req, err := s.client.NewRequestForNonJSON("POST", u, label)
	if err != nil {
		return nil, err
	}

	labelSynthetics := new(LabelSynthetics)
	resp, err := s.client.Do(ctx, req, labelSynthetics)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *LabelsSyntheticsService) DeleteLabelFromMonitor(ctx context.Context, monitorId string, label string) (*Response, error) {
	u := fmt.Sprintf("%v/labels/%v", monitorId, label)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, label)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
