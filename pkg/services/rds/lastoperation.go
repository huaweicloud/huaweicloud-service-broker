package rds

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *RDSBroker) LastOperation(instanceID string, operationData database.OperationDetails) (brokerapi.LastOperation, error) {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("lastoperation rds instance opts: instanceID: %s operationData: %v", instanceID, models.ToJson(operationData)))
	instance, err, serviceErr := SyncStatusWithService(b, instanceID, operationData.ServiceID,
		operationData.PlanID, operationData.TargetID, operationData.OperationType)

	if err != nil {
		return brokerapi.LastOperation{}, err
	}

	// Handle different cases
	if (operationData.OperationType == models.OperationProvisioning) ||
		(operationData.OperationType == models.OperationUpdating) {
		// OperationProvisioning || OperationUpdating
		if serviceErr != nil {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("get rds instance failed. Error: %s", serviceErr),
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
		if serviceErr != nil {
			e, ok := serviceErr.(golangsdk.ErrDefault404)
			if ok {
				return brokerapi.LastOperation{
					State:       brokerapi.Succeeded,
					Description: fmt.Sprintf("Status: %s", e.Error()),
				}, nil
			} else {
				return brokerapi.LastOperation{
					State:       brokerapi.Failed,
					Description: fmt.Sprintf("get rds instance failed. Error: %s", serviceErr),
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
