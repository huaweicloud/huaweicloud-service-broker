package obs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *OBSBroker) LastOperation(instanceID string, operationData models.OperationDatas) (brokerapi.LastOperation, error) {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("lastoperation obs bucket opts: instanceID: %s operationData: %v", instanceID, models.ToJson(operationData)))

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.LastOperation{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Invoke sdk get
	getOpts := &obs.GetBucketMetadataInput{}
	getOpts.Bucket = operationData.TargetID
	_, err = obsClient.GetBucketMetadata(getOpts)

	// Handle different cases
	if (operationData.OperationType == models.OperationProvisioning) ||
		(operationData.OperationType == models.OperationUpdating) {
		// OperationProvisioning || OperationUpdating
		if err != nil {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("get obs bucket failed. Error: %s", err),
			}, nil
		} else {
			return brokerapi.LastOperation{
				State:       brokerapi.Succeeded,
				Description: fmt.Sprintf("Bucket: %s", getOpts.Bucket),
			}, nil
		}
	} else if operationData.OperationType == models.OperationDeprovisioning {
		// OperationDeprovisioning
		if err != nil {
			e, ok := err.(golangsdk.ErrDefault404)
			if ok {
				return brokerapi.LastOperation{
					State:       brokerapi.Succeeded,
					Description: fmt.Sprintf("Status: %s", e.Error()),
				}, nil
			} else {
				return brokerapi.LastOperation{
					State:       brokerapi.Failed,
					Description: fmt.Sprintf("get obs bucket failed. Error: %s", err),
				}, nil
			}
		} else {
			return brokerapi.LastOperation{
				State:       brokerapi.InProgress,
				Description: fmt.Sprintf("Bucket: %s", getOpts.Bucket),
			}, nil
		}
	} else {
		b.Logger.Debug(fmt.Sprintf("unknown operation data type: %s", operationData.OperationType))
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("lastoperation obs bucket success: instanceID: %s", instanceID))

	return brokerapi.LastOperation{}, nil
}
