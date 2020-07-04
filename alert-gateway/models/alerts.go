package models

import (
	"time"

	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
)

type Alerts struct {
	Id              int64      `orm:"column(id);auto" json:"id,omitempty"`
	Rule            *Rules     `orm:"rel(fk)" json:"rule_id"`
	Labels          string     `orm:"column(labels);size(4095)" json:"labels"`
	Value           float64    `orm:"column(value)" json:"value"`
	Count           int        `json:"count"`
	Status          int8       `orm:"index" json:"status"`
	Summary         string     `orm:"column(summary);size(1023)" json:"summary"`
	Description     string     `orm:"column(description);size(1023)" json:"description"`
	Hostname        string     `orm:"column(hostname);size(255)" json:"hostname"`
	ConfirmedBy     string     `orm:"column(confirmed_by);size(1023)" json:"confirmed_by"`
	FiredAt         *time.Time `orm:"type(datetime)" json:"fired_at"`
	ConfirmedAt     *time.Time `orm:"null" json:"confirmed_at"`
	ConfirmedBefore *time.Time `orm:"null" json:"confirmed_before"`
	ResolvedAt      *time.Time `orm:"null" json:"resolved_at"`
}

func (u *Alerts) AlertsHandler(alert *common.Alerts) {

}
