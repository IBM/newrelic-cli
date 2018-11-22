package newrelic

import (
	"context"
)

type infraConditions service
type filterField map[string][]map[string]interface{}

//AlertsInfrastructureConditionList list of infrastructure conditions
type AlertsInfrastructureConditionList struct {
	AlertsInfrastructureConditions []*AlertsInfrastructureCondition `json:"data,omitempty"`
}

// AlertsInfrastructureCondition object for infrastructure condition
type AlertsInfrastructureCondition struct {
	AlertsInfraThreshold
	Comparison           *string      `json:"comparison,omitempty"`
	CreatedAtEpochMillis *int64       `json:"created_at_epoch_millis,omitempty"`
	Enabled              *bool        `json:"enabled,omitempty"`
	EventType            *string      `json:"event_type,omitempty"`
	Filter               *filterField `json:"filter,omitempty"`
	ID                   *int64       `json:"id,omitempty"`
	IntegrationProvider  *string      `json:"integration_provider,omitempty"`
	Name                 *string      `json:"name,omitempty"`
	PolicyID             *int64       `json:"policy_id,omitempty"`
	SelectValue          *string      `json:"select_value,omitempty"`
	Type                 *string      `json:"type,omitempty"`
	UpdatedAtEpochMillis *int64       `json:"updated_at_epoch_millis,omitempty"`
	WhereClause          *string      `json:"where_clause,omitempty"`
}

// AlertsInfraThreshold thresholds for alert conditions
type AlertsInfraThreshold struct {
	CriticalThreshold *AlertsInfraThresholdValues `json:"critical_threshold,omitempty"`
	WarningThreshold  *AlertsInfraThresholdValues `json:"warning_threshold,omitempty"`
}

// AlertsInfraThresholdValues threshold values for condition
type AlertsInfraThresholdValues struct {
	DurationMinutes *int64  `json:"duration_minutes,omitempty"`
	TimeFunction    *string `json:"time_function,omitempty"`
	Value           *int64  `json:"value,omitempty"`
}

func (s *infraConditions) listAll(ctx context.Context, list *AlertsConditionList, opt *AlertsConditionsOptions) (*Response, error) {
	u, err := addOptions("conditions", opt)
	if err != nil {
		return nil, err
	}
	u = infrastructureURL + u

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	list.AlertsInfrastructureConditionList = new(AlertsInfrastructureConditionList)
	resp, err := s.client.Do(ctx, req, list.AlertsInfrastructureConditionList)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
