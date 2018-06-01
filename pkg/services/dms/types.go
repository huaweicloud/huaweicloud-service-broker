package dms

import (
	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
)

// DMSBroker define
type DMSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// BindingCredential represent dms binding credential
type BindingCredential struct {
	Region    string `json:"region,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
	URL       string `json:"url,omitempty"`
	AK        string `json:"ak,omitempty"`
	SK        string `json:"sk,omitempty"`
	Type      string `json:"type,omitempty"`
}
