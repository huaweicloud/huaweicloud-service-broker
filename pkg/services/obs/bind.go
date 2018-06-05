package obs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *OBSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("bind obs bucket opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Invoke sdk
	bindOpts := &obs.GetBucketMetadataInput{}
	bindOpts.Bucket = instanceID
	obsResponse, err := obsClient.GetBucketMetadata(bindOpts)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("get obs bucket failed. Error: %s", err)
	}
	// Log result
	b.Logger.Debug(fmt.Sprintf("get obs bucket success: %v", obsResponse))

	// Find service
	service, err := b.Catalog.FindService(details.ServiceID)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("find obs service failed. Error: %s", err)
	}

	// Build Binding Credential
	credential, err := BuildBindingCredential(
		b.CloudCredentials.Region,
		obsClient.GetEndpoint(),
		bindOpts.Bucket,
		b.CloudCredentials.AccessKey,
		b.CloudCredentials.SecretKey,
		service.Name)
	if err != nil {
		return brokerapi.Binding{}, fmt.Errorf("build obs bucket binding credential failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("bind obs bucket success: %v", credential))

	// Return result
	return brokerapi.Binding{Credentials: credential}, nil
}
