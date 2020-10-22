package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/service"
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
		timer1 := time.NewTimer(time.Second * 6)
		<-timer1.C

		var ShopSdk service.UnionSDKAPI

		// starttime := time.Now().Add(-time.Minute * 60).Format("2006-01-02 15:04:05")
		// endtime := time.Now().Format("2006-01-02 15:04:05")
		//京东订单抓单
		ShopSdk = new(service.Jdsdk)
		JdOrders := ShopSdk.GetOrders("2020-10-15 11:04:05", "2020-10-15 12:00:05")
		logger.Info(fmt.Sprintf("response jd %+v", JdOrders))
		//拼多多订单抓单
		ShopSdk = new(service.Pddsdk)
		PddOrders := ShopSdk.GetOrders("2020-10-21 13:00:00", "2020-10-21 14:00:00")
		logger.Info(fmt.Sprintf("response pdd %+v", PddOrders))

	}
}
