package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *DCSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Check dcs instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("check dcs instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length == 0 {
		return brokerapi.UpdateServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// get InstanceDetails in back database
	ids := database.InstanceDetails{}
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		First(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("dcs instance in back database: %v", ids))

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Init updateOpts
	updateOpts := instances.UpdateOpts{}

	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &updateOpts)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update dcs instance opts: %v", updateOpts))

	// Invoke sdk
	updateResult := instances.Update(dcsClient, ids.TargetID, updateOpts)
	if updateResult.Err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("update dcs instance result: %v", updateResult))

	// Invoke sdk get
	freshInstance, err := instances.Get(dcsClient, ids.TargetID).Extract()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get dcs instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal dcs instance failed. Error: %s", err)
	}

	ids.TargetID = freshInstance.InstanceID
	ids.TargetName = freshInstance.Name
	ids.TargetStatus = freshInstance.Status
	ids.TargetInfo = string(targetinfo)

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("update dcs instance in back database opts: %v", ids))

	err = database.BackDBConnection.Save(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("update dcs instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDCS {
		// OperationDatas for OperationUpdating
		ods := models.OperationDatas{
			OperationType:  models.OperationUpdating,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       ids.TargetID,
			TargetName:     ids.TargetName,
			TargetStatus:   ids.TargetStatus,
			TargetInfo:     ids.TargetInfo,
			AdditionalInfo: ids.AdditionalInfo,
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("convert dcs instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dcs instance operation datas: %s", operationdata))

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: operationdata}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
