package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/util"
)

type JdItemsAllParam struct {
	Method       string `json:"method"`
	App_key      string `json:"app_key"`
	Timestamp    string `json:"timestamp"`
	Format       string `json:"format"`
	V            string `json:"v"`
	Sign_method  string `json:"sign_method"`
	Access_Token string `json:"access_token"`
	PageParam    string `json:"pageParam"`  //页码，返回第几页结果
	QueryParam   string `json:"queryParam"` //每页包含条数，上限为500
	Param_json   string `json:"param_json"`
}

type JdKplOpenUnionSearchByelitedResponse struct {
	JdKplOpenUnionSearchByelitedResponse struct {
		Data       []model.JdItemOriginal `json:"data"`
		Code       string                 `json:"code"`
		EliteId    int                    `json:"eliteId"`
		EliteName  string                 `json:"eliteName"`
		TotalCount int                    `json:"totalCount"`
	} `json:"jd_kpl_open_union_search_byelited_response"`
}

type JdItemReq struct {
	UnionSearchParam map[string]string `json:"unionSearchParam"` //页码，返回第几页结果
	// QueryParam map[string]string `json:"queryParam"` //每页包含条数，上限为500
	// OrderField string            `json:"orderField"`
}

type JdItemsdk struct {
	RequestParam JdItemsAllParam
	SignAndUri   string
}

func (J *JdItemsdk) GetParams(brand int) string {
	ParamStruct := JdItemReq{}
	ItemPage := map[string]string{"eliteId": strconv.Itoa(brand), "pageIndex": "1", "pageSize": "10", "sortType": "desc"}
	ParamStruct.UnionSearchParam = ItemPage

	// ItemQuery := map[string]string{"keywords": "phone", "skuId": ""}
	// ParamStruct.QueryParam = ItemQuery
	// ParamStruct.OrderField = "9"
	bytes, _ := json.Marshal(ParamStruct)
	return string(bytes)
}

//获取订单
func (J *JdItemsdk) FetchJdItems(brand int) interface{} {
	util.Config().Bind("conf", "thirdpartysdk", &conf)
	// logger.Info(fmt.Sprintf("get jd order %+v", conf))
	Param := J.GetParams(brand)
	J.SetSignJointUrlParam(Param)
	var urls strings.Builder
	urls.WriteString(conf.Jd.HOST)
	urls.WriteString(J.SignAndUri)
	body, _ := util.HttpGet(urls.String())
	fmt.Println(urls.String())
	// logger.Info(fmt.Sprintf("jd response %+v", string(body)))
	response := &JdKplOpenUnionSearchByelitedResponse{}
	e := json.Unmarshal([]byte(body), &response)
	if e != nil {
		panic(e)
	}
	// logger.Info(fmt.Sprintf("get jd item %+v", response.JdKplOpenUnionSearchByelitedResponse.Data))
	for _, item := range response.JdKplOpenUnionSearchByelitedResponse.Data {
		result := model.JdItemOriginal{}
		item.UpdateTime = strconv.FormatInt(time.Now().Unix(), 10)
		item.CreateTime = strconv.FormatInt(time.Now().Unix(), 10)
		item.EliteId = response.JdKplOpenUnionSearchByelitedResponse.EliteId
		model.DB.Table("tb_jd_original_items").Where("sku_id = ?", item.SkuId).First(&result)
		// logger.Info(fmt.Sprintf("dddddd  %+v", result.SkuId))
		// int_sku_id, _ := strconv.Atoi(result.SkuId)
		if result.SkuId > 0 {
			model.DB.Table("tb_jd_original_items").Where("sku_id = ?", result.SkuId).Updates(map[string]interface{}{
				"update_time":         item.UpdateTime,
				"sku_name":            item.SkuName,
				"pc_price":            item.PcPrice,
				"pc_commission":       item.PcCommission,
				"pc_commission_share": item.PcCommissionShare,
				"wl_commission":       item.WlCommission,
				"wl_commission_share": item.WlCommissionShare,
				"wl_price":            item.WlPrice,
				"image_url":           item.ImageUrl,
			})
			// logger.Info(fmt.Sprintf("item update  %+v", item))
			continue
		}
		model.DB.Table("tb_jd_original_items").Create(&item)
		// logger.Info(fmt.Sprintf("item create  %+v", item))
	}
	return urls.String()
}

//生成请求参数和签名
func (J *JdItemsdk) SetSignJointUrlParam(paramjson string) {
	J.RequestParam.App_key = conf.Jd.APPKEY
	J.RequestParam.Format = "json"
	J.RequestParam.V = "1.0"
	J.RequestParam.Method = conf.Jd.METHODITEMD
	J.RequestParam.Access_Token = conf.Jd.ACCESS_TOKEN

	J.RequestParam.Sign_method = "md5"
	J.RequestParam.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	J.RequestParam.Param_json = paramjson
	// J.RequestParam.PageParam = param.PageParam
	// J.RequestParam.QueryParam = param.QueryParam

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
	builder.WriteString(conf.Jd.APPSECRET)
	for _, person := range SortSlice {
		if person.Value == "" {
			continue
		}
		builder.WriteString(strings.ToLower(person.Key) + person.Value)
		u.Add(strings.ToLower(person.Key), person.Value)
	}
	builder.WriteString(conf.Jd.APPSECRET)
	//生成签名
	u.Add("sign", strings.ToUpper(util.Md5(builder.String())))
	//拼接参数
	J.SignAndUri = u.Encode()
}
