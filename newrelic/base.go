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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL     = "https://api.newrelic.com/v2/"
	syntheticsURL      = "https://synthetics.newrelic.com/synthetics/api/v3/monitors/"
	userAgent          = "go-newrelic"
	labelSyntheticsURL = "https://synthetics.newrelic.com/synthetics/api/v4/monitors/"
	insightsURL        = "https://insights-collector.newrelic.com/v1/accounts/"
	infrastructureURL  = "https://infra-api.newrelic.com/v2/alerts/"
	graphqlURL         = "https://api.newrelic.com/graphql"
)

type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string
	XApiKey   string
	ProxyAuth string
	Retries   int

	common service

	Users              *UsersService
	AlertsPolicies     *AlertsPoliciesService
	AlertsChannels     *AlertsChannelsService
	Labels             *LabelsService
	AlertsIncidents    *AlertsIncidentService
	AlertsViolations   *AlertsViolationService
	AlertsEvents       *AlertsEventService
	AlertsConditions   *AlertsConditionsService
	SyntheticsMonitors *SyntheticsService
	SyntheticsScript   *ScriptService
	LabelsSynthetics   *LabelsSyntheticsService
	Dashboards         *DashboardService
	CustomEvents       *CustomEventService
}

type service struct {
	client *Client
}

type Response struct {
	*http.Response

	NextPage  int
	PrePage   int
	FirstPage int
	LastPage  int
}

func (r *Response) String() string {
	if r == nil {
		return "nil"
	}
	return r.Status
}

type PageOptions struct {
	Page int `url:"page,omitempty"`
}

func NewClient(httpClient *http.Client, endpointType string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	var baseURL *url.URL

	switch endpointType {
	case "synthetics":
		baseURL, _ = url.Parse(syntheticsURL)
	case "labelSynthetics":
		baseURL, _ = url.Parse(labelSyntheticsURL)
	case "insights":
		baseURL, _ = url.Parse(insightsURL)
	case "infrastructure":
		baseURL, _ = url.Parse(infrastructureURL)
	case "graphql":
		baseURL, _ = url.Parse(graphqlURL)
	default:
		baseURL, _ = url.Parse(defaultBaseURL)
	}

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.common.client = c
	c.Users = (*UsersService)(&c.common)
	c.AlertsPolicies = (*AlertsPoliciesService)(&c.common)
	c.AlertsChannels = (*AlertsChannelsService)(&c.common)
	c.Labels = (*LabelsService)(&c.common)
	c.AlertsIncidents = (*AlertsIncidentService)(&c.common)
	c.AlertsViolations = (*AlertsViolationService)(&c.common)
	c.AlertsEvents = (*AlertsEventService)(&c.common)

	c.AlertsConditions = &AlertsConditionsService{}
	c.AlertsConditions.defaultConditions = (*defaultConditions)(&c.common)
	c.AlertsConditions.pluginsConditions = (*pluginsConditions)(&c.common)
	c.AlertsConditions.externalServiceConditions = (*externalServiceConditions)(&c.common)
	c.AlertsConditions.syntheticsConditions = (*syntheticsConditions)(&c.common)
	c.AlertsConditions.nrqlConditions = (*nrqlConditions)(&c.common)
	c.AlertsConditions.infraConditions = (*infraConditions)(&c.common)

	c.SyntheticsMonitors = (*SyntheticsService)(&c.common)
	c.SyntheticsScript = (*ScriptService)(&c.common)
	c.LabelsSynthetics = (*LabelsSyntheticsService)(&c.common)

	c.Dashboards = (*DashboardService)(&c.common)

	c.CustomEvents = (*CustomEventService)(&c.common)

	c.Retries = 3

	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	/* 	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	} */
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.XApiKey != "" {
		req.Header.Set("X-Api-Key", c.XApiKey)
	}
	if c.ProxyAuth != "" {
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.ProxyAuth))
		req.Header.Set("Proxy-Authorization", basic)
	}
	return req, nil
}

func (c *Client) NewRequestForNonJSON(method, urlStr string, body string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.XApiKey != "" {
		req.Header.Set("X-Api-Key", c.XApiKey)
	}
	if c.ProxyAuth != "" {
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.ProxyAuth))
		req.Header.Set("Proxy-Authorization", basic)
	}
	return req, nil
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	var retries = c.Retries

	var resp *http.Response
	var err error
	for retries > 0 {
		resp, err = c.client.Do(req)
		if err != nil || resp.StatusCode > 299 {
			if retries == 1 {
				if err == nil {
					bodyBytes, _ := ioutil.ReadAll(resp.Body)
					bodyString := truncateString(string(bodyBytes), 100)
					err = errors.New(fmt.Sprintf("%s %s returns status:%d body: %s", req.Method, req.URL.String(), resp.StatusCode, bodyString))
				}
				log.Println(err)
				break
			}
			time.Sleep(time.Duration(3) * time.Second)
			retries--
		} else {
			break
		}
	}
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

func (c *Client) DoWithBytes(ctx context.Context, req *http.Request) (*Response, []byte, error) {
	var retries = c.Retries

	var resp *http.Response
	var err error
	for retries > 0 {
		resp, err = c.client.Do(req)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(3) * time.Second)
			retries--
		} else {
			break
		}
	}
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}

		return nil, nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	var retBytes []byte
	retBytes, err = ioutil.ReadAll(resp.Body)

	return response, retBytes, err
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
