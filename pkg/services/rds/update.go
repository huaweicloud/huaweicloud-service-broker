package rds

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/datastores"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/flavors"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Check rds instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length == 0 {
		return brokerapi.UpdateServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// get InstanceDetails in back database
	ids := database.InstanceDetails{}
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		First(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("rds instance in back database: %v", models.ToJson(ids)))

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init updateParameters
	updateParameters := UpdateParameters{}
	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &updateParameters)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update rds instance opts: %v", models.ToJson(updateParameters)))

	// Invoke sdk
	if updateParameters.VolumeSize > 0 {
		// UpdateVolumeSize
		rdsInstance, err := instances.UpdateVolumeSize(
			rdsClient,
			instances.UpdateOps{
				Volume: map[string]interface{}{
					"size": updateParameters.VolumeSize,
				},
			},
			ids.TargetID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance volume size failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance volume size result: %v", models.ToJson(rdsInstance)))
	}

	// Invoke sdk
	if updateParameters.SpecCode != "" {

		// Find service plan
		servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
		}

		// Get parameters from service plan metadata
		metadataParameters := MetadataParameters{}
		if servicePlan.Metadata != nil {
			if len(servicePlan.Metadata.Parameters) > 0 {
				err := json.Unmarshal(servicePlan.Metadata.Parameters, &metadataParameters)
				if err != nil {
					return brokerapi.UpdateServiceSpec{},
						fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
				}
			}
		}

		// Get datastoresList
		datastoresList, err := datastores.List(rdsClient, metadataParameters.DatastoreType).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve datastores: %s", err)
		}
		if len(datastoresList) < 1 {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no datastore result")
		}
		b.Logger.Debug(fmt.Sprintf("update rds datastores opts: %v", models.ToJson(datastoresList)))

		// Get datastoreID
		var datastoreID string
		for _, datastore := range datastoresList {
			if datastore.Name == metadataParameters.DatastoreVersion {
				datastoreID = datastore.ID
				break
			}
		}
		if datastoreID == "" {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no datastore ID")
		}
		b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", datastoreID))

		// Get flavorsList
		flavorsList, err := flavors.List(rdsClient, datastoreID, b.CloudCredentials.Region).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve flavors: %s", err)
		}
		if len(flavorsList) < 1 {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no flavor result")
		}

		// Get flavorID
		var flavorID string
		for _, flavor := range flavorsList {
			if flavor.SpecCode == updateParameters.SpecCode {
				flavorID = flavor.ID
				break
			}
		}
		if flavorID == "" {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no flavor Id")
		}
		b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", flavorID))

		// UpdateFlavorRef
		rdsInstance, err := instances.UpdateFlavorRef(
			rdsClient,
			instances.UpdateFlavorOps{
				FlavorRef: flavorID,
			},
			ids.TargetID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance flavor failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance flavor result: %v", models.ToJson(rdsInstance)))
	}

	// Invoke sdk get
	freshInstance, err := instances.Get(rdsClient, ids.TargetID).Extract()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal rds instance failed. Error: %s", err)
	}

	ids.TargetID = freshInstance.ID
	ids.TargetName = freshInstance.Name
	ids.TargetStatus = freshInstance.Status
	ids.TargetInfo = string(targetinfo)

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("update rds instance in back database opts: %s", models.ToJson(ids)))

	err = database.BackDBConnection.Save(&ids).Error
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("update rds instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncRDS {
		// OperationDatas for OperationUpdating
		ods := models.OperationDatas{
			OperationType:  models.OperationUpdating,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       ids.TargetID,
			TargetName:     ids.TargetName,
			TargetStatus:   ids.TargetStatus,
			TargetInfo:     ids.TargetInfo,
			AdditionalInfo: ids.AdditionalInfo,
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("convert rds instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create rds instance operation datas: %s", operationdata))

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: operationdata}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
