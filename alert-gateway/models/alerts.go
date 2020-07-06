package models

import (
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
)

const (
	AlertsTable = "alert"
)

type Alerts struct {
	Model
	ID              int64 `gorm:"primary_key;auto" json:"id,omitempty"`
	Rule            Rules `gorm:"foreignkey:RuleID;save_associations:false:false" json:"rule_id"`
	RuleID          int64
	Labels          string     `gorm:"size:4095" json:"labels"`
	Value           float64    `json:"value"`
	Count           int        `json:"count"`
	Status          int8       `gorm:"index" json:"status"`
	Summary         string     `gorm:"size:1023" json:"summary"`
	Description     string     `gorm:"size:1023" json:"description"`
	Hostname        string     `gorm:"size:255" json:"hostname"`
	ConfirmedBy     string     `gorm:"size:1023" json:"confirmed_by"`
	FiredAt         *time.Time `gorm:"type:datetime" json:"fired_at"`
	ConfirmedAt     *time.Time `gorm:"null" json:"confirmed_at"`
	ConfirmedBefore *time.Time `gorm:"null" json:"confirmed_before"`
	ResolvedAt      *time.Time `gorm:"null" json:"resolved_at"`
}

type alertForQuery struct {
	*common.Alert
	label    string
	hostname string
	ruleId   int64
	firedAt  time.Time
}

func (Alerts) TableName() string {
	return AlertsTable
}

func (u *Alerts) AlertsHandler(alert *common.Alerts) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 16384)
			buf = buf[:runtime.Stack(buf, false)]
			log.Panic("Panic in AlertsHandler:", e, buf)
		}
	}()

	Cache := map[int64][]common.UserGroup{}
	log.Println("Cache", Cache)
	now := time.Now().Format("13:14:00")
	log.Println("now", now)
	todayZero, _ := time.ParseInLocation("0001-01-01", "0001-01-01 13:14:00", time.Local)

	for _, elemt := range *alert {
		var queryres []struct {
			Id     int64
			Status uint8
		}

		a := &alertForQuery{Alert: &elemt}
		a.setFields()

		log.Println("ruleId", a.ruleId, "labels", a.label)
		Ormer.Table(AlertsTable).Select("id,status").Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Find(&queryres)
		log.Println("len(queryres)", len(queryres))
		if len(queryres) > 0 {
			if queryres[0].Status != 0 {
				const AlertStatusOff = 0
				if elemt.State == AlertStatusOff {
					//rlist = append(rlist, queryres[0].Id)
					recoverInfo := struct {
						Id       int64
						Count    int
						Hostname string
					}{}
					log.Println("elemt.State", elemt.State)
					log.Println("recoverInfo", recoverInfo)
				}
			}
		} else {
			alert := &Alerts{
				// TODO: reset the "Id" to 0,which is very important:after a record is inserted,the value of "Id" will not be 0,but the auto primary key of the record
				ID: 0,
				// Rule:            Rules{ID: a.ruleId},
				RuleID:          a.ruleId,
				Labels:          a.label,
				FiredAt:         &a.firedAt,
				Description:     elemt.Annotations.Description,
				Summary:         elemt.Annotations.Summary,
				Count:           -1,
				Value:           elemt.Value,
				Status:          int8(elemt.State),
				Hostname:        a.hostname,
				ConfirmedAt:     &todayZero,
				ConfirmedBefore: &todayZero,
				ResolvedAt:      &todayZero,
			}
			Ormer.Table(AlertsTable).Create(alert)
		}
	}
}

/*
 set value for fields in alertForQuery
*/
func (a *alertForQuery) setFields() {
	var orderKey []string
	var labels []string

	// set ruleId
	a.ruleId, _ = strconv.ParseInt(a.Annotations.RuleId, 10, 64)
	for key := range a.Labels {
		orderKey = append(orderKey, key)
	}

	sort.Strings(orderKey)
	for _, i := range orderKey {
		// labels = append(labels, i+"\a"+a.Labels[i])
		labels = append(labels, i+":"+a.Labels[i])
	}

	// set label
	// a.label = strings.Join(labels, "\v")
	a.label = strings.Join(labels, ";")

	// set firedAt
	a.firedAt = a.FiredAt.Truncate(time.Second)

	// set hostname
	a.setHostname()
}

/*
 set hostname by instance label on data
*/
func (a *alertForQuery) setHostname() {
	h := ""
	if _, ok := a.Labels["instance"]; ok {
		h = a.Labels["instance"]
		boundary := strings.LastIndex(h, ":")
		if boundary != -1 {
			h = h[:boundary]
		}
	}
	a.hostname = h
}
