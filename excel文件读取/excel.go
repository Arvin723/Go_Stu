package main

//go get github.com/360EntSecGroup-Skylar/excelize

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	COL_VAL_OF_FILE = 0
	COL_VAL_OF_SRC  = 1
	COL_VAL_OF_DES  = 2
	COL_VAL_OF_LOG  = 3
)

const (
	LOG_DEFAULT = "cmt"
)

type data struct {
	file string
	src  string
	des  string
	log  string
}

func main() {
	f, err := excelize.OpenFile("config.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := f.GetRows("Sheet1")
	var tableDatas []data

	for _, row := range rows {
		var rowdata data

		//空单元格不会被读取到，默认值要提前赋值
		//但数据到结构体时，至少是有空字符串的
		//rows会读取第一到最后一行，中间即使没有数据也会空出来
		rowdata.log = LOG_DEFAULT

		for cv, colCell := range row {
			switch cv {
			case COL_VAL_OF_FILE:
				rowdata.file = colCell
			case COL_VAL_OF_SRC:
				rowdata.src = colCell
			case COL_VAL_OF_DES:
				rowdata.des = colCell
			case COL_VAL_OF_LOG:
				rowdata.log = colCell
			}
		}
		tableDatas = append(tableDatas, rowdata)
	}

	printTableDatas(tableDatas)
}

func printTableDatas(tableDatas []data) {
	lenOfDatas := len(tableDatas)
	for i := 0; i < lenOfDatas; i++ {
		fmt.Println(tableDatas[i].file + "\t" + tableDatas[i].src + "\t" + tableDatas[i].des + "\t" + tableDatas[i].log)
	}
}
