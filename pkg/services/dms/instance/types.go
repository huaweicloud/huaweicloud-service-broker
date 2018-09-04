package instance

import (
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"gopkg.in/mgo.v2/bson"
)

// DMSBroker define
type DMSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// BindingCredential represent dms binding credential
type BindingCredential struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URI      string `json:"uri,omitempty"`
	Type     string `json:"type,omitempty"`
}

// MetadataParameters represent plan metadata parameters in config
type MetadataParameters struct {
	Engine            string   `json:"engine,omitempty"`
	EngineVersion     string   `json:"engine_version,omitempty"`
	SpecCode          string   `json:"speccode,omitempty"`
	ChargingType      string   `json:"charging_type,omitempty"`
	VPCID             string   `json:"vpc_id,omitempty"`
	SubnetID          string   `json:"subnet_id,omitempty"`
	SecurityGroupID   string   `json:"security_group_id,omitempty"`
	AvailabilityZones []string `json:"availability_zones,omitempty"`
}

// ProvisionParameters represent provision parameters
type ProvisionParameters struct {
	VPCID             string                 `json:"vpc_id,omitempty"`
	SubnetID          string                 `json:"subnet_id,omitempty"`
	SecurityGroupID   string                 `json:"security_group_id,omitempty"`
	AvailabilityZones []string               `json:"availability_zones,omitempty"`
	Username          string                 `json:"username,omitempty"`
	Password          string                 `json:"password,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Description       string                 `json:"description,omitempty"`
	MaintainBegin     string                 `json:"maintain_begin,omitempty"`
	MaintainEnd       string                 `json:"maintain_end,omitempty"`
	UnknownFields     map[string]interface{} `json:"-" bson:",inline"`
}

// Collect unknown fields into "UnknownFields"
func (f *ProvisionParameters) UnmarshalJSON(b []byte) error {
	var j map[string]interface{}
	json.Unmarshal(b, &j)
	b, _ = bson.Marshal(&j)
	return bson.Unmarshal(b, f)
}

// UpdateParameters represent update parameters
type UpdateParameters struct {
	Name            string  `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	MaintainBegin   string  `json:"maintain_begin,omitempty"`
	MaintainEnd     string  `json:"maintain_end,omitempty"`
	SecurityGroupID string  `json:"security_group_id,omitempty"`
}

const (
	// AddtionalParamUsername for username
	AddtionalParamUsername string = "username"
	// AddtionalParamPassword for password
	AddtionalParamPassword string = "password"
)
