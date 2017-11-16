package brokerapi

import (
	"errors"
	"fmt"
)

type Catalog struct {
	Services []Service `json:"services"`
}

type Service struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	Bindable        bool             `json:"bindable,omitempty"`
	Tags            []string         `json:"tags,omitempty"`
	Metadata        *ServiceMetadata `json:"metadata,omitempty"`
	Requires        []string         `json:"requires,omitempty"`
	PlanUpdateable  bool             `json:"plan_updateable"`
	Plans           []ServicePlan    `json:"plans"`
	DashboardClient *DashboardClient `json:"dashboard_client,omitempty"`
}

type ServiceMetadata struct {
	DisplayName         string `json:"displayName,omitempty"`
	ImageURL            string `json:"imageUrl,omitempty"`
	LongDescription     string `json:"longDescription,omitempty"`
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
	DocumentationURL    string `json:"documentationUrl,omitempty"`
	SupportURL          string `json:"supportUrl,omitempty"`
}

type ServicePlan struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Metadata    *ServicePlanMetadata `json:"metadata,omitempty"`
	Free        bool                 `json:"free"`
}

type ServicePlanMetadata struct {
	Bullets     []string `json:"bullets,omitempty"`
	Costs       []Cost   `json:"costs,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
}

type DashboardClient struct {
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type Cost struct {
	Amount map[string]interface{} `json:"amount,omitempty"`
	Unit   string                 `json:"unit,omitempty"`
}

func (c Catalog) Validate() error {
	for _, service := range c.Services {
		if err := service.Validate(); err != nil {
			return fmt.Errorf("Validating Services configuration: %s", err)
		}
	}

	return nil
}

func (s Service) Validate() error {
	if s.ID == "" {
		return errors.New("Must provide a non-empty ID")
	}

	if s.Name == "" {
		return errors.New("Must provide a non-empty Name")
	}

	if s.Description == "" {
		return errors.New("Must provide a non-empty Description")
	}

	for _, servicePlan := range s.Plans {
		if err := servicePlan.Validate(); err != nil {
			return fmt.Errorf("Validating Plans configuration: %s", err)
		}
	}

	return nil
}

func (sp ServicePlan) Validate() error {
	if sp.ID == "" {
		return errors.New("Must provide a non-empty ID")
	}

	if sp.Name == "" {
		return errors.New("Must provide a non-empty Name")
	}

	if sp.Description == "" {
		return errors.New("Must provide a non-empty Description")
	}

	return nil
}

func (c Catalog) FindService(serviceID string) (service Service, found bool) {
	for _, service := range c.Services {
		if service.ID == serviceID {
			return service, true
		}
	}

	return service, false
}

func (c Catalog) FindServicePlan(planID string) (plan ServicePlan, found bool) {
	for _, service := range c.Services {
		for _, plan := range service.Plans {
			if plan.ID == planID {
				return plan, true
			}
		}
	}

	return plan, false
}
