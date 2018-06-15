package obs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *OBSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Check obs bucket length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("check obs bucket length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("obs bucket in back database: %v", ids))

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
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
	b.Logger.Debug(fmt.Sprintf("update obs bucket opts: %v", updateParameters))

	// Setting BucketAcl
	if updateParameters.BucketPolicy != "" {
		// Init aclOpts
		aclOpts := &obs.SetBucketAclInput{}
		aclOpts.Bucket = ids.TargetID
		aclOpts.ACL = obs.AclType(updateParameters.BucketPolicy)
		// Invoke sdk
		aclResponse, err := obsClient.SetBucketAcl(aclOpts)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("set obs bucket acl failed. Error: %s", err)
		}
		// Log result
		b.Logger.Debug(fmt.Sprintf("set obs bucket acl result: %v", aclResponse))
	}

	// Setting BucketVersioning
	if updateParameters.Status != "" {
		// Init versioningOpts
		versioningOpts := &obs.SetBucketVersioningInput{}
		versioningOpts.Bucket = ids.TargetID
		// Enabled && Suspended
		versioningOpts.Status = obs.VersioningStatusType(updateParameters.Status)
		// Invoke sdk
		versioningResponse, err := obsClient.SetBucketVersioning(versioningOpts)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("set obs bucket versioning failed. Error: %s", err)
		}
		// Log result
		b.Logger.Debug(fmt.Sprintf("set obs bucket versioning result: %v", versioningResponse))
	}

	// Invoke sdk get
	getOpts := &obs.GetBucketMetadataInput{}
	getOpts.Bucket = ids.TargetID
	freshBucket, err := obsClient.GetBucketMetadata(getOpts)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get obs bucket failed. Error: %s", err)
	}

	// Marshal bucket
	targetinfo, err := json.Marshal(freshBucket)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal obs bucket failed. Error: %s", err)
	}

	ids.TargetInfo = string(targetinfo)

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("update obs bucket in back database opts: %v", ids))

	err = database.BackDBConnection.Save(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update obs bucket in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("update obs bucket in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncOBS {
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
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("convert obs bucket operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create obs bucket operation datas: %s", operationdata))

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: operationdata}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
