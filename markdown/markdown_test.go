package markdown

import (
	"encoding/json"
	"fmt"
	"github.com/w3liu/gendoc"
	"testing"
)

type ReqCreateOrder struct {
	Pass       string  `json:"pass" doc:"交易密码不能为空，请用BASE64 进行转码 required"`
	Amount     float32 `json:"amount" doc:"支付金额，不能小于或等于0"`
	RandomNum  string  `json:"randomNum" doc:"随机字符串不能为空,最大长度为30"`
	TranBody   string  `json:"tranBody" doc:"交易描述不能为空,最大长度为30；"`
	OutTradeNo string  `json:"outTradeNo" doc:"三方交易唯一订单号，最大长度60"`
	CreateIp   string  `json:"createIp" doc:"IP地址"`
	StartTime  int64   `json:"startTime" doc:"交易开始时间搓,格式为yyyyMMddHHmmss"`
}

type RespCreateOrder struct {
	Code int32       `json:"code" doc:"错误码"`
	Msg  string      `json:"msg" doc:"错误信息"`
	Data interface{} `json:"data" doc:"业务数据"`
}

func TestGenDoc(t *testing.T) {
	doc := &gendoc.Document{}
	doc.AddItem("创建订单接口", "/v1/order/create", gendoc.POST, gendoc.Tomas, &ReqCreateOrder{}, &RespCreateOrder{Data: &ReqCreateOrder{}})
	doc.GenerateFields()
	list := doc.GetList()
	b, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestGenMarkDown(t *testing.T) {
	// 实例化文档
	doc := &gendoc.Document{
		Title:   "接口文档",
		Version: "v1.0.0",
		BaseUrl: "http://127.0.0.1:8080",
	}
	// 添加接口
	doc.AddItem("创建订单接口", "/v1/order/create", gendoc.POST, gendoc.Tomas, &ReqCreateOrder{}, &RespCreateOrder{Data: &ReqCreateOrder{}})
	// 生成字段
	doc.GenerateFields()
	// 实例化文档生成器
	md := New(doc)
	// 生成文档
	md.Generate("./doc.md")
}
