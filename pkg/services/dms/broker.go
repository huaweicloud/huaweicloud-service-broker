package dms

import (
	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"github.com/pivotal-cf/brokerapi"
)

// DMSBroker define
type DMSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// Provision implematation
func (b *DMSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{}, nil
}

// Deprovision implematation
func (b *DMSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, nil
}

// Bind implematation
func (b *DMSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}

// Unbind implematation
func (b *DMSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}

// Update implematation
func (b *DMSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

// LastOperation implematation
func (b *DMSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
