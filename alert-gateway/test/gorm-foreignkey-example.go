package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Users struct {
	// gorm.Model
	ID   int64 `gorm:"column:id;auto;" json:"id,omitempty"`
	Name string
}

// `Profile` 属于 `User`， 外键是`UserID`
type Profile struct {
	ID     int64 `gorm:"column:id;auto;" json:"id,omitempty"`
	UserID int
	Users  Users `gorm:"foreignkey:UserID" json:"user_id"`
	// Users Users `gorm:"association_foreignkey:ID" json:"user_id"`
	Name        string
	ConfirmedBy string
}

func main() {
	db, err := gorm.Open("mysql", "root:root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("failed to connect database")
	}
	defer db.Close()

	db.LogMode(true)

	db.DropTableIfExists(
		&Profile{},
		&Users{},
	)

	// db.AutoMigrate(
	// 	&Users{},
	// 	&Profile{},
	// )

	db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(
		&Profile{},
		&Users{},
	)

	// db.Model(&Profile{}).AddForeignKey("user_refer", "users(id)", "RESTRICT", "RESTRICT")

	profile := &Profile{
		Name: "profile-name",
		// UserID: 10,
		Users: Users{
			ID:   10,
			Name: "user-name",
		},
	}
	db.Create(profile)
}
