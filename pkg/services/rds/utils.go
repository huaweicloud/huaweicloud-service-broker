package rds

import (
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/jinzhu/gorm"

	// import mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// BuildBindingCredential from different rds instance
func BuildBindingCredential(
	host string,
	port int,
	name string,
	username string,
	password string,
	servicetype string,
	logger lager.Logger) (BindingCredential, error) {

	var uri string

	if servicetype == models.RDSPostgresqlServiceName {
		// Postgresql
		uri = fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	} else if servicetype == models.RDSMysqlServiceName {
		// Mysql
		name = "broker"
		err := CreateDatabase(username, password, host, port, "mysql", name)
		if err != nil {
			logger.Debug(fmt.Sprintf("create database failed: %v", err))
		}
		uri = fmt.Sprintf("%s:%s@%s:%d/%s", username, password, host, port, name)
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

// CreateDatabase create database by name
func CreateDatabase(DatabaseUsername string,
	DatabasePassword string,
	DatabaseHost string,
	DatabasePort int,
	DatabaseType string,
	DatabaseName string) error {
	// connect to back database
	connStr := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/mysql?charset=utf8&parseTime=True&loc=Local",
		DatabaseUsername,
		DatabasePassword,
		DatabaseHost,
		DatabasePort)

	dbConn, err := gorm.Open(DatabaseType, connStr)
	if err != nil {
		return err
	}

	CreateDataBaseSQL := fmt.Sprintf(`create database %s;`, DatabaseName)
	if err := dbConn.Exec(CreateDataBaseSQL).Error; err != nil {
		return err
	}

	return nil
}
