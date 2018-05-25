package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/logger"
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
	Swauth           bool   `json:"swauth"`
	TenantID         string `json:"tenant_id"`
	TenantName       string `json:"tenant_name"`
	Token            string `json:"token"`
	Username         string `json:"user_name"`
	UserID           string `json:"user_id"`

	CloudClient *golangsdk.ProviderClient
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

	return c.newCloudClient()
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

	// If using Swift Authentication, there's no need to validate authentication normally.
	if !c.Swauth {
		err = openstack.Authenticate(client, ao)
		if err != nil {
			return err
		}
	}

	c.CloudClient = client

	// TODO: add s3 and obs auth with ak/sk

	return nil
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

// OBSV1Client return obs v1 client
func (c *CloudCredentials) OBSV1Client() (*golangsdk.ServiceClient, error) {
	// TODO: If Swift Authentication is being used, return a swauth client.
	/*if c.Swauth {
		return swauth.NewObjectStorageV1(c.CloudClient, swauth.AuthOpts{
			User: c.Username,
			Key:  c.Password,
		})
	}*/

	return openstack.NewObjectStorageV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
}

// RDSV1Client return rds v1 client
func (c *CloudCredentials) RDSV1Client() (*golangsdk.ServiceClient, error) {
	return openstack.NewRdsServiceV1(c.CloudClient, golangsdk.EndpointOpts{
		Region:       c.Region,
		Availability: c.getEndpointType(),
	})
}
