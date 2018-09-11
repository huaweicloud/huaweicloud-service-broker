package instance

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dms instance
func BuildBindingCredential(
	host string,
	port int,
	username string,
	password string,
	servicetype string) (BindingCredential, error) {

	// Group uri
	uri := ""
	if servicetype == models.DMSRabbitMQServiceName {
		uri = fmt.Sprintf("%s:%d", host, port)
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
		URI:      uri,
		Type:     servicetype,
	}
	return bc, nil
}

func SyncStatusWithService(b *DMSBroker, instanceID string, serviceID string, planID string,
	targetID string) (database.InstanceDetails, error, error) {
	dbInstance := database.InstanceDetails{}
	// Log opts
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService dms instance opts: instanceID: %s serviceID: %s " +
		"planID: %s targetID: %s", instanceID, serviceID, planID, targetID))

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return dbInstance, fmt.Errorf("SyncStatusWithService create dms client failed. Error: %s", err), nil
	}
	// Invoke sdk get
	instance, serviceErr := instances.Get(dmsClient, targetID).Extract()
	if serviceErr != nil {
		return dbInstance, nil, serviceErr
	}
	// get InstanceDetails in back database
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, serviceID, planID).
		First(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService get dms instance in back database failed. Error: %s", err))
		return dbInstance, fmt.Errorf("SyncStatusWithService get dms instance failed. Error: %s", err), nil
	}
	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService dms instance in back database: %v", models.ToJson(dbInstance)))
	// update target info in back database
	targetInfo, err := json.Marshal(instance)
	if err != nil {
		return dbInstance, fmt.Errorf("SyncStatusWithService marshal dms instance failed. Error: %s", err), nil
	}

	dbInstance.TargetID = instance.InstanceID
	dbInstance.TargetName = instance.Name
	dbInstance.TargetStatus = instance.Status
	dbInstance.TargetInfo = string(targetInfo)

	err = database.BackDBConnection.Save(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update dms instance target status in back database failed. " +
			"Error: %s", err))
		return dbInstance, fmt.Errorf("SyncStatusWithService update dms instance target status failed. Error: %s", err), nil
	}
	// Sync target status success
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update dms instance target status succeed: %s", instanceID))
	return dbInstance, nil, nil
}
