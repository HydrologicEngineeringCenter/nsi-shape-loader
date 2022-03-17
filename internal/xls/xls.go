package xls

import "github.com/xuri/excelize/v2"

func NewXls(src string) (*excelize.File, error) {
	file, err := excelize.OpenFile(src)
	if err != nil {
		return &excelize.File{}, err
	}
	defer func() {
		// Close the spreadsheet.
		err = file.Close()
	}()
	return file, err
}
