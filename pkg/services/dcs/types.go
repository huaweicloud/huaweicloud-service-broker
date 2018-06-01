package dcs

import (
	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
)

// DCSBroker define
type DCSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// BindingCredential represent dcs binding credential
type BindingCredential struct {
	IP       string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
}
