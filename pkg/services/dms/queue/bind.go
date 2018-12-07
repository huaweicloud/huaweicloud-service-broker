package queue

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/groups"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *DMSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Check dms bind length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.BindDetails{}).
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check dms bind length in back database failed. Error: %s", err)
	}
	// ErrBindingAlreadyExistsSame
	if length > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExistsSame
	}

	// Check dms instance length in back database
	err = database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check dms instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("dms instance in back database: %v", models.ToJson(ids)))

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dms client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind dms instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk: the default includeDeadLetter value is false
	queue, err := queues.Get(dmsClient, ids.TargetID, false).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dms instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find dms service failed. Error: %s", err)
	}

	// Invoke sdk: the default includeDeadLetter value is false
	pages, err := groups.List(dmsClient, ids.TargetID, false).AllPages()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("list dms group failed. Error: %s", err)
	}
	grouplist, err := groups.ExtractGroups(pages)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("extract dms group failed. Error: %s", err)
	}
	var groupid string
	if len(grouplist) > 0 {
		groupid = grouplist[0].ID
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Get parameters from service plan metadata
	metadataParameters := MetadataParameters{}
	if servicePlan.Metadata != nil {
		if len(servicePlan.Metadata.Parameters) > 0 {
			err := json.Unmarshal(servicePlan.Metadata.Parameters, &metadataParameters)
			if err != nil {
				return brokerapi.Binding{},
					fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
			}
		}
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(
		metadataParameters.EndpointName,
		metadataParameters.EndpointPort,
		b.CloudCredentials.Region,
		dmsClient.ProjectID,
		dmsClient.Endpoint,
		b.CloudCredentials.AccessKey,
		b.CloudCredentials.SecretKey,
		queue.ID,
		groupid,
		service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build dms instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind dms instance success: %v", models.ToJson(credential)))

	// Constuct result
	result := brokerapi.Binding{Credentials: credential}

	// Marshal bind info
	bindinfo, err := json.Marshal(result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("marshal dms bind info failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("create dms bind in back database opts: %v", models.ToJson(bdsOpts)))

	err = database.BackDBConnection.Create(&bdsOpts).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dms bind in back database failed. Error: %s", err)
	}

	// Log BindDetails result
	b.Logger.Debug(fmt.Sprintf("create dms bind in back database succeed: %s", bindingID))

	// Return result
	return result, nil
}
