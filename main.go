package main

import (
	"os"
	"time"
	"fmt"
	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/service"
	"github.com/qiusnay/3dorderquery/model"
)

func main() {
	const logPath = "./log/3dorderquery.log"
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
	  logger.Fatalf("Failed to open log file: %v", err)
	}
	logger.Init("Logger", false, true, lf)
	model.DbStart()
	for {
		timer1 := time.NewTimer(time.Second * 60)
		<-timer1.C
		service := new(service.Jdsdk)
		// starttime := time.Now().Add(-time.Minute * 60).Format("2006-01-02 15:04:05")
		// endtime := time.Now().Format("2006-01-02 15:04:05")
		starttime := "2020-10-15 11:04:05"
		endtime := "2020-10-15 12:00:05"
		logger.Info(fmt.Sprintf("starttime %v %v", starttime, endtime))
		service.GetOrders(starttime, endtime)
	}
}