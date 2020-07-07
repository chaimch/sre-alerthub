package timer

import (
	"log"
	"math"
	"runtime"
	"strconv"
	"time"

	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/common"
	"gitlab.mobiuspace.net/mobiuspace/sre-team/sre-alerthub/models"
)

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

	}()
}
