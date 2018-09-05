package config

import (
	"errors"
	"fmt"
	"net/http"

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

func (c Catalog) ValidateOrgSpecGUID(organization_guid string, space_guid string) (error) {
	// Validate about parameters
	if organization_guid == "" {
		return brokerapi.NewFailureResponse(errors.New("organization_guid is empty"),
			http.StatusBadRequest, "organization_guid is empty")
	}
	if space_guid == "" {
		return brokerapi.NewFailureResponse(errors.New("space_guid is empty"),
			http.StatusBadRequest, "space_guid is empty")
	}
	return nil
}

// We get asyncAllowed in api.go of pivotal-cf/brokerapi like this:
// asyncAllowed := req.FormValue("accepts_incomplete") == "true"
func (c Catalog) ValidateAcceptsIncomplete(asyncAllowed bool) (error) {
	// Validate about parameters
	if ! asyncAllowed {
		return brokerapi.NewFailureResponse(errors.New("request doesn't have the accepts_incomplete parameter or it is false"),
			http.StatusUnprocessableEntity, "request doesn't have the accepts_incomplete parameter or it is false")
	}
	return nil
}
