package queue

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dms instance
func BuildBindingCredential(
	endpointname string,
	endpointport string,
	region string,
	projectid string,
	url string,
	ak string,
	sk string,
	queueid string,
	groupid string,
	servicetype string) (BindingCredential, error) {

	// Set ProtocolType
	var protocoltype string
	if servicetype == models.DMSStandardServiceName {
		protocoltype = ProtocolTypeHTTPS
	} else if servicetype == models.DMSActiveMQServiceName {
		protocoltype = ProtocolTypeTCP
	} else if servicetype == models.DMSKafkaServiceName {
		protocoltype = ProtocolTypeTCP
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Set url
	if url != "" {
		// Remove https://
		parts := strings.Split(url, "//")
		if len(parts) == 2 {
			// Remove last /
			url = strings.Replace(parts[1], "/", "", -1)
		}

		// Replace Endpoint Name
		if endpointname != "" {
			url = strings.Replace(url, "dms.", endpointname+".", -1)
		}
		// Add Endpoint Port
		if endpointport != "" {
			url = fmt.Sprintf("%s:%s", url, endpointport)
		}
	}
	// Init BindingCredential
	bc := BindingCredential{
		Region:       region,
		ProtocolType: protocoltype,
		ProjectID:    projectid,
		URL:          url,
		AK:           ak,
		SK:           sk,
		QueueID:      queueid,
		GroupID:      groupid,
		Type:         servicetype,
	}
	return bc, nil
}
