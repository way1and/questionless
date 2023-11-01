# QuestionLESS v1.0.0
基于Selenium的问卷识别/填写服务

大学牲做课题有问卷要求,让人家填烦了,索性花两天时间搞了这个

### 功能
1. 识别问卷星问卷
2. 自定义填写内容, 批量填写

### 使用方法
1. 下载对应版本的ChromeDriver
2. 命名为chromedriver.exe, 替换项目下chromedriver
3. go run main.go

### 运行参数
|  参数   |  默认值 | 描述 |
|  ----  | ----  |  ---- |
| -ip | localhost | gin服务器运行IP地址 |
| -port | 80 | gin服务器运行port |
| -dports | 9000-9010...-9090 | ChromeDriver服务端口 |

例: go run main.go -port=80 -dport=81
  
### 演示
![image](./face.png)


### 额外内容
1. 无需配置 输入问卷星 问卷网址一键扫描问卷
2. 可以为文本型答案自定义内容
3. 自定义每个选项选取概率
4. 支持多个任务同时进行, 仅需要打开多个网页
5. 支持分布式, 可以部署在服务器
   
浏览器打开 http://localhost/index 即可使用

### 未来
1. 等我闲了会更新其他问卷类型
2. 增加对选项附带文本的支持

我的微信 wayland_0916, 欢迎交流
