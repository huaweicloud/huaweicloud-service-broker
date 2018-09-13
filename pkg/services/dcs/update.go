package dcs

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	b.Logger.Debug(fmt.Sprintf("dcs instance in back database: %v", models.ToJson(ids)))

	// sync and check status whether allowed to update
	instance, err, serviceErr := SyncStatusWithService(b, instanceID, details.ServiceID, details.PlanID, ids.TargetID)

	if err != nil || serviceErr != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("sync status failed. error: %s, service error: %s", err, serviceErr)
	}
	if instance.Status != "RUNNING" {
		return brokerapi.UpdateServiceSpec{},
			brokerapi.NewFailureResponse(
				fmt.Errorf("Can only update dcs instance in RUNNING, but in: %s", instance.Status),
				http.StatusUnprocessableEntity, "Can only update dcs instance in RUNNING")
	}

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Init updateParameters
	updateParameters := UpdateParameters{}
	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &updateParameters)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update dcs instance opts: %v", models.ToJson(updateParameters)))

	// Init updateOpts
	updateOpts := instances.UpdateOpts{}
	// Name
	if updateParameters.Name != "" {
		updateOpts.Name = updateParameters.Name
	}
	// Description
	if updateParameters.Description != nil {
		updateOpts.Description = updateParameters.Description
	}
	// BackupStrategy
	if (updateParameters.BackupStrategySavedays > 0) &&
		(updateParameters.BackupStrategyBackupType != "") &&
		(updateParameters.BackupStrategyBeginAt != "") &&
		(updateParameters.BackupStrategyPeriodType != "") &&
		(len(updateParameters.BackupStrategyBackupAt) > 0) {
		updateOpts.InstanceBackupPolicy = &instances.InstanceBackupPolicy{
			SaveDays:   updateParameters.BackupStrategySavedays,
			BackupType: updateParameters.BackupStrategyBackupType,
			PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
				BeginAt:    updateParameters.BackupStrategyBeginAt,
				PeriodType: updateParameters.BackupStrategyPeriodType,
				BackupAt:   updateParameters.BackupStrategyBackupAt,
			},
		}
	}
	// MaintainBegin
	if updateParameters.MaintainBegin != "" {
		updateOpts.MaintainBegin = updateParameters.MaintainBegin
	}
	// MaintainEnd
	if updateParameters.MaintainEnd != "" {
		updateOpts.MaintainEnd = updateParameters.MaintainEnd
	}
	if updateParameters.SecurityGroupID != "" {
		updateOpts.SecurityGroupID = updateParameters.SecurityGroupID
	}

	// Invoke sdk update
	updateResult := instances.Update(dcsClient, ids.TargetID, updateOpts)
	if updateResult.Err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("update dcs instance result: %v", models.ToJson(updateResult)))

	// Update password
	if updateParameters.NewPassword != nil {
		updatePasswordResult := instances.UpdatePassword(
			dcsClient,
			ids.TargetID,
			instances.UpdatePasswordOpts{
				OldPassword: *updateParameters.OldPassword,
				NewPassword: *updateParameters.NewPassword,
			})
		if updatePasswordResult.Err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance password failed. Error: %s", updatePasswordResult.Err)
		}

		// Update back database
		addtionalparam := map[string]string{}
		if ids.AdditionalInfo != "" {
			err := json.Unmarshal([]byte(ids.AdditionalInfo), &addtionalparam)
			if err != nil {
				return brokerapi.UpdateServiceSpec{},
					fmt.Errorf("unmarshalling dcs addtional info failed. Error: %s", err)
			}
		}
		// Reset addtional info
		addtionalparam[AddtionalParamPassword] = *updateParameters.NewPassword
		// Marshal addtional info
		addtionalinfo, err := json.Marshal(addtionalparam)
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("marshal dcs addtional info failed. Error: %s", err)
		}
		ids.AdditionalInfo = string(addtionalinfo)

		// Log result
		b.Logger.Debug(fmt.Sprintf("update dcs instance password result: %v", models.ToJson(updatePasswordResult)))

	}

	// Extend capacity
	if updateParameters.NewCapacity > 0 {
		extendResult := instances.Extend(
			dcsClient,
			ids.TargetID,
			instances.ExtendOpts{
				NewCapacity: updateParameters.NewCapacity,
			})
		if extendResult.Err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("extend dcs instance failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("extend dcs instance result: %v", models.ToJson(extendResult)))
	}

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
	b.Logger.Debug(fmt.Sprintf("update dcs instance in back database opts: %v", models.ToJson(ids)))

	err = database.BackDBConnection.Save(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dcs instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("update dcs instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDCS {
		// OperationDatas for OperationUpdating
		ods := database.OperationDetails{
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

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: ""}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
