package brokerapi

import "errors"

type ServiceBroker interface {
	Services() CatalogResponse

	Provision(instanceID string, details ProvisionDetails, acceptsIncomplete bool) (ProvisioningResponse, bool, error)
	Update(instanceID string, details UpdateDetails, acceptsIncomplete bool) (bool, error)
	Deprovision(instanceID string, details DeprovisionDetails, acceptsIncomplete bool) (bool, error)

	Bind(instanceID string, bindingID string, details BindDetails) (BindingResponse, error)
	Unbind(instanceID string, bindingID string, details UnbindDetails) error

	LastOperation(instanceID string) (LastOperationResponse, error)
}

type ProvisionDetails struct {
	OrganizationGUID string                 `json:"organization_guid"`
	PlanID           string                 `json:"plan_id"`
	ServiceID        string                 `json:"service_id"`
	SpaceGUID        string                 `json:"space_guid"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type UpdateDetails struct {
	ServiceID      string                 `json:"service_id"`
	PlanID         string                 `json:"plan_id"`
	Parameters     map[string]interface{} `json:"parameters"`
	PreviousValues PreviousValues         `json:"previous_values"`
}

type PreviousValues struct {
	PlanID         string `json:"plan_id"`
	ServiceID      string `json:"service_id"`
	OrganizationID string `json:"organization_id"`
	SpaceID        string `json:"space_id"`
}

type DeprovisionDetails struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

type BindDetails struct {
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	AppGUID    string                 `json:"app_guid,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type UnbindDetails struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
}

var (
	ErrInstanceAlreadyExists = errors.New("instance already exists")
	ErrInstanceDoesNotExist  = errors.New("instance does not exist")
	ErrInstanceLimitMet      = errors.New("instance limit for this service has been reached")
	ErrInstanceNotUpdateable = errors.New("instance is not updateable")
	ErrInstanceNotBindable   = errors.New("instance is not bindable")
	ErrBindingAlreadyExists  = errors.New("binding already exists")
	ErrBindingDoesNotExist   = errors.New("binding does not exist")
	ErrAsyncRequired         = errors.New("This service plan requires client support for asynchronous service operations.")
	ErrAppGUIDRequired       = errors.New("This service supports generation of credentials through binding an application only.")
)
