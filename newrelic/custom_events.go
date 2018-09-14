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

type CustomEventService service

func (s *CustomEventService) Insert(ctx context.Context, insertKey string, accountID string, jsonData string) (*Response, []byte, error) {
	u := fmt.Sprintf("%v/events", accountID)
	req, err := s.client.NewRequestForNonJSON("POST", u, jsonData)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Insert-Key", insertKey)

	resp, bytes, err := s.client.DoWithBytes(ctx, req)
	if err != nil {
		return resp, nil, err
	}

	return resp, bytes, nil
}
