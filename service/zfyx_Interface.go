package service

//定义接口
type UnionSDKAPI interface {
	GetOrders(start string, end string) interface{}
	SetSignJointUrlParam(param string)
}

type Apiconfig struct {
	APPKEY    string
	APPSECRET string
	METHOD    string
	HOST      string
}
