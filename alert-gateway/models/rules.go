package models

const (
	RulesTable = "rule"
)

var RulesReceiver *Rules

type Rules struct {
	Model
	ID          int64  `gorm:"primary_key;auto;" json:"id,omitempty"`
	Expr        string `gorm:"size:1023" json:"expr"`
	Op          string `gorm:"size:31" json:"op"`
	Value       string `gorm:"size:1023" json:"value"`
	For         string `gorm:"size:1023" json:"for"`
	Summary     string `gorm:"size:1023" json:"summary"`
	Description string `gorm:"size:1023" json:"description"`
	Prom        Proms  `gorm:"foreignkey:PromID" json:"-"`
	PromID      int64  `json:"prom_id"`
	Plan        Plans  `gorm:"foreignkey:PlanID" json:"-"`
	PlanID      int64  `json:"plan_id"`
	//Labels      []*Labels `orm:"rel(m2m);rel_through(alert-gateway/models.RuleLabels)" json:"omitempty"`
}

func (Rules) TableName() string {
	return RulesTable
}

func (*Rules) Get(prom string, id string) []Rules {
	rules := []Rules{}
	if prom != "" {
		Ormer.Table(RulesTable).Where("prom_id = ?", prom).Find(&rules)
	} else if id != "" {
		Ormer.Table(RulesTable).Where("id = ?", id).Find(&rules)
	} else {
		Ormer.Table(RulesTable).Find(&rules)
	}
	return rules
}
