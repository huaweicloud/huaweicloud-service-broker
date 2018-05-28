package dms

import (
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *DMSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}
