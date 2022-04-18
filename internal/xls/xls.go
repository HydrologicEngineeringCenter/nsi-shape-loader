package xls

import (
	"github.com/usace/xlscellreader"
	"github.com/xuri/excelize/v2"
)

func NewXls(src string) (*xlscellreader.CellReader, error) {
	f, err := excelize.OpenFile(src)
	if err != nil {
		return &xlscellreader.CellReader{}, err
	}
	defer func() {
		// Close the spreadsheet.
		err = f.Close()
	}()
	reader := xlscellreader.CellReader{F: f}
	if err != nil {
		return &xlscellreader.CellReader{}, err
	}
	return &reader, err
}
