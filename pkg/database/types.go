package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// UpgradeTableName defines
var UpgradeTableName = "upgrades"

// UpgradeTableSQL matches with Upgrade Object
var UpgradeTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	upgrade_id int(10) DEFAULT NULL,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, UpgradeTableName)

// Upgrade defines for back database
type Upgrade struct {
	gorm.Model
	UpgradeID int
}
