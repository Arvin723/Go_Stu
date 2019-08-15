package main

import (
	"bufio"
	"fmt"
	"os/exec"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	UTF8     = "UTF-8"
	GB18030  = "GB18030"   //chcp 936
	SHIFTJIS = "Shift JIS" //chcp 932
)

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

func cmdRunner_GB18030(name string, arg ...string) (string, error) {
	return cmdRunner(GB18030, name, arg...)
}

func cmdRunner_SHIFTJIS(name string, arg ...string) (string, error) {
	return cmdRunner(SHIFTJIS, name, arg...)
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
