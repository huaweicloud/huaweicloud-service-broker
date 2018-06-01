package dms

import (
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *DMSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind dms instance opts: instanceID: %s details: %v", instanceID, details))

	// TODO do something for update

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind dms instance success: instanceID: %s details: %s", instanceID, details))

	// Return result
	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
