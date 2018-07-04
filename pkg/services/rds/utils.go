package rds

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
)

// BuildBindingCredential from different rds instance
func BuildBindingCredential(
	host string,
	port int,
	name string,
	username string,
	password string,
	servicetype string) (BindingCredential, error) {

	var uri string

	if servicetype == models.RDSPostgresqlServiceName {
		// Postgresql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSMysqlServiceName {
		// Mysql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSSqlserverServiceName {
		// Sqlserver
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSHwsqlServiceName {
		// Hwsql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init BindingCredential
	bc := BindingCredential{
		Host:     host,
		Port:     port,
		Name:     name,
		UserName: username,
		Password: password,
		URI:      uri,
		Type:     servicetype,
	}
	return bc, nil
}
