package rds

import (
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"gopkg.in/mgo.v2/bson"
)

// RDSBroker define
type RDSBroker struct {
	CloudCredentials config.CloudCredentials
	Catalog          config.Catalog
	Logger           lager.Logger
}

// BindingCredential represent rds binding credential
type BindingCredential struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URI      string `json:"uri,omitempty"`
	Type     string `json:"type,omitempty"`
}

// MetadataParameters represent plan metadata parameters in config
type MetadataParameters struct {
	DatastoreType    string `json:"datastore_type,omitempty"`
	DatastoreVersion string `json:"datastore_version,omitempty"`
	SpecCode         string `json:"speccode,omitempty"`
	VolumeType       string `json:"volume_type,omitempty"`
	VolumeSize       int    `json:"volume_size,omitempty"`
	AvailabilityZone string `json:"availability_zone,omitempty"`
	VPCID            string `json:"vpc_id,omitempty"`
	SubnetID         string `json:"subnet_id,omitempty"`
	SecurityGroupID  string `json:"security_group_id,omitempty"`
	DatabaseUsername string `json:"database_username,omitempty"`
}

// ProvisionParameters represent provision parameters
type ProvisionParameters struct {
	SpecCode                string                 `json:"speccode,omitempty"`
	VolumeType              string                 `json:"volume_type,omitempty"`
	VolumeSize              int                    `json:"volume_size,omitempty"`
	AvailabilityZone        string                 `json:"availability_zone,omitempty"`
	VPCID                   string                 `json:"vpc_id,omitempty"`
	SubnetID                string                 `json:"subnet_id,omitempty"`
	SecurityGroupID         string                 `json:"security_group_id,omitempty"`
	Name                    string                 `json:"name,omitempty"`
	DatabasePort            string                 `json:"database_port,omitempty"`
	DatabasePassword        string                 `json:"database_password,omitempty"`
	BackupStrategyStarttime string                 `json:"backup_strategy_starttime,omitempty"`
	BackupStrategyKeepdays  int                    `json:"backup_strategy_keepdays,omitempty"`
	HAEnable                bool                   `json:"ha_enable,omitempty"`
	HAReplicationMode       string                 `json:"ha_replicationmode,omitempty"`
	UnknownFields           map[string]interface{} `json:"-" bson:",inline"`
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
	VolumeSize int    `json:"volume_size,omitempty"`
	SpecCode   string `json:"speccode,omitempty"`
}

const (
	// AddtionalParamDBUsername for dbusername
	AddtionalParamDBUsername string = "dbusername"
	// AddtionalParamDBPassword for dbpassword
	AddtionalParamDBPassword string = "dbpassword"
)
