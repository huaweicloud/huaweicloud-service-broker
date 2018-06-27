package rds

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *RDSBroker) LastOperation(instanceID string, operationData models.OperationDatas) (brokerapi.LastOperation, error) {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("lastoperation rds instance opts: instanceID: %s operationData: %v", instanceID, models.ToJson(operationData)))

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.LastOperation{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Invoke sdk get
	instance, err := instances.Get(rdsClient, operationData.TargetID).Extract()

	// Handle different cases
	if (operationData.OperationType == models.OperationProvisioning) ||
		(operationData.OperationType == models.OperationUpdating) {
		// OperationProvisioning || OperationUpdating
		if err != nil {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("get rds instance failed. Error: %s", err),
			}, nil
		}
		// Status
		if instance.Status == "ACTIVE" {
			return brokerapi.LastOperation{
				State:       brokerapi.Succeeded,
				Description: fmt.Sprintf("Status: %s", instance.Status),
			}, nil
		} else if instance.Status == "FAILED" {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("Status: %s", instance.Status),
			}, nil
		} else {
			return brokerapi.LastOperation{
				State:       brokerapi.InProgress,
				Description: fmt.Sprintf("Status: %s", instance.Status),
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
					Description: fmt.Sprintf("get rds instance failed. Error: %s", err),
				}, nil
			}
		} else {
			return brokerapi.LastOperation{
				State:       brokerapi.InProgress,
				Description: fmt.Sprintf("Status: %s", instance.Status),
			}, nil
		}
	} else {
		b.Logger.Debug(fmt.Sprintf("unknown operation data type: %s", operationData.OperationType))
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("lastoperation rds instance success: instanceID: %s", instanceID))

	return brokerapi.LastOperation{}, nil
}
