package queue

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
	Region       string `json:"region,omitempty"`
	ProjectID    string `json:"projectid,omitempty"`
	ProtocolType string `json:"protocoltype,omitempty"`
	URL          string `json:"url,omitempty"`
	AK           string `json:"ak,omitempty"`
	SK           string `json:"sk,omitempty"`
	QueueID      string `json:"queueid,omitempty"`
	GroupID      string `json:"groupid,omitempty"`
	Type         string `json:"type,omitempty"`
}

// MetadataParameters represent plan metadata parameters in config
type MetadataParameters struct {
	QueueMode       string `json:"queue_mode,omitempty"`
	EndpointName    string `json:"endpoint_name,omitempty"`
	EndpointPort    string `json:"endpoint_port,omitempty"`
	RedrivePolicy   string `json:"redrive_policy,omitempty"`
	MaxConsumeCount int    `json:"max_consume_count,omitempty"`
	RetentionHours  int    `json:"retention_hours,omitempty"`
}

// ProvisionParameters represent provision parameters
type ProvisionParameters struct {
	RedrivePolicy   string                 `json:"redrive_policy,omitempty"`
	MaxConsumeCount int                    `json:"max_consume_count,omitempty"`
	RetentionHours  int                    `json:"retention_hours,omitempty"`
	QueueName       string                 `json:"queue_name,omitempty"`
	GroupName       string                 `json:"group_name,omitempty"`
	Description     string                 `json:"description,omitempty"`
	UnknownFields   map[string]interface{} `json:"-" bson:",inline"`
}

// Collect unknown fields into "UnknownFields"
func (f *ProvisionParameters) UnmarshalJSON(b []byte) error {
	var j map[string]interface{}
	json.Unmarshal(b, &j)
	b, _ = bson.Marshal(&j)
	return bson.Unmarshal(b, f)
}

const (
	// ProtocolTypeHTTPS for HTTPS
	ProtocolTypeHTTPS string = "HTTPS"
	// ProtocolTypeTCP for TCP
	ProtocolTypeTCP string = "TCP"
)
