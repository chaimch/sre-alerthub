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

		// 不存在指定 rule_id & labels & fired_ad 的 alert 则创建 alert 记录
		if len(queryres) <= 0 {
			Ormer.Table(AlertsTable).Create(&Alerts{
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
			})
			continue
		}

		// 第一条 DB 中报警如果已恢复, 则直接跳过, status 表示告警已恢复
		const AlertStatusOff = 0
		if queryres[0].Status == AlertStatusOff {
			continue
		}

		// 如果从 rule-engine 发送过来的不是告警恢复消息, 则更新对应的 summary & description & value
		if elemt.State != AlertStatusOff {
			Ormer.Table(AlertsTable).Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Update(map[string]interface{}{
				"summary":     elemt.Annotations.Summary,
				"description": elemt.Annotations.Description,
				"value":       elemt.Value,
			})
			continue
		}

		//  如果rule-engine 发送过来的是告警恢复消息
		recoverInfo := struct {
			Id       int64
			Count    int
			Hostname string
		}{}

		// 重新查询 alert 信息并加锁
		tx := Ormer.Begin()
		err := tx.Table(AlertsTable).Select("id,count,hostname").Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Set("gorm:query_option", "FOR UPDATE").Find(&recoverInfo)

		//  加锁超时时不能查询到对应的 recoverInfo, 则将 rule-engine 过来的告警恢复消息的 summary & description & value 更新到数据库
		if err != nil {
			Ormer.Table(AlertsTable).Where("rule_id =? AND labels=? AND fired_at=?", a.ruleId, a.label, a.firedAt).Update(map[string]interface{}{
				"summary":     elemt.Annotations.Summary,
				"description": elemt.Annotations.Description,
				"value":       elemt.Value,
			})
			tx.Commit()
			continue
		}

		// 如果未查询到任何待恢复的报警则 commit 并结束
		if recoverInfo.Id == 0 {
			tx.Commit()
			continue
		}

		// 从 alert 中获取到对应的待恢复报警时, 需更新 status & resolved_at
		err = tx.Table(AlertsTable).Where("id=?", recoverInfo.Id).Update(map[string]interface{}{
			"status":      elemt.State,
			"summary":     elemt.Annotations.Summary,
			"description": elemt.Annotations.Description,
			"value":       elemt.Value,
			"resolved_at": elemt.ResolvedAt,
		})
		// 更新恢复警报出错则回滚
		if err != nil {
			tx.Rollback()
			continue
		}

		common.Rw.RLock()

		//  如果机器在维护中, 则跳过发送告警恢复信息
		if _, ok := common.Maintain[a.hostname]; ok {
			common.Rw.RUnlock()
			tx.Commit()
		}

		// 根据 ruleId 找到对应的 planId, 根据 planId 获取对应的 plan_receiver
		var planId struct {
			PlanId  int64
			Summary string
		}
		var userGroupList []common.UserGroup
		Ormer.Table(RulesTable).Select("plan_id,summary").Where("id=?", a.ruleId).Find(&planId)

		// 当前 ruleId 对应的 planId 不在 Cache 中则添加
		if _, ok := Cache[planId.PlanId]; !ok {
			Ormer.Table(ReceiversTable).Select("id,start_time,end_time,start,period,reverse_polish_notation,user,`group`,duty_group,method").Where("plan_id=? AND (method='LANXIN' OR method LIKE 'HOOK %')", planId.PlanId).Find(&userGroupList)
			Cache[planId.PlanId] = userGroupList
		}

		for _, itemUserGroupList := range Cache[planId.PlanId] {
			// 如果当前组不可用或不在值日时候, 则跳过
			if !(itemUserGroupList.IsValid() && itemUserGroupList.IsOnDuty()) {
				continue
			}

			//  如果告警持续时间小于用户组的延迟报警时间, 则跳过
			if recoverInfo.Count < itemUserGroupList.Start {
				continue
			}

			// 如果(告警时间-告警组的延迟报警时间)仍小于告警组设置的报警周期时间, 则跳过
			if recoverInfo.Count-itemUserGroupList.Start < itemUserGroupList.Period {
				continue
			}

			itemRuleCount := [2]int64{a.ruleId, int64(itemUserGroupList.Start)}

			// 如果当前的 ruleId & userGroup 对应的 ruleCount 不存在则跳过
			if _, ok := common.RuleCount[itemRuleCount]; !ok {
				continue
			}

			// 如果当前的 ruleId & userGroup 对应的 ruleCount 小于 (告警时间 - 告警组延迟时间) 则跳过
			if common.RuleCount[itemRuleCount] < int64(recoverInfo.Count-itemUserGroupList.Start) {
				continue
			}

			// 如果 ruleCount 未到 userGroup 的告警周期, 则跳过
			if (common.RuleCount[itemRuleCount]-int64(recoverInfo.Count)+int64(itemUserGroupList.Start))%int64(itemUserGroupList.Period) != 0 {
				continue
			}

			// 如果 ruleCount 持续时间未到整数个 userGroup 的告警周期, 则跳过
			if common.RuleCount[itemRuleCount]-((common.RuleCount[itemRuleCount]-int64(recoverInfo.Count)+int64(itemUserGroupList.Start))/int64(itemUserGroupList.Period))*int64(itemUserGroupList.Period) < int64(itemUserGroupList.Period) {
				continue
			}

			// 如果 userGroup 的逆波兰已设置说明标签表达式校核合法则跳过
			if itemUserGroupList.ReversePolishNotation != "" {
				continue
			}

			// 如果标签匹配规则不符合逆波兰表达式, 说明表达式不合法, 则跳过
			if !common.CalculateReversePolishNotation(elemt.Labels, itemUserGroupList.ReversePolishNotation) {
				continue
			}

			common.Lock.Lock()
			// 如果该 userGroup 的告警恢复方法不存在, 则初始化该 userGroup 的告警恢复信息
			if _, ok := common.Recover2Send[itemUserGroupList.Method]; !ok {
				common.Recover2Send[itemUserGroupList.Method] = map[[2]int64]*common.Ready2Send{
					itemRuleCount: &common.Ready2Send{
						RuleId: a.ruleId,
						Start:  itemUserGroupList.Id,
						User: SendAlertsFor(&common.ValidUserGroup{
							User:      itemUserGroupList.User,
							Group:     itemUserGroupList.Group,
							DutyGroup: itemUserGroupList.DutyGroup,
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
				// 如果该 userGroup 对应的告警恢复方法的 ruleCount 不存在, 则初始化告警恢复信息
				if _, ok := common.Recover2Send[itemUserGroupList.Method][itemRuleCount]; !ok {
					common.Recover2Send[itemUserGroupList.Method][itemRuleCount] = &common.Ready2Send{
						RuleId: a.ruleId,
						Start:  itemUserGroupList.Id,
						User: SendAlertsFor(&common.ValidUserGroup{
							User:      itemUserGroupList.User,
							Group:     itemUserGroupList.Group,
							DutyGroup: itemUserGroupList.DutyGroup,
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
					// 如果该 userGroup 对应的告警恢复方法的 ruleCount 存在, 则添加当前告警恢复信息
					common.Recover2Send[itemUserGroupList.Method][itemRuleCount].Alerts = append(
						common.Recover2Send[itemUserGroupList.Method][itemRuleCount].Alerts,
						common.SingleAlert{
							Id:       recoverInfo.Id,
							Count:    recoverInfo.Count,
							Value:    elemt.Value,
							Summary:  elemt.Annotations.Summary,
							Hostname: recoverInfo.Hostname,
						})
				}
			}
			common.Lock.Unlock()
		}
		common.Rw.RUnlock()
		tx.Commit()
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
