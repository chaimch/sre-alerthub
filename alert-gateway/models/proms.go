package models

type Proms struct {
	Id   int64  `gorm:"auto" json:"id,omitempty"`
	Name string `gorm:"size(1023)" json:"name"`
	Url  string `gorm:"size(1023)" json:"url"`
	Model
}
