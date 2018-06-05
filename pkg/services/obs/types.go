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
