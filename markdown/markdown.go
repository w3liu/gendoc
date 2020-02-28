package markdown

import (
	"encoding/json"
	"fmt"
	"github.com/w3liu/gendoc"
	"os"
	"strings"
)

type SubTable struct {
	Title  string
	Fields []gendoc.Field
}

type Markdown struct {
	doc         *gendoc.Document
	id          int
	subReqList  []SubTable
	subRespList []SubTable
}

func New(doc *gendoc.Document) *Markdown {
	return &Markdown{doc: doc, subReqList: make([]SubTable, 0), subRespList: make([]SubTable, 0)}
}

func (m *Markdown) Generate(file string) {
	page := m.RenderPage(m.doc)
	_ = os.Remove(file)
	err := createFile(file, []byte(page))
	if err != nil {
		panic(err)
	}
}

func (m *Markdown) RenderPage(v *gendoc.Document) string {
	ts := TplPage
	ts = strings.Replace(ts, "{title}", v.Title, 1)
	ts = strings.Replace(ts, "{version}", v.Version, 1)
	ts = strings.Replace(ts, "{baseUrl}", v.BaseUrl, 1)
	body := ""
	list := v.GetList()
	if len(list) > 0 {
		for id, item := range list {
			tpl := m.RenderBody(id+1, item)
			body = fmt.Sprintf("%s%s", body, tpl)
		}
	}
	ts = strings.Replace(ts, "{body}", body, 1)
	return ts
}

func (m *Markdown) RenderBody(id int, v gendoc.DocItem) string {
	m.subReqList = make([]SubTable, 0)
	m.subRespList = make([]SubTable, 0)
	m.id = id
	ts := TplBody
	ts = strings.Replace(ts, "{id}", fmt.Sprintf("%v", id), 1)
	ts = strings.Replace(ts, "{name}", v.Title, 1)
	ts = strings.Replace(ts, "{author}", v.Author, 1)
	ts = strings.Replace(ts, "{method}", string(v.Method), 1)
	ts = strings.Replace(ts, "{url}", string(v.Url), 1)
	if len(v.ReqFields) > 0 {
		reqTable := m.RenderReqTable("", v.ReqFields)
		if len(m.subReqList) > 0 {
			subTable := ""
			for _, item := range m.subReqList {
				tpl := m.RenderReqTable(item.Title, item.Fields)
				subTable = fmt.Sprintf("%s%s", subTable, tpl)
			}
			reqTable = fmt.Sprintf("%s%s", reqTable, subTable)
		}
		ts = strings.Replace(ts, "{reqTable}", reqTable, 1)
	} else {
		ts = strings.Replace(ts, "{reqTable}", "", 1)
	}
	if len(v.RespFields) > 0 {
		respTable := m.RenderRespTable("", v.RespFields)
		if len(m.subRespList) > 0 {
			subTable := ""
			for _, item := range m.subRespList {
				tpl := m.RenderRespTable(item.Title, item.Fields)
				subTable = fmt.Sprintf("%s%s", subTable, tpl)
			}
			respTable = fmt.Sprintf("%s%s", respTable, subTable)
		}
		ts = strings.Replace(ts, "{respTable}", respTable, 1)
	} else {
		ts = strings.Replace(ts, "{respTable}", "", 1)
	}
	if v.RespParam != nil {
		respParam, _ := json.MarshalIndent(v.RespParam, "", "\t")
		ts = strings.Replace(ts, "{respParam}", fmt.Sprintf("```json\n %s \n```", string(respParam)), 1)
	} else {
		ts = strings.Replace(ts, "{respParam}", "", 1)
	}
	return ts

}

func (m *Markdown) RenderReqTable(title string, fields []gendoc.Field) string {
	ts := ""
	if title != "" {
		ts = fmt.Sprintf("\n<a id=\"%d.%s\"></a> \n##### %s \n %s ", m.id, title, title, TplReqTable)
	} else {
		ts = TplReqTable
	}

	params := ""
	for _, v := range fields {
		tpl := m.RenderReqParam(v)
		params = fmt.Sprintf("%s%s", params, tpl)
	}
	ts = strings.Replace(ts, "{params}", params, 1)
	return ts
}

func (m *Markdown) RenderReqParam(v gendoc.Field) string {
	ts := TplReqParam
	required := "是"
	if !v.Required {
		required = "否"
	}
	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{kind}", v.Kind, 1)
	ts = strings.Replace(ts, "{required}", required, 1)
	if len(v.List) > 0 {
		subTable := SubTable{
			Title:  v.Name,
			Fields: v.List,
		}
		m.subReqList = append(m.subReqList, subTable)
		v.Description = fmt.Sprintf("%s [点我](#%d.%s)", v.Description, m.id, v.Name)
	}
	ts = strings.Replace(ts, "{description}", v.Description, 1)
	return ts
}

func (m *Markdown) RenderRespTable(title string, fields []gendoc.Field) string {
	ts := ""
	if title != "" {
		ts = fmt.Sprintf("\n<a id=\"%d.%s\"></a> \n##### %s \n %s ", m.id, title, title, TplRespTable)
	} else {
		ts = TplRespTable
	}
	params := ""
	for _, v := range fields {
		tpl := m.RenderRespParam(v)
		params = fmt.Sprintf("%s%s", params, tpl)
	}
	ts = strings.Replace(ts, "{params}", params, 1)
	return ts
}

func (m *Markdown) RenderRespParam(v gendoc.Field) string {
	ts := TplRespParam
	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{kind}", v.Kind, 1)
	if len(v.List) > 0 {
		subTable := SubTable{
			Title:  v.Name,
			Fields: v.List,
		}
		m.subRespList = append(m.subRespList, subTable)
		v.Description = fmt.Sprintf("%s [点我](#%d.%s)", v.Description, m.id, v.Name)
	}
	ts = strings.Replace(ts, "{description}", v.Description, 1)
	return ts
}

func createFile(fileName string, val []byte) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(val)
	if err != nil {
		return err
	}
	return err
}
