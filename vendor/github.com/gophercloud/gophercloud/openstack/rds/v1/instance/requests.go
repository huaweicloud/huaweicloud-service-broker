package rds

import (
	"github.com/gophercloud/gophercloud"
)

type CreateOpsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

type CreateOps struct {
	Name string `json:"name" required:"true"`

	Datastore map[string]interface{} `json:"datastore,omitempty"`

	FlavorRef string `json:"flavorRef" required:"true"`

	Volume map[string]interface{} `json:"volume,omitempty"`

	Region string `json:"region,omitempty"`

	AvailabilityZone string `json:"availability,omitempty"`

	Vpc string `json:"vpc,omitempty"`

	Nics map[string]string `json:"nics,omitempty"`

	SecurityGroup map[string]string `json:"security,omitempty"`

	DbRtPd string `json:"dbRtPd,omitempty"`

	ReplicaOf string `json:"replicaOf" required:"true"`
}

func (ops CreateOps) ToInstanceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(ops, "instance")
}

//Create a Instance
func Create(client *gophercloud.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

//delete a instance
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

//get a instance
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

type UpdateVolumeBuilder interface {
	ToUpdateVolumeMap() (map[string]interface{}, error)
}

type UpdateVolumeOps struct {
	Volume UpdateVolumeSize `json:"volume"`
}

type UpdateVolumeSize struct {
	Size int `json:"size" required:"true"`
}

func (ops UpdateVolumeOps) ToUpdateVolumeMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(ops, "resize")
}

func UpdateVolume(client *gophercloud.ServiceClient, id string, ops UpdateVolumeBuilder) (r UpdateVolumeResult) {
	b, err := ops.ToUpdateVolumeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(updateVolumeURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
