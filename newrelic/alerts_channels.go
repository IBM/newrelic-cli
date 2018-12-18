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
	"encoding/json"
	"fmt"
	"strconv"
)

var channelTypes = make(map[AlertsChannelType]func() interface{})

func init() {
	channelTypes[ChannelByUser] = func() interface{} { return new(ChannelUserConfig) }
	channelTypes[ChannelByEmail] = func() interface{} { return new(ChannelEmailConfig) }
	channelTypes[ChannelBySlack] = func() interface{} { return new(ChannelSlackConfig) }
	channelTypes[ChannelByPagerDuty] = func() interface{} { return new(ChannelPagerDutyConfig) }
	channelTypes[ChannelByWebhook] = func() interface{} { return new(ChannelWebhookConfig) }
	channelTypes[ChannelByCampfire] = func() interface{} { return new(ChannelCampfireConfig) }
	channelTypes[ChannelByHipChat] = func() interface{} { return new(ChannelHipChatConfig) }
	channelTypes[ChannelByOpsGenie] = func() interface{} { return new(ChannelOpsGenieConfig) }
	channelTypes[ChannelByVictorOps] = func() interface{} { return new(ChannelVictorOpsConfig) }
	channelTypes[ChannelByXMatters] = func() interface{} { return new(ChannelXMattersConfig) }
}

type AlertsChannelType string
type AlertsChannelWebhookType string

const (
	ChannelByUser      AlertsChannelType = "user"
	ChannelByEmail     AlertsChannelType = "email"
	ChannelBySlack     AlertsChannelType = "slack"
	ChannelByPagerDuty AlertsChannelType = "pagerduty"
	ChannelByWebhook   AlertsChannelType = "webhook"
	ChannelByCampfire  AlertsChannelType = "campfire"
	ChannelByHipChat   AlertsChannelType = "hipchat"
	ChannelByOpsGenie  AlertsChannelType = "opsgenie"
	ChannelByVictorOps AlertsChannelType = "victorops"
	ChannelByXMatters  AlertsChannelType = "xmatters"
	// TODO: complete possible options: e.g. xMatters

	ChannelWebhookByJSON AlertsChannelWebhookType = "application/json"
	ChannelWebhookByForm AlertsChannelWebhookType = "application/x-www-form-urlencoded"
)

// ChannelUserConfig is system generated alerts channel, it's non-editable
type ChannelUserConfig struct {
	UserID *string `json:"user_id,omitempty"`
}

// ChannelEmailConfig is struct of Email Alert Channel
type ChannelEmailConfig struct {
	Recipients            *string `json:"recipients,omitempty"`
	IncludeJSONAttachment *bool   `json:"include_json_attachment,omitempty"`
}

// ChannelSlackConfig is struct of Slack Alert Channel
type ChannelSlackConfig struct {
	URL     *string `json:"url,omitempty"`
	Channel *string `json:"channel,omitempty"`
}

// ChannelPagerDutyConfig is struct of Pagerduty Alert Channel
type ChannelPagerDutyConfig struct {
	ServiceKey *string `json:"service_key,omitempty"`
}

// ChannelWebhookConfig is struct of Webhook Alert Channel
type ChannelWebhookConfig struct {
	BaseURL      *string                  `json:"base_url,omitempty"`
	AuthUsername *string                  `json:"auth_username,omitempty"`
	AuthPassword *string                  `json:"auth_password,omitempty"`
	PayloadType  AlertsChannelWebhookType `json:"payload_type,omitempty"`
	Payload      *map[string]interface{}  `json:"payload,omitempty"`
	Headers      *map[string]string       `json:"headers,omitempty"`
}

type ChannelHipChatConfig struct {
	AuthToken *string `json:"auth_token,omitempty"`
	RoomId    *string `json:"room_id,omitempty"`
}

type ChannelOpsGenieConfig struct {
	ApiKey     *string `json:"api_key,omitempty"`
	Teams      *string `json:"teams,omitempty"`
	Tags       *string `json:"tags,omitempty"`
	Recipients *string `json:"recipients,omitempty"`
}

type ChannelCampfireConfig struct {
	Subdomain *string `json:"subdomain,omitempty"`
	Token     *string `json:"token,omitempty"`
	Room      *string `json:"room,omitempty"`
}

type ChannelVictorOpsConfig struct {
	Key      *string `json:"key,omitempty"`
	RouteKey *string `json:"route_key,omitempty"`
}

type ChannelXMattersConfig struct {
	URL     *string `json:"url,omitempty"`
	Channel *string `json:"channel,omitempty"`
}

type AlertsChannelsService service

type AlertsChannel struct {
	ID            *int64              `json:"id,omitempty"`
	Name          *string             `json:"name,omitempty"`
	Type          AlertsChannelType   `json:"type,omitempty"`
	Configuration interface{}         `json:"configuration,omitempty"`
	Links         *AlertsChannelLinks `json:"links,omitempty"`
}

// UnmarshalJSON customizes UnmarshalJSON method so that Configuration field
// can be properly unmarshlled
func (c *AlertsChannel) UnmarshalJSON(data []byte) error {
	var envelope struct {
		ID            *int64              `json:"id,omitempty"`
		Name          *string             `json:"name,omitempty"`
		Type          AlertsChannelType   `json:"type,omitempty"`
		Configuration json.RawMessage     `json:"configuration,omitempty"`
		Links         *AlertsChannelLinks `json:"links,omitempty"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	configure, ok := channelTypes[envelope.Type]
	if !ok {
		return fmt.Errorf("type %q is not supported", envelope.Type)
	}
	configuration := configure()
	// in some cases, NewRelic mask or totally remove some fields
	// which results in the Configuration be a nil
	// so we should bypass the json decode err here
	json.Unmarshal(envelope.Configuration, configuration)
	// if err := json.Unmarshal(envelope.Configuration, configuration); err != nil {
	// 	return err
	// }
	c.Configuration = configuration
	c.ID = envelope.ID
	c.Name = envelope.Name
	c.Type = envelope.Type
	c.Links = envelope.Links
	return nil
}

// AlertsChannelLinks holds the links to AlertsPolicies
type AlertsChannelLinks struct {
	PolicyIDs []*int64 `json:"policy_ids,omitempty"`
}

type AlertsChannelEntity struct {
	AlertsChannel *AlertsChannel `json:"channel,omitempty"`
}

type AlertsChannelList struct {
	AlertsChannels []*AlertsChannel `json:"channels,omitempty"`
	AlertsLinks    *AlertsLinks     `json:"links,omitempty"`
}

type AlertsLinks struct {
	ChannelPolicyIDs *string `json:"channel.policy_ids,omitempty"`
}

type AlertsChannelListOptions struct {
	PageOptions
}

type PolicyChannelsAssociation struct {
	PolicyID      *int64   `json:"policyId,omitempty"`
	ChannelIDList []*int64 `json:"channels,omitempty"`
}

func (s *AlertsChannelsService) ListAll(ctx context.Context, opt *AlertsChannelListOptions) (*AlertsChannelList, *Response, error) {
	u, err := addOptions("alerts_channels.json", opt)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	alertsChannelList := new(AlertsChannelList)
	resp, err := s.client.Do(ctx, req, alertsChannelList)
	if err != nil {
		return nil, resp, err
	}

	return alertsChannelList, resp, nil
}

// Create POST a AlertsChannelEntity to create
//
// maybe a potential bug in NewRelic side
// from API document, Create returns a AlertsChannel json payload
// but actually it returns a AlertsChannelList json
// {"channels":[{"id":1104874,"name":"newrelic-cli-integration-test-channel","type":"email","configuration":{"recipients":""},"links":{"policy_ids":[]}}],"links":{"channel.policy_ids":"/v2/policies/{policy_id}"}}
// so change returned parameter *AlertsChannelEntity to *AlertsChannelList
func (s *AlertsChannelsService) Create(ctx context.Context, c *AlertsChannelEntity) (*AlertsChannelList, *Response, error) {
	u := "alerts_channels.json"
	req, err := s.client.NewRequest("POST", u, c)
	if err != nil {
		return nil, nil, err
	}

	channelList := new(AlertsChannelList)
	resp, err := s.client.Do(ctx, req, channelList)
	if err != nil {
		return nil, resp, err
	}

	return channelList, resp, nil
}

func (s *AlertsChannelsService) DeleteByID(ctx context.Context, id int64) (*Response, error) {
	u := fmt.Sprintf("alerts_channels/%v.json", id)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

func (s *AlertsChannelsService) UpdatePolicyChannels(ctx context.Context, policyId int64, channelIds []*int64) (*Response, error) {
	u := "alerts_policy_channels.json"

	var channels string
	var channelsLen = len(channelIds)
	for index, channelID := range channelIds {
		if index < (channelsLen - 1) {
			channels = channels + strconv.FormatInt(*channelID, 10) + ","
		} else {
			channels = channels + strconv.FormatInt(*channelID, 10)
		}
	}

	var c = "policy_id=" + strconv.FormatInt(policyId, 10) + "&channel_ids=" + channels

	u = u + "?" + c

	req, err := s.client.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, err
	}

	alertChannels := new(PolicyChannelsAssociation)
	resp, err := s.client.Do(ctx, req, alertChannels)
	if err != nil {
		return resp, err
	}

	return resp, err
}
