package instance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/availablezones"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/products"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *DMSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Check dms instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check dms instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dms client failed. Error: %s", err)
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
				brokerapi.NewFailureResponse(fmt.Errorf("Parameters are not following schema"),
					http.StatusBadRequest, "Parameters are not following schema")
		}
	}

	// List all the products
	ps, err := products.Get(dmsClient, metadataParameters.Engine).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dms products failed. Error: %s", err)
	}

	// Get ProductID and Storage Space
	productID := ""
	var storageSpace int
	var productParams []products.Parameter
	if metadataParameters.ChargingType == "Hourly" {
		productParams = ps.Hourly
	} else {
		productParams = ps.Monthly
	}
	// Go over
	for _, p := range productParams {
		for _, v := range p.Values {
			for _, d := range v.Details {
				// Single
				if d.SpecCode == metadataParameters.SpecCode {
					productID = d.ProductID
					storageSpace, err = strconv.Atoi(d.Storage)
					if err != nil {
						// Default 100
						storageSpace = 100
						// Log result
						b.Logger.Debug(fmt.Sprintf("get dms storage space failed. Error: %s", err))
					}
					break
				}
				// Cluster
				for _, pi := range d.ProductInfos {
					if pi.SpecCode == metadataParameters.SpecCode {
						productID = pi.ProductID
						storageSpace, err = strconv.Atoi(pi.Storage)
						if err != nil {
							// Default 200
							storageSpace = 200
							b.Logger.Debug(fmt.Sprintf("get dms storage space failed. Error: %s", err))
						}
					}
				}
			}
		}
	}

	// Init provisionOpts
	provisionOpts := instances.CreateOps{}
	// instance name
	provisionOpts.Name = provisionParameters.Name
	// instance description
	provisionOpts.Description = provisionParameters.Description
	// instance engine
	provisionOpts.Engine = metadataParameters.Engine
	// Default 3.7.0
	provisionOpts.EngineVersion = metadataParameters.EngineVersion
	// Default 100
	provisionOpts.StorageSpace = storageSpace
	// Password
	provisionOpts.Password = provisionParameters.Password
	// AccessUser
	provisionOpts.AccessUser = provisionParameters.Username
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
		azs, err := availablezones.Get(dmsClient).Extract()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dms availablezones failed. Error: %s", err)
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
	// ProductID
	provisionOpts.ProductID = productID
	// MaintainBegin
	if provisionParameters.MaintainBegin != "" {
		provisionOpts.MaintainBegin = provisionParameters.MaintainBegin
	}
	// MaintainEnd
	if provisionParameters.MaintainEnd != "" {
		provisionOpts.MaintainEnd = provisionParameters.MaintainEnd
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dms instance opts: %v", models.ToJson(provisionOpts)))

	// Invoke sdk
	dmsInstance, err := instances.Create(dmsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{},
			fmt.Errorf("provision dms instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision dms instance result: %v", models.ToJson(dmsInstance)))

	// Invoke sdk get
	freshInstance, err := instances.Get(dmsClient, dmsInstance.InstanceID).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dms instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal dms instance failed. Error: %s", err)
	}

	// Constuct addtional info
	addtionalparam := map[string]string{}
	addtionalparam[AddtionalParamUsername] = provisionOpts.AccessUser
	addtionalparam[AddtionalParamPassword] = provisionOpts.Password

	// Marshal addtional info
	addtionalinfo, err := json.Marshal(addtionalparam)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal dms addtional info failed. Error: %s", err)
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
	b.Logger.Debug(fmt.Sprintf("create dms instance in back database opts: %v", models.ToJson(idsOpts)))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dms instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create dms instance in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDMSInstance {
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
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("convert dms instance operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dms instance operation datas: %s", operationdata))

		// Create OperationDetails
		err = database.BackDBConnection.Create(&ods).Error
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create operation in back database failed. Error: %s", err)
		}

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: ""}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
