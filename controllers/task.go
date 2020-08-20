package controllers

import (
	"go-sghen/models"
	"strconv"
	"time"

	"github.com/robfig/cron"
)

var taskManager = cron.New()

// TaskController ...
type TaskController struct {
	BaseController
}

// InitTask ...
func InitTask() {
	taskManager.AddFunc("CRON_TZ=Asia/Shanghai 0 23 * * ?", func() {
		models.MConfig.MLogger.Info("每天23点定时任务")

		a := time.Now()
		b, _ := time.Parse("2006-01-02", "2020-03-17")
		d := a.Sub(b)

		smsCode := 3280 + int(d.Hours()/24)*20
		phones := []int64{models.MConfig.SmsMobile0, models.MConfig.SmsMobile1}
		for _, phone := range phones {
			if phone > 10000000000 {
				sendSmsCode(phone, strconv.Itoa(smsCode))
			}
		}
	})
	taskManager.AddFunc("CRON_TZ=Asia/Shanghai 0 0 * * ?", dynamicAPICacheTask)
	taskManager.Run()
}
