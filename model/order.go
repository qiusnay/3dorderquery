package model

type JdOriginalOrder struct {
	Id                int64   `gorm:"primary_key;AUTO_INCREMENT"`
	OrderId           int64   `gorm:"type:varchar(50);comment:'订单ID';not null;unique_index:IX_orderId" json:"orderId"`
	ParentId          int64   `gorm:"type:varchar(50);comment:'父单的订单ID，仅当发生订单拆分时返回， 0：未拆分，有值则表示此订单为子订单';not null" json:"parentId"`
	FinishTime        string  `gorm:"type:varchar(50);comment:'订单完成时间(时间戳，毫秒)';not null" json:"finishTime"`
	OrderEmt          int64   `gorm:"type:int;comment:'下单设备(1:PC,2:无线)';not null" json:"orderEmt"`
	OrderTime         string  `gorm:"type:varchar(50);comment:'下单时间(时间戳，毫秒)';not null" json:"orderTime"`
	PayMonth          int64   `gorm:"type:varchar(50);comment:'订单维度预估结算时间（格式：yyyyMMdd），0：未结算，订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。';not null" json:"payMonth"`
	Plus              int64   `gorm:"type:int;comment:'下单用户是否为PLUS会员 0：否，1：是';not null" json:"plus"`
	PopId             int64   `gorm:"type:int;comment:'商家ID';not null" json:"popId"`
	UnionId           int64   `gorm:"type:int;comment:'推客的联盟ID';not null" json:"unionId"`
	ValidCode         int64   `gorm:"type:varchar(50);comment:'订单维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效';not null" json:"validCode"`
	ActualCosPrice    float64 `gorm:"type:decimal(18,2);comment:'实际计算佣金的金额。订单完成后，会将误扣除的运费券金额更正。如订单完成后发生退款，此金额会更新。';not null" json:"actualCosPrice"`
	ActualFee         float64 `gorm:"type:decimal(18,2);comment:'推客获得的实际佣金（实际计佣金额*佣金比例*最终比例）。如订单完成后发生退款，此金额会更新。';not null" json:"actualFee"`
	Cid1              int64   `gorm:"type:int;comment:'一级类目ID';not null" json:"cid1"`
	Cid2              int64   `gorm:"type:int;comment:'二级类目ID';not null" json:"cid2"`
	Cid3              int64   `gorm:"type:int;comment:'三级类目ID';not null" json:"cid3"`
	CommissionRate    float64 `gorm:"type:decimal(18,2);comment:'佣金比例';not null" json:"commissionRate"`
	EstimateCosPrice  float64 `gorm:"type:decimal(18,2);comment:'预估计佣金额，即用户下单的金额(已扣除优惠券、白条、支付优惠、进口税，未扣除红包和京豆)，有时会误扣除运费券金额，完成结算时会在实际计佣金额中更正。如订单完成前发生退款，此金额不会更新';not null" json:"estimateCosPrice"`
	EstimateFee       float64 `gorm:"type:decimal(18,2);comment:'推客的预估佣金（预估计佣金额*佣金比例*最终比例），如订单完成前发生退款，此金额不会更新';not null" json:"estimateFee"`
	Ext1              string  `gorm:"type:varchar(50);comment:'推客生成推广链接时传入的扩展字段（需要联系运营开放白名单才能拿到数据）。订单行维度';not null" json:"ext1"`
	FinalRate         float64 `gorm:"type:decimal(18,2);comment:'最终比例（分成比例+补贴比例）';not null" json:"finalRate"`
	FrozenSkuNum      int64   `gorm:"type:int;comment:'商品售后中数量';not null" json:"frozenSkuNum"`
	Pid               string  `gorm:"type:varchar(50);comment:'联盟子站长身份标识，格式：子站长ID_子站长网站ID_子站长推广位ID';not null" json:"pid"`
	PositionId        int64   `gorm:"type:int;comment:'推广位ID,0代表无推广位';not null" json:"positionId"`
	Price             float64 `gorm:"type:decimal(18,2);comment:'商品单价';not null" json:"price"`
	SiteId            int64   `gorm:"type:int;comment:'网站ID，0：无网站';not null" json:"siteId"`
	SkuId             int64   `gorm:"type:varchar(50);comment:'商品ID';not null" json:"skuId"`
	SkuName           string  `gorm:"type:varchar(500);comment:'商品名称';not null" json:"skuName"`
	SkuNum            int64   `gorm:"type:int;comment:'商品数量';not null" json:"skuNum"`
	SkuReturnNum      int64   `gorm:"type:int;comment:'商品已退货数量';not null" json:"skuReturnNum"`
	SubSideRate       float64 `gorm:"type:decimal(18,2);comment:'分成比例';not null" json:"subSideRate"`
	SubUnionId        string  `gorm:"type:varchar(50);comment:'子联盟ID(需要联系运营开放白名单才能拿到数据)';not null" json:"subUnionId"`
	SubsidyRate       float64 `gorm:"type:decimal(18,2);comment:'补贴比例';not null" json:"subsidyRate"`
	TraceType         int64   `gorm:"type:int;comment:'2：同店；3：跨店';not null" json:"traceType"`
	UnionAlias        string  `gorm:"type:varchar(50);comment:'PID所属母账号平台名称（原第三方服务商来源）';not null" json:"unionAlias"`
	UnionTag          string  `gorm:"type:varchar(50);comment:'联盟标签数据（整型的二进制字符串(32位)，目前只返回8位：00000001。数据从右向左进行，每一位为1表示符合联盟的标签特征，第1位：京喜红包，第2位：组合推广订单，第3位：拼购订单，第5位：有效首次购订单（00011XXX表示有效首购，最终奖励活动结算金额会结合订单状态判断，以联盟后台对应活动效果数据报表https://union.jd.com/active为准）。例如：00000001:京喜红包订单，00000010:组合推广订单，00000100:拼购订单，00011000:有效首购，00000111：京喜红包+组合推广+拼购等）';not null" json:"unionTag"`
	UnionTrafficGroup int64   `gorm:"type:int;comment:'渠道组 1：1号店，其他：京东';not null" json:"unionTrafficGroup"`
	HasMore           int64   `gorm:"type:int;comment:'是否还有下一页';not null" json:"has_more"`
}

func (JdOriginalOrder) TableName() string {
	return "tb_jd_original_order"
}

type PddOriginalOrder struct {
	Id                    int64  `gorm:"primary_key;AUTO_INCREMENT"`
	OrderSn               string `gorm:"type:varchar(50);comment:'推广订单编号';not null;unique_index:IX_order_sn" json:"order_sn"`
	GoodsId               int    `gorm:"type:varchar(50);comment:'商品ID';not null" json:"goods_id"`
	GoodsName             string `gorm:"type:varchar(500);comment:'商品标题';not null;" json:"goods_name"`
	GoodsThumbnailUrl     string `gorm:"type:varchar(100);comment:'商品缩略图';not null;" json:"goods_thumbnail_url"`
	GoodsQuantity         int    `gorm:"type:int;comment:'购买商品的数量';not null;" json:"goods_quantity"`
	GoodsPrice            int    `gorm:"type:decimal(18,2);comment:'订单中sku的单件价格，单位为分';not null;" json:"goods_price"`
	OrderAmount           int    `gorm:"type:decimal(18,2);comment:'实际支付金额，单位为分';not null;" json:"order_amount"`
	OrderCreateTime       int    `gorm:"type:varchar(50);comment:'订单生成时间，UNIX时间戳';not null;" json:"order_create_time"`
	OrderSettleTime       int    `gorm:"type:varchar(50);comment:'结算时间';not null;" json:"order_settle_time"`    // 结算时间
	OrderVerifyTime       int    `gorm:"type:varchar(50);comment:'审核时间';not null;" json:"order_verify_time"`    // 审核时间
	OrderReceiveTime      int    `gorm:"type:varchar(50);comment:'确认收货时间';not null;" json:"order_receive_time"` // 收货时间
	OrderPayTime          int    `gorm:"type:varchar(50);comment:'支付时间';not null;" json:"order_pay_time"`
	PromotionRate         int    `gorm:"type:decimal(18,2);comment:'佣金比例，千分比';not null;" json:"promotion_rate"`
	PromotionAmount       int    `gorm:"type:decimal(18,2);comment:'佣金金额，单位为分';not null;" json:"promotion_amount"`
	BatchNo               string `gorm:"type:varchar(50);comment:'结算批次号';default null;" json:"batch_no"`
	OrderStatus           int    `gorm:"type:int;comment:'订单状态： -1 未支付,0-已支付,1-已成团,2-确认收货,3-审核成功,4-审核失败（不可提现）,5-已经结算,8-非多多进宝商品（无佣金订单）';not null;" json:"order_status"`
	OrderStatusDesc       string `gorm:"type:varchar(500);comment:'订单状态描述';not null;" json:"order_status_desc"`
	OrderGroupSuccessTime int    `gorm:"type:varchar(50);comment:'成团时间';not null;" json:"order_group_success_time"`
	OrderModifyAt         int    `gorm:"type:varchar(50);comment:'最后更新时间';not null;" json:"order_modify_at"`
	Type                  int    `gorm:"type:int;comment:'订单推广类型';not null;" json:"type"`
	GroupId               string `gorm:"type:varchar(50);comment:'成团编号';not null;" json:"group_id"`
	AuthDuoId             int    `gorm:"type:varchar(50);comment:'订多多客工具id';not null;" json:"auth_duo_id"`
	ZsDuoId               int    `gorm:"type:varchar(50);comment:'招商多多客id';not null;" json:"zs_duo_id"`
	CustomParameters      string `gorm:"type:varchar(50);comment:'自定义参数';not null;" json:"custom_parameters"`
	PId                   string `gorm:"type:varchar(50);comment:'推广位ID';not null;" json:"p_id"`
	Pid                   string `gorm:"type:varchar(50);comment:'未知ID';not null;" json:"pid"`
	VerifyTime            int    `gorm:"type:varchar(50);comment:'订单审核时间';not null;" json:"verify_time"` // 审核时间
	MatchChannel          int    `gorm:"type:varchar(50);comment:'匹配渠道';not null;" json:"match_channel"`
	CpaNew                int    `gorm:"type:varchar(50);comment:'是否是 cpa 新用户，1表示是，0表示否';default null;" json:"cpa_new"`
	DuoCouponAmount       int    `gorm:"type:decimal(18,2);comment:'优惠券金额';not null;" json:"duo_coupon_amount"`
}

func (PddOriginalOrder) TableName() string {
	return "tb_pdd_original_order"
}
