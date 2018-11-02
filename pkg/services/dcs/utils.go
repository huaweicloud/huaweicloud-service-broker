package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dcs instance
func BuildBindingCredential(
	ip string,
	port int,
	username string,
	password string,
	name string,
	servicetype string) (BindingCredential, error) {

	if servicetype == models.DCSRedisServiceName {
		username = ""
	} else if servicetype == models.DCSMemcachedServiceName {

	} else if servicetype == models.DCSIMDGServiceName {
		port = 0
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		IP:       ip,
		Port:     port,
		UserName: username,
		Password: password,
		Name:     name,
		Type:     servicetype,
	}
	return bc, nil
}

func SyncStatusWithService(b *DCSBroker, instanceID string, serviceID string, planID string,
	targetID string) (*instances.Instance, error, error) {
	dbInstance := database.InstanceDetails{}
	// Log opts
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService dcs instance opts: instanceID: %s serviceID: %s "+
		"planID: %s targetID: %s", instanceID, serviceID, planID, targetID))

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return nil, fmt.Errorf("SyncStatusWithService create dcs client failed. Error: %s", err), nil
	}
	// Invoke sdk get
	instance, serviceErr := instances.Get(dcsClient, targetID).Extract()
	if serviceErr != nil {
		return nil, nil, serviceErr
	}
	// get InstanceDetails in back database
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, serviceID, planID).
		First(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService get dcs instance in back database failed. Error: %s", err))
		return instance, fmt.Errorf("SyncStatusWithService get dcs instance failed. Error: %s", err), nil
	}
	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService dcs instance in back database: %v", models.ToJson(dbInstance)))
	// update target info in back database
	targetInfo, err := json.Marshal(instance)
	if err != nil {
		return instance, fmt.Errorf("SyncStatusWithService marshal dcs instance failed. Error: %s", err), nil
	}

	dbInstance.TargetID = instance.InstanceID
	dbInstance.TargetName = instance.Name
	dbInstance.TargetStatus = instance.Status
	dbInstance.TargetInfo = string(targetInfo)

	err = database.BackDBConnection.Save(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update dcs instance target status in back database failed. "+
			"Error: %s", err))
		return instance, fmt.Errorf("SyncStatusWithService update dcs instance target status failed. Error: %s", err), nil
	}
	// Sync target status success
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update dcs instance target status succeed: %s", instanceID))
	return instance, nil, nil
}
