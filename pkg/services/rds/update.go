package rds

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init rawParameters
	rawParameters := map[string]string{}

	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &rawParameters)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update rds instance opts: %v", rawParameters))

	// Invoke sdk
	if volumesize, ok := rawParameters["volumesize"]; ok {
		// UpdateVolumeSize
		rdsInstance, err := instances.UpdateVolumeSize(
			rdsClient,
			instances.UpdateOps{
				Volume: map[string]interface{}{
					"size": volumesize,
				},
			},
			instanceID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance volume size failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance volume size result: %v", rdsInstance))
	}

	// Invoke sdk
	if flavorname, ok := rawParameters["flavorname"]; ok {
		// UpdateFlavorRef
		rdsInstance, err := instances.UpdateFlavorRef(
			rdsClient,
			instances.UpdateFlavorOps{
				// TODO convert flavor name to flavor ref
				FlavorRef: flavorname,
			},
			instanceID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance flavor failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance flavor result: %v", rdsInstance))
	}

	// Return result
	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
