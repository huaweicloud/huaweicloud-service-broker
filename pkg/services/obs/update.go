package obs

import (
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *OBSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind obs bucket opts: instanceID: %s details: %v", instanceID, details))

	// TODO do something for update

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind obs bucket success: instanceID: %s details: %s", instanceID, details))

	// Return result
	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
