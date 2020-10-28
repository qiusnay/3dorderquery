package service

import (
	"github.com/qiusnay/3dorderquery/model"
)

//定义接口
type UnionSDKAPI interface {
	GetOrders(start string, end string) interface{}
	SetSignJointUrlParam(param string)
}

type Apiconfig struct {
	APPKEY       string
	APPSECRET    string
	METHOD       string
	HOST         string
	METHODITEMD  string
	ACCESS_TOKEN string
}

type JdOrderResult struct {
	Code      int    `json:"code"`
	HasMore   bool   `json:"hasMore"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
	Data      []model.JdOriginalOrder
}
