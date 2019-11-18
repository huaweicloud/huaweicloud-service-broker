package rds

import (
	"fmt"
	"net/http"

	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *RDSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Check accepts_incomplete if this service support async
	if models.OperationAsyncRDS {
		e := b.Catalog.ValidateAcceptsIncomplete(asyncAllowed)
		if e != nil {
			return brokerapi.DeprovisionServiceSpec{}, e
		}
	}

	// Check rds instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length == 0 {
		return brokerapi.DeprovisionServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// get InstanceDetails in back database
	ids := database.InstanceDetails{}
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		First(&ids).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("rds instance in back database: %v", models.ToJson(ids)))

	// sync and check status whether allowed to update
	instance, err, serviceErr := SyncStatusWithService(b, instanceID, details.ServiceID, details.PlanID, ids.TargetID)

	if err != nil || serviceErr != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("sync status failed. error: %s, service error: %s", err, serviceErr)
	}
	if instance.Status != "ACTIVE" && instance.Status != "FAILED" {
		return brokerapi.DeprovisionServiceSpec{},
			brokerapi.NewFailureResponse(
				fmt.Errorf("Can only delete rds instance in ACTIVE or FAILED, but in: %s", instance.Status),
				http.StatusUnprocessableEntity, "Can only delete rds instance in ACTIVE or FAILED")
	}

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV3Client()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision rds instance opts: %s", instanceID))

	// Invoke sdk
	result := instances.Delete(rdsClient, ids.TargetID)
	if result.Err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision rds instance failed. Error: %s", result.Err)
	}

	// Delete InstanceDetails in back database
	err = database.BackDBConnection.Delete(&ids).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("delete rds instance in back database failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision rds instance success: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncRDS {
		// OperationDatas for OperationDeprovisioning
		ods := database.OperationDetails{
			OperationType:  models.OperationDeprovisioning,
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
			return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("convert rds instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create rds instance operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: ""}, nil
	}

	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
