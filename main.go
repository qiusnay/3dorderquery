package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/service"
	"github.com/qiusnay/3dorderquery/util"
)

func main() {
	Log := util.NewLogger("main")
	Log.Info(fmt.Sprintf("作业启动" + strconv.Itoa(runtime.NumGoroutine())))
	model.DbStart()
	model.RedisStart()

	var wg sync.WaitGroup
	wg.Add(6)
	c := ZfyxSdk{}
	//京东订单抓取 3分钟 - 过去1小时
	go c.JdOrderFetchBy3min(3)
	//拼多多订单抓取 3分钟 - 过去1小时
	go c.PddOrderFetchBy3min(3)
	//京东商城抓商品 30分钟 - 最新 50条
	go c.FetchJdGoodsBy30Min(30)
	//拼多多商城抓商品  30分钟 - 最新500条
	go c.FetchPddGoodsBy30Min(30)
	//订单同步 - 京东 1分钟
	go c.SyncJdOrderFromOriginal(1)
	//订单同步 - 拼多多 1分钟
	go c.SyncPddOrderFromOriginal(1)
	wg.Wait()
}

type ZfyxSdk struct{}

//京东订单同步 1分钟
func (z *ZfyxSdk) SyncJdOrderFromOriginal(TimeMinutes int) {
	log := util.NewLogger("sync_order_jd")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("SyncJdOrderFromOriginal", TimeMinutes)
			log.Error(fmt.Sprintf("京东订单同步异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始京东订单同步,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C
		//京东订单同步tb_dingdan,tb_dingdan_items
		syncJdOrderModel := new(service.JdOrderCreate)
		syncJdOrderModel.Sync()
	}
}

//拼多多订单同步 1分钟
func (z *ZfyxSdk) SyncPddOrderFromOriginal(TimeMinutes int) {
	log := util.NewLogger("sync_order_pdd")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("SyncPddOrderFromOriginal", TimeMinutes)
			log.Error(fmt.Sprintf("拼多多订单同步异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始拼多多订单同步,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C
		//拼多多订单同步tb_dingdan,tb_dingdan_items
		syncPddOrderModel := new(service.PddOrderCreate)
		syncPddOrderModel.Sync()
	}
}

//拼多多抓商品 30min
//频道id：
//0-1.9包邮, 1-今日爆款, 2-品牌清仓,3-相似商品推荐,4-猜你喜欢,5-实时热销,6-实时收益,7-今日畅销,8-高佣榜单，默认1
func (z *ZfyxSdk) FetchPddGoodsBy30Min(TimeMinutes int) {
	log := util.NewLogger("fetch_pdd_goods")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("FetchPddGoodsBy30Min", TimeMinutes)
			log.Error(fmt.Sprintf("拼多多商品抓取异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始拼多多商品抓取,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C
		ShopSdk := new(service.PddItemsdk)
		for _, brand := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8} {
			FetchUrl := ShopSdk.FetchPddItems(brand)
			log.Info(fmt.Sprintf("当前抓取地址: %s", FetchUrl))
		}
	}
}

//京东商城抓商品
//频道id：
//1-好券商品,2-超级大卖场,10-9.9专区,22-热销爆品,23-为你推荐,24-数码家电,25-超市,26-母婴玩具,27-家具日用,28-美妆穿搭,
//29-医药保健,30-图书文具,31-今日必推,32-品牌好货,33-秒杀商品,34-拼购商品,109-新品首发,110-自营,125-首购商品,129-高佣榜单,130-视频商品
func (z *ZfyxSdk) FetchJdGoodsBy30Min(TimeMinutes int) {
	log := util.NewLogger("fetch_jd_goods")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("FetchJdGoodsBy30Min", TimeMinutes)
			log.Error(fmt.Sprintf("京东商品抓取异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始京东商品抓取,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C
		ShopSdk := new(service.JdItemsdk)
		for _, brand := range []int{1, 2, 10, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 109, 110, 125, 129, 130} {
			FetchUrl := ShopSdk.FetchJdItems(brand, log)
			log.Info(fmt.Sprintf("当前抓取地址: %s", FetchUrl))
		}
	}
}

//京东订单抓取 3分钟
func (z *ZfyxSdk) JdOrderFetchBy3min(TimeMinutes int) {
	log := util.NewLogger("fetch_jd_order")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("JdOrderFetchBy3min", TimeMinutes)
			log.Error(fmt.Sprintf("京东订单抓取异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始京东订单抓取,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C

		//京东订单抓单
		// starttime := time.Now().Add(-time.Minute * 60).Format("2006-01-02 15:04:05")
		// endtime := time.Now().Format("2006-01-02 15:04:05")
		ShopSdk := new(service.Jdsdk)
		FetchUrl := ShopSdk.FetchOrders("2020-10-15 11:04:05", "2020-10-15 12:00:05")
		log.Info(fmt.Sprintf("当前抓取地址: %s", FetchUrl))
	}
}

//拼多多订单抓取 3分钟
func (z *ZfyxSdk) PddOrderFetchBy3min(TimeMinutes int) {
	log := util.NewLogger("fetch_pdd_order")
	defer func() {
		if err := recover(); err != nil {
			z.RestartRoutine("PddOrderFetchBy3min", TimeMinutes)
			log.Error(fmt.Sprintf("拼多多订单抓取异常,错误信息 : %v, 自动恢复 : %s", err, time.Now().Format("2006-01-02 15:04:05")))
		}
	}()
	for {
		log.Info(fmt.Sprintf("开始拼多多订单抓取,当前时间 : " + time.Now().Format("2006-01-02 15:04:05")))
		timer1 := time.NewTimer(time.Second * time.Duration(60*TimeMinutes))
		<-timer1.C

		//拼多多订单抓单
		ShopSdk := new(service.Pddsdk)
		FetchUrl := ShopSdk.FetchOrders("2020-10-30 13:00:00", "2020-10-31 12:00:00")
		log.Info(fmt.Sprintf("当前抓取地址: %s", FetchUrl))
	}
}

//反射调用
func (z *ZfyxSdk) RestartRoutine(Method string, TimeMinutes int) {
	v := reflect.ValueOf(&ZfyxSdk{})
	argsFunc := v.MethodByName(Method)
	args := []reflect.Value{reflect.ValueOf(TimeMinutes)}
	argsFunc.Call(args)
}
