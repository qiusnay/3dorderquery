package service

import (
	"bytes"
	"log"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/qiusnay/3dorderquery/model"
	"github.com/qiusnay/3dorderquery/util"
)

//定义接口
type UnionSDKAPI interface {
	GetOrders(start string, end string) interface{}
	// SetSignJointUrlParam(param string)
}

type Apiconfig struct {
	APPKEY                     string
	APPSECRET                  string
	METHOD                     string
	HOST                       string
	METHODITEMD                string
	METHOD_GERNERATE_PROMOTION string
	ACCESS_TOKEN               string
}

type Configpdd struct {
	Pdd Apiconfig `toml:"pdd"`
}

var PddConf Configpdd

type JdOrderResult struct {
	Code      int    `json:"code"`
	HasMore   bool   `json:"hasMore"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	Data      []model.JdOriginalOrder
}

type PddUrlReq struct {
	CustomParameters     map[string]string `json:"custom_parameters"`      //自定义参数，为链接打上自定义标签；自定义参数最长限制64个字节；格式为： {"uid":"11111","sid":"22222"} ，其中 uid 用户唯一标识，可自行加密后传入，每个用户仅且对应一个标识，必填； sid 上下文信息标识，例如sessionId等，非必填。该json字符串中也可以加入其他自定义的key
	GenerateSchemaUrl    bool              `json:"generate_schema_url"`    //是否返回 schema URL
	GenerateShortUrl     bool              `json:"generate_short_url"`     //是否生成短链接，true-是，false-否
	GenerateWeappWebview bool              `json:"generate_weapp_webview"` //是否生成唤起微信客户端链接，true-是，false-否，默认false
	GenerateWeApp        bool              `json:"generate_we_app"`        //是否生成小程序推广
	GoodsIdList          string            `json:"goods_id_list"`          //商品ID，仅支持单个查询
	PId                  string            `json:"p_id"`                   //推广位ID
	Type                 string            `json:"type"`
}

//生成请求参数和签名
func SetSignJointUrlParam(param interface{}) string {
	defer func() {
		if err := recover(); err != nil {
			// logger.Error(err)
		}
	}()
	pddParams := make(map[string]interface{})
	pddParams["data_type"] = "json"
	pddParams["client_id"] = PddConf.Pdd.APPKEY
	pddParams["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)

	if pddStruct, isin := param.(PddUrlReq); isin {
		t := reflect.TypeOf(pddStruct)
		v := reflect.ValueOf(pddStruct)
		for k := 0; k < t.NumField(); k++ {
			// fmt.Println(fmt.Sprintf("key show  : %+v", t.Field(k).Name))
			pddParams[t.Field(k).Name] = v.Field(k).Interface()
		}
	}
	if pddStruct, isin := param.(PddOrderReq); isin {
		t := reflect.TypeOf(pddStruct)
		v := reflect.ValueOf(pddStruct)
		for k := 0; k < t.NumField(); k++ {
			// fmt.Println(fmt.Sprintf("key show  : %+v", t.Field(k).Name))
			pddParams[t.Field(k).Name] = v.Field(k).Interface()
		}
	}
	if pddStruct, isin := param.(PddItemReq); isin {
		t := reflect.TypeOf(pddStruct)
		v := reflect.ValueOf(pddStruct)
		for k := 0; k < t.NumField(); k++ {
			// fmt.Println(fmt.Sprintf("key show  : %+v", t.Field(k).Name))
			pddParams[t.Field(k).Name] = v.Field(k).Interface()
		}
	}
	// fmt.Println(fmt.Sprintf("pddparams show  : %+v", pddParams))
	// values := reflect.ValueOf(pddParams)
	// keys := reflect.TypeOf(pddParams)
	// count := values.NumField()
	SortSlice := util.Items{}
	for key, item := range pddParams {
		switch value := item.(type) {
		case string:
			SortSlice = append(SortSlice, util.Onestruct{Camel2Case(key), value})
		case int:
			SortSlice = append(SortSlice, util.Onestruct{Camel2Case(key), strconv.Itoa(value)})
		case bool:
			SortSlice = append(SortSlice, util.Onestruct{Camel2Case(key), strconv.FormatBool(value)})
		}
	}
	sort.Sort(SortSlice)
	var builder strings.Builder
	u := url.Values{}
	builder.WriteString(PddConf.Pdd.APPSECRET)
	for _, person := range SortSlice {
		if person.Value == "" {
			continue
		}
		builder.WriteString(strings.ToLower(person.Key) + person.Value)
		u.Add(strings.ToLower(person.Key), person.Value)
	}
	builder.WriteString(PddConf.Pdd.APPSECRET)

	//生成签名
	u.Add("sign", strings.ToUpper(util.Md5(builder.String())))
	//拼接参数
	return u.Encode()
}

func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}
