package obs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *OBSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Check obs bind length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.BindDetails{}).
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check obs bind length in back database failed. Error: %s", err)
	}
	// ErrBindingAlreadyExistsSame
	if length > 0 {
		// Get BindDetails in back database
		bddetail := database.BindDetails{}
		err = database.BackDBConnection.
			Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
			First(&bddetail).Error
		if err != nil {
			return brokerapi.Binding{}, fmt.Errorf("get binding in back database failed. Error: %s", err)
		}

		// Get additional info from InstanceDetails
		bindingdetail := brokerapi.Binding{}
		err = bddetail.GetBindInfo(&bindingdetail)
		if err != nil {
			return brokerapi.Binding{}, fmt.Errorf("get binding info failed. Error: %s", err)
		}
		return bindingdetail, brokerapi.ErrBindingAlreadyExistsSame
	}

	// Check obs instance length in back database
	err = database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check obs instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length == 0 {
		return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
	}

	// get InstanceDetails in back database
	ids := database.InstanceDetails{}
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		First(&ids).Error
	if err != nil {
		return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
	}

	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("obs instance in back database: %v", models.ToJson(ids)))

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind obs bucket opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	bindOpts := &obs.GetBucketMetadataInput{}
	bindOpts.Bucket = ids.TargetID
	obsResponse, err := obsClient.GetBucketMetadata(bindOpts)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get obs bucket failed. Error: %s", err)
	}
	// Log result
	b.Logger.Debug(fmt.Sprintf("get obs bucket success: %v", models.ToJson(obsResponse)))

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find obs service failed. Error: %s", err)
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(
		b.CloudCredentials.Region,
		obsClient.GetEndpoint(),
		bindOpts.Bucket,
		b.CloudCredentials.AccessKey,
		b.CloudCredentials.SecretKey,
		service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build obs bucket binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind obs bucket success: %v", models.ToJson(credential)))

	// Constuct result
	result := brokerapi.Binding{Credentials: credential}

	// Marshal bind info
	bindinfo, err := json.Marshal(result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("marshal obs bind info failed. Error: %s", err)
	}

	// create BindDetails in back database
	bdsOpts := database.BindDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		BindID:         bindingID,
		BindInfo:       string(bindinfo),
		AdditionalInfo: "",
	}

	// log BindDetails opts
	b.Logger.Debug(fmt.Sprintf("create obs bind in back database opts: %v", models.ToJson(bdsOpts)))

	err = database.BackDBConnection.Create(&bdsOpts).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create obs bind in back database failed. Error: %s", err)
	}

	// Log BindDetails result
	b.Logger.Debug(fmt.Sprintf("create obs bind in back database succeed: %s", bindingID))

	// Return result
	return result, nil
}
