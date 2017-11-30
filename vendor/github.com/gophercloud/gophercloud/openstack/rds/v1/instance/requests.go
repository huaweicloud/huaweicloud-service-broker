package rds

import (
	"github.com/gophercloud/gophercloud"
)

var RequestOpts gophercloud.RequestOpts = gophercloud.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

//CreateOpsBuilder is used for creating instance parameters.
//any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

type UpdateOpsBuilder interface {
	ToInstanceUpdateMap() (map[string]interface{}, error)
}

type UpdateOps struct {
	Volume map[string]interface{} `json:"volume"`
}

//CreateOps is a struct that contains all the parameters.
type CreateOps struct {

	//instance name
	Name string `json:"name" required:"true"`

	//data store
	Datastore map[string]string `json:"datastore,omitempty"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Volume map[string]interface{} `json:"volume,omitempty"`

	Region string `json:"region,omitempty"`

	AvailabilityZone string `json:"availabilityZone,omitempty"`

	Vpc string `json:"vpc,omitempty"`

	Nics map[string]string `json:"nics,omitempty"`

	SecurityGroup map[string]string `json:"securityGroup,omitempty"`

	DbPort string `json:"dbPort,omitempty"`

	BackupStrategy map[string]interface{} `json:"backupStrategy,omitempty"`

	DbRtPd string `json:"dbRtPd,omitempty"`

	//ReplicaOf string `json:"replicaOf" required:"true"`
}

func (ops CreateOps) ToInstanceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(ops, "instance")
}

func (ops UpdateOps) ToInstanceUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(ops, "resize")
}

//Create a instance with given parameters.
func Create(client *gophercloud.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	RequestOpts.OkCodes = []int{202}
	RequestOpts.JSONBody = nil
	_, r.Err = client.Post(createURL(client), b, &r.Body, &RequestOpts)

	return
}

func UpdateVolumeSize(client *gophercloud.ServiceClient, ops UpdateOpsBuilder, id string) (r UpdateResult) {
	b, err := ops.ToInstanceUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	RequestOpts.OkCodes = []int{202}
	RequestOpts.JSONBody = nil
	_, r.Err = client.Post(updateURL(client, id), b, &r.Body, &RequestOpts)

	return
}

//delete a instance via id
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	RequestOpts.OkCodes = []int{202}
	RequestOpts.JSONBody = nil
	_, r.Err = client.Delete(deleteURL(client, id), &RequestOpts)
	return
}

//get a instance with detailed information by id
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	RequestOpts.OkCodes = []int{200}
	RequestOpts.JSONBody = nil
	_, r.Err = client.Get(getURL(client, id), &r.Body, &RequestOpts)
	return
}

//list all the instances
func List(client *gophercloud.ServiceClient) (r ListResult) {
	RequestOpts.OkCodes = []int{200}
	RequestOpts.JSONBody = nil
	_, r.Err = client.Get(listURL(client), &r.Body, &RequestOpts)
	return
}
