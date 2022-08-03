package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/tebeka/selenium"
)

const (
	//设置常量 分别设置chromedriver.exe的地址和本地调用端口
	seleniumPath = `C:\Projects\chromedriver_win32\chromedriver.exe`
	cookiePath   = `./cookies/cookies.json`
	port         = 9515
)

var (
	action = flag.String("a", "search", "selenium action")
)

func SearchPage(wd selenium.WebDriver, sheet *xlsx.Sheet) {
	time.Sleep(time.Second)
	//找到商品wes
	wes, err := wd.FindElements(selenium.ByXPATH, "//*[@id=\"orderAppContainer\"]/div/div[2]/div[3]/div[3]/div/div/div/div[2]/div/div")
	if err != nil {
		fmt.Println("FindElements filed: ", err)
	}
	for _, we := range wes {
		row := sheet.AddRow()
		cellOrderID := row.AddCell()
		cellOrderTime := row.AddCell()
		cellPhone := row.AddCell()
		cellGoodsMan := row.AddCell()

		time.Sleep(time.Duration(5+rand.Intn(10)*100) * time.Millisecond)

		// 显示按钮
		showButtonWe, err := we.FindElement(selenium.ByXPATH, "./div[2]/div[3]/div[2]/div/div/div/div/span[2]/a")
		if err == nil {
			if err := showButtonWe.Click(); err == nil {
				// 尝试读取号码
				for i := 0; i < 20; i++ {
					phoneNumWe, err := we.FindElement(selenium.ByXPATH, "./div[2]/div[3]/div[2]/div/div/div/div/span[1]")
					if err == nil {
						if txt, err := phoneNumWe.Text(); err == nil && !strings.Contains(txt, "***") {
							cellPhone.Value = txt
							break
						} else {
							time.Sleep(100 * time.Millisecond)
						}
					}
				}
			}
		} else {
			phoneNumWe, err := we.FindElement(selenium.ByXPATH, "./div[2]/div[3]/div[2]/div/div/div/div/span[1]")
			if err == nil {
				if txt, err := phoneNumWe.Text(); err == nil {
					cellPhone.Value = txt
				}
			}
		}

		// 订单号
		orderIdWe, err := we.FindElement(selenium.ByXPATH, "./div[1]/div/div[1]/span[1]/div/div")
		if err == nil {
			if txt, err := orderIdWe.Text(); err == nil {
				words := strings.Split(txt, " ")
				if len(words) == 2 {
					cellOrderID.Value = words[1]
				}
			}
		}

		// 下单时间
		orderTimeWe, err := we.FindElement(selenium.ByXPATH, "./div[1]/div/div[1]/span[2]")
		if err == nil {
			if txt, err := orderTimeWe.Text(); err == nil {
				words := strings.Split(txt, " ")
				if len(words) == 3 {
					cellOrderTime.Value = words[1] + "" + words[2]
				}
			}
		}

		// 带货达人
		goodsManWe, err := we.FindElement(selenium.ByXPATH, "./div[2]/div[1]/div/div[1]/div/div/div[2]/div[3]/div/div")
		if err == nil {
			if txt, err := goodsManWe.Text(); err == nil {
				if strings.Contains(txt, "带货达人") {
					cellGoodsMan.Value = strings.Replace(txt, "带货达人：", "", -1)
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	//1.开启selenium服务
	//设置selium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//延迟关闭服务
	defer service.Stop()

	//2.调用浏览器
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		panic(err)
	}
	//延迟退出chrome
	defer wd.Quit()

	if *action == "login" {
		if err := wd.Get("https://fxg.jinritemai.com/ffa/morder/order/list?btm_ppre=a2427.b6931.c0.d0&btm_pre=a2427.b08003.c0.d0&btm_show_id=af77892e-11a1-4f86-a3b7-94437675b125"); err != nil {
			panic(err)
		}

		time.Sleep(160 * time.Second)
		cs, err := wd.GetCookies()
		if err != nil {
			panic(err)
		}
		bytes, err := json.Marshal(cs)
		if err != nil {
			fmt.Println("Marshal cookies failed: ", err)
		}
		err = ioutil.WriteFile(cookiePath, bytes, 0644)
		if err != nil {
			fmt.Println("WriteFile cookies failed: ", err)
		}
	} else {
		bytes, err := ioutil.ReadFile(cookiePath)
		if err != nil {
			fmt.Println("ReadFile cookies failed: ", err)
		}
		var cookies []selenium.Cookie
		err = json.Unmarshal(bytes, &cookies)
		if err != nil {
			fmt.Println("Unmarshal cookies failed: ", err)
		}
		if err := wd.Get("https://fxg.jinritemai.com/ffa/morder/order/list?btm_ppre=a2427.b6931.c0.d0&btm_pre=a2427.b08003.c0.d0&btm_show_id=af77892e-11a1-4f86-a3b7-94437675b125"); err != nil {
			panic(err)
		}
		for _, cookie := range cookies {
			err = wd.AddCookie(&cookie)
			if err != nil {
				fmt.Println("Add cookie failed: ", err)
			}
		}
		if err := wd.Get("https://fxg.jinritemai.com/ffa/morder/order/list?btm_ppre=a2427.b6931.c0.d0&btm_pre=a2427.b08003.c0.d0&btm_show_id=af77892e-11a1-4f86-a3b7-94437675b125"); err != nil {
			panic(err)
		}

		file := xlsx.NewFile()
		sheet, err := file.AddSheet("信息收集")
		if err != nil {
			panic(err.Error())
		}
		defer file.Save("data.xlsx")

		row := sheet.AddRow()
		tabs := []string{"id", "下单时间", "联系方式", "带货达人"}
		for _, tab := range tabs {
			cell := row.AddCell()
			cell.Value = tab
		}

		var cmd string
		fmt.Scan(&cmd)
		for cmd != "export" {
			for {
				// 扫描当前页
				SearchPage(wd, sheet)
				// 下一页
				lis, err := wd.FindElements(selenium.ByXPATH, "//*[@id=\"orderAppContainer\"]/div/div[2]/div[3]/div[3]/div/div/div/ul/li")
				if err == nil {
					nextButtonWe, err := lis[len(lis)-2].FindElement(selenium.ByXPATH, "./button")
					if err == nil {
						enabled, _ := nextButtonWe.IsEnabled()
						if enabled && nextButtonWe.Click() == nil {
							time.Sleep(1 * time.Second)
							continue
						}
					}
				}
				break
			}
			fmt.Scan(&cmd)
		}

		fmt.Println("Exporting...")
	}
}
