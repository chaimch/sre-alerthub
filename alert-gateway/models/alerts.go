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
	ID              int64      `gorm:"primary_key;auto" json:"id,omitempty"`
	Rule            Rules      `gorm:"foreignkey:RuleID;save_associations:false:false"`
	RuleID          int64      `json:"rule_id"`
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
					tx := Ormer.Begin()
					tx.Table(AlertsTable).Select("id,count,hostname").Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Set("gorm:query_option", "FOR UPDATE").Find(&recoverInfo)
					if recoverInfo.Id != 0 {
						tx.Table(AlertsTable).Where("id=?", recoverInfo.Id).Updates(map[string]interface{}{
							"status":      elemt.State,
							"summary":     elemt.Annotations.Summary,
							"description": elemt.Annotations.Description,
							"value":       elemt.Value,
							"resolved_at": elemt.ResolvedAt,
						})
						common.Rw.RLock()
						if _, ok := common.Maintain[a.hostname]; !ok {
							var userGroupList []common.UserGroup
							var planId struct {
								PlanId  int64
								Summary string
							}
							Ormer.Table(RulesTable).Select("plan_id,summary").Where("id=?", a.ruleId).Find(&planId)
							if _, ok := Cache[planId.PlanId]; !ok {
								Ormer.Table(ReceiversTable).Select("id,start_time,end_time,start,period,reverse_polish_notation,user,`group`,duty_group,method").Where("plan_id=? AND (method='LANXIN' OR method LIKE 'HOOK %')", planId.PlanId).Find(&userGroupList)
								Cache[planId.PlanId] = userGroupList
							}
							for _, element := range Cache[planId.PlanId] {
								if element.IsValid() && element.IsOnDuty() {
									if recoverInfo.Count >= element.Start {
										sendFlag := false
										if recoverInfo.Count-element.Start >= element.Period {
											sendFlag = true
										} else {
											ruleCountA := [2]int64{a.ruleId, int64(element.Start)}
											if _, ok := common.RuleCount[ruleCountA]; ok {
												if common.RuleCount[ruleCountA] >= int64(recoverInfo.Count-element.Start) {
													if (common.RuleCount[ruleCountA]-int64(recoverInfo.Count)+int64(element.Start))%int64(element.Period) == 0 || common.RuleCount[ruleCountA]-((common.RuleCount[ruleCountA]-int64(recoverInfo.Count)+int64(element.Start))/int64(element.Period))*int64(element.Period) >= int64(element.Period) {
														sendFlag = true
													}
												}
											}
										}
										if sendFlag {
											if element.ReversePolishNotation == "" || common.CalculateReversePolishNotation(elemt.Labels, element.ReversePolishNotation) {
												common.Lock.Lock()
												if _, ok := common.Recover2Send[element.Method]; !ok {
													common.Recover2Send[element.Method] = map[[2]int64]*common.Ready2Send{[2]int64{a.ruleId, element.Id}: &common.Ready2Send{
														RuleId: a.ruleId,
														Start:  element.Id,
														User: SendAlertsFor(&common.ValidUserGroup{
															User:      element.User,
															Group:     element.Group,
															DutyGroup: element.DutyGroup,
														}),
														Alerts: []common.SingleAlert{common.SingleAlert{
															Id:       recoverInfo.Id,
															Count:    recoverInfo.Count,
															Value:    elemt.Value,
															Summary:  elemt.Annotations.Summary,
															Hostname: recoverInfo.Hostname,
														}},
													}}
												} else {
													if _, ok := common.Recover2Send[element.Method][[2]int64{a.ruleId, element.Id}]; !ok {
														common.Recover2Send[element.Method][[2]int64{a.ruleId, element.Id}] = &common.Ready2Send{
															RuleId: a.ruleId,
															Start:  element.Id,
															User: SendAlertsFor(&common.ValidUserGroup{
																User:      element.User,
																Group:     element.Group,
																DutyGroup: element.DutyGroup,
															}),
															Alerts: []common.SingleAlert{common.SingleAlert{
																Id:       recoverInfo.Id,
																Count:    recoverInfo.Count,
																Value:    elemt.Value,
																Summary:  elemt.Annotations.Summary,
																Hostname: recoverInfo.Hostname,
																Labels:   elemt.Labels,
															}},
														}
													} else {
														common.Recover2Send[element.Method][[2]int64{a.ruleId, element.Id}].Alerts = append(common.Recover2Send[element.Method][[2]int64{a.ruleId, element.Id}].Alerts, common.SingleAlert{
															Id:       recoverInfo.Id,
															Count:    recoverInfo.Count,
															Value:    elemt.Value,
															Summary:  elemt.Annotations.Summary,
															Hostname: recoverInfo.Hostname,
														})
													}
												}
												//logs.Panic.Debug("[%s] %v",common.Recover2Send["LANXIN"])
												common.Lock.Unlock()
											}
										}
									}
								}
							}
						}
						common.Rw.RUnlock()
						tx.Commit()
					}
					tx.Commit()
				} else {
					//send the recover message
					Ormer.Table(AlertsTable).Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Update(map[string]interface{}{
						"summary":     elemt.Annotations.Summary,
						"description": elemt.Annotations.Description,
						"value":       elemt.Value,
					})
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
