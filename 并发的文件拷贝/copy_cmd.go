package main

import "os"

const (
	COPY  = "COPY"
	COVER = "FUGAI"
)

func copyFile(filename, src, des string) (bool, string, error) {
	src, des, err := checkLastByte(src, des)
	if err != nil {
		return false, "copyFile: checkLastByte failed!", err
	}

	filepath := src + "\\" + filename

	if fileOk, err := isFileExist(filepath); fileOk == false {
		return false, "copyFile: src File Not Exit!", err
	}

	//msg, err := cmdRunner_GB18030("xcopy", filepath, des, "/y")
	msg, err := cmdRunner_CHCP("xcopy", filepath, des, "/y")
	if err != nil {
		return false, msg + "\ncopyFile: cmdRunner error!", err
	}

	return true, "copyFile: \n" + msg + "copyFile end\n", nil
}

func copyFileIfDesExist(filename, src, des string) (bool, string, error) {
	despath := des + "\\" + filename
	if fileOk, err := isFileExist(despath); fileOk == false {
		return false, "copyFileIfDesExist: des File Not Exit!", err
	}
	return copyFile(filename, src, des)
}

func isFileExist(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkLastByte(src, des string) (string, string, error) {
	var srcOk bool = false
	var desOk bool = false

	rs := []rune(src)
	lenSrc := len(rs)
	lastByte := string(rs[lenSrc-1:])
	if lastByte == "\\" || lastByte == "/" {
		src = string(rs[:lenSrc-1])
	} else {
		srcOk = true
	}

	rs = []rune(des)
	lenSrc = len(rs)
	lastByte = string(rs[lenSrc-1:])
	if lastByte == "\\" || lastByte == "/" {
		des = string(rs[:lenSrc-1])
	} else {
		desOk = true
	}

	if srcOk && desOk {
		return src, des, nil
	} else {
		src, des, _ = checkLastByte(src, des)
	}

	return src, des, nil
}
