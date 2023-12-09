package conf

import (
	"fmt"
	"os"
)

func SetAPI(ip string, port int) {
	path := "./static/api.js"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777)

	if err != nil {
		fmt.Println("修改 api.js 错误 无法创建或打开文件")
		return
	}

	_, err = f.WriteString(fmt.Sprintf("var API = 'http://%s:%d'", ip, port))
	if err != nil {
		fmt.Println("修改 api.js 错误 无法写入文字")
		return
	}

}
