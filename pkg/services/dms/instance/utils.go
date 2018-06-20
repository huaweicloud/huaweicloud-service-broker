package instance

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dms instance
func BuildBindingCredential(
	host string,
	port int,
	username string,
	password string,
	servicetype string) (BindingCredential, error) {

	// Group uri
	uri := ""
	if servicetype == models.DMSRabbitMQServiceName {
		uri = fmt.Sprintf("%s:%d", host, port)
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
		URI:      uri,
		Type:     servicetype,
	}
	return bc, nil
}
