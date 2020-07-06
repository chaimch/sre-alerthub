package models

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
)

const (
	GroupsTable = "group"
)

type Groups struct {
	Id   int64  `gorm:"auto" json:"id,omitempty"`
	Name string `gorm:"unique;size:255" json:"name"`
	User string `gorm:"size:1023" json:"user"`
}

type HttpRes struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   []struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Mobile  string `json:"mobile"`
		Email   string `json:"email"`
		AddTime string `json:"add_time"`
		Account string `json:"account"`
	} `json:"data"`
}

func (*Groups) TableName() string {
	return "group"
}

func SendAlertsFor(VUG *common.ValidUserGroup) []string {
	var userList []string
	if VUG.User != "" {
		userList = strings.Split(VUG.User, ",")
	}
	if VUG.Group != "" {
		var groups []*Groups
		Ormer.Table(GroupsTable).Where("name IN (?)", strings.Split(VUG.Group, ",")).Find(&groups)
		for _, v := range groups {
			userList = append(userList, strings.Split(v.User, ",")...)
		}
	}
	if VUG.DutyGroup != "" {
		date := time.Now().Format("2006-1-2")
		idList := strings.Split(VUG.DutyGroup, ",")
		for _, id := range idList {
			// res, _ := common.HttpGet(beego.AppConfig.String("DutyGroupUrl"), map[string]string{"teamId": id, "day": date}, nil)
			res, _ := common.HttpGet("http://dutygroupurl:8000/Api/getDutyUser", map[string]string{"teamId": id, "day": date}, nil)
			info := HttpRes{}
			jsonDataFromHttp, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(jsonDataFromHttp, &info)
			for _, i := range info.Data {
				userList = append(userList, i.Account)
			}
		}
	}
	hashMap := map[string]bool{}
	for _, name := range userList {
		hashMap[name] = true
	}
	res := []string{}
	for key := range hashMap {
		res = append(res, key)
	}
	return res
}
