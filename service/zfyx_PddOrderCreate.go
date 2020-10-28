package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/model"
)

type PddOrderCreate struct{}

func (s *PddOrderCreate) Sync() {
	//获取当前扫表的索引值
	index, err := model.Redis.Get("pdd_order_scan_index").Result()
	if err != nil {
		logger.Info(fmt.Sprintf("redis key not exist,init"))
		model.Redis.Set("pdd_order_scan_index", 0, 0).Err()
	}
	Int64Index, _ := strconv.ParseInt(index, 10, 64)
	orders := []model.PddOriginalOrder{}
	//查询拼多多订单数据
	model.DB.Where("id > ?", index).Order("id").Limit(1).Find(&orders)
	insertOrderChannel := make(chan int64)
	for _, order := range orders {
		if Int64Index < order.Id {
			Int64Index = order.Id
		}

		go s.CreateOrder(order, insertOrderChannel)
		Did := <-insertOrderChannel
		if Did > 0 {
			go s.CreateOrderItem(order, Did)
		}
	}
	model.Redis.Set("pdd_order_scan_index", Int64Index, 0).Err()
}

//订单表同步
func (s *PddOrderCreate) CreateOrder(order model.PddOriginalOrder, cha chan int64) {
	userid, _ := strconv.ParseInt(order.CustomParameters, 10, 64)
	InserOrder := &model.TbDingdan{
		Ordernum:       order.OrderSn,
		OrdernumParent: order.OrderSn,
		ShopId:         10001,
		Userid:         userid,
		Buydate:        time.Unix(order.OrderCreateTime, 0).Format("2006-01-02 15:04:05"),
		Amount:         s.getRmbYuan(order.OrderAmount),
		InputDate:      time.Now().Format("2006-01-02 15:04:05"),
		OrderState:     order.OrderStatus,
		FanliState:     0,
		PreCommission:  s.getRmbYuan(order.PromotionAmount),
		Fanli:          0,
		TrackingCode:   "",
		ShopTitle:      "",
		PayDate:        time.Unix(order.OrderGroupSuccessTime, 0).Format("2006-01-02 15:04:05"),
		Commission:     s.getRmbYuan(order.PromotionAmount),
		Yujifanli:      0.00,
		Remark:         "",
		ModifyDate:     time.Now().Format("2006-01-02 15:04:05"),
	}
	model.DB.Create(InserOrder)
	cha <- InserOrder.Did
}

//商品表同步
func (s *PddOrderCreate) CreateOrderItem(order model.PddOriginalOrder, Did int64) {
	userid, _ := strconv.ParseInt(order.CustomParameters, 10, 64)
	model.DB.Create(&model.TbDingdanItems{
		Did:            Did,
		Userid:         userid,
		ShopId:         10001,
		Pid:            strconv.FormatInt(order.GoodsId, 10),
		Pnum:           order.GoodsQuantity,
		Price:          s.getRmbYuan(order.GoodsPrice * float64(order.GoodsQuantity)),
		Unitprice:      s.getRmbYuan(order.GoodsPrice),
		Comm:           0,
		Cid:            "",
		Fanli:          0,
		Category_id:    0,
		Shop_title:     "",
		Product_title:  order.GoodsName,
		Product_status: order.OrderStatus,
		Related_pid:    strconv.FormatInt(order.GoodsId, 10),
		Itempic:        order.GoodsThumbnailUrl,
		Remark:         "",
		ModifyDate:     time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (s *PddOrderCreate) getRmbYuan(cash float64) float64 {
	AmountYuan, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cash/100), 64)
	return AmountYuan
}
