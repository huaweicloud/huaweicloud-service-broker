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

func SyncStatusWithService(b *RDSBroker, instanceID string, serviceID string, planID string,
	targetID string) (*instances.Instance, error, error) {
	dbInstance := database.InstanceDetails{}
	// Log opts
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService rds instance opts: instanceID: %s serviceID: %s "+
		"planID: %s targetID: %s", instanceID, serviceID, planID, targetID))

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return nil, fmt.Errorf("SyncStatusWithService create rds client failed. Error: %s", err), nil
	}
	// Invoke sdk get
	instance, serviceErr := instances.Get(rdsClient, targetID).Extract()
	if serviceErr != nil {
		return nil, nil, serviceErr
	}

	// Check other update operation
	if instance.Status == "ACTIVE" {
		// Check OperationDetails length in back database
		var length int
		err := database.BackDBConnection.
			Model(&database.OperationDetails{}).
			Where("instance_id = ? and operation_type = ?", instanceID, models.OperationUpdating).
			Count(&length).Error
		if err != nil {
			return nil, nil, err
		}

		// last OperationDetails is exist
		if length > 0 {
			// Log length
			b.Logger.Debug(fmt.Sprintf("instance: %s operation: %s length: %d", instanceID, models.OperationUpdating, length))

			// get last OperationDetails in back database
			ods := database.OperationDetails{}
			err := database.BackDBConnection.
				Where("instance_id = ? and operation_type = ?", instanceID, models.OperationUpdating).
				Last(&ods).Error
			if err != nil {
				return nil, nil, err
			}
			// Get additional info from OperationDetails
			addtionalparamdetail := map[string]string{}
			err = ods.GetAdditionalInfo(&addtionalparamdetail)
			if err != nil {
				return nil, nil, err
			}
			// Check AddtionalParamFlavorID exist
			if _, ok := addtionalparamdetail[AddtionalParamFlavorID]; ok {
				if (addtionalparamdetail[AddtionalParamFlavorID] != "") &&
					(addtionalparamdetail[AddtionalParamFlavorID] != instance.Flavor.Id) {
					// Log begin
					b.Logger.Debug(fmt.Sprintf("update rds instance %s flavor %s begin.",
						targetID, addtionalparamdetail[AddtionalParamFlavorID]))

					// UpdateFlavorRef
					rdsInstance, err := instances.UpdateFlavorRef(
						rdsClient,
						instances.UpdateFlavorOps{
							FlavorRef: addtionalparamdetail[AddtionalParamFlavorID],
						},
						targetID).Extract()
					if err != nil {
						return nil, nil, fmt.Errorf("update rds instance flavor failed. Error: %s", err)
					}

					// Log result
					b.Logger.Debug(fmt.Sprintf("update rds instance flavor result: %v", models.ToJson(rdsInstance)))

					// Invoke sdk get again
					instance, serviceErr = instances.Get(rdsClient, targetID).Extract()
					if serviceErr != nil {
						return nil, nil, serviceErr
					}

					// Remove AddtionalParamFlavorID
					delete(addtionalparamdetail, AddtionalParamFlavorID)
					// Marshal addtional info
					addtionalinfo, err := json.Marshal(addtionalparamdetail)
					if err != nil {
						// Just log
						b.Logger.Debug(fmt.Sprintf("marshal rds addtional info failed. Error: %s", err))
					}
					ods.AdditionalInfo = string(addtionalinfo)
					err = database.BackDBConnection.Save(&ods).Error
					if err != nil {
						// Just log
						b.Logger.Debug(fmt.Sprintf("Remove %s from OperationDetails in back database failed.", AddtionalParamFlavorID))
					}
				}
			}
		}
	}

	// get InstanceDetails in back database
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, serviceID, planID).
		First(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService get rds instance in back database failed. Error: %s", err))
		return instance, fmt.Errorf("SyncStatusWithService get rds instance failed. Error: %s", err), nil
	}
	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService rds instance in back database: %v", models.ToJson(dbInstance)))
	// update target info in back database
	targetInfo, err := json.Marshal(instance)
	if err != nil {
		return instance, fmt.Errorf("SyncStatusWithService marshal rds instance failed. Error: %s", err), nil
	}

	dbInstance.TargetID = instance.ID
	dbInstance.TargetName = instance.Name
	dbInstance.TargetStatus = instance.Status
	dbInstance.TargetInfo = string(targetInfo)

	err = database.BackDBConnection.Save(&dbInstance).Error
	if err != nil {
		b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update rds instance target status in back database failed. "+
			"Error: %s", err))
		return instance, fmt.Errorf("SyncStatusWithService update rds instance target status failed. Error: %s", err), nil
	}
	// Sync target status success
	b.Logger.Debug(fmt.Sprintf("SyncStatusWithService update rds instance target status succeed: %s", instanceID))
	return instance, nil, nil
}
