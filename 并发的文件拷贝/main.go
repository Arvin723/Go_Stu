package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Unknwon/goconfig"
)

const (
	GoroutineStartLog = "\n----------------------Start.----------------------\n"
	GoroutineEndLog   = "---------------------- End. ----------------------\n\n"
)

type param struct {
	copymode string
	filename string
	src      string
	des      string
	log      string
	idx      int
}

//copy原则：
//0.不能有相同的des+filenale

func main() {
	fmt.Println("======================================================")
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		panic("配置文件读取错误")
	}
	fmt.Println("配置文件读取成功")
	fmt.Println("======================================================")

	var count int32 = 0
	var mutex sync.Mutex
	var wg sync.WaitGroup
	files := cfg.GetSectionList()
	lenOfFiles := len(files)
	//copyAttrs := make([]param, 0, lenOfFiles+8)
	var copyAttrs []param
	wg.Add(lenOfFiles)
	var currentWorkPath string
	getCurrentWorkPath()
	if currentWorkPath, err = getCurrentWorkPath(); err != nil {
		fmt.Println("getCurrentWorkPath() returned error!")
		panic(err)
	}
	fmt.Println("___ CurrentWorkPath ___: ", currentWorkPath)

	for i := 0; i < lenOfFiles; i++ {
		copyAttr, err := getParam(i, cfg, files)
		if err != nil {
			fmt.Println(err)
			wg.Done()
			continue
		}

		go func(copyAttr param) {
			defer wg.Done()
			var msg string
			idxStr := strconv.Itoa(copyAttr.idx)

			msg += GoroutineStartLog
			msg += "$ index [ " + idxStr + " ]. $\n"
			msg += "*** configs: ***\n"
			msg += "copymode:\t" + copyAttr.copymode + "\n"
			msg += "file:\t" + copyAttr.filename + "\n"
			msg += "src:\t" + copyAttr.src + "\n"
			msg += "des:\t" + copyAttr.des + "\n"
			msg += "log:\t" + copyAttr.log + "\n\n"

			if copyAttr.filename == "NO_FILE" {
				msg += "filename == NO_FILE\n"
				msg += GoroutineEndLog
				fmt.Println(msg)
				return
			}

			var desPathIsUnique bool
			mutex.Lock()
			copyAttrs, desPathIsUnique = uniqueFilePathCheck(copyAttrs, copyAttr)
			mutex.Unlock()
			if desPathIsUnique == false {
				msg += copyAttr.filename + ": desPathIsUnique == false"
				fmt.Println(msg)
				return
			}

			msg += "commands]:\n"
			var cpmsg string
			var copyOk bool
			if copyAttr.copymode == COPY {
				copyOk, cpmsg, _ = copyFile(copyAttr.filename, copyAttr.src, copyAttr.des)
			} else if copyAttr.copymode == COVER {
				copyOk, cpmsg, _ = copyFileIfDesExist(copyAttr.filename, copyAttr.src, copyAttr.des)
			}
			msg += cpmsg + "\n"
			if copyOk == false {
				msg += "#### copy failed! ####\n"
				msg += GoroutineEndLog
				fmt.Println(msg)
				return
			}
			atomic.AddInt32(&count, 1)

			msg += "#### copy success! ####\n"
			msg += GoroutineEndLog
			fmt.Println(msg)
		}(copyAttr)
	}

	wg.Wait()
	printCopyedFile(copyAttrs)
	fmt.Println("\n**************************************************************************")
	fmt.Println("----\t\tcopy count:", count, "\t\t---")
	fmt.Println("**************************************************************************")
	return
}

func getParam(index int, cfg *goconfig.ConfigFile, files []string) (param, error) {
	copymode, err := cfg.GetValue(files[index], "copymode")
	if err != nil {
		copymode = COPY
		fmt.Println("index[", index, "]\t", err)
	}
	filename, err := cfg.GetValue(files[index], "filename")
	if err != nil {
		filename = "NO_FILE"
		fmt.Println("index[", index, "]\t", err)
	}
	src, err := cfg.GetValue(files[index], "src")
	if err != nil {
		fmt.Println("index[", index, "]\t", err)
		return param{}, err
	}
	des, err := cfg.GetValue(files[index], "des")
	if err != nil {
		fmt.Println("index[", index, "]\t", err)
		return param{}, err
	}
	log, err := cfg.GetValue(files[index], "log")
	if err != nil {
		fmt.Println("index[", index, "]\t", err)
	}

	return param{
		copymode: copymode,
		filename: filename,
		src:      src,
		des:      des,
		log:      log,
		idx:      index,
	}, nil

}

func uniqueFilePathCheck(copyAttrs []param, copyAttr param) ([]param, bool) {
	lenOfAttrs := len(copyAttrs)
	desPath, err := filepath.Abs(copyAttr.des)
	if err != nil {
		panic(err)
	}
	copyAttr.des = desPath
	for i := 0; i < lenOfAttrs; i++ {
		if copyAttrs[i].des == copyAttr.des &&
			copyAttrs[i].filename == copyAttr.filename {
			return copyAttrs, false
		}
	}
	copyAttrs = append(copyAttrs, copyAttr)
	return copyAttrs, true
}

func getCurrentWorkPath() (string, error) {
	currentWorkPath := ".\\"
	currentWorkPath, err := filepath.Abs(currentWorkPath)
	if err != nil {
		return "", err
	}
	return currentWorkPath, nil
}

func printCopyedFile(copyAttrs []param) {
	fmt.Println("####################################################################")
	lenOfAttrs := len(copyAttrs)
	for i := 0; i < lenOfAttrs; i++ {
		fmt.Println("\n~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("index:[", copyAttrs[i].idx, "].")
		fmt.Println("\tcopymode:\t", copyAttrs[i].copymode)
		fmt.Println("\tfilename:\t", copyAttrs[i].filename)
		fmt.Println("\t  src  :\t", copyAttrs[i].src)
		fmt.Println("\t  des  :\t", copyAttrs[i].des)
		fmt.Println("\t  log  :\t", copyAttrs[i].log)
		fmt.Println("\n~~~~~~~~~~~~~~~~~~~~~")
	}
	fmt.Println("####################################################################")
}
