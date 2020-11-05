package util

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log map[string]*logrus.Logger

func NewLogger(Logname string) *logrus.Logger {
	cLog := Log[Logname]
	if cLog != nil {
		return cLog
	}
	pathMap := lfshook.PathMap{
		logrus.InfoLevel: "./log/" + Logname + "." + time.Now().Format("2006-01-02") + ".log",
	}
	cLog = logrus.New()
	cLog.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return cLog
}

//生成 MD5
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func HttpGet(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	//设置请求头部信息
	//req.Header.Add("Authorization", q.Token)
	//发送请求
	response, _ := http.DefaultClient.Do(req)
	//关闭请求
	defer Close(response)
	//解析返回结果
	bytes, err := ioutil.ReadAll(response.Body)
	return bytes, err
}

func Close(response *http.Response) {
	e := response.Body.Close()
	if e != nil {
		panic(e)
	}
}

type Onestruct struct {
	Key   string
	Value string
}
type Items []Onestruct

// Len()方法和Swap()方法不用变化
// 获取此 slice 的长度
func (p Items) Len() int { return len(p) }

// 交换数据
func (p Items) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p Items) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}
