package timer

import (
	"encoding/json"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

type Record struct {
	Id              int64
	RuleId          int64
	Value           float64
	Count           int
	Summary         string
	Description     string
	Hostname        string
	ConfirmedBefore *time.Time
	FiredAt         *time.Time
	Labels          string
}

func UpdateMaintainlist() {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 16384)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println("Panic in UpdateMaintainlist:", e, buf)
		}
	}()

	delta, _ := time.ParseDuration("30s")
	datetime := time.Now().Add(delta)
	now := datetime.Format("15:04")
	maintainIds := []struct {
		Id int64
	}{}
	// Format 必须为 2006-01-02 15:04:05, 6-1-2-3-4-5, go 诞生之日
	models.Ormer.Table(models.MaintainsTable).Select("id").Where("valid>=? AND day_start<=? AND day_end>=? AND (flag=true AND (time_start<=? OR time_end>=?) OR flag=false AND time_start<=? AND time_end>=?) AND month&"+strconv.Itoa(int(math.Pow(2, float64(time.Now().Month()))))+">0", datetime.Format("2006-01-02 15:04:05"), datetime.Day(), datetime.Day(), now, now, now, now).Find(&maintainIds)

	m := map[string]bool{}
	for _, mid := range maintainIds {
		hosts := []struct {
			Hostname string
		}{}
		models.Ormer.Table(models.HostsTable).Select("hostname").Where("mid=?", mid.Id).Find(&hosts)
		for _, name := range hosts {
			m[name.Hostname] = true
		}
	}

	common.Rw.Lock()
	common.Maintain = m
	common.Rw.Unlock()

	log.Println("Maintain: ", common.Maintain)
}

func Filter(alerts map[int64][]Record, maxCount map[int64]int) map[string][]common.Ready2Send {
	log.Println("Filter")
	SendClass := map[string][]common.Ready2Send{
		"SMS":    []common.Ready2Send{},
		"LANXIN": []common.Ready2Send{},
		"CALL":   []common.Ready2Send{},
	}
	Cache := map[int64][]common.UserGroup{}
	NewRuleCount := map[[2]int64]int64{}
	for key := range alerts {
		var usergroupList []common.UserGroup
		var planId struct {
			PlanId  int64
			Summary string
		}
		AlertsMap := map[int][]common.SingleAlert{}
		models.Ormer.Table(models.RulesTable).Select("plan_id,summary").Where("id=?", key).Find(&planId)
		if _, ok := Cache[planId.PlanId]; !ok {
			models.Ormer.Table(models.ReceiversTable).Select("id,start_time,end_time,start,period,reverse_polish_notation,user,`group`,duty_group,method").Where("plan_id=?", planId.PlanId).Find(&usergroupList)
			Cache[planId.PlanId] = usergroupList
		}

		for _, element := range Cache[planId.PlanId] {
			if !(element.IsValid() && element.IsOnDuty()) {
				continue
			}

			if maxCount[key] < element.Start {
				continue
			}

			itemRuleCount := [2]int64{key, int64(element.Start)}
			NewRuleCount[itemRuleCount] = -1
			if _, ok := common.RuleCount[itemRuleCount]; ok {
				itemRuleCount = [2]int64{key, int64(element.Start)}
				NewRuleCount[itemRuleCount] = common.RuleCount[itemRuleCount]
			}
			NewRuleCount[itemRuleCount] += 1

			if NewRuleCount[itemRuleCount]%int64(element.Period) != 0 {
				continue
			}

			if _, ok := AlertsMap[element.Start]; !ok {
				AlertsMap[element.Start] = []common.SingleAlert{}
			} else {
				if len(AlertsMap[element.Start]) > 0 {
					if element.ReversePolishNotation == "" {
						SendClass[element.Method] = append(SendClass[element.Method], common.Ready2Send{
							RuleId: key,
							Start:  element.Id,
							User: models.SendAlertsFor(&common.ValidUserGroup{
								User:      element.User,
								Group:     element.Group,
								DutyGroup: element.DutyGroup,
							}),
							Alerts: AlertsMap[element.Start],
						})
					} else {
						filteredAlerts := []common.SingleAlert{}
						for _, alert := range AlertsMap[element.Start] {
							if common.CalculateReversePolishNotation(alert.Labels, element.ReversePolishNotation) {
								filteredAlerts = append(filteredAlerts, alert)
							}
						}
						if len(filteredAlerts) > 0 {
							SendClass[element.Method] = append(SendClass[element.Method], common.Ready2Send{
								RuleId: key,
								Start:  element.Id,
								User: models.SendAlertsFor(&common.ValidUserGroup{
									User:      element.User,
									Group:     element.Group,
									DutyGroup: element.DutyGroup,
								}),
								Alerts: filteredAlerts,
							})
						}
					}
				}
				continue
			}

			for _, alert := range alerts[key] {
				if alert.Count < element.Start {
					continue
				}

				if _, ok := common.Maintain[alert.Hostname]; ok {
					continue
				}

				labelMap := map[string]string{}
				if alert.Labels != "" {
					for _, j := range strings.Split(alert.Labels, "\v") {
						kv := strings.Split(j, "\a")
						labelMap[kv[0]] = kv[1]
					}
				}
				AlertsMap[element.Start] = append(AlertsMap[element.Start], common.SingleAlert{
					Id:       alert.Id,
					Count:    alert.Count,
					Value:    alert.Value,
					Summary:  alert.Summary,
					Hostname: alert.Hostname,
					Labels:   labelMap,
				})
			}

			if len(AlertsMap[element.Start]) <= 0 {
				continue
			}

			if element.ReversePolishNotation == "" {
				SendClass[element.Method] = append(SendClass[element.Method], common.Ready2Send{
					RuleId: key,
					Start:  element.Id,
					User: models.SendAlertsFor(&common.ValidUserGroup{
						User:      element.User,
						Group:     element.Group,
						DutyGroup: element.DutyGroup,
					}),
					Alerts: AlertsMap[element.Start],
				})
			} else {
				filteredAlerts := []common.SingleAlert{}
				for _, alert := range AlertsMap[element.Start] {
					if common.CalculateReversePolishNotation(alert.Labels, element.ReversePolishNotation) {
						filteredAlerts = append(filteredAlerts, alert)
					}
				}
				if len(filteredAlerts) > 0 {
					SendClass[element.Method] = append(SendClass[element.Method], common.Ready2Send{
						RuleId: key,
						Start:  element.Id,
						User: models.SendAlertsFor(&common.ValidUserGroup{
							User:      element.User,
							Group:     element.Group,
							DutyGroup: element.DutyGroup,
						}),
						Alerts: filteredAlerts,
					})
				}
			}

		}
	}
	common.RuleCount = NewRuleCount
	return SendClass
}

func Sender(SendClass map[string][]common.Ready2Send, now string) {
	for k, v := range SendClass {
		switch k {
		case "SMS":
			log.Println("SMS send ...", now, k, v)
		case "LANXIN":
			log.Println("LANXIN send ...", now, k, v)
		case "CALL":
			log.Println("CALL send ...", now, k, v)
		default:
			log.Println("Hook send ...", now, k, v)
			go Send2Hook(v, now, "alert", k[5:])
		}
	}
}

func RecoverSender(SendClass map[string]map[[2]int64]*common.Ready2Send, now string) {
	lanxin := []common.Ready2Send{}
	for _, v := range SendClass["LANXIN"] {
		lanxin = append(lanxin, *v)
	}
	go SendRecover("http://lanxinurl:8000/api/v1/lanxin/text", "StreeAlert", map[string]string{"key": "6E358A78-0A5B-49D2-A12F-6A4EB07A9671"}, lanxin, now)
	delete(SendClass, "LANXIN")
	for k := range SendClass {
		hook := []common.Ready2Send{}
		for _, u := range SendClass[k] {
			hook = append(hook, *u)
		}
		go Send2Hook(hook, now, "recover", k[5:])
	}
}

func SendRecover(url string, from string, param map[string]string, content []common.Ready2Send, now string) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 16384)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println("Panic in SendRecover:", e, buf)
		}
	}()
	for _, i := range content {
		msg := []string{"[故障恢复:" + strconv.FormatInt(int64(len(i.Alerts)), 10) + "条] " + i.Alerts[0].Summary}
		for _, j := range i.Alerts {
			duration, _ := time.ParseDuration(strconv.FormatInt(int64(j.Count), 10) + "m")

			id := strconv.FormatInt(j.Id, 10)
			value := strconv.FormatFloat(j.Value, 'f', 2, 64)
			msg = append(msg, "["+duration.String()+"][ID:"+id+"] "+j.Hostname+" 当前值:"+value)
		}
		msg = append(msg, "[时间] "+now)
		data, _ := json.Marshal(common.Msg{
			Content: strings.Join(msg, "\n"),
			From:    from,
			Title:   "Alerts",
			To:      i.User})
		common.HttpPost(url, param, nil, data)
	}
}

func Send2Hook(content []common.Ready2Send, now string, t string, url string) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 16384)
			buf = buf[:runtime.Stack(buf, false)]
			log.Println("Panic in Send2Hook:", e, buf)
		}
	}()
	if t == "recover" {
		for _, i := range content {
			data, _ := json.Marshal(
				struct {
					Type   string               `json:"type"`
					Time   string               `json:"time"`
					RuleId int64                `json:"rule_id"`
					To     []string             `json:"to"`
					Alerts []common.SingleAlert `json:"alerts"`
				}{
					Type:   t,
					RuleId: i.RuleId,
					Time:   now,
					To:     i.User,
					Alerts: i.Alerts,
				})
			common.HttpPost(url, nil, common.GenerateJsonHeader(), data)
		}
	} else {
		for _, i := range content {
			data, _ := json.Marshal(
				struct {
					Type        string               `json:"type"`
					Time        string               `json:"time"`
					RuleId      int64                `json:"rule_id"`
					To          []string             `json:"to"`
					ConfirmLink string               `json:"confirm_link"`
					Alerts      []common.SingleAlert `json:"alerts"`
				}{
					Type:        t,
					RuleId:      i.RuleId,
					Time:        now,
					ConfirmLink: "http://172.30.31.126:32000" + "/alerts_confirm/" + strconv.FormatInt(i.RuleId, 10) + "?start=" + strconv.FormatInt(i.Start, 10),
					To:          i.User,
					Alerts:      i.Alerts,
				})
			common.HttpPost(url, nil, common.GenerateJsonHeader(), data)
		}
	}
}

func init() {
	go func() {
		for {
			timeSleep := time.Duration(90-time.Now().Second()) * time.Second
			log.Println("UpdateMaintainlist time sleep: ", timeSleep)
			time.Sleep(timeSleep)
			UpdateMaintainlist()
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(60-time.Now().Second()) * time.Second)
			now := time.Now().Format("2006-01-02 15:04:05")

			go func() {
				defer func() {
					if e := recover(); e != nil {
						buf := make([]byte, 16384)
						buf = buf[:runtime.Stack(buf, false)]
						log.Println("Panic in timer:", e, buf)
					}
				}()

				var info []Record

				models.Ormer.Table(models.AlertsTable).Where("status=1 AND confirmed_before<?", now).Update("status", 2)

				tx := models.Ormer.Begin()
				tx.Table(models.AlertsTable).Where("status!=0").Update("count", gorm.Expr("count + 1"))
				tx.Table(models.AlertsTable).Select("id,rule_id,value,count,summary,description,hostname,confirmed_before,fired_at,labels").Where("status = ?", 2).Find(&info)

				aggregation := map[int64][]Record{}
				maxCount := map[int64]int{}
				for _, i := range info {
					aggregation[i.RuleId] = append(aggregation[i.RuleId], i)

					// 如果 db 中的 ruleId 在 maxCount 中存在且小于或等于 maxCount 的报警时间则跳过
					if val, ok := maxCount[i.RuleId]; ok && (i.Count <= val) {
						continue
					}

					// 如果 max 中不存在 ruleId 或小于 ruleId 对应的报警时间, 则更新 maxCount
					maxCount[i.RuleId] = i.Count

				}

				common.Rw.RLock()
				ready2send := Filter(aggregation, maxCount)
				common.Rw.RUnlock()
				tx.Commit()

				Sender(ready2send, now)
				common.Lock.Lock()
				recover2send := common.Recover2Send
				common.Recover2Send = map[string]map[[2]int64]*common.Ready2Send{
					"LANXIN": map[[2]int64]*common.Ready2Send{},
				}
				common.Lock.Unlock()
				RecoverSender(recover2send, now)
			}()
		}
	}()
}
