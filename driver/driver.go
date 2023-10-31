package driver

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

var PORT = 0

func Driver(url string) (*selenium.Service, *selenium.WebDriver) {

	var opts []selenium.ServiceOption
	var page *selenium.Service
	var err error
	for {
		page, err = selenium.NewChromeDriverService("./chromedriver.exe", PORT, opts...)
		if err != nil {
			fmt.Println("网络问题正在重启")
			duration, _ := time.ParseDuration("5s")
			time.Sleep(duration)
			continue
		}
		// 无问题
		break
	}
	// 浏览器设置
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args: []string{
			"--headless", //设置Chrome无头模式
		},
	}
	caps.AddChrome(chromeCaps)
	//
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", PORT))
	if err != nil {
		fmt.Println(err)
		fmt.Println("测试")
		return nil, nil
	}
	//

	if err := wd.Get(url); err != nil {
		fmt.Println(err)
		return nil, nil
	}

	for {
		d, _ := time.ParseDuration("2s")
		time.Sleep(d)
		ele, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[1]/form/div[13]/div[10]/div[3]/div/div/div")
		if err != nil || ele == nil {
			continue
		}
		break
	}

	return page, &wd
}
