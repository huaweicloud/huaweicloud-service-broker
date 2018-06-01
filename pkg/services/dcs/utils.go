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

	// TODO confirm from different dcs instance
	if servicetype == models.DCSRedisServiceName ||
		servicetype == models.DCSMemcachedServiceName ||
		servicetype == models.DCSIMDGServiceName {

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
