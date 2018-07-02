package obs

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *OBSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Check obs instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("check obs instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("obs instance in back database: %v", models.ToJson(ids)))

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision obs bucket opts: %s", instanceID))

	// Invoke sdk
	obsResponse, err := obsClient.DeleteBucket(ids.TargetID)
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision obs bucket failed. Error: %s", err)
	}

	// Delete InstanceDetails in back database
	err = database.BackDBConnection.Delete(&ids).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("delete obs bucket in back database failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision obs bucket %s success: %v", instanceID, models.ToJson(obsResponse)))

	// Return result
	if asyncAllowed && models.OperationAsyncOBS {
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
			return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("convert obs bucket operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create obs bucket operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.DeprovisionServiceSpec{IsAsync: true, OperationData: ""}, nil
	}

	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
