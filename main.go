package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	PostQYWX("测试消息。")
	//RegularRemind()
}

func RegularRemind() {
	fmt.Println("钉钉喝水提醒小助手 ver_0.3")
	fmt.Println("最后更新时间:2022年8月30日09:08:44")
	weekDay := time.Now().Weekday()
	if weekDay == 6 || weekDay == 7 {
		fmt.Println("周末就别加班了啊靓仔！")
		select {}
	}
	var timeArr = [...]int{9, 10, 11, 13, 14, 15, 16, 17}
	var maxNum = len(timeArr) - 1
	nowNum := -1
	fmt.Println("Now:", time.Now().Hour(), time.Now().Minute())
	nowHour := time.Now().Hour()
	if time.Now().Minute() == 59 {
		nowHour++
	}
	for num, v := range timeArr {
		if nowHour <= v {
			nowNum = num
			nowMinute := time.Now().Minute()
			if nowHour == v && ((nowHour < 13 && nowMinute > 0) || (nowHour >= 13 && nowMinute > 30)) {
				nowNum++
			}
			break
		}
	}
	if nowNum == -1 {
		return
	}
	if nowNum > maxNum {
		return
	}
	_time := timeArr[nowNum]
	c := cron.New(cron.WithSeconds())
	var (
		minute string
		hour   string
	)
	if _time < 12 {
		minute = "0"
	} else {
		minute = "30"
	}
	hour = strconv.Itoa(_time)
	// spec := "*/1 53 14 * * 1-5" //每周1-5 14:53开始 每隔1秒运行一次
	fmt.Println("Next:", hour, minute)
	c.AddFunc("00 * * * * ?", func() {
		now := time.Now()
		fmt.Println("现在时间：", now.Hour(), ":", now.Minute())
		text := ""
		if strconv.Itoa(now.Minute()) == minute &&
			strconv.Itoa(now.Hour()) == hour {
			//text := "测试消息。"
			//atAll := "true"
			//PostDingDing(text, atAll)
			text = "喝水提醒小助手提醒您，该喝今天的第" + strconv.Itoa(nowNum+1) + "杯水了。"
			PostQYWX(text)
			fmt.Println("已发送")
			if nowNum == maxNum {
				c.Stop()
			}
			nowNum++
			_time = timeArr[nowNum]
			hour = strconv.Itoa(_time)
			if _time < 12 {
				minute = "0"
			} else {
				minute = "30"
			}
			fmt.Println("Next:", hour, minute)
		}
	})
	c.Start()
	select {}
}

// JsonType
var JsonType = "application/json;charset=utf-8"

// 钉钉Url
var DingDingUrl = "https://oapi.dingtalk.com/robot/send?access_token=c74ed8c3a605f20838dd6278daa43f733a6fa077166c10bb05e80d173ec32613"

// 向钉钉发送消息
func PostDingDing(text string, atAll string) {
	fmt.Println()
	fmt.Println(text)
	body := "{\"at\":{\"isAtAll\":" + atAll + "},\"text\":{\"content\":\"" + text + "\"},\"msgtype\":\"text\"}"
	res, err := http.Post(
		DingDingUrl,
		JsonType,
		bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	fmt.Println("接口返回值：", string(content))
}

// 企业微信Url
var QYWXTestUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3a5066fa-8218-4f2a-b0e7-e3d0eb4351bc"

// 向企业微信发送消息
func PostQYWX(text string) {
	fmt.Println()
	fmt.Println(text)
	body := "{\"text\":{\"content\":\"" + text + "\"},\"msgtype\":\"text\"}"
	res, err := http.Post(
		QYWXTestUrl,
		JsonType,
		bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	fmt.Println("接口返回值：", string(content))
}
