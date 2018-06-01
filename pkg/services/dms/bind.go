package dms

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *DMSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create dms client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind dms instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk: the default includeDeadLetter value is false
	_, err = queues.Get(dmsClient, instanceID, false).Extract()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get dms instance failed. Error: %s", err)
	}

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find dms service failed. Error: %s", err)
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(
		b.CloudCredentials.Region,
		dmsClient.Endpoint,
		b.CloudCredentials.TenantID,
		b.CloudCredentials.AccessKey,
		b.CloudCredentials.SecretKey,
		service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build dms instance binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind dms instance success: %v", credential))

	// Return result
	return brokerapi.Binding{Credentials: credential}, nil
}
