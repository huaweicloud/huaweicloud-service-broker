package obs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *OBSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Init obs client
	obsClient, err := b.CloudCredentials.OBSClient()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create obs client failed. Error: %s", err)
	}
	// Close obs client
	if obsClient != nil {
		defer obsClient.Close()
	}

	// Init provisionOpts
	var params map[string]string
	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &params)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error unmarshalling rawParameters: %s", err)
		}
	}

	provisionOpts := &obs.CreateBucketInput{}
	// Setting Bucket
	provisionOpts.Bucket = params["bucket_name"]
	if provisionOpts.Bucket == "" {
		provisionOpts.Bucket = instanceID
	}
	// Setting StorageClass
	provisionOpts.StorageClass = obs.StorageClassType(params["storage_class"])
	if provisionOpts.StorageClass == "" {
		provisionOpts.StorageClass = obs.StorageClassStandard
	}
	// Setting ACL
	provisionOpts.ACL = obs.AclType(params["acl"])
	if provisionOpts.ACL == "" {
		provisionOpts.ACL = obs.AclPrivate
	}
	// Setting Location
	provisionOpts.Location = params["location"]

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dcs bucket opts: %v", provisionOpts))

	// Invoke sdk
	obsResponse, err := obsClient.CreateBucket(provisionOpts)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision obs bucket failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision obs bucket result: %v", obsResponse))

	// Return result
	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
