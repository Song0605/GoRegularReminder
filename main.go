package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	//Modify:2022年8月26日14:05:48 ver 0.1
	fmt.Println("钉钉喝水提醒小助手 ver_0.1")
	weekDay := time.Now().Weekday()
	if weekDay == 6 || weekDay == 7 {
		fmt.Println("周末就别加班了啊靓仔！~")
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
			if (nowHour < 13 && nowMinute > 0) || (nowHour >= 13 && nowMinute > 30) {
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
	// if nowNum == maxNum && time.Now().Minute() > 30 {
	// 	return
	// }
	_time := timeArr[nowNum]
	// fmt.Println("_time =", _time)
	c := cron.New(cron.WithSeconds())
	var (
		// frequency string,
		minute string
		hour   string
		// back string,
		// space string
	)
	// frequency = "*/60"
	if _time < 12 {
		minute = "0"
	} else {
		minute = "30"
	}
	hour = strconv.Itoa(_time)
	// back = "* * 1-5"
	// space = " "
	// spec := frequency + space + minute + space + hour + space + back
	// spec := "*/1 53 14 * * 1-5" //每周1-5 14:53开始 每隔1秒运行一次
	fmt.Println("Next:", hour, minute)
	c.AddFunc("00 * * * * ?", func() {
		now := time.Now()
		// fmt.Println(now.Hour(), ":", now.Minute())
		fmt.Printf(". ")
		if strconv.Itoa(now.Minute()) == minute && strconv.Itoa(now.Hour()) == hour {
			PostDingDing(nowNum + 1)
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

func PostDingDing(Num int) {
	fmt.Println()
	text := "是时候了！饮下今天的第" + strconv.Itoa(Num) + "杯水吧！"
	fmt.Println(text)
	body := "{\"at\":{\"isAtAll\":true},\"text\":{\"content\":\"" + text + "\"},\"msgtype\":\"text\"}"
	res, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=c74ed8c3a605f20838dd6278daa43f733a6fa077166c10bb05e80d173ec32613", "application/json;charset=utf-8", bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error!\n", err.Error())
	}
	fmt.Println(string(content))
}
