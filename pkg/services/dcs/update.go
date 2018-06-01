package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *DCSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Init updateOpts
	updateOpts := instances.UpdateOpts{}

	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &updateOpts)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update dcs instance opts: %v", updateOpts))

	// Invoke sdk
	updateResult := instances.Update(dcsClient, instanceID, updateOpts)
	if updateResult.Err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("update dcs instance result: %v", updateResult))

	// Return result
	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
