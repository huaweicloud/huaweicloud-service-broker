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

	// Init provisionOpts
	provisionOpts := instances.CreateOps{}
	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &provisionOpts)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error unmarshalling rawParameters: %s", err)
		}
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Setting provisionOpts
	if provisionOpts.Name == "" {
		provisionOpts.Name = instanceID
	}

	// TODO need to confirm different engine name
	if servicePlan.Name == models.DCSRedisServiceName {
		provisionOpts.Engine = "Redis"
	} else if servicePlan.Name == models.DCSMemcachedServiceName {
		provisionOpts.Engine = "Memcached"
	} else if servicePlan.Name == models.DCSIMDGServiceName {
		provisionOpts.Engine = "IMDG"
	} else {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("unknown service name: %s", servicePlan.Name)
	}

	// TODO get default vpc
	if provisionOpts.VPCID == "" {

	}

	// TODO get default security group
	if provisionOpts.SecurityGroupID == "" {

	}

	// TODO get default Subnet
	if provisionOpts.SubnetID == "" {

	}

	// Get default AvailableZones
	if len(provisionOpts.AvailableZones) == 0 {
		// List all the azs in this region
		azs, err := availablezones.Get(dcsClient).Extract()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dcs availablezones failed. Error: %s", err)
		}

		// Choose the first one Still have available resources in this az
		for _, az := range azs.AvailableZones {
			if az.ResourceAvailability == "true" {
				provisionOpts.AvailableZones = []string{az.ID}
				break
			}
		}
	}

	// 	Get default Product
	if provisionOpts.ProductID == "" {
		// List all the products
		ps, err := products.Get(dcsClient).Extract()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dcs products failed. Error: %s", err)
		}

		// Choose the first one
		for _, p := range ps.Products {
			if p.ProductID != "" {
				provisionOpts.ProductID = p.ProductID
				break
			}
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dcs instance opts: %v", provisionOpts))

	// Invoke sdk
	dcsInstance, err := instances.Create(dcsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision dcs instance result: %v", dcsInstance))

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
	addtionalparam["password"] = provisionOpts.Password

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
	b.Logger.Debug(fmt.Sprintf("create dcs instance in back database opts: %v", idsOpts))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dcs instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create dcs instance in back database succeed: %s", instanceID))

	// Return result
	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
