package driver

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/tebeka/selenium/log"
)

var caps selenium.Capabilities
var PORTUsingMap = make(map[int]bool) // false 空闲

// init 配置 chrome driver 参数
func init() {
	caps = selenium.Capabilities{
		"browserName": "chrome",
	}
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Args: []string{
			"--headless",    //设置Chrome无头模式
			"--log-level=3", // 日志等级
			"--silent",      // 无输出
		},
		ExcludeSwitches: []string{
			"enable-logging", // 关闭 dev listen on
		},
	}

	caps.AddChrome(chromeCaps) // 设置日志等级
	mode := log.Severe
	caps.AddLogging(log.Capabilities{
		log.Server:      mode,
		log.Driver:      mode,
		log.Browser:     mode,
		log.Client:      mode,
		log.Profiler:    mode,
		log.Performance: mode,
	})
}

func UsePORT(port int) {
	PORTUsingMap[port] = true
}

// GetPORT 获得空闲端口
func GetPORT() int {
	for port, using := range PORTUsingMap {
		if !using {
			UsePORT(port)
			return port
		}
	}
	return 0
}

// CheckPort 检查端口占用 true 可用 false 不可用
func CheckPort(port int) bool {

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(port)))
	if err != nil {
		return false
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)
	return true
}

// FreePORT 释放端口
func FreePORT(port int) {
	PORTUsingMap[port] = false
}

// AddPORT 添加新端口
func AddPORT(port int) {
	PORTUsingMap[port] = false
}

// Driver 获得 driver 对象
func Driver(url string) (*selenium.Service, *selenium.WebDriver, int) {

	var opts []selenium.ServiceOption
	var page *selenium.Service
	var err error
	var PORT int
	for {
		PORT = GetPORT()
		page, err = selenium.NewChromeDriverService("./chromedriver.exe", PORT, opts...)
		if err != nil || PORT == 0 {
			fmt.Println("网络问题正在重启")
			duration, _ := time.ParseDuration("5s")
			time.Sleep(duration)
			continue
		}
		// 无问题
		break
	}

	//
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", PORT))
	if err != nil {
		fmt.Println(err)
		fmt.Println("测试")
		return nil, nil, 0
	}
	//

	if err := wd.Get(url); err != nil {
		fmt.Println(err)
		return nil, nil, 0
	}

	script := `
    Object.defineProperty(navigator, 'webdriver', {
        get: () => undefined
    });
	`
	_, err2 := wd.ExecuteScript(script, nil)
	if err2 != nil {
		fmt.Println("Failed to execute JavaScript:", err)
		return nil, nil, 0
	}

	for {
		d, _ := time.ParseDuration("2s")
		time.Sleep(d)
		ele, err := wd.FindElement(selenium.ByXPATH, "//*[@id='ctlNext']")
		if err != nil || ele == nil {
			// 问卷已关闭
			ele, err := wd.FindElement(selenium.ByID, "divWorkError")
			if err != nil || ele != nil {
				return nil, nil, 0
			}

			continue
		}
		break
	}
	return page, &wd, PORT
}
