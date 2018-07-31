package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/logger"

	"github.com/gophercloud/gophercloud"
	nativeopenstack "github.com/gophercloud/gophercloud/openstack"
)

// CloudCredentials define
type CloudCredentials struct {
	AccessKey        string `json:"access_key"`
	SecretKey        string `json:"secret_key"`
	CACertFile       string `json:"cacert_file"`
	ClientCertFile   string `json:"cert"`
	ClientKeyFile    string `json:"key"`
	DomainID         string `json:"domain_id"`
	DomainName       string `json:"domain_name"`
	EndpointType     string `json:"endpoint_type"`
	IdentityEndpoint string `json:"auth_url"`
	Insecure         bool   `json:"insecure"`
	Password         string `json:"password"`
	Region           string `json:"region"`
	TenantID         string `json:"tenant_id"`
	TenantName       string `json:"tenant_name"`
	Token            string `json:"token"`
	Username         string `json:"user_name"`
	UserID           string `json:"user_id"`

	OpenStackClient *gophercloud.ProviderClient
	CloudClient     *golangsdk.ProviderClient
}

// Validate CloudCredentials
func (c *CloudCredentials) Validate() error {
	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	err := c.newOpenStackClient()
	if err != nil {
		return err
	}

	return c.newCloudClient()
}

// newOpenStackClient returns new native openstack client
func (c *CloudCredentials) newOpenStackClient() error {
	ao := gophercloud.AuthOptions{
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
		IdentityEndpoint: c.IdentityEndpoint,
		Password:         c.Password,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		TokenID:          c.Token,
		Username:         c.Username,
		UserID:           c.UserID,
		// allow to renew tokens
		AllowReauth: true,
	}

	client, err := nativeopenstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := ReadContents(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := ReadContents(c.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := ReadContents(c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &logger.LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = nativeopenstack.Authenticate(client, ao)
	if err != nil {
		return err
	}

	c.OpenStackClient = client

	return nil
}

// newCloudClient returns new cloud client
func (c *CloudCredentials) newCloudClient() error {
	ao := golangsdk.AuthOptions{
		DomainID:         c.DomainID,
		DomainName:       c.DomainName,
		IdentityEndpoint: c.IdentityEndpoint,
		Password:         c.Password,
		TenantID:         c.TenantID,
		TenantName:       c.TenantName,
		TokenID:          c.Token,
		Username:         c.Username,
		UserID:           c.UserID,
		// allow to renew tokens
		AllowReauth: true,
	}

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := ReadContents(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure {
		config.InsecureSkipVerify = true
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := ReadContents(c.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := ReadContents(c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	// if OS_DEBUG is set, log the requests and responses
	var osDebug bool
	if os.Getenv("OS_DEBUG") != "" {
		osDebug = true
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: &logger.LogRoundTripper{
			Rt:      transport,
			OsDebug: osDebug,
		},
	}

	err = openstack.Authenticate(client, ao)
	if err != nil {
		return err
	}

	c.CloudClient = client

	return nil
}

// getNativeEndpointType returns native openstack endpoint type
func (c *CloudCredentials) getNativeEndpointType() gophercloud.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return gophercloud.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return gophercloud.AvailabilityAdmin
	}
	return gophercloud.AvailabilityPublic
}

// getEndpointType returns cloud endpoint type
func (c *CloudCredentials) getEndpointType() golangsdk.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return golangsdk.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return golangsdk.AvailabilityAdmin
	}
	return golangsdk.AvailabilityPublic
}

// NetworkingV2Client return native networking v2 client
func (c *CloudCredentials) NetworkingV2Client() (*gophercloud.ServiceClient, error) {
	return nativeopenstack.NewNetworkV2(c.OpenStackClient, gophercloud.EndpointOpts{
		Region:       c.Region,
		Availability: c.getNativeEndpointType(),
	})
}

// RDSV1Client return rds v1 client
func (c *CloudCredentials) RDSV1Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewRdsServiceV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
}

// DCSV1Client return dcs v1 client
func (c *CloudCredentials) DCSV1Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewDCSServiceV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
}

// DMSV1Client return dms v1 client
func (c *CloudCredentials) DMSV1Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewDMSServiceV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
}

// OBSClient return obs client
func (c *CloudCredentials) OBSClient() (*obs.ObsClient, error) {
	sc, err := openstack.NewOBSService(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
	if err != nil {
		return nil, err
	}

	return obs.New(c.AccessKey, c.SecretKey, sc.ServiceURL())
}
