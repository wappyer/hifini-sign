package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const UrlIndex = "https://www.hifiti.com/"
const UrlSignIn = "https://www.hifiti.com/user-login.htm"

func main() {
	var token = ""
	flag.StringVar(&token, "token", "", "登录后拿cookie中「bbs_token」的值")
	flag.Parse()
	if token == "" {
		token = os.Getenv("token")
		if token == "" {
			log.Printf("token is empty")
			return
		}
	}

	sign := getSign(token)
	SignIn(token, sign)
}

// 获取sign(sign为签到接口重要传参)
func getSign(token string) (sign string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", UrlIndex, nil)
	req.Header.Set("Cookie", fmt.Sprintf("bbs_token=%s", token))
	indexResp, err := client.Do(req)
	if err != nil {
		log.Printf("[getSign]请求首页失败，err: %v\n", err)
		return
	}
	defer indexResp.Body.Close()

	indexHtmlByte, err := io.ReadAll(indexResp.Body)
	if err != nil {
		log.Printf("[getSign]读取首页内容失败，err: %v\n", err)
		return
	}

	reg, err := regexp.Compile(`var sign = "(.*)"`)
	if err != nil {
		log.Printf("[getSign]正则匹配表达式错误，err: %v\n", err)
		return
	}
	match := reg.FindStringSubmatch(string(indexHtmlByte))
	if len(match) < 2 {
		log.Printf("[getSign]正则匹配未获取到sign，匹配结果:%v\n", match)
		return
	}
	sign = match[1]
	return
}

func SignIn(token, sign string) {
	client := http.Client{}
	body := strings.NewReader(fmt.Sprintf("sign=%s", sign))
	req, err := http.NewRequest("POST", UrlSignIn, body)
	req.Header.Set("Cookie", fmt.Sprintf("bbs_token=%s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[SignIn]签到请求失败，err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	log.Println("[SignIn]签到请求成功")
	respBodyByte, _ := io.ReadAll(resp.Body)
	log.Printf("[SignIn]签到请求返回内容：%s\n", respBodyByte)
	return
}
