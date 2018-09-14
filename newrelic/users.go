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

// UsersService handles communication with the user related
// methods of the NewRelic API.
//
// NewRelic API docs: https://docs.newrelic.com/docs/apis/rest-api-v2/account-examples-v2/listing-users-your-account
type UsersService service

// User represents a NewRelic user.
type User struct {
	ID        *int64  `json:"id,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	Role      *string `json:"role,omitempty"`
}

// UserEntity corresponds to a JSON payload returned by NewRelic API
type UserEntity struct {
	User *User `json:"user,omitempty"`
}

// UserList represents a collection of NewRelic users.
type UserList struct {
	Users []*User `json:"users,omitempty"`
}

// UserListOptions specifies optional parameters to the UsersService.ListAll
// method.
type UserListOptions struct {
	IDOptions    string `url:"filter[ids],omitempty"`
	EmailOptions string `url:"filter[email],omitempty"`

	PageOptions
}

// ListAll retruns all NewRelic users under current account/NEW_RELIC_API_KEY
// When given `opt`, it returns fitlered resultset
func (s *UsersService) ListAll(ctx context.Context, opt *UserListOptions) (*UserList, *Response, error) {
	u, err := addOptions("users.json", opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	userList := new(UserList)
	resp, err := s.client.Do(ctx, req, userList)
	if err != nil {
		return nil, resp, err
	}

	return userList, resp, nil
}

// GetByID returns specfic NewRelic user by given `id`
func (s *UsersService) GetByID(ctx context.Context, id int64) (*UserEntity, *Response, error) {
	u := fmt.Sprintf("users/%v.json", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	userEntity := new(UserEntity)
	resp, err := s.client.Do(ctx, req, userEntity)
	if err != nil {
		return nil, resp, err
	}

	return userEntity, resp, nil
}
