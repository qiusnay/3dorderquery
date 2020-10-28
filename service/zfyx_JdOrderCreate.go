package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/model"
)

type JdOrderCreate struct {
}

func (s *JdOrderCreate) Sync() {
	//获取当前扫表的索引值
	index, err := model.Redis.Get("jd_order_scan_index").Result()
	if err != nil {
		logger.Info(fmt.Sprintf("redis key not exist,init"))
		model.Redis.Set("jd_order_scan_index", 0, 0).Err()
	}
	Int64Index, _ := strconv.ParseInt(index, 10, 64)
	orders := []model.JdOriginalOrder{}
	//查询京东订单数据
	model.DB.Where("oid > ?", index).Order("oid").Limit(1).Find(&orders)
	insertOrderChannel := make(chan int64)
	for _, order := range orders {
		if Int64Index < order.Oid {
			Int64Index = order.Oid
		}
		go s.CreateOrder(order, insertOrderChannel)
		Did := <-insertOrderChannel
		if Did > 0 {
			go s.CreateOrderItem(order, Did)
		}
	}
	model.Redis.Set("jd_order_scan_index", Int64Index, 0).Err()
}

//订单表同步
func (s *JdOrderCreate) CreateOrder(order model.JdOriginalOrder, cha chan int64) {
	userid, _ := strconv.ParseInt(order.Ext1, 10, 64)
	InserOrder := &model.TbDingdan{
		Ordernum:       strconv.FormatInt(order.OrderId, 10),
		OrdernumParent: strconv.FormatInt(order.ParentId, 10),
		ShopId:         10000,
		Userid:         userid,
		Buydate:        order.OrderTime,
		Amount:         order.ActualCosPrice,
		InputDate:      time.Now().Format("2006-01-02 15:04:05"),
		OrderState:     order.ValidCode,
		FanliState:     0,
		PreCommission:  order.EstimateFee,
		Fanli:          0,
		TrackingCode:   "",
		ShopTitle:      "",
		PayDate:        strconv.FormatInt(order.PayMonth, 10),
		Commission:     order.ActualFee,
		Yujifanli:      0.00,
		Remark:         "",
		ModifyDate:     time.Now().Format("2006-01-02 15:04:05"),
	}
	model.DB.Create(InserOrder)
	cha <- InserOrder.Did
}

//商品表同步
func (s *JdOrderCreate) CreateOrderItem(order model.JdOriginalOrder, Did int64) {
	userid, _ := strconv.ParseInt(order.Ext1, 10, 64)
	model.DB.Create(&model.TbDingdanItems{
		Did:            Did,
		Userid:         userid,
		ShopId:         10000,
		Pid:            strconv.FormatInt(order.SkuId, 10),
		Pnum:           order.SkuNum,
		Price:          order.Price * float64(order.SkuNum),
		Unitprice:      order.Price,
		Comm:           0,
		Cid:            "",
		Fanli:          0,
		Category_id:    0,
		Shop_title:     "",
		Product_title:  order.SkuName,
		Product_status: order.ValidCode,
		Related_pid:    strconv.FormatInt(order.SkuId, 10),
		Itempic:        "",
		Remark:         "",
		ModifyDate:     time.Now().Format("2006-01-02 15:04:05"),
	})
}
