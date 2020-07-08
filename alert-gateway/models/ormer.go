package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Ormer *gorm.DB

func init() {
	var err error
	// TODO: 待加入配置解析
	Ormer, err = gorm.Open("mysql", "root:root@/alerthub?charset=utf8&parseTime=true")
	if err != nil {
		log.Panic(err)
	}

	// TODO: 上线前删除
	Ormer.LogMode(true)
	Ormer.AutoMigrate(
		&Alerts{},
		&Rules{},
		&Receivers{},
		&Groups{},
		&Plans{},
		&Proms{},
		&Maintains{},
		&Hosts{},
	)
}

// refer to	gorm.Model
type Model struct {
	// TODO: CreatedAt & UpdatedAt & DeletedAt 会增加多余的 select
	// CreatedAt time.Time  `gorm:"type:timestamp; default: NOW(); not null" json:"created_at"`
	// UpdatedAt time.Time  `gorm:"type:timestamp; default: NOW(); not null" json:"updated_at"`
	// DeletedAt *time.Time `gorm:"type:timestamp; null" sql:"index" json:"deleted_at"`
}
