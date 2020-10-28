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
	model.RedisStart()
	for {
		timer1 := time.NewTimer(time.Second * 2)
		<-timer1.C

		var ShopSdk service.UnionSDKAPI

		// starttime := time.Now().Add(-time.Minute * 60).Format("2006-01-02 15:04:05")
		// endtime := time.Now().Format("2006-01-02 15:04:05")
		//京东订单抓单
		ShopSdk = new(service.Jdsdk)
		JdOrders := ShopSdk.GetOrders("2020-10-15 11:04:05", "2020-10-15 12:00:05")
		logger.Info(fmt.Sprintf("response jd %+v", JdOrders))
		// //拼多多订单抓单
		ShopSdk = new(service.Pddsdk)
		PddOrders := ShopSdk.GetOrders("2020-10-21 13:00:00", "2020-10-21 14:00:00")
		logger.Info(fmt.Sprintf("response pdd %+v", PddOrders))

		//京东商城抓商品
		//频道id：
		//1-好券商品,2-超级大卖场,10-9.9专区,22-热销爆品,23-为你推荐,24-数码家电,25-超市,26-母婴玩具,27-家具日用,28-美妆穿搭,
		//29-医药保健,30-图书文具,31-今日必推,32-品牌好货,33-秒杀商品,34-拼购商品,109-新品首发,110-自营,125-首购商品,129-高佣榜单,130-视频商品
		ShopSdkS := new(service.JdItemsdk)
		for _, brand := range []int{1, 2, 10, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 109, 110, 125, 129, 130} {
			ShopSdkS.GetJdItems(brand)
		}

		//京东订单同步tb_dingdan,tb_dingdan_items
		syncJdOrderModel := new(service.JdOrderCreate)
		syncJdOrderModel.Sync()

		//拼多多订单同步tb_dingdan,tb_dingdan_items
		syncPddOrderModel := new(service.PddOrderCreate)
		syncPddOrderModel.Sync()
		// logger.Info(fmt.Sprintf("response jd %+v", JdOrders))
	}
}
