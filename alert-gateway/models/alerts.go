package models

import (
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
)

const (
	AlertsTable = "alert"
)

type Alerts struct {
	gorm.Model
	Id              int        `gorm:"column(id);auto" json:"id,omitempty"`
	Rule            *Rules     `gorm:"foreignkey(fk)" json:"rule_id"`
	Labels          string     `gorm:"column(labels);size(4095)" json:"labels"`
	Value           float64    `gorm:"column(value)" json:"value"`
	Count           int        `json:"count"`
	Status          int8       `gorm:"index" json:"status"`
	Summary         string     `gorm:"column(summary);size(1023)" json:"summary"`
	Description     string     `gorm:"column(description);size(1023)" json:"description"`
	Hostname        string     `gorm:"column(hostname);size(255)" json:"hostname"`
	ConfirmedBy     string     `gorm:"column(confirmed_by);size(1023)" json:"confirmed_by"`
	FiredAt         *time.Time `gorm:"type(datetime)" json:"fired_at"`
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
	log.Println("todayZero", todayZero)

	for _, elemt := range *alert {
		var queryres []struct {
			Id     int64
			Status uint8
		}
		log.Println("queryres", queryres)

		a := &alertForQuery{Alert: &elemt}
		log.Println("a", a)
		a.setFields()
		log.Println("a", a)

		Ormer.Table(AlertsTable).Select("id,status").Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Find(&queryres)
		log.Println("queryres", queryres)
		if len(queryres) > 0 {
			if queryres[0].Status != 0 {
				log.Println("Get something")
			}
		} else {
			Ormer.Table(AlertsTable).Create(&Alerts{
				// TODO: reset the "Id" to 0,which is very important:after a record is inserted,the value of "Id" will not be 0,but the auto primary key of the record
				Id:              0,
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
			})
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
		labels = append(labels, i+"\a"+a.Labels[i])
	}

	// set label
	a.label = strings.Join(labels, "\v")

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
