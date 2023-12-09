package driver

import (
	"filler/models"
	"filler/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/tebeka/selenium"
)

type Page struct {
	Url string
}

type Operator struct {
	driver *selenium.WebDriver
}

func (o *Operator) QuestionGet(index int) models.Question {
	if o.QuestionExist(index) {
		question := models.Question{
			Type:     o.QuestionType(index),
			Title:    o.QuestionTitle(index),
			Options:  o.OptionGets(index),
			Required: o.QuestionRequire(index),
			Index:    index,
		}
		return question
	}

	return models.Question{Index: 0}
}

func (o *Operator) QuestionExist(index int) bool {
	driver := *o.driver
	_, err := driver.FindElement(selenium.ByCSSSelector, fmt.Sprintf("#div%d", index))
	if err != nil {
		return false
	}
	return true
}

func (o *Operator) OptionGets(index int) []models.Option {
	var options []models.Option
	driver := *o.driver
	eles, _ := driver.FindElements(selenium.ByCSSSelector, fmt.Sprintf("#div%d>.ui-controlgroup>div", index))
	var extended = false
	for optionIndex, ele := range eles {
		name, _ := ele.Text()
		desc := ""
		// 可能级联
		if _, err := driver.FindElement(selenium.ByCSSSelector, fmt.Sprintf("#tqq%d_%d", index, optionIndex+1)); err == nil {
			// 级联元素
			extended = true
			desc = "勾选后产生文本框"
		}

		// 将选项加入列表
		options = append(options, models.Option{
			Name:     name,
			Extended: extended,
			Desc:     desc,
			Index:    optionIndex + 1,
		})
	}
	return options
}

// QuestionInput 文字问题
func (o *Operator) QuestionInput(index int, text string) bool {
	driver := *o.driver
	ele, err := driver.FindElement(selenium.ByID, fmt.Sprintf("q%d", index+1))
	if err != nil || ele == nil {
		fmt.Println(err, ele)
		return false
	}

	err = ele.SendKeys(text)
	if err != nil {
		return false
	}
	return true
}

// QuestionSubmit 提交
func (o *Operator) QuestionSubmit() bool {
	driver := *o.driver
	ele, _ := driver.FindElement(selenium.ByXPATH, "//*[@id='ctlNext']")
	d1, _ := time.ParseDuration("1s")
	time.Sleep(d1)
	err := ele.Click()
	d2, _ := time.ParseDuration("0.5s")
	time.Sleep(d2)

	// 人机验证弹窗
	verifyPopup, _ := driver.FindElement(selenium.ByXPATH, "//*[@id='layui-layer1']/div[3]/a")
	if verifyPopup != nil {
		err := verifyPopup.Click()
		if err != nil {
			return false
		}
		d, _ := time.ParseDuration("2s")
		time.Sleep(d)
	}

	// 人机验证按钮
	verifyButton, _ := driver.FindElement(selenium.ByXPATH, "//*[@id='SM_BTN_1']")
	if verifyButton != nil {
		err := verifyButton.Click()
		if err != nil {
			return false
		}
		d, _ := time.ParseDuration("4s")
		time.Sleep(d)
	}

	// 人机验证滑动条
	verifySlider, _ := driver.FindElement(selenium.ByXPATH, "//*[@id='nc_1_n1z']")
	if verifySlider != nil {
		// err = slideSlider(driver, )
		err = slideSlider(driver, verifySlider)
		if err != nil {
			fmt.Println("Failed to slide slider:", err)
			return false
		}
	}
	if err != nil {
		return false
	}
	return true
}

func slideSlider(wd selenium.WebDriver, sliderElement selenium.WebElement) error {
	// 获取滑块的位置
	location, err := sliderElement.Location()
	if err != nil {
		return err
	}

	if wd.ButtonDown() != nil {
		fmt.Println(err)
	}

	// 计算滑动目标位置（这里示意性地滑动到右侧，具体根据实际情况调整）
	targetX := location.X + 100
	targetY := location.Y + 50

	// 模拟鼠标点击滑块并保持
	err = sliderElement.MoveTo(targetX, targetY)
	if err != nil {
		fmt.Println(err)
	}

	// 使用JavaScript执行滑动操作
	// script := fmt.Sprintf("arguments[0].style.left='%dpx';", targetX)
	// _, err = wd.ExecuteScript(script, []interface{}{sliderElement})
	// if err != nil {
	// 	return err
	// }

	// 等待一段时间以确保滑块移动完成（具体等待时间根据实际情况调整）
	time.Sleep(2 * time.Second)

	return nil
}

func (o *Operator) OptionInput(index, optionIndex int, text string) bool {
	driver := *o.driver
	ele, err := driver.FindElement(selenium.ByID, fmt.Sprintf("tqq%d_%d", index, optionIndex))
	if err != nil {
		return false
	}
	err = ele.SendKeys(text)
	if err != nil {
		return false
	}
	return true
}

func (o *Operator) QuestionSelect(idx int, selectIdx int) bool {
	driver := *o.driver
	elems, err := driver.FindElements(selenium.ByCSSSelector, fmt.Sprintf("#div%d>.ui-controlgroup>.ui-radio", idx+1))
	if err != nil || elems == nil || len(elems) < selectIdx+1 {
		fmt.Println(elems, err, len(elems), fmt.Sprintf("#div%d>.ui-controlgroup>.ui-radio", idx+1))
		fmt.Println("错误1")
		return false
	}
	err = elems[selectIdx].Click()
	if err != nil {
		fmt.Println("错误2")
		return false
	}
	return true
}

func (o *Operator) QuestionSelects(idx int, selectIdxArr []int) bool {
	driver := *o.driver
	elems, err := driver.FindElements(selenium.ByCSSSelector, fmt.Sprintf("#div%d>.ui-controlgroup>.ui-checkbox", idx+1))
	if err != nil || elems == nil || len(elems) < len(selectIdxArr) {
		return false
	}
	for _, idx := range selectIdxArr {
		if len(elems) < idx+1 {
			return false
		}
		err := elems[idx].Click()
		if err != nil {
			return false
		}
	}
	return true
}

func (o *Operator) QuestionType(index int) string {
	driver := *o.driver
	ele, _ := driver.FindElement(selenium.ByCSSSelector, fmt.Sprintf("#div%d", index))
	t, _ := ele.GetAttribute("type")

	// 问题类型
	var questionType string
	switch t {
	case "1":
		questionType = "文本"
	case "2":
		questionType = "文本"
	case "3":
		questionType = "单选"
	case "4":
		questionType = "多选"
	case "8":
		questionType = "文本"
	default:
		questionType = "未知"
	}
	return questionType
}

func (o *Operator) QuestionTitle(index int) string {
	driver := *o.driver
	ele, _ := driver.FindElement(selenium.ByCSSSelector, fmt.Sprintf("#div%d", index)+">.field-label")
	title, err := ele.Text()
	if err != nil {
		return ""
	}
	// 去符号
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "*", "")
	title = strings.ReplaceAll(title, "\n", "")
	// 去修饰符号 及其中内容
	end := strings.Index(title, "】")
	if end != -1 && end < len(title)-1 {
		title = title[0:end]
	}
	start := strings.Index(title, "【")
	if start != -1 && start < end {
		title = title[0:start]
	}
	start = strings.Index(title, ".")
	if start != -1 && start < len(title)-1 {
		title = title[start+1:]
	}
	return title
}

func (o *Operator) QuestionRequire(index int) bool {
	driver := *o.driver
	ele, _ := driver.FindElement(selenium.ByCSSSelector, fmt.Sprintf("#div%d", index))

	var required = false
	req, _ := ele.GetAttribute("req")
	if req == "1" {
		required = true
	}
	return required
}

func (p *Page) GetStructure() []models.Question {
	chrome, page, port := Driver(p.Url)
	if chrome == nil || page == nil || port == 0 {
		return nil
	}

	operator := Operator{driver: page}
	defer func(chrome *selenium.Service) {
		_ = chrome.Stop()
		FreePORT(port)
	}(chrome)
	defer func(driver selenium.WebDriver) {
		_ = driver.Close()
	}(*page)

	var questions []models.Question
	var question models.Question
	var i = 0

	for {
		i += 1
		question = operator.QuestionGet(i)

		if question.Index == 0 { // 未找到
			break
		}
		questions = append(questions, question)
	}

	return questions
}

func (p *Page) Submit(data []any) bool {
	chrome, page, port := Driver(p.Url)

	// fmt.Printf("chrome: %v, page: %v, port: %v\n", chrome, page, port)
	operator := Operator{driver: page}

	defer func(chrome *selenium.Service) {
		_ = chrome.Stop()
		FreePORT(port)
	}(chrome)
	defer func(driver selenium.WebDriver) {
		_ = driver.Close()
	}(*page)

	for idx := 0; idx < len(data); idx++ {
		item := cast.ToStringMap(data[idx])

		if item["type"] == "文本" {
			texts := strings.Split(cast.ToString(item["text"]), "\n")
			if len(texts) == 0 {
				texts = append(texts, "无")
			}

			success := operator.QuestionInput(idx, texts[rand.Intn(len(texts))])
			if !success {
				fmt.Println("问题", idx, "填写失败")
				return false
			}
			continue
		}

		if item["type"] == "单选" || item["type"] == "多选" {
			options := cast.ToSlice(item["options"])
			length := len(options)
			if length == 0 {
				return false
			}
			// 获取选项 rate array
			var rateArray []float64
			var count = 0.0
			for odx := 0; odx < length; odx++ {
				option := cast.ToStringMap(options[odx])
				rateArray = append(rateArray, cast.ToFloat64(option["rate"]))
				count += rateArray[odx]
			}
			// 总概率
			if (count == 0.0) && cast.ToBool(item["required"]) {
				fmt.Println("必选选项 选取率和为 0")
				return false
			}
			// 执行
			if item["type"] == "单选" {
				// 随机 选择
				odx := utils.RateChoose(rateArray)
				success := operator.QuestionSelect(idx, odx)
				if !success {
					fmt.Println("选项", idx, "-", odx, "选择失败")
					return false
				}
				continue

			}
			if item["type"] == "多选" {
				//
				odxArr := utils.RateSample(rateArray, cast.ToBool(item["required"]))
				success := operator.QuestionSelects(idx, odxArr)
				if !success {
					fmt.Println("选项", idx, "选择失败")
					return false
				}
				continue
			}
		}

		if item["type"] == "未知" {
			return false
		}
		return false
	}

	if !operator.QuestionSubmit() {
		fmt.Println("点击提交失败")
		return false
	}

	return true
}

func NewPage(url string) *Page {
	return &Page{Url: url}
}

func Exec(task *models.Task, data []any) {
	rNum := task.Info.Num
	startNum := task.Info.CurrentNum
	for r := startNum; r < rNum; r++ {
		// 检测有没有终止
		if task.Info.Running {
			fmt.Println("[", task.ID, "]", "第", r+1, "次提交...")
			task.Info.CurrentNum = r + 1
			if NewPage(task.Info.Url).Submit(data) {
				fmt.Println("[", task.ID, "]", "第", r+1, " 次提交成功")
				task.Info.SuccessCount += 1
			} else {
				task.Info.FailedCount += 1
				if task.Info.FailedCount >= 5 {
					// 超过 5次 自动终止
					fmt.Println("[", task.ID, "]", "提交失败次数过多 自动终止")
					task.Info.Running = false
				}
			}
		} else {
			// 已终止
			return
		}
	}
	task.Finished = true
}
