package pkg

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcelFromMemory(data []byte, processor func(lineCount int64, lineText []string) bool) (int64, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("open excel failed: %w", err)
	}
	defer f.Close()

	rows, err := f.Rows(f.GetSheetList()[0]) // 流式读取
	if err != nil {
		return 0, err
	}

	var lineCount int64
	for rows.Next() {
		row, _ := rows.Columns()
		if !processor(lineCount, row) {
			break
		}
		lineCount++
	}
	return lineCount, nil
}
