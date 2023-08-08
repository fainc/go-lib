package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type mapExport struct{}

var mapExportVar = mapExport{}

func MapExport() *mapExport {
	return &mapExportVar
}

type TitleValue struct {
	TitleText string `dc:"标题文字"`
	Column    int    `dc:"指定竖向合并行数(横向合并列根据数据列自动计算)"`
	Align     string `dc:"文字水平对齐方式"`
}
type KeysValue struct {
	Name  string `dc:"标题名称"`
	Index string `dc:"数据索引"`
}
type SheetsValue struct {
	SheetName string                   `dc:"表名称，不设置默认顺序Sheet1,2,3"`
	Title     []TitleValue             `dc:"表内标题参数"`
	Keys      []KeysValue              `dc:"数据键"`
	Data      []map[string]interface{} `dc:"数据集"`
}
type FileValue struct {
	SavePath    string `dc:"文件保存路径"`
	ActiveSheet int    `dc:"设置默认表，默认0"`
}

// SaveFile  保存到本地
func (rec *mapExport) SaveFile(savePath string, activeSheet int, sheets []*SheetsValue) (err error) {
	f := excelize.NewFile()
	defer func() {
		if err = f.Close(); err != nil {
			return
		}
	}()
	err = rec.writer(f, sheets)
	if err != nil {
		return err
	}
	f.SetActiveSheet(activeSheet)
	if err = f.SaveAs(savePath); err != nil {
		return
	}
	return
}

// IOWrite  数据流写入
func (rec *mapExport) IOWrite(w io.Writer, activeSheet int, sheets []*SheetsValue) (err error) {
	f := excelize.NewFile()
	defer func() {
		if err = f.Close(); err != nil {
			return
		}
	}()
	err = rec.writer(f, sheets)
	if err != nil {
		return err
	}
	f.SetActiveSheet(activeSheet)
	if err = f.Write(w); err != nil {
		return
	}
	return
}

func (rec *mapExport) writer(f *excelize.File, sheets []*SheetsValue) (err error) {
	for i, sheet := range sheets {
		if i == 0 {
			if sheet.SheetName != "" { // 重命名Sheet1
				err = f.SetSheetName("Sheet1", sheet.SheetName)
				if err != nil {
					return err
				}
			} else {
				sheet.SheetName = "Sheet1"
			}
		}
		if i != 0 { // 新增表
			_, err = f.NewSheet(sheet.SheetName)
		}
	}
	for _, sheet := range sheets {
		writeRowID := 1 // 表内写入行起始位置标记
		// 标题写入
		for _, title := range sheet.Title {
			cell := fmt.Sprintf("A%v", writeRowID)
			err = f.SetCellValue(sheet.SheetName, cell, title.TitleText)
			if err != nil {
				return
			}
			// 合并标题单元格
			var endRowID int
			if title.Column <= 1 { // 默认占一行
				endRowID = writeRowID
			} else {
				endRowID = writeRowID + title.Column - 1
			}
			endCell := fmt.Sprintf("%v%v", getCell(len(sheet.Keys)-1), endRowID)
			_ = f.MergeCell(sheet.SheetName, cell, endCell)

			// 定义样式
			style, _ := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{Horizontal: title.Align, Vertical: "center", WrapText: true}})
			if err != nil {
				return
			}
			_ = f.SetCellStyle(sheet.SheetName, cell, cell, style)
			// 更新行号
			if title.Column <= 1 { // 默认占一行
				writeRowID++
			} else {
				writeRowID += title.Column
			}
		}
		// 数据键写入
		for i, key := range sheet.Keys {
			cell := fmt.Sprintf("%v%v", getCell(i), writeRowID)
			err = f.SetCellValue(sheet.SheetName, cell, key.Name)
			if err != nil {
				return
			}
		}
		writeRowID++
		// 数据行写入
		for _, data := range sheet.Data { // 行
			for c, key := range sheet.Keys { // 列
				cell := fmt.Sprintf("%v%v", getCell(c), writeRowID)
				err = f.SetCellValue(sheet.SheetName, cell, data[key.Index])
				if err != nil {
					return
				}
			}
			writeRowID++
		}
	}
	return
}
func getCell(i int) string {
	tables := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	return tables[i]
}
