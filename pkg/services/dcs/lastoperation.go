package dcs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *DCSBroker) LastOperation(instanceID string, operationData models.OperationDatas) (brokerapi.LastOperation, error) {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("lastoperation dcs instance opts: instanceID: %s operationData: %v", instanceID, models.ToJson(operationData)))

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.LastOperation{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Invoke sdk get
	instance, err := instances.Get(dcsClient, operationData.TargetID).Extract()

	// Handle different cases
	if (operationData.OperationType == models.OperationProvisioning) ||
		(operationData.OperationType == models.OperationUpdating) {
		// OperationProvisioning || OperationUpdating
		if err != nil {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("get dcs instance failed. Error: %s", err),
			}, nil
		}
		// ErrorCode
		if instance.ErrorCode != "" {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("ErrorCode: %s", instance.ErrorCode),
			}, nil
		}
		// Status
		if instance.Status == "RUNNING" {
			return brokerapi.LastOperation{
				State:       brokerapi.Succeeded,
				Description: fmt.Sprintf("Status: %s", instance.Status),
			}, nil
		} else if (instance.Status == "CREATEFAILED") ||
			(instance.Status == "ERROR") {
			return brokerapi.LastOperation{
				State:       brokerapi.Failed,
				Description: fmt.Sprintf("ErrorCode: %s", instance.ErrorCode),
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
					Description: fmt.Sprintf("get dcs instance failed. Error: %s", err),
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
	b.Logger.Debug(fmt.Sprintf("lastoperation dcs instance success: instanceID: %s", instanceID))

	return brokerapi.LastOperation{}, nil
}
