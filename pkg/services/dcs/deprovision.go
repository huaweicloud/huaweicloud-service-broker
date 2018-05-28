package dcs

import (
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *DCSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	return brokerapi.DeprovisionServiceSpec{}, nil
}
