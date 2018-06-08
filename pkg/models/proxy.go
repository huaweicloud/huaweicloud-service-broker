package models

import (
	"github.com/pivotal-cf/brokerapi"
)

// ServiceBrokerProxy is used to implement details
type ServiceBrokerProxy interface {
	Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error)

	Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error)

	Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error)

	Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error

	Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error)

	LastOperation(instanceID string, operationData OperationDatas) (brokerapi.LastOperation, error)
}
