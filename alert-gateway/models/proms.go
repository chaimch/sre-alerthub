package models

const (
	PromsTable = "prom"
)

var PromsReceiver *Proms

type Proms struct {
	Id   int64  `gorm:"auto" json:"id,omitempty"`
	Name string `gorm:"size(1023)" json:"name"`
	Url  string `gorm:"size(1023)" json:"url"`
	Model
}

func (*Proms) TableName() string {
	return PromsTable
}

func (p *Proms) GetAllProms() []Proms {
	proms := []Proms{}
	Ormer.Table(PromsTable).Find(&proms)
	return proms
}
