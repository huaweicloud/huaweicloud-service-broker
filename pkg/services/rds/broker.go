package rds

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"github.com/pivotal-cf/brokerapi"
)

// RDSBroker define
type RDSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// Provision implematation
func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init provisionOpts
	provisionOpts := instances.CreateOps{
		/*Name:             instanceID,
		Datastore:        map[string]string{"type": rds_prop.DatastoreType, "version": rds_prop.DatastoreVersion},
		FlavorRef:        rds_prop.FlavorId,
		Volume:           map[string]interface{}{"type": rds_prop.VolumeType, "size": rds_prop.VolumeSize},
		Region:           rds_prop.Region,
		AvailabilityZone: rds_prop.AvailabilityZone,
		Vpc:              rds_prop.VpcId,
		Nics:             map[string]string{"subnetId": rds_prop.SubnetId},
		SecurityGroup:    map[string]string{"id": rds_prop.SecurityGroupId},
		DbPort:           rds_prop.Dbport,
		BackupStrategy:   map[string]interface{}{"startTime": rds_prop.BackupStrategyStarttime, "keepDays": rds_prop.BackupStrategyKeepdays},
		DbRtPd:           rds_prop.Dbpassword,*/
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision rds instance opts: %v", provisionOpts))

	// Invoke sdk
	rdsInstance, err := instances.Create(rdsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision rds instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision rds instance result: %v", rdsInstance))

	// Return result
	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}

// Deprovision implematation
func (b *RDSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision rds instance opts: %s", instanceID))

	// Invoke sdk
	result := instances.Delete(rdsClient, instanceID)
	if result.Err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision rds instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision rds instance success: %s", instanceID))

	// Return result
	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}

// Bind implematation
func (b *RDSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind rds instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	instance, err := instances.Get(rdsClient, instanceID).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find rds service failed. Error: %s", err)
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(instance.HostName, instance.DbPort, "DbName", "Dbusername", "Dbpassword", service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build rds instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind rds instance success: %v", credential))

	// Return result
	return brokerapi.Binding{Credentials: credential}, nil
}

// Unbind implematation
func (b *RDSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind rds instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// TODO do something for unbind

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind rds instance success: instanceID: %s bindingID: %s", instanceID, bindingID))

	return nil
}

// Update implematation
func (b *RDSBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.UpdateServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init rawParameters
	rawParameters := map[string]string{}

	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &rawParameters)
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("Error unmarshalling parameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("update rds instance opts: %v", rawParameters))

	// Invoke sdk
	if volumesize, ok := rawParameters["volumesize"]; ok {
		// UpdateVolumeSize
		rdsInstance, err := instances.UpdateVolumeSize(
			rdsClient,
			instances.UpdateOps{
				Volume: map[string]interface{}{
					"size": volumesize,
				},
			},
			instanceID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance volume size failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance volume size result: %v", rdsInstance))
	}

	// Invoke sdk
	if flavorname, ok := rawParameters["flavorname"]; ok {
		// UpdateFlavorRef
		rdsInstance, err := instances.UpdateFlavorRef(
			rdsClient,
			instances.UpdateFlavorOps{
				// TODO convert flavor name to flavor ref
				FlavorRef: flavorname,
			},
			instanceID).Extract()
		if err != nil {
			return brokerapi.UpdateServiceSpec{}, fmt.Errorf("update rds instance flavor failed. Error: %s", err)
		}

		// Log result
		b.Logger.Debug(fmt.Sprintf("update rds instance flavor result: %v", rdsInstance))
	}

	// Return result
	return brokerapi.UpdateServiceSpec{IsAsync: false, OperationData: ""}, nil
}

// LastOperation implematation
func (b *RDSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
