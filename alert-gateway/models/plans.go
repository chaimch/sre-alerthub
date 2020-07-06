package models

type Plans struct {
	Id          int64  `gorm:"auto" json:"id,omitempty"`
	RuleLabels  string `gorm:"column(rule_labels);size(255)" json:"rule_labels"`
	Description string `gorm:"column(description);size(1023)" json:"description"`
	Model
}
