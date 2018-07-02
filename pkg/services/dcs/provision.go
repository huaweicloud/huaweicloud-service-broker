package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/availablezones"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/products"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *DCSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Check dcs instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check dcs instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find dcs service failed. Error: %s", err)
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

	// List all the products
	ps, err := products.Get(dcsClient).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dcs products failed. Error: %s", err)
	}

	// Get ProductID
	productID := ""
	for _, p := range ps.Products {
		if (p.SpecCode == metadataParameters.SpecCode) &&
			(p.ChargingType == metadataParameters.ChargingType) {
			productID = p.ProductID
			break
		}
	}

	// Init provisionOpts
	provisionOpts := instances.CreateOps{}
	provisionOpts.Name = provisionParameters.Name
	if provisionParameters.Description != "" {
		provisionOpts.Description = provisionParameters.Description
	}
	provisionOpts.Engine = metadataParameters.Engine
	if (metadataParameters.EngineVersion != "") &&
		(service.Name == models.DCSRedisServiceName) {
		provisionOpts.EngineVersion = metadataParameters.EngineVersion
	} else {
		provisionOpts.EngineVersion = ""
	}

	// Default Capacity
	provisionOpts.Capacity = metadataParameters.Capacity
	if provisionParameters.Capacity > 0 {
		provisionOpts.Capacity = provisionParameters.Capacity
	}
	// Username
	if provisionParameters.Username != "" {
		if (service.Name == models.DCSMemcachedServiceName) ||
			(service.Name == models.DCSIMDGServiceName) {
			provisionOpts.AccessUser = provisionParameters.Username
			provisionOpts.NoPasswordAccess = "false"
		}
	}
	provisionOpts.Password = provisionParameters.Password
	// Default VPCID
	provisionOpts.VPCID = metadataParameters.VPCID
	if provisionParameters.VPCID != "" {
		provisionOpts.VPCID = provisionParameters.VPCID
	}
	// Default SubnetID
	provisionOpts.SubnetID = metadataParameters.SubnetID
	if provisionParameters.SubnetID != "" {
		provisionOpts.SubnetID = provisionParameters.SubnetID
	}
	// Default SecurityGroupID
	provisionOpts.SecurityGroupID = metadataParameters.SecurityGroupID
	if provisionParameters.SecurityGroupID != "" {
		provisionOpts.SecurityGroupID = provisionParameters.SecurityGroupID
	}
	// Default AvailabilityZones
	provisionOpts.AvailableZones = metadataParameters.AvailabilityZones
	if len(provisionParameters.AvailabilityZones) > 0 {
		provisionOpts.AvailableZones = provisionParameters.AvailabilityZones
	}
	// Convert AvailabilityZones from code to id
	if len(provisionOpts.AvailableZones) > 0 {
		// List all the azs in this region
		azs, err := availablezones.Get(dcsClient).Extract()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dcs availablezones failed. Error: %s", err)
		}
		// Convert code to id
		azIDs := []string{}
		for _, az := range azs.AvailableZones {
			for _, azCode := range provisionOpts.AvailableZones {
				if az.Code == azCode {
					azIDs = append(azIDs, az.ID)
				}
			}
		}
		provisionOpts.AvailableZones = azIDs
	}
	provisionOpts.ProductID = productID
	// BackupStrategy only supported in Master/standby
	if (provisionParameters.BackupStrategySavedays > 0) &&
		(provisionParameters.BackupStrategyBackupType != "") &&
		(provisionParameters.BackupStrategyBeginAt != "") &&
		(provisionParameters.BackupStrategyPeriodType != "") &&
		(len(provisionParameters.BackupStrategyBackupAt) > 0) {
		provisionOpts.InstanceBackupPolicy = &instances.InstanceBackupPolicy{
			SaveDays:   provisionParameters.BackupStrategySavedays,
			BackupType: provisionParameters.BackupStrategyBackupType,
			PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
				BeginAt:    provisionParameters.BackupStrategyBeginAt,
				PeriodType: provisionParameters.BackupStrategyPeriodType,
				BackupAt:   provisionParameters.BackupStrategyBackupAt,
			},
		}
	}
	// MaintainBegin
	if provisionParameters.MaintainBegin != "" {
		provisionOpts.MaintainBegin = provisionParameters.MaintainBegin
	}
	// MaintainEnd
	if provisionParameters.MaintainEnd != "" {
		provisionOpts.MaintainEnd = provisionParameters.MaintainEnd
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dcs instance opts: %v", models.ToJson(provisionOpts)))

	// Invoke sdk
	dcsInstance, err := instances.Create(dcsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision dcs instance result: %v", models.ToJson(dcsInstance)))

	// Invoke sdk get
	freshInstance, err := instances.Get(dcsClient, dcsInstance.InstanceID).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dcs instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal dcs instance failed. Error: %s", err)
	}

	// Constuct addtional info
	addtionalparam := map[string]string{}
	addtionalparam[AddtionalParamPassword] = provisionOpts.Password

	// Marshal addtional info
	addtionalinfo, err := json.Marshal(addtionalparam)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal dcs addtional info failed. Error: %s", err)
	}

	// create InstanceDetails in back database
	idsOpts := database.InstanceDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		TargetID:       freshInstance.InstanceID,
		TargetName:     freshInstance.Name,
		TargetStatus:   freshInstance.Status,
		TargetInfo:     string(targetinfo),
		AdditionalInfo: string(addtionalinfo),
	}

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("create dcs instance in back database opts: %v", models.ToJson(idsOpts)))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dcs instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create dcs instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDCS {
		// OperationDatas for OperationProvisioning
		ods := database.OperationDetails{
			OperationType:  models.OperationProvisioning,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       freshInstance.InstanceID,
			TargetName:     freshInstance.Name,
			TargetStatus:   freshInstance.Status,
			TargetInfo:     string(targetinfo),
			AdditionalInfo: string(addtionalinfo),
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("convert dcs instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dcs instance operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: ""}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
