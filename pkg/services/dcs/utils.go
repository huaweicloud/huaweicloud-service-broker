package dcs

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different dcs instance
func BuildBindingCredential(
	ip string,
	port int,
	username string,
	password string,
	name string,
	servicetype string) (BindingCredential, error) {

	if servicetype == models.DCSRedisServiceName {
		username = ""
	} else if servicetype == models.DCSMemcachedServiceName {

	} else if servicetype == models.DCSIMDGServiceName {
		port = 0
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		IP:       ip,
		Port:     port,
		UserName: username,
		Password: password,
		Name:     name,
		Type:     servicetype,
	}
	return bc, nil
}
