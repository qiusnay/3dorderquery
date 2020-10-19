package service

import (
	"fmt"
	"time"
	"strings"
	"reflect"
	"sort"
	"net/url"
	"github.com/qiusnay/3dorderquery/util"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"github.com/google/logger"
	// "github.com/jinzhu/gorm"
)

type OrderParam struct {
	OrderReq OrderReq `json:"orderReq"`
}

type OrderReq struct {
	PageNo       int    `json:"pageNo"`                 //页码，返回第几页结果
	PageSize     int    `json:"pageSize"`               //每页包含条数，上限为500
	Type         int    `json:"type"`                   //订单时间查询类型(1：下单时间，2：完成时间，3：更新时间)
	StartTime         string `json:"startTime"`                   //查询时间，建议使用分钟级查询，格式：yyyyMMddHH、yyyyMMddHHmm或yyyyMMddHHmmss，如201811031212 的查询范围从12:12:00--12:12:59
	EndTime         string `json:"endTime"`                   //查询时间，建议使用分钟级查询，格式：yyyyMMddHH、yyyyMMddHHmm或yyyyMMddHHmmss，如201811031212 的查询范围从12:12:00--12:12:59
	ChildUnionId int64  `json:"childUnionId,omitempty"` //子站长ID（需要联系运营开通PID账户权限才能拿到数据），childUnionId和key不能同时传入
	Key          string `json:"key,omitempty"`          //其他推客的授权key，查询工具商订单需要填写此项，childUnionid和key不能同时传入
}

type JdUnionOpenOrderRowQueryResponse struct {
	JdUnionOpenOrderRowQueryResponse struct {
		Result string `json:"result"`
		Code   string `json:"code"`
	} `json:"jd_union_open_order_row_query_response"`
}

type OrderResult struct {
	Code      int    `json:"code"`
	HasMore   bool   `json:"hasMore"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	Data      []struct {
		FinishTime string         `json:"finishTime"` //订单完成时间(时间戳，毫秒)
		OrderEmt   int           `json:"orderEmt"`   //下单设备(1:PC,2:无线)
		OrderId    int64         `json:"orderId"`    //订单ID
		OrderTime  string `json:"orderTime"`  //下单时间(时间戳，毫秒)
		ParentId   int64         `json:"parentId"`   //父单的订单ID，仅当发生订单拆分时返回， 0：未拆分，有值则表示此订单为子订单
		PayMonth   int           `json:"payMonth"`   //订单维度预估结算时间（格式：yyyyMMdd），0：未结算，订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
		Plus       int           `json:"plus"`       //下单用户是否为PLUS会员 0：否，1：是
		PopId      int           `json:"popId"`      //商家ID
		SkuList    []struct {
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
			UnionTrafficGroup int     `json:"unionTrafficGroup"` //渠道组 1：1号店，其他：京东
			ValidCode         int     `json:"validCode"`         //sku维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
			SubUnionId        string  `json:"subUnionId"`        //子联盟ID(需要联系运营开放白名单才能拿到数据)
			TraceType         int     `json:"traceType"`         //2：同店；3：跨店
			PayMonth          int     `json:"payMonth"`          //订单行维度预估结算时间（格式：yyyyMMdd） ，0：未结算。订单的预估结算时间仅供参考。账号未通过资质审核或订单发生售后，会影响订单实际结算时间。
			PopId             int64   `json:"popId"`             //商家ID，订单行维度
			Ext1              string  `json:"ext1"`              //推客生成推广链接时传入的扩展字段（需要联系运营开放白名单才能拿到数据）。&lt;订单行维度&gt;
			Price             float64 `json:"price"`             //商品单价
			SkuName           string  `json:"skuName"`           //商品名称
			SubSideRate       float64 `json:"subSideRate"`       //分成比例
			SubsidyRate       float64 `json:"subsidyRate"`       //补贴比例
		}
		UnionId   int64  `json:"unionId"`   //推客的联盟ID
		Ext1      string `json:"ext1"`      //推客生成推广链接时传入的扩展字段，订单维度（需要联系运营开放白名单才能拿到数据）
		ValidCode int    `json:"validCode"` //订单维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,10.无效-开增值税专用发票订单,11.无效-乡村推广员下单,12.无效-自己推广自己下单,13.无效-违规订单,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成,18.已结算（5.9号不再支持结算状态回写展示））注：自2018/7/13起，自己推广自己下单已经允许返佣，故12无效码仅针对历史数据有效
		HasMore   bool   `json:"hasMore"`   //是否还有更多,true：还有数据；false:已查询完毕，没有数据
	}
}

type apiconfig struct {
	Appkey   string
	APPsecret   string
	Method    string
}

type tomlConfig struct {
	JD   apiconfig
}

type SysParam struct {
	Method       string `json:"method"`
	App_key      string `json:"app_key"`
	Access_token string `json:"access_token"`
	Timestamp    string `json:"timestamp"`
	Format       string `json:"format"`
	V            string `json:"v"`
	Sign_method  string `json:"sign_method"`
	Param_json   string `json:"param_json"`
}

type Jdsdk struct {
	RequestParam SysParam
	SignAndUri   string
	RequestRul   string
}

func (J *Jdsdk) GetParams() string {
	ParamStruct := OrderParam{}
	ParamStruct.OrderReq.StartTime = "2020-10-15 11:04:05"
	ParamStruct.OrderReq.EndTime = "2020-10-15 12:00:05"
	ParamStruct.OrderReq.PageNo = 1
	ParamStruct.OrderReq.PageSize = 10
	ParamStruct.OrderReq.Type = 1
	bytes, _ := json.Marshal(ParamStruct)
	return string(bytes)
}

//获取订单
func (J *Jdsdk) GetOrders() (result *OrderResult) {
	Param := J.GetParams()
	J.SetSignJointUrlParam(Param)
	var urls strings.Builder
	urls.WriteString(util.JD_HOST)
	urls.WriteString(J.SignAndUri)
	body, _ := util.HttpGet(urls.String())
	// fmt.Println(string(body))
	response := &JdUnionOpenOrderRowQueryResponse{}
	e := json.Unmarshal([]byte(body), &response)
	logger.Info(fmt.Sprintf("get jd order %+v", response))
	// fmt.Println(response)
	if e != nil {
		panic(e)
	}
	result = &OrderResult{}
	e = json.Unmarshal([]byte(response.JdUnionOpenOrderRowQueryResponse.Result), result)
	if e != nil {
		panic(e)
	}
	return result
}

//生成请求参数和签名
func (J *Jdsdk) SetSignJointUrlParam(paramjson string) *Jdsdk {
	var config tomlConfig
	if _, err := toml.DecodeFile("./conf/conf.toml", &config); err != nil {
		fmt.Println(err)
		return J
	}
	J.RequestParam.App_key = config.JD.Appkey
	J.RequestParam.Format  = "json"
	J.RequestParam.V       = "1.0"
	J.RequestParam.Method = config.JD.Method

	J.RequestParam.Sign_method = "md5"
	J.RequestParam.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	J.RequestParam.Param_json = paramjson
	
	values := reflect.ValueOf(J.RequestParam)
	keys := reflect.TypeOf(J.RequestParam)
	count := values.NumField()
	SortSlice := util.Items{}
	for i := 0; i < count; i++ {
		value := values.Field(i)
		key := keys.Field(i)
		switch value.Kind() {
		case reflect.String:
			SortSlice = append(SortSlice, util.Onestruct{strings.ToLower(key.Name), value.String()})
		case reflect.Int:
			SortSlice = append(SortSlice, util.Onestruct{strings.ToLower(key.Name), value.String()})
		}
	}
	sort.Sort(SortSlice)
	
	var builder strings.Builder
	u := url.Values{}
	builder.WriteString(config.JD.APPsecret)
	for _, person := range SortSlice {
		if person.Value == "" {
			continue
		}
		builder.WriteString(strings.ToLower(person.Key) + person.Value)
		u.Add(strings.ToLower(person.Key), person.Value)
	}
	builder.WriteString(config.JD.APPsecret)

	//生成签名
	u.Add("sign", strings.ToUpper(util.Md5(builder.String())))
	//拼接参数
	J.SignAndUri = u.Encode()
	return J
}

