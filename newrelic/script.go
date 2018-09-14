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

type Script struct {
	ScriptText *string `json:"scriptText,omitempty"`
}

type ScriptService service

func (s *ScriptService) GetByID(ctx context.Context, id string) (*Script, *Response, error) {
	u := fmt.Sprintf("%v/script/", id)
	//fmt.Println("path : ", u)

	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	script := new(Script)
	resp, err := s.client.Do(ctx, req, script)
	if err != nil {
		return nil, resp, err
	}

	return script, resp, nil
}

func (s *ScriptService) UpdateByID(ctx context.Context, scriptText *Script, id string) (*Response, error) {
	u := fmt.Sprintf("%v/script/", id)
	req, err := s.client.NewRequest("PUT", u, scriptText)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, scriptText)
	return resp, err
}
