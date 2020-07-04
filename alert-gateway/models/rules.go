package models

type Rules struct {
	Model
	Id          int64  `gorm:"column(id);auto;" json:"id,omitempty"`
	Expr        string `gorm:"column(expr);size(1023)" json:"expr"`
	Op          string `gorm:"column(op);size(31)" json:"op"`
	Value       string `gorm:"column(value);size(1023)" json:"op"`
	For         string `gorm:"column(for);size(1023)" json:"for"`
	Summary     string `gorm:"column(summary);size(1023)" json:"summary"`
	Description string `gorm:"column(description);size(1023)" json:"description"`
	Prom        *Proms `gorm:"foreignkey(Proms)" json:"prom_id"`
	Plan        *Plans `gorm:"foreignkey(Plans)" json:"plan_id"`
	//Labels      []*Labels `orm:"rel(m2m);rel_through(alert-gateway/models.RuleLabels)" json:"omitempty"`
}
