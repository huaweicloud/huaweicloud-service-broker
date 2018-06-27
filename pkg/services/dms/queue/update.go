package queue

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation is not necessary for DMSStandardServiceName, DMSActiveMQServiceName and DMSKafkaServiceName
func (b *DMSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Check dms instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("check dms instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("dms instance in back database: %v", models.ToJson(ids)))

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create dms client failed. Error: %s", err)
	}

	// Invoke sdk get
	freshQueue, err := queues.Get(dmsClient, ids.TargetID, false).Extract()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get dms queue failed. Error: %s", err)
	}

	// Marshal queue
	targetinfo, err := json.Marshal(freshQueue)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal dms queue failed. Error: %s", err)
	}

	ids.TargetID = freshQueue.ID
	ids.TargetName = freshQueue.Name
	ids.TargetInfo = string(targetinfo)

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("update dms instance in back database opts: %v", models.ToJson(ids)))

	err = database.BackDBConnection.Save(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update dms queue in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("update dms queue in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDMS {
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
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("convert dms queue operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dms queue operation datas: %s", operationdata))

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: operationdata}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
