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
	taskManager.AddFunc("CRON_TZ=Asia/Shanghai 0 22 * * ?", func() {
		models.MConfig.MLogger.Info("每天22点定时任务")

		a := time.Now()
		b, _ := time.Parse("2006-01-02", "2020-03-17")
		d := a.Sub(b)

		smsCode := 3280 + int(d.Hours()/24)*20
		phones := []int64{15625045984, 13570578655}
		for _, phone := range phones {
			sendSmsCode(phone, strconv.Itoa(smsCode))
		}
	})
	taskManager.Run()
}
