package dcs

import (
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *DCSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	return brokerapi.ProvisionedServiceSpec{}, nil
}
