package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Unknwon/goconfig"
)

const (
	GoroutineStartLog = "\n----------------------Start.----------------------\n"
	GoroutineEndLog   = "---------------------- End. ----------------------\n\n"
)

func main() {
	var count int32 = 0
	fmt.Println("======================================================")
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic("配置文件读取错误")
	}
	fmt.Println("配置文件读取成功")
	fmt.Println("======================================================")

	files := cfg.GetSectionList()
	lenOfFiles := len(files)

	var wg sync.WaitGroup
	wg.Add(lenOfFiles)
	for i := 0; i < lenOfFiles; i++ {
		copymode, err := cfg.GetValue(files[i], "copymode")
		if err != nil {
			copymode = COPY
			fmt.Println("index[", i, "]\t", err)
		}
		filename, err := cfg.GetValue(files[i], "filename")
		if err != nil {
			filename = "NO_FILE"
			fmt.Println("index[", i, "]\t", err)
		}
		src, err := cfg.GetValue(files[i], "src")
		if err != nil {
			fmt.Println("index[", i, "]\t", err)
		}
		des, err := cfg.GetValue(files[i], "des")
		if err != nil {
			fmt.Println("index[", i, "]\t", err)
		}
		log, err := cfg.GetValue(files[i], "log")
		if err != nil {
			fmt.Println("index[", i, "]\t", err)
		}

		go func(copymode, filename, src, des, log string, idx int) {
			defer wg.Done()
			var msg string
			idxStr := strconv.Itoa(idx)
			msg += GoroutineStartLog
			msg += "$ index [ " + idxStr + " ]. $\n"
			msg += "*** configs: ***\n"
			msg += "copymode:\t" + copymode + "\n"
			msg += "file:\t" + filename + "\n"
			msg += "src:\t" + src + "\n"
			msg += "des:\t" + des + "\n"
			msg += "log:\t" + log + "\n\n"

			if filename == "NO_FILE" {
				msg += "filename == NO_FILE\n"
				msg += GoroutineEndLog
				fmt.Println(msg)
				return
			}

			msg += "commands]:\n"
			var cpmsg string
			var copyOk bool
			if copymode == COPY {
				copyOk, cpmsg, _ = copyFile(filename, src, des)
			} else if copymode == COVER {
				copyOk, cpmsg, _ = copyFileIfDesExist(filename, src, des)
			}
			msg += cpmsg + "\n"
			if copyOk == false {
				msg += "#### copy failed! ####\n"
				msg += GoroutineEndLog
				fmt.Println(msg)
				return
			}
			count++

			msg += "#### copy success! ####\n"
			msg += GoroutineEndLog
			fmt.Println(msg)
		}(copymode, filename, src, des, log, i)
	}

	wg.Wait()
	fmt.Println("***************************************************************")
	fmt.Println("----\t\tcopy count:", count, "\t\t---")
	fmt.Println("***************************************************************")
}
