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

	// TODO confirm different uri from different rds instance
	if servicetype == models.RDSPostgresqlServiceName {
		// Postgresql
		uri = fmt.Sprintf("mysql://%s:%s@%s:%d/%s?reconnect=true", username, password, host, port, name)
	} else if servicetype == models.RDSMysqlServiceName {
		// Mysql
		uri = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?reconnect=true", username, password, host, port, name)
	} else if servicetype == models.RDSSqlserverServiceName {
		// Sqlserver
		uri = fmt.Sprintf("sqlserver://%s:%s@%s:%d/%s?reconnect=true", username, password, host, port, name)
	} else if servicetype == models.RDSHwsqlServiceName {
		// Hwsql
		uri = fmt.Sprintf("hwsql://%s:%s@%s:%d/%s?reconnect=true", username, password, host, port, name)
	} else {
		return BindingCredential{}, fmt.Errorf("unknown service type: %s", servicetype)
	}

	// Init provisionOpts
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
