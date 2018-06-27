package queue

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *DMSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Check dms instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("check dms instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("dms instance in back database: %v", models.ToJson(ids)))

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create dms client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision dms queue opts: %s", instanceID))

	// Invoke sdk
	result := queues.Delete(dmsClient, ids.TargetID)
	if result.Err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision dms instance failed. Error: %s", err)
	}

	// Delete InstanceDetails in back database
	err = database.BackDBConnection.Delete(&ids).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("delete dms instance in back database failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision dms instance success: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDMS {
		// OperationDatas for OperationDeprovisioning
		ods := models.OperationDatas{
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
			return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("convert dms instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dms instance operation datas: %s", operationdata))

		return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: operationdata}, nil
	}

	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
