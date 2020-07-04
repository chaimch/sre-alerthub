package models

type Rules struct {
	Id          int64  `orm:"column(id);auto" json:"id,omitempty"`
	Expr        string `orm:"column(expr);size(1023)" json:"expr"`
	Op          string `orm:"column(op);size(31)" json:"op"`
	Value       string `orm:"column(value);size(1023)" json:"op"`
	For         string `orm:"column(for);size(1023)" json:"for"`
	Summary     string `orm:"column(summary);size(1023)" json:"summary"`
	Description string `orm:"column(description);size(1023)" json:"description"`
	Prom        *Proms `orm:"rel(fk)" json:"prom_id"`
	Plan        *Plans `orm:"rel(fk)" json:"plan_id"`
	//Labels      []*Labels `orm:"rel(m2m);rel_through(alert-gateway/models.RuleLabels)" json:"omitempty"`
}
