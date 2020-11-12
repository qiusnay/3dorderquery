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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/qiusnay/3dorderquery/service"
	"github.com/qiusnay/3dorderquery/util"
)

func main() {
	Log := util.NewLogger("api")
	mux := http.NewServeMux()
	mux.Handle("/", &myAPIHandler{})
	mux.HandleFunc("/api/getpddurl", getUrlPinduoduo)
	mux.HandleFunc("/api/getjdurl", getUrlJd)
	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 30, //设置30秒的写超时
		Handler:      mux,
	}
	server.ListenAndServe() // 开启监听
	defer func() {
		if err := recover(); err != nil {
			Log.Info(err)
		}
	}()
}

type GoodsPromotionUrlGenerateResponse struct {
	GoodsPromotionUrlGenerateResponse struct {
		Code      int64       `json:"code"`
		Message   string      `json:"message"`
		Data      interface{} `json:"data"`
		RequestId string      `json:"request_id"`
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
	Log := util.NewLogger("api")
	defer func() {
		if err := recover(); err != nil {
			Log.Error(err)
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
	newbody := strings.Replace(string(body), "goods_promotion_url_list", "data", 1)
	response := &GoodsPromotionUrlGenerateResponse{}
	e := json.Unmarshal([]byte(newbody), &response)
	if e != nil {
		panic(e)
	}
	Log.Info("pdd response : %+v", response.GoodsPromotionUrlGenerateResponse.Data)
	// fmt.Println(fmt.Sprintf("get form : %+v", response.GoodsPromotionUrlGenerateResponse.GoodsPromotionUrlList))
	// 睡眠4秒  上面配置了3秒写超时，所以访问 “/bye“路由会出现没有响应的现象
	time.Sleep(100 * time.Microsecond)
	response.GoodsPromotionUrlGenerateResponse.Code = 200
	response.GoodsPromotionUrlGenerateResponse.Message = "success"
	recommandUrl, _ := json.Marshal(response.GoodsPromotionUrlGenerateResponse)
	w.Write([]byte(recommandUrl))
}

type configjd struct {
	Jd service.Apiconfig `toml:"jd"`
}

var ApiJdConf configjd

func (h *myAPIHandler) GetJdParams(r *http.Request, Conf configjd) string {
	ParamStruct := service.JdUrlReq{}
	ParamStruct.PromotionCodeReq.MaterialId = r.URL.Query().Get("materialId")
	ParamStruct.PromotionCodeReq.SiteId = Conf.Jd.SITEID
	bytes, _ := json.Marshal(ParamStruct)
	return string(bytes)
}

type JdUnionOpenPromotionCommonGetResponse struct {
	JdUnionOpenPromotionCommonGetResponse struct {
		Result string `json:"result"`
		Code   string `json:"code"`
	} `json:"jd_union_open_promotion_common_get_response"`
}

func getUrlJd(w http.ResponseWriter, r *http.Request) {
	Log := util.NewLogger("api")
	defer func() {
		if err := recover(); err != nil {
			Log.Error(err)
		}
	}()
	h := &myAPIHandler{}
	util.Config().Bind("conf", "thirdpartysdk", &ApiJdConf)
	Param := h.GetJdParams(r, ApiJdConf)
	JdUrl := SetJdSignJointUrlParam(Param)
	var urls strings.Builder
	urls.WriteString(ApiJdConf.Jd.HOST)
	urls.WriteString(JdUrl)
	// fmt.Println(urls.String())
	body, _ := util.HttpGet(urls.String())
	response := &JdUnionOpenPromotionCommonGetResponse{}
	e := json.Unmarshal([]byte(body), &response)
	if e != nil {
		panic(e)
	}
	Log.Info(fmt.Sprintf("jd url response %+v", response.JdUnionOpenPromotionCommonGetResponse.Result))
	// fmt.Println(fmt.Sprintf("get form : %+v", response.GoodsPromotionUrlGenerateResponse.GoodsPromotionUrlList))
	// 睡眠4秒  上面配置了3秒写超时，所以访问 “/bye“路由会出现没有响应的现象
	time.Sleep(100 * time.Microsecond)
	// recommandUrl, _ := json.Marshal()
	w.Write([]byte(response.JdUnionOpenPromotionCommonGetResponse.Result))

}

//生成请求参数和签名
func SetJdSignJointUrlParam(paramjson string) string {
	J := service.JdSysItemUrlParam{}
	J.App_key = ApiJdConf.Jd.APPKEY
	J.Format = "json"
	J.V = "1.0"
	J.Method = ApiJdConf.Jd.METHOD_GERNERATE_PROMOTION

	J.Sign_method = "md5"
	J.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	J.Param_json = paramjson

	values := reflect.ValueOf(J)
	keys := reflect.TypeOf(J)
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
	builder.WriteString(ApiJdConf.Jd.APPSECRET)
	for _, person := range SortSlice {
		if person.Value == "" {
			continue
		}
		builder.WriteString(strings.ToLower(person.Key) + person.Value)
		u.Add(strings.ToLower(person.Key), person.Value)
	}
	builder.WriteString(ApiJdConf.Jd.APPSECRET)

	//生成签名
	u.Add("sign", strings.ToUpper(util.Md5(builder.String())))
	//拼接参数
	return u.Encode()
}
