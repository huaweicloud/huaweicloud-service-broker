package main

import (
	"crypto/tls"
	"net/http"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	trusts "github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/rds/v1/instance"
	netutil "k8s.io/apimachinery/pkg/util/net"
	certutil "k8s.io/client-go/util/cert"
)

func main() {
	keystoneEndpoint := "https://iam.eu-de.otc.t-systems.com/v3" //TODO set the endpoint of keystone.
	pc, err := openstack.NewClient(keystoneEndpoint)
	if err != nil {
		fmt.Println("error 01")
		return
	}

	roots, err := certutil.NewPool("/root/ca/ca.crt")
	if err != nil {
		fmt.Println("error 001")
		return
	}
	config := &tls.Config{}
	config.RootCAs = roots
	pc.HTTPClient.Transport = netutil.SetOldTransportDefaults(&http.Transport{TLSClientConfig: config})

	eo := gophercloud.EndpointOpts{Region: "eu-de", Availability: gophercloud.AvailabilityPublic} //TODO change the two parameters
	opts := tokens3.AuthOptions{
		IdentityEndpoint: "https://iam.eu-de.otc.t-systems.com/v3",
		Username:         "swx414799",
		//UserID:           "cb3cfa12219b47f5809e864b3d511ff5",
		Password:         "Huawei@1234",
		DomainID:         "",
		DomainName:       "swx414799",
		//TenantName:       "eu-de",
		//Scope:            '{"project": {"name": "eu-de"}}',
		Scope:            tokens3.Scope{ProjectName: "eu-de", DomainName: "swx414799"},
		AllowReauth:      true,
	}
	authOptsExt := trusts.AuthOptsExt{
		TrustID:            "", //TODO config the trust id
		AuthOptionsBuilder: &opts,
	}
	//authenticate
	err = openstack.AuthenticateV3(pc, authOptsExt, gophercloud.EndpointOpts{})

	if err != nil {
		fmt.Println(err)
		fmt.Println("error 02")
		return
	}

	sc, err33 := openstack.NewRdsServiceV1(pc, eo)
	fmt.Println("sc '%s' #####", sc)
	if err33 != nil{
		fmt.Println(err33)
		fmt.Println("error 33")
		return
	}
	r2 := rds.List(sc)

	fmt.Println("r2 '%s' call list", r2)

	r := rds.Create(sc, rds.CreateOps{
		Name: "rds_name_666",
		Datastore: map[string]string{"type": "PostgreSQL", "version": "9.5.5"},
		FlavorRef: "7fbf27c5-07e5-43dc-cf13-ad7a0f1c5d9a",
		Volume: map[string]interface{}{"type": "COMMON", "size": 100},
		Region: "eu-de",
		AvailabilityZone: "eu-de-01",
		Vpc: "bf693499-6d9e-49b7-83b1-b2e5e156c7f0",
		Nics: map[string]string{"subnetId": "1ea3b3a0-9689-4f9e-88b0-7a82fc538d4d"},
		SecurityGroup: map[string]string{"id": "dc3ec145-9029-4b39-b5a3-ace5a01f772b"},
		DbPort: "8635",
		BackupStrategy: map[string]interface{}{"startTime": "00:00:00", "keepDays": 0},
		DbRtPd: "Huangwei!120521"})
	if r.Err != nil {
	        //TODO log the error
	        fmt.Println("03 ERR '%s' not found", r)
	        fmt.Println("error 03")
	        return
	}

	fmt.Println("r '%s' call create", r)

	//var instance rds.Instance
	instance, err3 := r.Extract()
	if err3 != nil {
		fmt.Println(err3)
		fmt.Println("error 04")
		return
	}

	fmt.Println("instance '%s' call create", instance)


}
