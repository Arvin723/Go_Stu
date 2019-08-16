package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	UTF8     = "UTF-8"     //chcp 65001
	GB18030  = "GB18030"   //chcp 936
	SHIFTJIS = "Shift JIS" //chcp 932
)

func cmdRunner_CHCP(name string, arg ...string) (string, error) {
	out, err := exec.Command("chcp").Output()
	if err != nil {
		log.Fatal(err)
	}
	outString := string(out)
	//outString, _ := cmdRunner_GB18030("chcp")

	if strings.Contains(outString, "936") {
		outString = "936"
	} else if strings.Contains(outString, "932") {
		outString = "932"
	} else if strings.Contains(outString, "65001") {
		outString = "65001"
	} else {
		outString = "65001"
	}

	switch outString {
	case "936":
		return cmdRunner(GB18030, name, arg...)
	case "932":
		return cmdRunner(SHIFTJIS, name, arg...)
	case "65001":
		return cmdRunner(UTF8, name, arg...)
	default:
		return cmdRunner(UTF8, name, arg...)
	}
}

func cmdRunner_GB18030(name string, arg ...string) (string, error) {
	return cmdRunner(GB18030, name, arg...)
}

func cmdRunner_SHIFTJIS(name string, arg ...string) (string, error) {
	return cmdRunner(SHIFTJIS, name, arg...)
}

func cmdRunner(encode, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd run failed!")
		return "outPipe create failed", err
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("cmd run failed!")
		return "cmd.Start failed", err
	}

	ioScanner := bufio.NewScanner(outPipe)
	var msg string
	for ioScanner.Scan() {
		msg += ConvertByte2String(ioScanner.Bytes(), encode) + "\n"
	}
	if err := cmd.Wait(); err != nil {
		msg += "cmd.Wait error!"
		return msg, err
	}

	return msg, nil
}

func ConvertByte2String(byte []byte, charset string) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case SHIFTJIS:
		var decodeBytes, _ = japanese.ShiftJIS.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
