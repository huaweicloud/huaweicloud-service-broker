package dcs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *DCSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind dcs instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	instance, err := instances.Get(dcsClient, instanceID).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dcs instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find dcs service failed. Error: %s", err)
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(instance.IP, instance.Port, instance.UserName, "password", instance.Name, service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build dcs instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind dcs instance success: %v", credential))

	// Return result
	return brokerapi.Binding{Credentials: credential}, nil
}
