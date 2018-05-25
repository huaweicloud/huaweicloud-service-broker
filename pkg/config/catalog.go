package config

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Catalog define
type Catalog struct {
	Services []brokerapi.Service `json:"services"`
}

// Validate catalog
func (c Catalog) Validate() error {
	for _, service := range c.Services {

		if err := c.ValidateService(service); err != nil {
			return fmt.Errorf("Validating Services configuration: %s", err)
		}

		for _, serviceplan := range service.Plans {
			if err := c.ValidateServicePlan(serviceplan); err != nil {
				return fmt.Errorf("Validating ServicesPlans configuration: %s", err)
			}
		}
	}

	return nil
}

// ValidateService from file
func (c Catalog) ValidateService(s brokerapi.Service) error {
	if s.ID == "" {
		return fmt.Errorf("Service must provide a non-empty ID (%+v)", s)
	}

	if s.Name == "" {
		return fmt.Errorf("Service must provide a non-empty Name (%+v)", s)
	}

	if s.Description == "" {
		return fmt.Errorf("Service must provide a non-empty Description (%+v)", s)
	}

	return nil
}

// ValidateServicePlan from file
func (c Catalog) ValidateServicePlan(sp brokerapi.ServicePlan) error {
	if sp.ID == "" {
		return fmt.Errorf("ServicePlan must provide a non-empty ID (%+v)", sp)
	}

	if sp.Name == "" {
		return fmt.Errorf("ServicePlan must provide a non-empty Name (%+v)", sp)
	}

	if sp.Description == "" {
		return fmt.Errorf("ServicePlan must a non-empty Description (%+v)", sp)
	}

	return nil
}

// FindService from config
func (c Catalog) FindService(serviceid string) (brokerapi.Service, error) {

	// Validate about parameters
	if serviceid == "" {
		return brokerapi.Service{}, errors.New("serviceid is empty")
	}

	// Find service
	var service brokerapi.Service
	var existingservice bool
	for _, s := range c.Services {
		if s.ID == serviceid {
			service = s
			existingservice = true
			break
		}
	}

	// Return results
	if !existingservice {
		return brokerapi.Service{}, fmt.Errorf("service is not existing. serviceid: %s", serviceid)
	}
	return service, nil
}

// FindServicePlan from config
func (c Catalog) FindServicePlan(serviceid string, planid string) (brokerapi.ServicePlan, error) {

	// Validate about parameters
	if serviceid == "" {
		return brokerapi.ServicePlan{}, errors.New("serviceid is empty")
	}
	if planid == "" {
		return brokerapi.ServicePlan{}, errors.New("planid is empty")
	}

	// Find service plan
	var serviceplan brokerapi.ServicePlan
	var existingplan bool
	for _, service := range c.Services {
		if service.ID == serviceid {
			for _, plan := range service.Plans {
				if plan.ID == planid {
					serviceplan = plan
					existingplan = true
					break
				}
			}
		}
	}

	// Return results
	if !existingplan {
		return brokerapi.ServicePlan{}, fmt.Errorf("service plan is not existing. serviceid: %s planid: %s", serviceid, planid)
	}
	return serviceplan, nil
}
