package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *DCSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Check dcs bind length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.BindDetails{}).
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check dcs bind length in back database failed. Error: %s", err)
	}
	// ErrBindingAlreadyExists
	if length > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	}

	// Check dcs instance length in back database
	err = database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check dcs instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("dcs instance in back database: %v", ids))

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind dcs instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	instance, err := instances.Get(dcsClient, ids.TargetID).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dcs instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find dcs service failed. Error: %s", err)
	}

	// Get additional info from InstanceDetails
	addtionalparam := map[string]string{}
	err = ids.GetAdditionalInfo(&addtionalparam)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dcs instance additional info failed. Error: %s", err)
	}

	// Get specified parameters
	password := addtionalparam["password"]

	// Build Binding Credential
	credential, err := BuildBindingCredential(instance.IP, instance.Port, instance.UserName, password, instance.Name, service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build dcs instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind dcs instance success: %v", credential))

	// Constuct result
	result := brokerapi.Binding{Credentials: credential}

	// Marshal bind info
	bindinfo, err := json.Marshal(result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("marshal dcs bind info failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("create dcs bind in back database opts: %v", bdsOpts))

	err = database.BackDBConnection.Create(&bdsOpts).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dcs bind in back database failed. Error: %s", err)
	}

	// Log BindDetails result
	b.Logger.Debug(fmt.Sprintf("create dcs bind in back database succeed: %s", bindingID))

	// Return result
	return result, nil
}
