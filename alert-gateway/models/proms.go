package models

type Proms struct {
	Id   int64  `orm:"auto" json:"id,omitempty"`
	Name string `orm:"size(1023)" json:"name"`
	Url  string `orm:"size(1023)" json:"url"`
}
