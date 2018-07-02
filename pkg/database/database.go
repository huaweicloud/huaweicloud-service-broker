package database

import (
	"fmt"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"github.com/jinzhu/gorm"

	// import mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// BackDBConnection is used to operate on database
var BackDBConnection *gorm.DB

// once is used to init BackDBConnection
var once sync.Once

// New Database Connection
func New(logger lager.Logger, config config.Config) error {
	once.Do(func() {
		// Connect Data
		backDBConnection, err := ConnectBackDatabase(logger, config)
		if err != nil {
			logger.Error("Error connecting to back database", err)
			panic(err)
		}
		err = UpgradeBackDatabase(logger, backDBConnection)
		if err != nil {
			logger.Error("Error upgrading back database", err)
			panic(err)
		}
		BackDBConnection = backDBConnection
	})
	return nil
}

// ConnectBackDatabase connect back database
func ConnectBackDatabase(logger lager.Logger, config config.Config) (*gorm.DB, error) {
	// connect to back database
	connStr := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.BackDatabase.DatabaseUsername,
		config.BackDatabase.DatabasePassword,
		config.BackDatabase.DatabaseHost,
		config.BackDatabase.DatabasePort,
		config.BackDatabase.DatabaseName)

	return gorm.Open(config.BackDatabase.DatabaseType, connStr)
}

// UpgradeBackDatabase upgrades tables and datas
func UpgradeBackDatabase(logger lager.Logger, backdatabase *gorm.DB) error {

	// upgrades defines
	upgrades := make([]func() error, 1)
	upgrades[0] = func() error {
		// table upgrades
		if err := backdatabase.Exec(UpgradesTableSQL).Error; err != nil {
			return err
		}
		// table instance_details
		if err := backdatabase.Exec(InstanceDetailsTableSQL).Error; err != nil {
			return err
		}
		// table bind_details
		if err := backdatabase.Exec(BindDetailsTableSQL).Error; err != nil {
			return err
		}
		// table operation_details
		if err := backdatabase.Exec(OperationDetailsTableSQL).Error; err != nil {
			return err
		}
		return nil
	}

	// Get UpgradeID
	var lastUpgrade = -1
	if backdatabase.HasTable(UpgradesTableName) {
		var ups []Upgrades
		err := backdatabase.Order("upgrade_id desc").Find(&ups).Error
		if err != nil {
			logger.Error("Error getting upgrades", err)
			return err
		}

		// Get value
		if len(ups) > 0 {
			lastUpgrade = ups[0].UpgradeID
		}
	}

	// Exec upgrades one by one
	for index := lastUpgrade + 1; index < len(upgrades); index++ {
		// Begin transaction
		tx := backdatabase.Begin()

		// Exec upgrades
		err := upgrades[index]()
		if err != nil {
			tx.Rollback()
			return err
		}

		// Save last upgrade
		err = backdatabase.Save(&Upgrades{
			UpgradeID: index,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// End transaction
		tx.Commit()
	}

	return nil
}
