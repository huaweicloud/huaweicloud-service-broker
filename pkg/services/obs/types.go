package obs

import (
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"gopkg.in/mgo.v2/bson"
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
	BucketName    string                 `json:"bucket_name,omitempty" bson:"bucket_name,omitempty"`
	BucketPolicy  string                 `json:"bucket_policy,omitempty" bson:"bucket_policy,omitempty"`
	UnknownFields map[string]interface{} `json:"-" bson:",inline"`
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
	b, _ = bson.Marshal(&j)
	return bson.Unmarshal(b, f)
}

// UpdateParameters represent update parameters
type UpdateParameters struct {
	BucketPolicy string `json:"bucket_policy,omitempty"`
	Status       string `json:"status,omitempty"`
}

const (
	// AddtionalParamRequest for request
	AddtionalParamRequest string = "request"
)
