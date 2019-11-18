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

// Provision implematation
func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	// Check accepts_incomplete if this service support async
	if models.OperationAsyncRDS {
		e := b.Catalog.ValidateAcceptsIncomplete(asyncAllowed)
		if e != nil {
			return brokerapi.ProvisionedServiceSpec{}, e
		}
	}

	// Check rds instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		// Get InstanceDetails in back database
		iddetail := database.InstanceDetails{}
		err = database.BackDBConnection.
			Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
			First(&iddetail).Error
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get instance in back database failed. Error: %s", err)
		}

		// Get additional info from InstanceDetails
		addtionalparamdetail := map[string]string{}
		err = iddetail.GetAdditionalInfo(&addtionalparamdetail)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get instance additional info failed. Error: %s", err)
		}

		// Check AddtionalParamRequest exist
		if _, ok := addtionalparamdetail[AddtionalParamRequest]; ok {
			if (addtionalparamdetail[AddtionalParamRequest] != "") &&
				(addtionalparamdetail[AddtionalParamRequest] == string(details.RawParameters)) {
				return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExistsSame
			}
		}

		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV3Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds v3 client failed. Error: %s", err)
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Get parameters from service plan metadata
	metadataParameters := MetadataParameters{}
	if servicePlan.Metadata != nil {
		if len(servicePlan.Metadata.Parameters) > 0 {
			err := json.Unmarshal(servicePlan.Metadata.Parameters, &metadataParameters)
			if err != nil {
				return brokerapi.ProvisionedServiceSpec{},
					fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
			}
		}
	}

	// Get parameters from details
	provisionParameters := ProvisionParameters{}
	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &provisionParameters)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{},
				brokerapi.NewFailureResponse(fmt.Errorf("Error unmarshalling rawParameters from details: %s", err),
					http.StatusBadRequest, "Error unmarshalling rawParameters")
		}
		// Exist other unknown fields,
		if len(provisionParameters.UnknownFields) > 0 {
			return brokerapi.ProvisionedServiceSpec{},
				brokerapi.NewFailureResponse(
					fmt.Errorf("Parameters are not following schema: %+v", provisionParameters.UnknownFields),
					http.StatusBadRequest, "Parameters are not following schema")
		}
	}

	// Get datastoresList
	pages, err := datastores.List(rdsClient, metadataParameters.DatastoreType).AllPages()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve datastores: %s", err)
	}

	datastoresList, err := datastores.ExtractDataStores(pages)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve datastores: %s", err)
	}
	if len(datastoresList.DataStores) < 1 {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no datastore result")
	}
	b.Logger.Debug(fmt.Sprintf("provision rds instance datastoresList: %v", models.ToJson(datastoresList)))

	// Get datastoreID
	var datastoreID string
	for _, datastore := range datastoresList.DataStores {
		if datastore.Name == metadataParameters.DatastoreVersion {
			datastoreID = datastore.Id
			break
		}
	}
	if datastoreID == "" {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no datastore ID")
	}
	b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", datastoreID))

	// Get flavorsList
	dbflavorOpt := flavors.DbFlavorsOpts{}
	dbflavorOpt.Versionname = metadataParameters.DatastoreVersion
	flavorPages, err := flavors.List(rdsClient, dbflavorOpt, metadataParameters.DatastoreType).AllPages()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve flavors: %s", err)
	}

	flavorsList, err := flavors.ExtractDbFlavors(flavorPages)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve flavors: %s", err)
	}
	if len(flavorsList.Flavorslist) < 1 {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no flavor result")
	}
	b.Logger.Debug(fmt.Sprintf("provision rds instance flavorsList: %v", models.ToJson(flavorsList)))

	// Get flavorID
	var flavorID string
	// Default SpecCode
	specCode := metadataParameters.SpecCode
	if provisionParameters.SpecCode != "" {
		specCode = provisionParameters.SpecCode
	}
	for _, flavor := range flavorsList.Flavorslist {
		if flavor.Speccode == specCode {
			flavorID = flavor.Speccode
			break
		}
	}
	if flavorID == "" {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no flavor Id")
	}
	b.Logger.Debug(fmt.Sprintf("Received flavor ID: %s", flavorID))

	var replicaOpts instances.CreateReplicaOpts
	var provisionOpts instances.CreateRdsOpts
	if provisionParameters.ReplicaOfId != "" {
		volumeType := metadataParameters.VolumeType
		if provisionParameters.VolumeType != "" {
			volumeType = provisionParameters.VolumeType
		}

		volumeSize := metadataParameters.VolumeSize
		if provisionParameters.VolumeSize > 0 {
			volumeSize = provisionParameters.VolumeSize
		}
		replicaOpts = instances.CreateReplicaOpts{
			Name:        provisionParameters.Name,
			ReplicaOfId: provisionParameters.ReplicaOfId,
			FlavorRef:   flavorID,
			Volume: &instances.Volume{
				Type: volumeType,
				Size: volumeSize,
			},
			Region:           b.CloudCredentials.Region,
			AvailabilityZone: metadataParameters.AvailabilityZone,
		}
	} else {
		provisionOpts = getPrivisionOpts(provisionParameters, metadataParameters,
										 flavorID, b.CloudCredentials.Region)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision rds instance opts: %v", models.ToJson(provisionOpts)))
	// Create rds instance
	var rdsInstance *instances.CreateRds
	if provisionParameters.ReplicaOfId == "" {
		rdsInstance, err = instances.Create(rdsClient, provisionOpts).Extract()
	} else {
		rdsInstance, err = instances.CreateReplica(rdsClient, replicaOpts).Extract()
	}

	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision rds instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision rds instance result: %v", models.ToJson(rdsInstance)))

	listInstanceOpts := instances.ListRdsInstanceOpts{}
	listInstanceOpts.Id = rdsInstance.Instance.Id
	// Invoke sdk get
	instancePages, err := instances.List(rdsClient, listInstanceOpts).AllPages()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve instance: %s", err)
	}

	freshInstances, err := instances.ExtractRdsInstances(instancePages)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	if len(freshInstances.Instances) != 1 {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("The rds instance not exist or more than one.")
	}
	freshInstance := freshInstances.Instances[0]
	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal rds instance failed. Error: %s", err)
	}

	// Constuct addtional info
	addtionalparam := map[string]string{}
	addtionalparam[AddtionalParamDBUsername] = metadataParameters.DatabaseUsername
	addtionalparam[AddtionalParamDBPassword] = provisionOpts.Password
	addtionalparam[AddtionalParamRequest] = string(details.RawParameters)

	// Marshal addtional info
	addtionalinfo, err := json.Marshal(addtionalparam)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal rds addtional info failed. Error: %s", err)
	}

	// create InstanceDetails in back database
	idsOpts := database.InstanceDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		TargetID:       freshInstance.Id,
		TargetName:     freshInstance.Name,
		TargetStatus:   freshInstance.Status,
		TargetInfo:     string(targetinfo),
		AdditionalInfo: string(addtionalinfo),
	}

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("create rds instance in back database opts: %v", models.ToJson(idsOpts)))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create rds instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncRDS {
		// OperationDatas for OperationProvisioning
		ods := database.OperationDetails{
			OperationType:  models.OperationProvisioning,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       freshInstance.Id,
			TargetName:     freshInstance.Name,
			TargetStatus:   freshInstance.Status,
			TargetInfo:     string(targetinfo),
			AdditionalInfo: string(addtionalinfo),
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("convert rds instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create rds instance operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: ""}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}

func getPrivisionOpts(provisionParameters ProvisionParameters,
					  metadataParameters MetadataParameters,
					  flavorID string, region string) (instances.CreateRdsOpts) {
	provisionOpts := instances.CreateRdsOpts{}
	provisionOpts.Name = provisionParameters.Name
	provisionOpts.Datastore = &instances.Datastore{
		Type:    metadataParameters.DatastoreType,
		Version: metadataParameters.DatastoreVersion,
	}

	provisionOpts.FlavorRef = flavorID
	// Default VolumeType
	volumeType := metadataParameters.VolumeType
	if provisionParameters.VolumeType != "" {
		volumeType = provisionParameters.VolumeType
	}
	// Default VolumeSize
	volumeSize := metadataParameters.VolumeSize
	if provisionParameters.VolumeSize > 0 {
		volumeSize = provisionParameters.VolumeSize
	}

	provisionOpts.Volume = &instances.Volume{
		Type: volumeType,
		Size: volumeSize,
	}
	provisionOpts.Region = region
	// Default AvailabilityZone
	availabilityZone := metadataParameters.AvailabilityZone
	if provisionParameters.AvailabilityZone != "" {
		availabilityZone = provisionParameters.AvailabilityZone
	}
	provisionOpts.AvailabilityZone = availabilityZone
	// Default VPCID
	vpcID := metadataParameters.VPCID
	if provisionParameters.VPCID != "" {
		vpcID = provisionParameters.VPCID
	}
	provisionOpts.VpcId = vpcID
	// Default SubnetID
	subnetID := metadataParameters.SubnetID
	if provisionParameters.SubnetID != "" {
		subnetID = provisionParameters.SubnetID
	}
	provisionOpts.SubnetId = subnetID
	// Default SecurityGroupID
	securitygroupID := metadataParameters.SecurityGroupID
	if provisionParameters.SecurityGroupID != "" {
		securitygroupID = provisionParameters.SecurityGroupID
	}
	provisionOpts.SecurityGroupId = securitygroupID
	if provisionParameters.DatabasePort != "" {
		provisionOpts.Port = provisionParameters.DatabasePort
	}

	provisionOpts.BackupStrategy = &instances.BackupStrategy{}
	if provisionParameters.BackupStrategyStarttime != "" {
		provisionOpts.BackupStrategy.StartTime = provisionParameters.BackupStrategyStarttime
		if provisionParameters.BackupStrategyKeepdays > 0 {
			provisionOpts.BackupStrategy.KeepDays = provisionParameters.BackupStrategyKeepdays
		}
	} else {
		// Default Value
		provisionOpts.BackupStrategy.StartTime = "23:00-00:00"
		provisionOpts.BackupStrategy.KeepDays = 0
	}

	provisionOpts.Password = provisionParameters.DatabasePassword

	if provisionParameters.HAEnable == true {
		provisionOpts.Ha = &instances.Ha{}
		provisionOpts.Ha.Mode = "ha"
		if provisionParameters.HAReplicationMode != "" {
			provisionOpts.Ha.ReplicationMode = provisionParameters.HAReplicationMode
		}
	}

	return provisionOpts
}
