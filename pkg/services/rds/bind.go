package rds

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *RDSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Check rds bind length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.BindDetails{}).
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check rds bind length in back database failed. Error: %s", err)
	}
	// ErrBindingAlreadyExists
	if length > 0 {
		return brokerapi.Binding{}, brokerapi.ErrBindingAlreadyExists
	}

	// Check rds instance length in back database
	err = database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("rds instance in back database: %v", ids))

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind rds instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	instance, err := instances.Get(rdsClient, instanceID).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find rds service failed. Error: %s", err)
	}

	// Get additional info from InstanceDetails
	addtionalparam := map[string]string{}
	err = ids.GetAdditionalInfo(&addtionalparam)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get rds instance additional info failed. Error: %s", err)
	}

	// Get specified parameters
	dbrtpd := addtionalparam["dbrtpd"]

	// Build Binding Credential: Default database user name is root
	credential, err := BuildBindingCredential(instance.HostName, instance.DbPort, instance.Name, "root", dbrtpd, service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build rds instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind rds instance success: %v", credential))

	// Constuct result
	result := brokerapi.Binding{Credentials: credential}

	// Marshal bind info
	bindinfo, err := json.Marshal(result)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("marshal rds bind info failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("create rds bind in back database opts: %v", bdsOpts))

	err = database.BackDBConnection.Create(&bdsOpts).Error
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create rds bind in back database failed. Error: %s", err)
	}

	// Log BindDetails result
	b.Logger.Debug(fmt.Sprintf("create rds bind in back database succeed: %s", bindingID))

	// Return result
	return result, nil
}
