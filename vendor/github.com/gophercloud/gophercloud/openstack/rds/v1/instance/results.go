package rds

import "github.com/gophercloud/gophercloud"

type Instance struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Name             string `json:"name"`
	Created          string `json:"created"`
	HostName         string `json:"hostname"`
	Type             string `json:"type"`
	Region           string `json:"region"`
	Updated          string `json:"updated"`
	AvailabilityZone string `json:"availabilityZone"`
	Vpc              string `json:"vpc"`
	Nics             map[string]string `json:"nics""`
	SecurityGroup    map[string]string `json:"securityGroup""`
	Flavor           map[string]string `json:"flavor""`
	Volume           map[string]interface{} `json:"volume""`
	DbPort           int64  `json:"dbPort"`
	DataStoreInfo    map[string]string `json:"dataStoreInfo""`
	Extendparam      map[string]interface{} `json:"extendparam"`
	BackupStrategy   map[string]interface{} `json:"backupStrategy"`
	//SlaveId          string `json:"slaveId"`
	//ReplicationMode  string `json:"ha.replicationMode"`
	//ReplicaOf        string `json:"replica_of"`
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

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type ListResult struct {
	gophercloud.Result
}

func (lr ListResult) Extract() ([]Instance, error) {
	var a struct {
		Instances []Instance `json:"instances"`
	}
	err := lr.Result.ExtractInto(&a)
	return a.Instances, err
}
