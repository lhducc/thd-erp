package utils

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportField struct {
	Header    string
	FieldName string
	GetValue  func(item interface{}) any
}

type ExcelExporter struct {
	SheetName string
	Fields    map[string]ExportField
}

func NewExcelExporter(sheetName string) *ExcelExporter {
	return &ExcelExporter{
		SheetName: sheetName,
		Fields:    make(map[string]ExportField),
	}
}

func (e *ExcelExporter) RegisterField(key string, header string, fieldName string, valueFunc func(interface{}) any) {
	e.Fields[key] = ExportField{
		Header:    header,
		FieldName: fieldName,
		GetValue:  valueFunc,
	}
}

func (e *ExcelExporter) Export(items interface{}, selectedFields []string) ([]byte, string, error) {
	itemsValue := reflect.ValueOf(items)
	if itemsValue.Kind() != reflect.Slice {
		return nil, "", fmt.Errorf("dữ liệu đầu vào phải là một slice")
	}

	var fields []ExportField
	for _, name := range selectedFields {
		if f, ok := e.Fields[name]; ok {
			fields = append(fields, f)
		}
	}

	if len(fields) == 0 {
		return nil, "", fmt.Errorf("không có trường nào hợp lệ để export")
	}

	f := excelize.NewFile()
	index, err := f.NewSheet(e.SheetName)
	if err != nil {
		return nil, "", fmt.Errorf("không thể tạo sheet mới: %w", err)
	}
	f.SetActiveSheet(index)

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("không thể tạo style: %w", err)
	}

	for i, field := range fields {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(e.SheetName, cell, field.Header)
		f.SetCellStyle(e.SheetName, cell, cell, headerStyle)
		f.SetColWidth(e.SheetName, col, col, 20)
	}

	for i := 0; i < itemsValue.Len(); i++ {
		item := itemsValue.Index(i).Interface()
		row := i + 2
		for j, field := range fields {
			col := string(rune('A' + j))
			cell := fmt.Sprintf("%s%d", col, row)
			f.SetCellValue(e.SheetName, cell, field.GetValue(item))
		}
	}

	f.DeleteSheet("Sheet1")

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("không thể ghi file Excel: %w", err)
	}

	filename := e.SheetName + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	return buf.Bytes(), filename, nil
}

type ExcelReader struct {
	file *excelize.File
}

func NewExcelReader(reader io.Reader) (*ExcelReader, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	return &ExcelReader{file: file}, nil
}

// GetSheetName gets the name of a sheet by its index.
// It's exported.
func (er *ExcelReader) GetSheetName(idx int) string {
	return er.file.GetSheetName(idx)
}

// GetRows gets all rows from a sheet.
// It's exported.
func (er *ExcelReader) GetRows(sheetName string) ([][]string, error) {
	return er.file.GetRows(sheetName)
}
