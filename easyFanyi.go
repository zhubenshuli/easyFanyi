package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var url string = "http://fanyi.youdao.com/openapi.do?keyfrom=zhubenshuli&key=1532828825&type=data&doctype=json&version=1.1&q="
var word WordFanyi

type WordFanyi struct {
	ErrorCode int
	Basic     interface{}
}

func main() {
	startTime := time.Now().Unix()

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "使用方法：easy_fanyi 需要翻译的txt文件")
		os.Exit(-1)
	}
	fmt.Fprintln(os.Stderr, "正在翻译...")

	// 读取文件
	inFilepath := os.Args[1]
	inFile, err := os.Open(inFilepath)
	defer inFile.Close()
	checkErr(err)

	// 翻译文件保存的字符串
	fanyiStr := ""

	fileContent, err := ioutil.ReadAll(inFile)
	wordArr := strings.Fields(string(fileContent))
	for _, v := range wordArr {
		searchUrl := url + v

		// 调用有道翻译的API
		resp, err := http.Get(searchUrl)
		defer resp.Body.Close()
		checkErr(err)
		body, err := ioutil.ReadAll(resp.Body)
		checkErr(err)

		// 对返回的json字符串做处理
		err = json.Unmarshal(body, &word)
		checkErr(err)

		if word.ErrorCode == 0 {
			wordBasic := word.Basic.(map[string]interface{})
			var phonetic string
			var explains []interface{}
			if wordBasic["phonetic"] != nil {
				phonetic = wordBasic["phonetic"].(string)
			}
			if wordBasic["explains"] != nil {
				explains = wordBasic["explains"].([]interface{})
			}

			fanyiStr += v + " [" + phonetic + "] "
			for _, val := range explains {
				explainStr := val.(string)
				fanyiStr += explainStr
			}
			fanyiStr += "\r\n"
		}
	}

	// 输出文件
	outFile, err := os.Create("newFanyi.txt")
	defer outFile.Close()
	checkErr(err)
	_, err = outFile.WriteString(fanyiStr)
	checkErr(err)

	endTime := time.Now().Unix()
	fmt.Fprintf(os.Stderr, "翻译耗时：%d秒，翻译文件为newFanyi.txt", endTime-startTime)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
