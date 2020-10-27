package service

// "encoding/json"

type OrderResult struct {
	Code      int    `json:"code"`
	HasMore   bool   `json:"hasMore"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	Data      []JdOrders
}

type JdOrders struct {
	FinishTime string `json:"finishTime"` //订单完成时间(时间戳，毫秒)
	OrderEmt   int64  `json:"orderEmt"`   //下单设备(1:PC,2:无线)
	OrderId    int64  `json:"orderId"`    //订单ID
	OrderTime  string `json:"orderTime"`  //下单时间(时间戳，毫秒)
	ParentId   int64  `json:"parentId"`   //父单的订单ID，仅当发生订单拆分时返回， 0：未拆分，有值则表示此订单为子订单
	PayMonth   int64  `json:"payMonth"`   //订单维度预估结算时间（格式：yyyyMMdd），0：未结算，订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
	Plus       int64  `json:"plus"`       //下单用户是否为PLUS会员 0：否，1：是
	PopId      int64  `json:"popId"`      //商家ID
	//订单包含的商品信息列表
	ActualCosPrice    float64 `json:"actualCosPrice"`    //实际计算佣金的金额。订单完成后，会将误扣除的运费券金额更正。如订单完成后发生退款，此金额会更新。
	ActualFee         float64 `json:"actualFee"`         //推客获得的实际佣金（实际计佣金额*佣金比例*最终比例）。如订单完成后发生退款，此金额会更新。
	CommissionRate    float64 `json:"commissionRate"`    //佣金比例
	EstimateCosPrice  float64 `json:"estimateCosPrice"`  //预估计佣金额，即用户下单的金额(已扣除优惠券、白条、支付优惠、进口税，未扣除红包和京豆)，有时会误扣除运费券金额，完成结算时会在实际计佣金额中更正。如订单完成前发生退款，此金额不会更新
	EstimateFee       float64 `json:"estimateFee"`       //推客的预估佣金（预估计佣金额*佣金比例*最终比例），如订单完成前发生退款，此金额不会更新
	FinalRate         float64 `json:"finalRate"`         //最终比例（分成比例+补贴比例）
	Cid1              int64   `json:"cid1"`              //一级类目ID
	FrozenSkuNum      int64   `json:"frozenSkuNum"`      //商品售后中数量
	Pid               string  `json:"pid"`               //联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID
	PositionId        int64   `json:"positionId"`        //推广位ID,0代表无推广位
	Cid2              int64   `json:"cid2"`              //二级类目ID
	SiteId            int64   `json:"siteId"`            //网站ID，0：无网站
	SkuId             int64   `json:"skuId"`             //商品ID
	SkuNum            int64   `json:"skuNum"`            //商品数量
	SkuReturnNum      int64   `json:"skuReturnNum"`      //商品已退货数量
	Cid3              int64   `json:"cid3"`              //三级类目ID
	UnionAlias        string  `json:"unionAlias"`        //PID所属母账号平台名称（原第三方服务商来源）
	UnionTag          string  `json:"unionTag"`          //联盟标签数据（整型的二进制字符串(32位)，目前只返回8位：00000001。数据从右向左进行，每一位为1表示符合联盟的标签特征，第1位：京喜红包，第2位：组合推广订单，第3位：拼购订单，第5位：有效首次购订单（00011XXX表示有效首购，最终奖励活动结算金额会结合订单状态判断，以联盟后台对应活动效果数据报表https://union.jd.com/active为准）。例如：00000001:京喜红包订单，00000010:组合推广订单，00000100:拼购订单，00011000:有效首购，00000111：京喜红包+组合推广+拼购等）
	UnionTrafficGroup int64   `json:"unionTrafficGroup"` //渠道组 1：1号店，其他：京东
	SubUnionId        string  `json:"subUnionId"`        //子联盟ID(需要联系运营开放白名单才能拿到数据)
	TraceType         int64   `json:"traceType"`         //2：同店；3：跨店
	Price             float64 `json:"price"`             //商品单价
	SkuName           string  `json:"skuName"`           //商品名称
	SubSideRate       float64 `json:"subSideRate"`       //分成比例
	SubsidyRate       float64 `json:"subsidyRate"`       //补贴比例
	UnionId           int64   `json:"unionId"`           //推客的联盟ID
	Ext1              string  `json:"ext1"`              //推客生成推广链接时传入的扩展字段，订单维度（需要联系运营开放白名单才能拿到数据）
	ValidCode         int64   `json:"validCode"`         //订单维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
	HasMore           bool    `json:"hasMore"`           //是否还有更多,true：还有数据；false:已查询完毕，没有数据
}

type JdItemOriginal struct {
	SkuId                 int     `json:"skuId"`
	Pid                   int     `json:"pid"`
	WareId                int     `json:"wareId"`
	SkuName               string  `json:"skuName"`
	Cid1                  int     `json:"cid1"`
	Cid1Name              string  `json:"cid1Name"`
	Cid2                  int     `json:"cid2"`
	Cid2Name              string  `json:"cid2Name"`
	Cid3                  int     `json:"cid3`
	Cid3Name              string  `json:"cid3Name"`
	BrandCode             int     `json:"brandCode"`
	BrandName             string  `json:"brandName"`
	Owner                 string  `json:"owner"`
	ImageUrl              string  `json:"imageUrl"`
	ImgList               string  `json:"imgList"`
	Vid                   int     `json:"vid`
	PcPrice               float64 `json:"pcPrice"`
	WlPrice               float64 `json:"wlPrice"`
	PcCommissionShare     float64 `json:"pcCommissionShare"`
	WlCommissionShare     float64 `json:"wlCommissionShare"`
	PcCommission          float64 `json:"pcCommission"`
	WlCommission          float64 `json:"wlCommission"`
	HasCoupon             int     `json:"hasCoupon"`
	IsHot                 int     `json:"isHot"`
	CouponId              int     `json:"couponId"`
	CouponLink            string  `json:"couponLink"`
	RfId                  int     `json:"rfId`
	Comments              int     `json:"comments"`
	GoodComments          int     `json:"goodComments"`
	GoodCommentsShare     float64 `json:"goodCommentsShare"`
	VenderName            string  `json:"venderName"`
	InOrderCount30Days    int     `json:"inOrderCount30Days"`
	InOrderCount30DaysSku int     `json:"inOrderCount30DaysSku"`
	IsPinGou              int     `json:"isPinGou"`
	PingouActiveId        int     `json:"pingouActiveId"`
	PingouPrice           float64 `json:"pingouPrice"`
	PingouTmCount         int     `json:"pingouTmCount"`
	EliteId               int     `json:"eliteId"`
	CreateTime            string  `json:"create_time"`
	UpdateTime            string  `json:"update_time"`
}

type PddOrder struct {
	OrderSn               string `json:"order_sn"`
	GoodsId               int    `json:"goods_id"`
	GoodsName             string `json:"goods_name"`
	GoodsThumbnailUrl     string `json:"goods_thumbnail_url"`
	GoodsQuantity         int    `json:"goods_quantity"`
	GoodsPrice            int    `json:"goods_price"`
	OrderAmount           int    `json:"order_amount"`
	OrderCreateTime       int    `json:"order_create_time"`
	OrderSettleTime       int    `json:"order_settle_time"`  // 结算时间
	OrderVerifyTime       int    `json:"order_verify_time"`  // 审核时间
	OrderReceiveTime      int    `json:"order_receive_time"` // 收货时间
	OrderPayTime          int    `json:"order_pay_time"`
	PromotionRate         int    `json:"promotion_rate"`
	PromotionAmount       int    `json:"promotion_amount"`
	BatchNo               string `json:"batch_no"`
	OrderStatus           int    `json:"order_status"`
	OrderStatusDesc       string `json:"order_status_desc"`
	VerifyTime            int    `json:"verify_time"`
	OrderGroupSuccessTime int    `json:"order_group_success_time"`
	OrderModifyAt         int    `json:"order_modify_at"`
	Type                  int    `json:"type"`
	GroupId               int    `json:"group_id"`
	AuthDuoId             int    `json:"auth_duo_id"`
	ZsDuoId               int    `json:"zs_duo_id"`
	CustomParameters      string `json:"custom_parameters"`
	Pid                   string `json:"pid"` // common use pid
	PId                   string `json:"p_id"`
	MatchChannel          int    `json:"match_channel"`
	DuoCouponAmount       int    `json:"duo_coupon_amount"`
}
