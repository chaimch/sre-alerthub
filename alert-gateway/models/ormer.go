package models

import (
	"log"
	"time"

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
}

type Model struct {
	gorm.Model
	CreatedAt time.Time  `gorm:"column:created_at; type:timestamp; default: NOW(); not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:created_at; type:timestamp; default: NOW(); not null" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:created_at; type:timestamp; default: NOW(); null" sql:"index" json:"deleted_at"`
}
