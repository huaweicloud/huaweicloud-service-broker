package obs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *OBSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Check obs instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check obs instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Get parameters from service plan metadata
	metadataParameters := MetadataParameters{}
	if servicePlan.Metadata != nil {
		if len(servicePlan.Metadata.Parameters) > 0 {
			err := json.Unmarshal(servicePlan.Metadata.Parameters, &metadataParameters)
			if err != nil {
				return brokerapi.ProvisionedServiceSpec{},
					fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
			}
		}
	}

	// Get parameters from details
	provisionParameters := ProvisionParameters{}
	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &provisionParameters)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{},
				fmt.Errorf("Error unmarshalling rawParameters from details: %s", err)
		}
	}

	// Init provisionOpts
	provisionOpts := &obs.CreateBucketInput{}
	// Setting Bucket
	provisionOpts.Bucket = provisionParameters.BucketName
	// Setting Default StorageClass
	provisionOpts.StorageClass = obs.StorageClassType(metadataParameters.StorageClass)
	// Setting ACL
	provisionOpts.ACL = obs.AclType(metadataParameters.BucketPolicy)
	if provisionParameters.BucketPolicy != "" {
		provisionOpts.ACL = obs.AclType(provisionParameters.BucketPolicy)
	}
	// Setting Location
	provisionOpts.Location = b.CloudCredentials.Region

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision obs bucket opts: %v", provisionOpts))

	// Invoke sdk
	obsResponse, err := obsClient.CreateBucket(provisionOpts)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision obs bucket failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision obs bucket result: %v", obsResponse))

	// Invoke sdk get
	getOpts := &obs.GetBucketMetadataInput{}
	getOpts.Bucket = provisionOpts.Bucket
	freshBucket, err := obsClient.GetBucketMetadata(getOpts)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get obs bucket failed. Error: %s", err)
	}

	// Marshal bucket
	targetinfo, err := json.Marshal(freshBucket)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal obs bucket failed. Error: %s", err)
	}

	// create InstanceDetails in back database
	idsOpts := database.InstanceDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		TargetID:       provisionOpts.Bucket,
		TargetName:     provisionOpts.Bucket,
		TargetStatus:   "",
		TargetInfo:     string(targetinfo),
		AdditionalInfo: "",
	}

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("create obs bucket in back database opts: %v", idsOpts))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create obs bucket in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create obs bucket in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncOBS {
		// OperationDatas for OperationProvisioning
		ods := models.OperationDatas{
			OperationType:  models.OperationProvisioning,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       provisionOpts.Bucket,
			TargetName:     provisionOpts.Bucket,
			TargetStatus:   "",
			TargetInfo:     string(targetinfo),
			AdditionalInfo: "",
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("convert obs bucket operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create obs bucket operation datas: %s", operationdata))

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: operationdata}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
