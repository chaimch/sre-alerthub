package models

const (
	RulesTable = "rule"
)

type Rules struct {
	Model
	ID          int64  `gorm:"primary_key;auto;" json:"id,omitempty"`
	Expr        string `gorm:"size:1023" json:"expr"`
	Op          string `gorm:"size:31" json:"op"`
	Value       string `gorm:"size:1023" json:"value"`
	For         string `gorm:"size:1023" json:"for"`
	Summary     string `gorm:"size:1023" json:"summary"`
	Description string `gorm:"size:1023" json:"description"`
	Prom        *Proms `gorm:"foreignkey:PromID" json:"prom_id"`
	PromID      int64
	Plan        Plans `gorm:"foreignkey:PlanID" json:"plan_id"`
	PlanID      int64
	//Labels      []*Labels `orm:"rel(m2m);rel_through(alert-gateway/models.RuleLabels)" json:"omitempty"`
}

func (Rules) TableName() string {
	return RulesTable
}
