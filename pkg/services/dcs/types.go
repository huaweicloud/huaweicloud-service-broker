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

// MetadataParameters represent plan metadata parameters in config
type MetadataParameters struct {
	Engine            string   `json:"engine,omitempty"`
	EngineVersion     string   `json:"engine_version,omitempty"`
	SpecCode          string   `json:"speccode,omitempty"`
	ChargingType      string   `json:"charging_type,omitempty"`
	Capacity          int      `json:"capacity,omitempty"`
	VPCID             string   `json:"vpc_id,omitempty"`
	SubnetID          string   `json:"subnet_id,omitempty"`
	SecurityGroupID   string   `json:"security_group_id,omitempty"`
	AvailabilityZones []string `json:"availability_zones,omitempty"`
}

// ProvisionParameters represent provision parameters
type ProvisionParameters struct {
	Capacity                 int      `json:"capacity,omitempty"`
	VPCID                    string   `json:"vpc_id,omitempty"`
	SubnetID                 string   `json:"subnet_id,omitempty"`
	SecurityGroupID          string   `json:"security_group_id,omitempty"`
	AvailabilityZones        []string `json:"availability_zones,omitempty"`
	Username                 string   `json:"username,omitempty"`
	Password                 string   `json:"password,omitempty"`
	Name                     string   `json:"name,omitempty"`
	Description              string   `json:"description,omitempty"`
	BackupStrategySavedays   int      `json:"backup_strategy_savedays,omitempty"`
	BackupStrategyBackupType string   `json:"backup_strategy_backup_type,omitempty"`
	BackupStrategyBackupAt   []int    `json:"backup_strategy_backup_at,omitempty"`
	BackupStrategyBeginAt    string   `json:"backup_strategy_begin_at,omitempty"`
	BackupStrategyPeriodType string   `json:"backup_strategy_period_type,omitempty"`
	MaintainBegin            string   `json:"maintain_begin,omitempty"`
	MaintainEnd              string   `json:"maintain_end,omitempty"`
}

// UpdateParameters represent update parameters
type UpdateParameters struct {
	Name                     string  `json:"name,omitempty"`
	Description              *string `json:"description,omitempty"`
	BackupStrategySavedays   int     `json:"backup_strategy_savedays,omitempty"`
	BackupStrategyBackupType string  `json:"backup_strategy_backup_type,omitempty"`
	BackupStrategyBackupAt   []int   `json:"backup_strategy_backup_at,omitempty"`
	BackupStrategyBeginAt    string  `json:"backup_strategy_begin_at,omitempty"`
	BackupStrategyPeriodType string  `json:"backup_strategy_period_type,omitempty"`
	MaintainBegin            string  `json:"maintain_begin,omitempty"`
	MaintainEnd              string  `json:"maintain_end,omitempty"`
	SecurityGroupID          string  `json:"security_group_id,omitempty"`
	NewCapacity              int     `json:"new_capacity,omitempty"`
	OldPassword              *string `json:"old_password,omitempty"`
	NewPassword              *string `json:"new_password,omitempty"`
}

const (
	// AddtionalParamPassword for password
	AddtionalParamPassword string = "password"
)
