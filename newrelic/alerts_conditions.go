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

const (
	// https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts
	ConditionDefault         ConditionCategory = "condition"
	ConditionPlugins         ConditionCategory = "plugins_condition"
	ConditionExternalService ConditionCategory = "external_service_condition"
	ConditionSynthetics      ConditionCategory = "synthetics_condition"
	ConditionNRQL            ConditionCategory = "nrql_condition"
	ConditionInfrastructure  ConditionCategory = "infrastructure_condition"
)

type ConditionCategory string

type AlertsConditionsService struct {
	*defaultConditions
	*pluginsConditions
	*externalServiceConditions
	*syntheticsConditions
	*nrqlConditions
	*infraConditions
}

type AlertsConditionsOptions struct {
	PolicyIDOptions string `url:"policy_id,omitempty"`

	PageOptions
}

type AlertsConditionList struct {
	*AlertsDefaultConditionList
	*AlertsExternalServiceConditionList
	*AlertsNRQLConditionList
	*AlertsPluginsConditionList
	*AlertsSyntheticsConditionList
	*AlertsInfrastructureConditionList
}

type AlertsConditionEntity struct {
	*AlertsDefaultConditionEntity
	*AlertsExternalServiceConditionEntity
	*AlertsNRQLConditionEntity
	*AlertsPluginsConditionEntity
	*AlertsSyntheticsConditionEntity
	*AlertsInfrastructureConditionEntity
}

func (s *AlertsConditionsService) ListAll(ctx context.Context, opt *AlertsConditionsOptions) (*AlertsConditionList, error) {
	if opt == nil || opt.PolicyIDOptions == "" {
		return nil, fmt.Errorf("policy_id is required")
	}

	list := new(AlertsConditionList)
	cats := []ConditionCategory{ConditionDefault, ConditionPlugins, ConditionExternalService, ConditionSynthetics, ConditionNRQL, ConditionInfrastructure}
	for _, cat := range cats {
		// TODO: paralleize and use ctx.Done() to cancel the parent context
		listFunc := s.listByCategory(cat)
		resp, err := listFunc(ctx, list, opt)
		if err != nil || resp.StatusCode >= 400 {
			return nil, fmt.Errorf("%v.Response: %v. Error: %v.", cat, resp, err)
		}
	}

	return list, nil
}

func (s *AlertsConditionsService) List(ctx context.Context, opt *AlertsConditionsOptions, cat ConditionCategory) (*AlertsConditionList, *Response, error) {
	if opt == nil || opt.PolicyIDOptions == "" {
		return nil, nil, fmt.Errorf("policy_id is required")
	}

	list := new(AlertsConditionList)
	listFunc := s.listByCategory(cat)
	resp, err := listFunc(ctx, list, opt)
	return list, resp, err
}

func (s *AlertsConditionsService) Create(ctx context.Context, cat ConditionCategory, c *AlertsConditionEntity, conditionID int64) (*AlertsConditionEntity, *Response, error) {
	createFunc := s.createByCategory(cat)
	condition, resp, err := createFunc(ctx, c, conditionID)
	return condition, resp, err
}

func (s *AlertsConditionsService) Update(ctx context.Context, cat ConditionCategory, c *AlertsConditionEntity, conditionID int64) (*AlertsConditionEntity, *Response, error) {
	updateFunc := s.updateByCategory(cat)
	condition, resp, err := updateFunc(ctx, c, conditionID)
	return condition, resp, err
}

func (s *AlertsConditionsService) DeleteByID(ctx context.Context, cat ConditionCategory, conditionID int64) (*Response, error) {
	deleteFunc := s.deleteByCategory(cat)
	resp, err := deleteFunc(ctx, conditionID)
	return resp, err
}

func (s *AlertsConditionsService) listByCategory(cat ConditionCategory) func(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	switch cat {
	case ConditionDefault:
		return s.defaultConditions.listAll
	case ConditionExternalService:
		return s.externalServiceConditions.listAll
	case ConditionNRQL:
		return s.nrqlConditions.listAll
	case ConditionPlugins:
		return s.pluginsConditions.listAll
	case ConditionSynthetics:
		return s.syntheticsConditions.listAll
	case ConditionInfrastructure:
		return s.infraConditions.listAll
	default:
		return func(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
			return nil, fmt.Errorf("unsupported category %q", cat)
		}
	}
}

func (s *AlertsConditionsService) createByCategory(cat ConditionCategory) func(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
	switch cat {
	case ConditionDefault:
		return s.defaultConditions.create
	case ConditionExternalService:
		return s.externalServiceConditions.create
	case ConditionNRQL:
		return s.nrqlConditions.create
	case ConditionPlugins:
		return s.pluginsConditions.create
	case ConditionSynthetics:
		return s.syntheticsConditions.create
	case ConditionInfrastructure:
		return s.infraConditions.create
	default:
		return func(ctx context.Context, c *AlertsConditionEntity, policyID int64) (*AlertsConditionEntity, *Response, error) {
			return nil, nil, fmt.Errorf("unsupported category %q", cat)
		}
	}
}

func (s *AlertsConditionsService) updateByCategory(cat ConditionCategory) func(ctx context.Context, c *AlertsConditionEntity, conditionID int64) (*AlertsConditionEntity, *Response, error) {
	switch cat {
	case ConditionDefault:
		return s.defaultConditions.update
	case ConditionExternalService:
		return s.externalServiceConditions.update
	case ConditionNRQL:
		return s.nrqlConditions.update
	case ConditionPlugins:
		return s.pluginsConditions.update
	case ConditionSynthetics:
		return s.syntheticsConditions.update
	case ConditionInfrastructure:
		return s.infraConditions.update
	default:
		return func(ctx context.Context, c *AlertsConditionEntity, conditionID int64) (*AlertsConditionEntity, *Response, error) {
			return nil, nil, fmt.Errorf("unsupported category %q", cat)
		}
	}
}

func (s *AlertsConditionsService) deleteByCategory(cat ConditionCategory) func(ctx context.Context, conditionID int64) (*Response, error) {
	switch cat {
	case ConditionDefault:
		return s.defaultConditions.deleteByID
	case ConditionExternalService:
		return s.externalServiceConditions.deleteByID
	case ConditionNRQL:
		return s.nrqlConditions.deleteByID
	case ConditionPlugins:
		return s.pluginsConditions.deleteByID
	case ConditionSynthetics:
		return s.syntheticsConditions.deleteByID
	default:
		return func(ctx context.Context, conditionID int64) (*Response, error) {
			return nil, fmt.Errorf("unsupported category %q", cat)
		}
	}
}

type AlertsConditionTerm struct {
	Duration     *string `json:"duration,omitempty"`
	Operator     *string `json:"operator,omitempty"`
	Priority     *string `json:"priority,omitempty"`
	Threshold    *string `json:"threshold,omitempty"`
	TimeFunction *string `json:"time_function,omitempty"`
}
