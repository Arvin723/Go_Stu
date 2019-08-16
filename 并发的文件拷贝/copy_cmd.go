package main

const (
	COPY  = "COPY"
	COVER = "FUGAI"
	NEW   = "NEW"
)

func copyFile_MODE(mode, filename, src, des string) (bool, string, error) {
	switch mode {
	case COPY:
		return copyFile_COPY(filename, src, des)
	case COVER:
		return copyFile_COVER(filename, src, des)
	case NEW:
		return copyFile_NEW(filename, src, des)
	default:
		return copyFile(filename, src, des)
	}
}

//copy命令的调用
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

//任何情况(除非源处没有文件) -> 拷贝
func copyFile_COPY(filename, src, des string) (bool, string, error) {
	return copyFile(filename, src, des)
}

//目标处没有文件 -> 不拷贝
func copyFile_COVER(filename, src, des string) (bool, string, error) {
	despath := des + "\\" + filename
	if fileOk, err := isFileExist(despath); fileOk == false {
		return false, "copyFileIfDesExist: des File Not Exit!", err
	}
	return copyFile(filename, src, des)
}

//目标处已有文件 -> 不拷贝
func copyFile_NEW(filename, src, des string) (bool, string, error) {
	despath := des + "\\" + filename
	if fileOk, err := isFileExist(despath); fileOk == true {
		return false, "copyFileIfDesExist: des File Not Exit!", err
	}
	return copyFile(filename, src, des)
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
