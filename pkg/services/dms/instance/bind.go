package instance

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
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
	// ErrBindingAlreadyExists
	if length > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
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

	// Invoke sdk
	instance, err := instances.Get(dmsClient, ids.TargetID).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dms instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find dms service failed. Error: %s", err)
	}

	// Get additional info from InstanceDetails
	addtionalparam := map[string]string{}
	err = ids.GetAdditionalInfo(&addtionalparam)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get rds instance additional info failed. Error: %s", err)
	}

	// Get specified parameters
	password := addtionalparam[AddtionalParamPassword]

	// Build Binding Credential
	credential, err := BuildBindingCredential(
		instance.ConnectAddress,
		instance.Port,
		instance.UserName,
		password,
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
