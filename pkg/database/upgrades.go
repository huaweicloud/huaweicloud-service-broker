package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// UpgradesTableName defines
var UpgradesTableName = "upgrades"

// UpgradesTableSQL matches with Upgrades Object
var UpgradesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	upgrade_id int(10) DEFAULT NULL,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, UpgradesTableName)

// Upgrades defines for back database
type Upgrades struct {
	gorm.Model
	UpgradeID int
}
