package main

/**
请求地址:http://localhost:1210/api/getpddurl?itemid=11111&userid=123&pid=11272100_178340985
itemid : 商品ID
userid : 支付优选的用户ID
pid : 拼多多分配的推广位,用来标识当前商品是用哪一个推广位来推广的.这个对生成的推广链接有影响,目前有这样几个可以分配 :
<a val="11272100_178367015">支付优选公众号推文</a>
<a val="11272100_178341000">支付优选APP-开发测试3</a>
<a val="11272100_178340985">支付优选APP-开发测试2</a>
<a val="11272100_178340976">支付优选APP-开发测试1</a>
<a val="11272100_178340875">支付优选APP-线上优惠-拼多多</a>
<a val="11272100_178340361">支付优选APP-开屏1</a>
<a val="11272100_178340767">支付优选APP-开屏2</a>
<a val="11272100_178340731">支付优选APP-购物返利频道</a>
<a val="11272100_178340639">支付优选APP-首页轮播banner3</a>
<a val="11272100_178340633">支付优选APP-首页轮播banner2</a>
<a val="11272100_178340620">支付优选APP-首页轮播banner1</a>
<a val="11272100_176024774">zfyx</a>
<a val="11272100_148477208">A</a>

返回值说明 :
mobile_short_url : 唤醒拼多多app的推广短链接
mobile_url : 唤醒拼多多app的推广长链接
schema_url:schema的链接
short_url:推广短链接
url:推广长链接
we_app_info: []小程序信息
we_app_web_view_short_url:唤起微信app推广短链接
we_app_web_view_url:唤起微信app推广链接
*/

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/service"
	"github.com/qiusnay/3dorderquery/util"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &myAPIHandler{})
	mux.HandleFunc("/api/getpddurl", getUrlPinduoduo)
	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 30, //设置30秒的写超时
		Handler:      mux,
	}
	server.ListenAndServe() // 开启监听
}

type GoodsPromotionUrlGenerateResponse struct {
	GoodsPromotionUrlGenerateResponse struct {
		GoodsPromotionUrlList []interface{} `json:"goods_promotion_url_list"`
		RequestId             string        `json:"request_id"`
	} `json:"goods_promotion_url_generate_response"`
}

func (*myAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("支付优选接口服务器v1\nIP : 127.0.0.1:1210\n/ 首页相关介绍\n/ getbroadUrlForpdd 获取拼多多推广链接\n email : qiusnay@gmail.com"))
}

type myAPIHandler struct{}

func (h *myAPIHandler) GetParams(r *http.Request, Conf service.Configpdd) service.PddUrlReq {
	ParamStruct := service.PddUrlReq{
		Type:                 Conf.Pdd.METHOD_GERNERATE_PROMOTION,
		GoodsIdList:          "[" + r.URL.Query().Get("itemid") + "]",
		PId:                  r.URL.Query().Get("pid"),
		CustomParameters:     map[string]string{"uid": r.URL.Query().Get("userid")},
		GenerateSchemaUrl:    true,
		GenerateShortUrl:     true,
		GenerateWeappWebview: true,
		GenerateWeApp:        true,
	}
	return ParamStruct
}

func getUrlPinduoduo(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	h := &myAPIHandler{}
	util.Config().Bind("conf", "thirdpartysdk", &service.PddConf)
	Param := h.GetParams(r, service.PddConf)
	pddUrl := service.SetSignJointUrlParam(Param)
	var urls strings.Builder
	urls.WriteString(service.PddConf.Pdd.HOST)
	urls.WriteString(pddUrl)
	// fmt.Println(urls.String())
	body, _ := util.HttpGet(urls.String())
	response := &GoodsPromotionUrlGenerateResponse{}
	e := json.Unmarshal([]byte(body), &response)
	if e != nil {
		panic(e)
	}
	// fmt.Println(fmt.Sprintf("get form : %+v", response.GoodsPromotionUrlGenerateResponse.GoodsPromotionUrlList))
	// 睡眠4秒  上面配置了3秒写超时，所以访问 “/bye“路由会出现没有响应的现象
	time.Sleep(100 * time.Microsecond)
	recommandUrl, _ := json.Marshal(response.GoodsPromotionUrlGenerateResponse.GoodsPromotionUrlList)
	w.Write([]byte(recommandUrl))
}
