# 接口文档
> 版本号：v1.0.0<br>
> BaseUrl: http://127.0.0.1:8080


### 1. 创建订单接口
> 作者：Tomas

#### 请求说明
> 请求方式：POST<br>
请求URL ：[/v1/order/create](#)

#### 请求参数

|字段      |字段类型       |必填     |字段说明    |
|---------|--------------|--------|-----------|
|pass   |string        |是      |交易密码不能为空，请用BASE64 进行转码     |
|amount   |float32        |否      |支付金额，不能小于或等于0    |
|randomNum   |string        |否      |随机字符串不能为空,最大长度为30    |
|tranBody   |string        |否      |交易描述不能为空,最大长度为30；    |
|outTradeNo   |string        |否      |三方交易唯一订单号，最大长度60    |
|createIp   |string        |否      |IP地址    |
|startTime   |int64        |否      |交易开始时间搓,格式为yyyyMMddHHmmss    |


#### 返回结果
```json
 {
	"code": 0,
	"msg": "",
	"data": {
		"pass": "",
		"amount": 0,
		"randomNum": "",
		"tranBody": "",
		"outTradeNo": "",
		"createIp": "",
		"startTime": 0
	}
} 
```
#### 返回参数

|字段      |字段类型       |字段说明    |
|---------|--------------|-----------|
|code   |int32        |错误码    |
|msg   |string        |错误信息    |
|data   |interface        |业务数据 [点我](#1.data)    |

<a id="1.data"></a> 
##### data 
 
|字段      |字段类型       |字段说明    |
|---------|--------------|-----------|
|pass   |string        |交易密码不能为空，请用BASE64 进行转码     |
|amount   |float32        |支付金额，不能小于或等于0    |
|randomNum   |string        |随机字符串不能为空,最大长度为30    |
|tranBody   |string        |交易描述不能为空,最大长度为30；    |
|outTradeNo   |string        |三方交易唯一订单号，最大长度60    |
|createIp   |string        |IP地址    |
|startTime   |int64        |交易开始时间搓,格式为yyyyMMddHHmmss    |
 


