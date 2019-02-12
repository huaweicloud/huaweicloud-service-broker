package instance

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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
	VPCID             string                 `json:"vpc_id,omitempty" bson:"vpc_id,omitempty"`
	SubnetID          string                 `json:"subnet_id,omitempty" bson:"subnet_id,omitempty"`
	SecurityGroupID   string                 `json:"security_group_id,omitempty" bson:"security_group_id,omitempty"`
	AvailabilityZones []string               `json:"availability_zones,omitempty" bson:"availability_zones,omitempty"`
	Username          string                 `json:"username,omitempty" bson:"username,omitempty"`
	Password          string                 `json:"password,omitempty" bson:"password,omitempty"`
	Name              string                 `json:"name,omitempty" bson:"name,omitempty"`
	Description       string                 `json:"description,omitempty" bson:"description,omitempty"`
	MaintainBegin     string                 `json:"maintain_begin,omitempty" bson:"maintain_begin,omitempty"`
	MaintainEnd       string                 `json:"maintain_end,omitempty" bson:"maintain_end,omitempty"`
	UnknownFields     map[string]interface{} `json:"-" bson:",inline"`
}

func (f *ProvisionParameters) MarshalJSON() ([]byte, error) {
	var j interface{}
	b, _ := bson.Marshal(f)
	bson.Unmarshal(b, &j)
	return json.Marshal(&j)
}

// Collect unknown fields into "UnknownFields"
func (f *ProvisionParameters) UnmarshalJSON(b []byte) error {
	var j map[string]interface{}
	json.Unmarshal(b, &j)
	fmt.Printf("DMS UnmarshalJSON ProvisionParameters: %v\n", j)
	// Compatibles Array and String for availability_zones
	if j["availability_zones"] != nil {
		t := reflect.TypeOf(j["availability_zones"]).Kind()
		fmt.Printf("DMS UnmarshalJSON availability_zones type: %v\n", t)
		if t == reflect.String {
			str := FormatStr(j["availability_zones"].(string))
			if str != "" {
				j["availability_zones"] = strings.Split(str, ",")
				fmt.Printf("DMS UnmarshalJSON availability_zones value: %v\n", j["availability_zones"])
			}
		}
	}

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

func FormatStr(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, `"`, "", -1)
	str = strings.TrimPrefix(str, "[")
	str = strings.TrimSuffix(str, "]")
	return str
}

const (
	// AddtionalParamUsername for username
	AddtionalParamUsername string = "username"
	// AddtionalParamPassword for password
	AddtionalParamPassword string = "password"
	// AddtionalParamRequest for request
	AddtionalParamRequest string = "request"
)
