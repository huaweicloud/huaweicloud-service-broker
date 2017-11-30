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
        keystoneEndpoint := "https://iam.huaweiclouds.com/v3" //TODO set the endpoint of keystone.
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
                Username:         "Username",
                //UserID:           "cb3cfa12219b47f5809e864b3d511ff5",
                Password:         "Password",
                DomainID:         "",
                DomainName:       "DomainName",
                //TenantName:       "eu-de",
                //Scope:            '{"project": {"name": "eu-de"}}',                 
                Scope:            tokens3.Scope{ProjectName: "eu-de", DomainName: "DomainName"},    
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

        //r := rds.Create(sc, rds.CreateOps{Name: "rds_name", FlavorRef: "c1", ReplicaOf: "test"})
        //if r.Err != nil {
        //        //TODO log the error
        //        fmt.Println("03 ERR '%s' not found", r)
        //        fmt.Println("error 03")
        //        return
        //}

        //var instance rds.Instance
        //err3 := r.ExtractInto(&instance)
        //if err3 != nil {
        //        //TODO deal with the error
        //        fmt.Println("error 04")
        //        return 
        //}
}
