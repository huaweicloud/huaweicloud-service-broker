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

// Provision implematation
func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

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
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
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
				fmt.Errorf("Error unmarshalling rawParameters from details: %s", err)
		}
	}

	// Get datastoresList
	datastoresList, err := datastores.List(rdsClient, metadataParameters.DatastoreType).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve datastores: %s", err)
	}
	if len(datastoresList) < 1 {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no datastore result")
	}
	b.Logger.Debug(fmt.Sprintf("provision rds instance datastoresList: %v", models.ToJson(datastoresList)))

	// Get datastoreID
	var datastoreID string
	for _, datastore := range datastoresList {
		if datastore.Name == metadataParameters.DatastoreVersion {
			datastoreID = datastore.ID
			break
		}
	}
	if datastoreID == "" {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no datastore ID")
	}
	b.Logger.Debug(fmt.Sprintf("Received datastore ID: %s", datastoreID))

	// Get flavorsList
	flavorsList, err := flavors.List(rdsClient, datastoreID, b.CloudCredentials.Region).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("Unable to retrieve flavors: %s", err)
	}
	if len(flavorsList) < 1 {
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
	for _, flavor := range flavorsList {
		if flavor.SpecCode == specCode {
			flavorID = flavor.ID
			break
		}
	}
	if flavorID == "" {
		return brokerapi.ProvisionedServiceSpec{},
			errors.New("Returned no flavor Id")
	}
	b.Logger.Debug(fmt.Sprintf("Received flavor ID: %s", flavorID))

	// Init provisionOpts
	provisionOpts := instances.CreateOps{}
	provisionOpts.Name = provisionParameters.Name
	provisionOpts.DataStore = instances.DataStoreOps{
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
	provisionOpts.Volume = instances.VolumeOps{
		Type: volumeType,
		Size: volumeSize,
	}
	provisionOpts.Region = b.CloudCredentials.Region
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
	provisionOpts.Vpc = vpcID
	// Default SubnetID
	subnetID := metadataParameters.SubnetID
	if provisionParameters.SubnetID != "" {
		subnetID = provisionParameters.SubnetID
	}
	provisionOpts.Nics = instances.NicsOps{
		SubnetId: subnetID,
	}
	// Default SecurityGroupID
	securitygroupID := metadataParameters.SecurityGroupID
	if provisionParameters.SecurityGroupID != "" {
		securitygroupID = provisionParameters.SecurityGroupID
	}
	provisionOpts.SecurityGroup = instances.SecurityGroupOps{
		Id: securitygroupID,
	}
	if provisionParameters.DatabasePort != "" {
		provisionOpts.DbPort = provisionParameters.DatabasePort
	}
	if provisionParameters.BackupStrategyStarttime != "" {
		provisionOpts.BackupStrategy = instances.BackupStrategyOps{}
		provisionOpts.BackupStrategy.StartTime = provisionParameters.BackupStrategyStarttime
		if provisionParameters.BackupStrategyKeepdays > 0 {
			provisionOpts.BackupStrategy.KeepDays = provisionParameters.BackupStrategyKeepdays
		}
	} else {
		// Default Value
		provisionOpts.BackupStrategy.StartTime = "00:00:00"
		provisionOpts.BackupStrategy.KeepDays = 0
	}
	provisionOpts.DbRtPd = provisionParameters.DatabasePassword
	provisionOpts.Ha = instances.HaOps{}
	provisionOpts.Ha.Enable = provisionParameters.HAEnable
	if provisionParameters.HAReplicationMode != "" {
		provisionOpts.Ha.ReplicationMode = provisionParameters.HAReplicationMode
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision rds instance opts: %v", models.ToJson(provisionOpts)))

	// Invoke sdk
	rdsInstance, err := instances.Create(rdsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision rds instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision rds instance result: %v", models.ToJson(rdsInstance)))

	// Invoke sdk get
	freshInstance, err := instances.Get(rdsClient, rdsInstance.ID).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal rds instance failed. Error: %s", err)
	}

	// Constuct addtional info
	addtionalparam := map[string]string{}
	addtionalparam[AddtionalParamDBUsername] = metadataParameters.DatabaseUsername
	addtionalparam[AddtionalParamDBPassword] = provisionOpts.DbRtPd

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
		TargetID:       freshInstance.ID,
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
		ods := models.OperationDatas{
			OperationType:  models.OperationProvisioning,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       freshInstance.ID,
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

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: operationdata}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
