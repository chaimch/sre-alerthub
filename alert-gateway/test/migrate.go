package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

func main() {
	// models.Ormer.DropTableIfExists(
	// 	&Profile{},
	// 	&Users{},
	// )

	models.Ormer.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(
		&models.Rules{},
		&models.Alerts{},
	)

	// models.Ormer.AutoMigrate(
	// 	&models.Alerts{},
	// 	&models.Rules{},
	// )

}
