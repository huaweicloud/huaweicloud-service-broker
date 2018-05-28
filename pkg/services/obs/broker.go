package obs

import (
	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"github.com/pivotal-cf/brokerapi"
)

// OBSBroker define
type OBSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// Provision implematation
func (b *OBSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{}, nil
}

// Deprovision implematation
func (b *OBSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, nil
}

// Bind implematation
func (b *OBSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}

// Unbind implematation
func (b *OBSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}

// Update implematation
func (b *OBSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}

// LastOperation implematation
func (b *OBSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
