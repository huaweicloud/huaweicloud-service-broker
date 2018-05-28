package rds

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

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
