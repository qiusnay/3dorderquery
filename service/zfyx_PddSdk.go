package service

import (
	"encoding/json"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/util"
)

type OrderListGetResponse struct {
	OrderListGetResponse struct {
		OrderList  []PddOrder `json:"order_list"`
		TotalCount int        `json:"total_count"`
	} `json:"order_list_get_response"`
}

type Pddsdk struct {
	SignAndUri string
}

type PddOrderReq struct {
	Type              string `json:"type"`
	Client_id         string `json:"client_id"`
	Timestamp         string `json:"timestamp"`
	Data_type         string `json:"data_type"`
	Page              string `json:"pageNo"`    //页码，返回第几页结果
	Page_size         string `json:"pageSize"`  //每页包含条数，上限为500
	End_update_time   string `json:"startTime"` //查询时间，建议使用分钟级查询，格式：yyyyMMddHH、yyyyMMddHHmm或yyyyMMddHHmmss，如201811031212 的查询范围从12:12:00--12:12:59
	Start_update_time string `json:"endTime"`   //查询时间，建议使用分钟级查询，格式：yyyyMMddHH、yyyyMMddHHmm或yyyyMMddHHmmss，如201811031212 的查询范围从12:12:00--12:12:59
}

func (J *Pddsdk) GetParams(start string, end string) PddOrderReq {
	ParamStruct := PddOrderReq{}
	startUnix, _ := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	ParamStruct.Start_update_time = strconv.FormatInt(startUnix.Unix(), 10)
	endUnix, _ := time.ParseInLocation("2006-01-02 15:04:05", end, time.Local)
	ParamStruct.End_update_time = strconv.FormatInt(endUnix.Unix(), 10)
	ParamStruct.Page = strconv.Itoa(1)
	ParamStruct.Page_size = strconv.Itoa(10)
	return ParamStruct
}

type configpdd struct {
	Pdd Apiconfig `toml:"pdd"`
}

var pddConf configpdd

//获取订单
func (J *Pddsdk) GetOrders(start string, end string) interface{} {
	util.Config().Bind("conf", "thirdpartysdk", &pddConf)
	// logger.Info(fmt.Sprintf("get pdd conf %+v", pddConf))
	Param := J.GetParams(start, end)
	paramsString, _ := json.Marshal(Param)
	J.SetSignJointUrlParam(string(paramsString))
	var urls strings.Builder
	urls.WriteString(pddConf.Pdd.HOST)
	urls.WriteString(J.SignAndUri)
	body, _ := util.HttpGet(urls.String())
	response := &OrderListGetResponse{}
	e := json.Unmarshal([]byte(body), &response)

	if e != nil {
		panic(e)
	}
	for _, ord := range response.OrderListGetResponse.OrderList {
		// logger.Info(fmt.Sprintf("response %+v", ord))
		model.DB.Table("tb_pdd_original_order").Create(&ord)
	}
	return *response
}

//生成请求参数和签名
func (J *Pddsdk) SetSignJointUrlParam(paramstr string) {
	param := PddOrderReq{}
	json.Unmarshal([]byte(paramstr), &param)
	param.Client_id = pddConf.Pdd.APPKEY
	param.Data_type = "json"
	param.Type = pddConf.Pdd.METHOD

	param.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	values := reflect.ValueOf(param)
	keys := reflect.TypeOf(param)
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
	builder.WriteString(pddConf.Pdd.APPSECRET)
	for _, person := range SortSlice {
		if person.Value == "" {
			continue
		}
		builder.WriteString(strings.ToLower(person.Key) + person.Value)
		u.Add(strings.ToLower(person.Key), person.Value)
	}
	builder.WriteString(pddConf.Pdd.APPSECRET)

	//生成签名
	u.Add("sign", strings.ToUpper(util.Md5(builder.String())))
	//拼接参数
	J.SignAndUri = u.Encode()
}
