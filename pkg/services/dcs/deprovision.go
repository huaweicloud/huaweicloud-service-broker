package dcs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *DCSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision dcs instance opts: %s", instanceID))

	// Invoke sdk
	result := instances.Delete(dcsClient, instanceID)
	if result.Err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision dcs instance success: %s", instanceID))

	// Return result
	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
