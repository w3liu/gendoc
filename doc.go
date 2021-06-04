package gendoc

import (
	"fmt"
	"reflect"
	"strings"
)

type MethodType string
type UrlType string
type RequestType interface{}
type ResponseType interface{}

const (
	POST   MethodType = "POST"
	GET    MethodType = "GET"
	PUT    MethodType = "PUT"
	DELETE MethodType = "DELETE"
)

const (
	Tomas string = "Tomas"
)

type Document struct {
	Title   string    `json:"title"`   // 文档标题
	Version string    `json:"version"` // 版本号
	BaseUrl string    `json:"baseUrl"` // BaseUrl
	List    []DocItem `json:"list"`    // 文档列表
}

// 字段
type Field struct {
	Name        string  `json:"name"`        // 字段名称
	Kind        string  `json:"kind"`        // 字段类型
	Required    bool    `json:"required"`    // 是否必填
	Description string  `json:"description"` // 字段说明
	List        []Field `json:"list"`        // 字段列表
}

// 文档对象
type DocItem struct {
	Title      string       `json:"title"`     // 标题
	Url        UrlType      `json:"url"`       // 接口地址
	Method     MethodType   `json:"method"`    // 请求类型
	ReqParam   RequestType  `json:"reqParam"`  // 请求参数
	RespParam  ResponseType `json:"respParam"` // 返回参数
	Author     string       `json:"author"`    // 作者
	ReqFields  []Field      `json:"fields"`    // 字段列表
	RespFields []Field      `json:"fields"`    // 字段列表
}

/*
添加文档对象
title 接口名称
url 相对路径
method http请求方法
author 作者
req 请求结构体（需要实例化）
resp 返回结果结构体（需要实例化）
*/
func (d *Document) AddItem(title string, url UrlType, method MethodType, author string, req, resp interface{}) {
	v := DocItem{
		Title:     title,
		Url:       url,
		Method:    method,
		ReqParam:  req,
		RespParam: resp,
		Author:    author,
	}
	if len(d.List) == 0 {
		d.List = make([]DocItem, 0)
		d.List = append(d.List, v)
	} else {
		d.List = append(d.List, v)
	}
}

// 生成接口列表
func (d *Document) GenerateFields() {
	if len(d.List) > 0 {
		for i, docItem := range d.List {
			docItem.ReqFields = createFields(docItem.ReqParam)
			docItem.RespFields = createFields(docItem.RespParam)
			d.List[i] = docItem
		}
	}
}

// 获取接口列表
func (d *Document) GetList() []DocItem {
	return d.List
}

// 创建字段
func createFields(param interface{}) []Field {
	if param == nil {
		return nil
	}
	fields := make([]Field, 0)
	val := reflect.ValueOf(param)
	if !val.IsValid() {
		panic("not valid")
	}
	fmt.Println(val.Kind())
	if val.Kind() == reflect.Slice {
		if val.Len() > 0 {
			val = val.Index(0)
		} else {
			return nil
		}
	}
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	typ := val.Type()
	cnt := val.NumField()
	for i := 0; i < cnt; i++ {
		fd := val.Field(i)
		kd := fd.Kind()
		ty := typ.Field(i)
		field := Field{
			Name:        capitalize(ty.Name),
			Kind:        kd.String(),
			Required:    getRequired(ty),
			Description: getDescription(ty),
			List:        nil,
		}
		if field.Kind == "interface" || field.Kind == "struct" {
			subFields := createFields(fd.Interface())
			field.List = subFields
		}
		if field.Kind == "slice" {
			subFields := createFields(fd.Slice(0, 1).Interface())
			field.List = subFields
		}
		// 如果是数字型字符串 例 Id int `json:"id,string"`
		if field.Kind == "int" && strings.Contains(ty.Tag.Get("json"), ",string") {
			field.Kind = "string"
		}
		//如果是内嵌结构体
		if ty.Anonymous {
			subFields := createFields(fd.Interface())
			fields = append(fields, subFields...)
		} else {
			fields = append(fields, field)
		}
	}
	return fields
}

// 获取tag
func getDescription(field reflect.StructField) string {
	tag := field.Tag.Get("doc")
	desc := strings.Trim(tag, " ")
	desc = strings.TrimRight(desc, "required")
	return desc
}

// 判断是否必填
func getRequired(field reflect.StructField) bool {
	tag := field.Tag.Get("doc")
	return strings.Contains(tag, "required")
}

// 首字母小写
func capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
