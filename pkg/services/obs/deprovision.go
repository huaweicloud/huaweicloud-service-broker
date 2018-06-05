package obs

import (
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *OBSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision obs bucket opts: %s", instanceID))

	// Invoke sdk
	obsResponse, err := obsClient.DeleteBucket(instanceID)
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision obs bucket failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision obs bucket %s success: %v", instanceID, obsResponse))

	// Return result
	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
