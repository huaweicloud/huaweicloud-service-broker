package rds

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/huaweicloud/golangsdk/openstack/rds/v3/datastores"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/flavors"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Update implematation
func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Check accepts_incomplete if this service support async
	if models.OperationAsyncRDS {
		e := b.Catalog.ValidateAcceptsIncomplete(asyncAllowed)
		if e != nil {
			return brokerapi.UpdateServiceSpec{}, e
		}
	}

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

	// sync and check status whether allowed to update
	instance, err, serviceErr := SyncStatusWithService(b, instanceID, details.ServiceID, details.PlanID, ids.TargetID)

	if err != nil || serviceErr != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("sync status failed. error: %s, service error: %s", err, serviceErr)
	}
	if instance.Status != "ACTIVE" {
		return brokerapi.UpdateServiceSpec{},
			brokerapi.NewFailureResponse(
				fmt.Errorf("Can only update rds instance in ACTIVE, but in: %s", instance.Status),
				http.StatusUnprocessableEntity, "Can only update rds instance in ACTIVE")
	}

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV3Client()
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

	// additionalInfo
	odAdditionalInfo := ids.AdditionalInfo

	// Invoke sdk
	isUpdateVolumeSize := false
	if updateParameters.VolumeSize > 0 {
		// Update in case of difference
		if instance.Volume.Size != updateParameters.VolumeSize {
			// UpdateVolumeSize
			rdsInstance, err := instances.EnlargeVolume(
				rdsClient,
				instances.EnlargeVolumeRdsOpts{
					EnlargeVolume: &instances.EnlargeVolumeSize{
						Size: updateParameters.VolumeSize,
					},
				},
				ids.TargetID).Extract()
			if err != nil {
				return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance volume size failed. Error: %s", err)
			}

			// Log result
			b.Logger.Debug(fmt.Sprintf("update rds instance volume size result: %v", models.ToJson(rdsInstance)))
			isUpdateVolumeSize = true
		}
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
		pages, err := datastores.List(rdsClient, metadataParameters.DatastoreType).AllPages()
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve datastores: %s", err)
		}

		datastoresList, err := datastores.ExtractDataStores(pages)
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve datastores: %s", err)
		}
		if len(datastoresList.DataStores) < 1 {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no datastore result")
		}
		b.Logger.Debug(fmt.Sprintf("update rds datastores opts: %v", models.ToJson(datastoresList.DataStores)))

		// Get datastoreID
		var datastoreID string
		for _, datastore := range datastoresList.DataStores {
			if datastore.Name == metadataParameters.DatastoreVersion {
				datastoreID = datastore.Id
				break
			}
		}
		if datastoreID == "" {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no datastore ID")
		}
		b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", datastoreID))

		// Get flavorsList
		dbflavorOpt := flavors.DbFlavorsOpts{}
		dbflavorOpt.Versionname = metadataParameters.DatastoreVersion
		flavorPages, err := flavors.List(rdsClient, dbflavorOpt, metadataParameters.DatastoreType).AllPages()
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve flavors: %s", err)
		}

		flavorsList, err := flavors.ExtractDbFlavors(flavorPages)
		if err != nil {
			return brokerapi.UpdateServiceSpec{},
				fmt.Errorf("Unable to retrieve flavors: %s", err)
		}
		if len(flavorsList.Flavorslist) < 1 {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no flavor result")
		}

		// Get flavorID
		var flavorID string
		for _, flavor := range flavorsList.Flavorslist {
			if flavor.Speccode == updateParameters.SpecCode {
				flavorID = flavor.Speccode
				break
			}
		}
		if flavorID == "" {
			return brokerapi.UpdateServiceSpec{},
				errors.New("Returned no flavor Id")
		}
		b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", flavorID))

		// Update in case of difference
		if instance.FlavorRef != flavorID {
			// If Update Volume Size is running
			if isUpdateVolumeSize {
				// Get additional info from InstanceDetails
				addtionalparamdetail := map[string]string{}
				err = ids.GetAdditionalInfo(&addtionalparamdetail)
				if err != nil {
					return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get instance additional info failed in update. Error: %s", err)
				}
				// Add AddtionalParamFlavorID
				addtionalparamdetail[AddtionalParamFlavorID] = flavorID
				// Marshal addtional info
				addtionalinfo, err := json.Marshal(addtionalparamdetail)
				if err != nil {
					return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal instance addtional info failed in update. Error: %s", err)
				}
				odAdditionalInfo = string(addtionalinfo)
			} else {
				// UpdateFlavorRef
				rdsInstance, err := instances.Resize(
					rdsClient,
					instances.ResizeFlavorOpts{
						ResizeFlavor: &instances.SpecCode{
							Speccode: flavorID,
						},
					},
					ids.TargetID).Extract()
				if err != nil {
					return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance flavor failed. Error: %s", err)
				}

				// Log result
				b.Logger.Debug(fmt.Sprintf("update rds instance flavor result: %v", models.ToJson(rdsInstance)))
			}
		}
	}

	// Invoke sdk list
	listInstanceOpts := instances.ListRdsInstanceOpts{}
	listInstanceOpts.Id = ids.TargetID
	instancePages, err := instances.List(rdsClient, listInstanceOpts).AllPages()
	if err != nil {
		return brokerapi.UpdateServiceSpec{},
			fmt.Errorf("Unable to retrieve instance: %s", err)
	}

	freshInstances, err := instances.ExtractRdsInstances(instancePages)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	if len(freshInstances.Instances) != 1 {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("The rds instance not exist or more than one.")
	}
	freshInstance := freshInstances.Instances[0]

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("marshal rds instance failed. Error: %s", err)
	}

	ids.TargetID = freshInstance.Id
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
		ods := database.OperationDetails{
			OperationType:  models.OperationUpdating,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       ids.TargetID,
			TargetName:     ids.TargetName,
			TargetStatus:   ids.TargetStatus,
			TargetInfo:     ids.TargetInfo,
			AdditionalInfo: odAdditionalInfo,
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("convert rds instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create rds instance operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.UpdateServiceSpec{IsAsync: true, OperationData: ""}, nil
	}

	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}
