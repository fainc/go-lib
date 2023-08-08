package excel

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type testData struct {
	Id   int
	Name interface{}
}

func TestCreate_ToLocal(t *testing.T) {
	var sqldata []testData
	sqldata = append(sqldata, testData{Id: 1, Name: "1"})
	var Sheets []*SheetsValue
	data := gconv.Maps(sqldata)
	keys := []KeysValue{{Name: "年度", Index: "Year"}, {Name: "月份", Index: "Name"}, {Name: "收入金额", Index: "Income"}, {Name: "支出金额", Index: "Out"}, {Name: "结余", Index: "Val"}}
	Sheets = append(Sheets, &SheetsValue{SheetName: "测试", Title: []TitleValue{{TitleText: "2023年度结算计算表", Column: 3}, {TitleText: "制表人：FAInc.", Column: 1}, {TitleText: "制表时间：" + gtime.Now().String(), Column: 1}}, Keys: keys, Data: data})
	err := MapExport().SaveFile("test.xlsx", 0, Sheets)
	if err != nil {
		g.Dump(err)
		return
	}
}
