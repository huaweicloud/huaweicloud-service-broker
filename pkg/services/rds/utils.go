package rds

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different rds instance
func BuildBindingCredential(
	host string,
	port int,
	name string,
	username string,
	password string,
	servicetype string) (BindingCredential, error) {

	var uri string

	if servicetype == models.RDSPostgresqlServiceName {
		// Postgresql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSMysqlServiceName {
		// Mysql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSSqlserverServiceName {
		// Sqlserver
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSHwsqlServiceName {
		// Hwsql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		Host:     host,
		Port:     port,
		Name:     name,
		UserName: username,
		Password: password,
		URI:      uri,
		Type:     servicetype,
	}
	return bc, nil
}

func SyncStatusWithReal(b *RDSBroker, instanceID string, serviceID string, planID string,
	targetID string) (database.InstanceDetails, error, error) {
	ids := database.InstanceDetails{}
	// Log opts
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithReal rds instance opts: instanceID: %s serviceID: %s " +
		"planID: %s targetID: %s", instanceID, serviceID, planID, targetID))

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return ids, fmt.Errorf("SyncStatusWithReal create rds client failed. Error: %s", err), nil
	}
	// Invoke sdk get
	instance, err := instances.Get(rdsClient, targetID).Extract()
	if err != nil {
		return ids, nil, err
	}
	// get InstanceDetails in back database
	findErr := database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, serviceID, planID).
		First(&ids).Error
	if findErr != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithReal get rds instance in back database failed. Error: %s", findErr))
		return ids, fmt.Errorf("SyncStatusWithReal get rds instance failed. Error: %s", findErr), nil
	}
	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithReal rds instance in back database: %v", models.ToJson(ids)))
	// update target info in back database
	targetinfo, err := json.Marshal(instance)
	if err != nil {
		return ids, fmt.Errorf("SyncStatusWithReal marshal rds instance failed. Error: %s", err), nil
	}

	ids.TargetID = instance.ID
	ids.TargetName = instance.Name
	ids.TargetStatus = instance.Status
	ids.TargetInfo = string(targetinfo)

	updateErr := database.BackDBConnection.Save(&ids).Error
	if updateErr != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithReal update rds instance target status in back database failed. " +
			"Error: %s", updateErr))
		return ids, fmt.Errorf("SyncStatusWithReal update rds instance target status failed. Error: %s", updateErr), nil
	}
	// Sync target status success
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithReal update rds instance target status succeed: %s", instanceID))
	return ids, nil, nil
}
