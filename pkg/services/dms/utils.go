package dms

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dms instance
func BuildBindingCredential(
	region string,
	projectid string,
	url string,
	ak string,
	sk string,
	servicetype string) (BindingCredential, error) {

	// TODO confirm url from different dms instance
	if servicetype == models.DMSStandardServiceName ||
		servicetype == models.DMSActiveMQServiceName ||
		servicetype == models.DMSKafkaServiceName ||
		servicetype == models.DMSRabbitMQServiceName {

	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		Region:    region,
		ProjectID: projectid,
		URL:       url,
		AK:        ak,
		SK:        sk,
		Type:      servicetype,
	}
	return bc, nil
}
