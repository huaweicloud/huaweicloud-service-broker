package obs

import (
	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
)

// OBSBroker define
type OBSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// BindingCredential represent obs binding credential
type BindingCredential struct {
	Region     string `json:"region,omitempty"`
	URL        string `json:"url,omitempty"`
	BucketName string `json:"bucketname,omitempty"`
	AK         string `json:"ak,omitempty"`
	SK         string `json:"sk,omitempty"`
	Type       string `json:"type,omitempty"`
}

// MetadataParameters represent plan metadata parameters in config
type MetadataParameters struct {
	StorageClass string `json:"storage_class,omitempty"`
	BucketPolicy string `json:"bucket_policy,omitempty"`
}

// ProvisionParameters represent provision parameters
type ProvisionParameters struct {
	BucketName   string `json:"bucket_name,omitempty"`
	BucketPolicy string `json:"bucket_policy,omitempty"`
}

// UpdateParameters represent update parameters
type UpdateParameters struct {
	BucketPolicy string `json:"bucket_policy,omitempty"`
	Status       string `json:"status,omitempty"`
}
