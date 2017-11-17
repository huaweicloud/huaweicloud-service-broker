package main

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	trusts "github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/rds/v1/instance"
)

func main() {
	keystoneEndpoint := "" //TODO set the endpoint of keystone.
	pc, err := openstack.NewClient(keystoneEndpoint)
	if err != nil {
		return
	}
	eo := gophercloud.EndpointOpts{Region: "ReionOne", Availability: gophercloud.AvailabilityPublic} //TODO change the two parameters
	opts := tokens3.AuthOptions{
		IdentityEndpoint: "",
		Username:         "",
		UserID:           "",
		Password:         "",
		DomainID:         "",
		DomainName:       "",
		AllowReauth:      true,
	}
	authOptsExt := trusts.AuthOptsExt{
		TrustID:            "", //TODO config the trust id
		AuthOptionsBuilder: &opts,
	}
	//authenticate
	err = openstack.AuthenticateV3(pc, authOptsExt, eo)

	if err != nil {
		return
	}

	sc, _ := openstack.NewRdsServiceV1(pc, eo)

	r := rds.Create(sc, rds.CreateOps{Name: "", FlavorRef: "", ReplicaOf: ""})
	if r.Err != nil {
		//TODO log the error
		return
	}

	var instance rds.Instance
	err3 := r.ExtractInto(&instance)
	if err3 != nil {
		//TODO deal with the error
		return
	}

	//TODO deal with the instance object
}
