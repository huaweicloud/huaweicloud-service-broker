package rds

import "github.com/gophercloud/gophercloud"

type Instance struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Name             string `json:"name"`
	Created          string `json:"created"`
	HostName         string `json:"hostname"`
	Type             string `json:"type"`
	Region           string `json:"string"`
	Updated          string `json:"updated"`
	AvailabilityZone string `json:"availabilityZone"`
	Vpc              string `json:"vpc"`
	SubnetId         string `json:"nics.subnetId"`
	SecurityGroupId  string `json:"securityGroup.id"`
	FlavorId         string `json:"flavor.id"`
	Volume           string `json:"volume"`
	DbPort           int    `json:"dbPort"`
	DataStoreInfo    string `json:"dataStoreInfo"`
	JobsId           string `json:"extendParam.jobs.id"`
	BackupStrategy   string `json:"backupStrategy"`
	SlaveId          string `json:"slaveId"`
	ReplicationMode  string `json:"ha.replicationMode"`
	ReplicaOf        string `json:"replica_of"`
}

// Extract will get the Volume object out of the commonResult object.
func (r commonResult) Extract() (*Instance, error) {
	var s Instance
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "instance")
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateVolumeResult struct {
	commonResult
}
