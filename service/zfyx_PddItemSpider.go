package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/util"
)

type PddItemsAllParam struct {
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

type GoodsBasicDetailResponse struct {
	GoodsBasicDetailResponse struct {
		List      []model.PddItemOriginal `json:"list"`
		RequestId string                  `json:"request_id"`
		SearchId  string                  `json:"search_id"`
		ListId    string                  `json:"list_id"`
		Total     int                     `json:"total"`
	} `json:"goods_basic_detail_response"`
}

type PddItemReq struct {
	Type        string `json:"type"`
	Limit       int64  `json:"limit"`
	ChannelType int    `json:"channel_type"`
}

type PddItemsdk struct {
	RequestParam PddItemsAllParam
	SignAndUri   string
}

func (J *PddItemsdk) GetParams(brand int) PddItemReq {
	ParamStruct := PddItemReq{}
	ParamStruct.Limit = 400
	ParamStruct.ChannelType = brand
	ParamStruct.Type = PddConf.Pdd.METHODITEMD
	return ParamStruct
}

//获取拼多多商品
func (J *PddItemsdk) FetchPddItems(brand int) interface{} {
	util.Config().Bind("conf", "thirdpartysdk", &PddConf)

	Param := J.GetParams(brand)
	SignAndUri := SetSignJointUrlParam(Param)
	var urls strings.Builder
	urls.WriteString(PddConf.Pdd.HOST)
	urls.WriteString(SignAndUri)
	body, _ := util.HttpGet(urls.String())
	response := &GoodsBasicDetailResponse{}
	e := json.Unmarshal([]byte(body), &response)
	if e != nil {
		panic(e)
	}
	// logger.Info(fmt.Sprintf("get jd item %+v", response.JdKplOpenUnionSearchByelitedResponse.Data))
	for _, item := range response.GoodsBasicDetailResponse.List {
		result := model.PddItemOriginal{}
		item.UpdateTime = strconv.FormatInt(time.Now().Unix(), 10)
		item.CreateTime = strconv.FormatInt(time.Now().Unix(), 10)
		item.EliteId = brand
		model.DB.Table("tb_pdd_original_items").Where("goods_id = ?", item.GoodsId).First(&result)
		// logger.Info(fmt.Sprintf("dddddd  %+v", result.SkuId))
		// int_sku_id, _ := strconv.Atoi(result.SkuId)
		if result.GoodsId > 0 {
			model.DB.Table("tb_pdd_original_items").Where("goods_id = ?", result.GoodsId).Updates(map[string]interface{}{
				"update_time":      item.UpdateTime,
				"goods_name":       item.GoodsName,
				"min_group_price":  item.MinGroupPrice,
				"min_normal_price": item.MinNormalPrice,
				"promotion_rate":   item.PromotionRate,
				"goods_image_url":  item.GoodsImageUrl,
				"category_id":      item.CategoryId,
			})
			// logger.Info(fmt.Sprintf("item update  %+v", item))
			continue
		}
		model.DB.Table("tb_pdd_original_items").Create(&item)
		// logger.Info(fmt.Sprintf("item create  %+v", item))
	}
	return urls.String()
}
