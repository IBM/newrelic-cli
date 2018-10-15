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

type Dashboard struct {
	ID          *int64  `json:"id,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	//created_at
	//updated_at
	Visibility *string `json:"visibility,omitempty"`
	Editable   *string `json:"editable,omitempty"`
	UIURL      *string `json:"ui_url,omitempty"`
	APIURL     *string `json:"api_url,omitempty"`
	OwnerEmail *string `json:"owner_email,omitempty"`
}

type DashboardList struct {
	Dashboards []*Dashboard `json:"dashboards,omitempty"`
}

type DashboardListOptions struct {
	PageOptions
}

type DashboardService service

type CreateDashboardResponse struct {
	Dashboard *Dashboard `json:"dashboard,omitempty"`
}

func (s *DashboardService) ListAll(ctx context.Context, opt *DashboardListOptions) (*Response, []byte, error) {
	u, err := addOptions("dashboards.json", opt)

	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}

func (s *DashboardService) GetByID(ctx context.Context, id int64) (*Response, []byte, error) {
	u := fmt.Sprintf("dashboards/%v.json", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}

func (s *DashboardService) DeleteByID(ctx context.Context, id int64) (*Response, []byte, error) {
	u := fmt.Sprintf("dashboards/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}

func (s *DashboardService) Create(ctx context.Context, dashboard string) (*Response, []byte, error) {
	u := "dashboards.json"
	req, err := s.client.NewRequestForNonJSON("POST", u, dashboard)
	if err != nil {
		return nil, nil, err
	}

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}

func (s *DashboardService) Update(ctx context.Context, dashboard string, id int64) (*Response, []byte, error) {
	u := fmt.Sprintf("dashboards/%v.json", id)
	req, err := s.client.NewRequestForNonJSON("PUT", u, dashboard)
	if err != nil {
		return nil, nil, err
	}

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}
