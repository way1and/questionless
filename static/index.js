const API = "http://localhost:80"
const UrlInput = document.getElementById("url-input")
const FormList = document.getElementById("form-list")
const ButtonSubmit = document.getElementById("btn-submit")
const ButtonScan = document.getElementById("btn-scan")
const SelectType = document.getElementById("select-type")
const NumInput = document.getElementById("num-input")



var Scanned = false
var Scanning = false
var Submitting = false
var Data = []
var UrlScanned = ""
var TaskID = ""
var TaskSecret = ""
var StatusInfo = {}
var IntervalID = 0

function ClickScanUrl(){
    // 防止重复
    if (Scanning){
        return
    }

    let url = UrlInput.value
    if (!/^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\*\+,;=.]+$/.test(url)){
       layer.msg("URL 填写错误")
        return false
    }

    layer.msg("开始扫描, 请耐心等待...")
    Scanning = true
    $.ajax({
        url:API + "/scan?url=" + url,
        async:true,
        method:"GET",
        success(res){
            console.log(res)
            let data

            if (res.data == null){
                layer.msg(res["msg"])
               return
            }
            data = res.data
            layer.msg("扫描成功")
            Scanning = false // 设置状态
            Scanned = true // 设置状态
            UrlScanned = url
            $(ButtonSubmit).addClass("layui-bg-green")
            $(ButtonScan).removeClass("layui-bg-blue")
            $(NumInput).val("60")
                .attr("disabled", false)
            if (data!== null && data.length !== 0){
                refreshFormList(data)
            }
        },
        error(res){
            Scanning = false // 设置状态
            let msg = res["responseJSON"].msg
            layer.msg(`扫描失败: ${msg}`)
            $(ButtonSubmit).removeClass("layui-bg-green")
            refreshFormList(null)
        }
    })
}

function refreshFormList(data){
    $(FormList).empty()
    if (!data){
        return
    }
    Data = data // 设置
    GenerateAllData()
    let title, option,  card, top, bottom
    let item
    let tag;
    let req, idx;
    for (let i = 0; i < data.length; i++) {
        item = data[i]
        // 卡片
        card = document.createElement("div")
        $(card).addClass("layui-card")
        $(card).attr("id", `form-${i}`)
        // 上半
        top = document.createElement("div")
        $(top).addClass("layui-card-header").css("overflow", "hidden")
        {
            // 标签
            tag = document.createElement("span")
            $(tag).addClass("layui-badge")

            switch (item["type"]){
                case "单选":
                    $(tag).addClass("layui-bg-blue")
                    break
                case "多选":
                    $(tag).addClass("layui-bg-green")
                    break
                case "文本":
                    $(tag).addClass("layui-bg-orange")
                    break
                default:
                    $(tag).addClass("layui-bg-cyan")
            }

            $(tag).text(item["type"])
            top.appendChild(tag)



            // 题号
            idx = document.createElement("span")
            $(idx).text(`  ${item["index"]}.`)
            top.appendChild(idx)

            // 必填
            if (item["required"]){
                req = document.createElement("span")
                $(req).text("*")
                    .css("color", "#ff5722")
                top.appendChild(req)
            }

            // 标题
            title = document.createElement("span")
            $(title).text(`  ${item["title"]}`) // 标题
            top.appendChild(title)

        }
        // 下半
        bottom = document.createElement("div")
        $(bottom).addClass("layui-card-body")
        {
            if (item["options"] != null) {
                for (let optionData of item["options"]) {
                    option = document.createElement("div")

                    $(option).text(`[${optionData["index"]}]  ${optionData["name"]}`)
                    if (optionData["extended"]) {
                        option.innerHTML +=  `
                        <div class="layui-input-inline layui-input-wrap" >
                            <input class="layui-input layui-padding-1 layui-bg-gray" id="url-input" value="附加文本内容" disabled/>
                        </div>`
                    }

                    option.innerHTML += `<span style="color: #2f363c" id="rate-${i}-${optionData["index"]-1}">${"  "+optionData["rate"]}%</span>`
                    bottom.appendChild(option)
                }
            }else if (item["type"] === "文本") {
                let div;
                div = document.createElement("div")
                $(div).html(`  
        <div class="layui-form" style="">
            <div class="layui-form-item">
                    <label class="layui-form-label"> 描述 </label>
                    <div class="layui-input-inline layui-input-wrap" >
                        <input placeholder="文本域" class="layui-input" disabled/>
                    </div>
            </div></div>
        `)
                bottom.appendChild(div)

            }

        }
        // 控制按钮
        let config = document.createElement("button")
        $(config).addClass("layui-btn layui-btn-primary")
            .text("设置")
            .attr("index" , `${item["index"]}`)
            .css("position","absolute")
            .css("right", "20px")
            .css("bottom", "20px").css("z-index", "100")
            .on("click", ()=>{
            ClickConfigShow($(config).attr("index"))
        })
        bottom.appendChild(config)
        // 将选项加入options
        card.appendChild(top)
        card.appendChild(bottom)
        $(FormList).append(card)

    }
}

function GenerateAllData(){
    let item
    let option;
    let data = Data
    for (let idx in data) {
        GenerateDefaultRow(idx)
    }
    Data = data
}

function GenerateDefaultRow(idx){
    let item = Data[idx]
    if (item["type"] === "文本") {
        item["action"] = "input"
        item["ai"] = false
        item["text"] = "无"
    }
    if (item["type"] === "未知") {
        item["action"] = "-"
    }
    let rate
    // 有 options 循环遍历
    if (item["options"] && item["options"].length !== 0) {
        item["action"] =  "单选"? "select":"selects"
        for (let optionIdx in item["options"]) {
            if (item["options"][optionIdx]["extended"]){
                item["options"][optionIdx]["text"] = "无"
            }

            if (item["type"] === "单选") {
                rate = ((1 / (item["options"].length) )* 100) .toFixed(2) // 单选
            }else {


                rate = 50 // 多选
            }
           item["options"][optionIdx]["rate"] = rate
        }
    }
    Data[idx] = item
}

function onRateInput(form){
    let index = $(form).attr("index")
    let optionIndex = $(form).attr("option-index")
    let val = $(`#rate-input-${optionIndex}`).val()
    $(`#rate-${index}-${optionIndex}`).text(`  ${val}%`)
    Data[parseInt(index)]["options"][optionIndex]["rate"] = val
}

function onTextInput(input){
    let index = $(input).attr("index")
    let optionIndex = $(input).attr("option-index")

    if (!optionIndex){
        Data[parseInt(index)]["text"] = $(input).val()
    }else{
        Data[parseInt(index)]["options"][parseInt(optionIndex)]["text"] = $(input).val()
    }
}

function onConfigReset(btn){
    let idx = $(btn).attr("index")
    layer.closeAll("page")
    GenerateDefaultRow(idx)
    refreshFormList(Data)
}

function ClickConfigShow(index){
    if (!index || index === 0 || !Scanned || Scanning){
        return
    }
    let div = document.createElement("div")
    $(div).attr("id", "config-popup")
        .css("padding-top", "20px")


    let item = Data[index - 1]
    // 题目
    let title = document.createElement("p")
    $(title).html(`<p style="margin-top: 20px;margin-bottom: 20px;margin-left:20px;"><span style="color: #2f363c;">${item["title"]}</span></p>`).css("margin-left", "20px")

    div.appendChild(title)


    let textUpload;
    if (item["type"] === "文本") {
        let form = document.createElement("form")
        $(form).addClass("layui-form")
            .css("padding-right", "18px")
            .css("overflow", "hidden")
        let text = document.createElement("div")

        $(text).addClass("layui-form-item")
            .html(`<textarea id="text-${index - 1}" index="${index - 1}" oninput="onTextInput(this)" placeholder="文本" style="margin-left: 10px;" class="layui-textarea">${item["text"]}</textarea>`)

        form.appendChild(text)

        textUpload = document.createElement("button")
        $(textUpload).html(`<button index="${index - 1}" class="layui-btn layui-btn-normal layui-bg-gray" style="color: #16baaa;" onclick="chooseText(this)">添加文本文件</button>`)
        div.appendChild(textUpload)
        div.appendChild(form)

        let desc = document.createElement("p")
        $(desc).html(`
            <p style="margin-left: 20px;" class="layui-font-12">1. 每一行是一个答案, 对每行的选取是均等的</p><br/>
            <p style="margin-left: 20px;" class="layui-font-12">2. 如果需要提高某种答案的选取率请将此答案重复输入多行</p>`)

        div.appendChild(desc)

    } else {
        for (let optionIdx in item.options) {
            let option = item.options[optionIdx]
            //
            if (item["type"] === "单选" || item["type"] === "多选") {
                console.log("123")
                let form = document.createElement("form")
                $(form).addClass("layui-form")

                let name = document.createElement("p")
                $(name).html(`<span>${option["index"]}. ${option["name"]}</span>`)
                $(name).css("margin-left", "20px")

                form.appendChild(name)
                { // 选项配置表单
                    let rate = document.createElement("div")
                    $(rate).addClass("layui-form-item")

                    $(rate).html(`
                    <label class="layui-form-label"> 选取率 </label>
                        <div class="layui-input-inline  layui-input-wrap">
                            <input class="layui-input" value="${option['rate']}" id="rate-input-${optionIdx}" oninput="onRateInput(form)"/>
                        <div class="layui-input-split layui-input-suffix"> % </div>
                        </div>`)
                    form.appendChild(rate)
                }
                $(form).attr("index", item["index"] - 1)
                $(form).attr("option-index", optionIdx)

                div.appendChild(form)
            }
            // 其他情况
            else if (item["type"] === "未知") {

            }

        }
    }

    // 重置
    let reset = document.createElement("button")

    $(reset)
        .html(`<button class="layui-btn layui-bg-red" index="${index - 1}" onclick="onConfigReset(this)" style="position:absolute;bottom: 20px;right: 20px;z-index: 100;">重置</button>`)
    div.appendChild(reset)
    // 页面层
    layer.open({
        type: 1,
        area: ['480px', '420px'], // 宽高
        shadeClose: true,
        title:`填写设置 - 问题 ${index}  [ ${item["type"]} ]`,
        content: div.innerHTML
    });
}

async function chooseText(btn){
    let index = $(btn).attr("index")
    let optionIndex = $(btn).attr("option-index")
    let id = ("text-" + index) + (optionIndex?`-${optionIndex}`:"")

   let [handler] = await window.showOpenFilePicker()
    let f = await handler.getFile()
    let reader = new FileReader()
    // 传入需要被转换的文本流 file,这个是转字符串的关键方法
    reader.readAsText(f)
    // onload是异步的,封装的话可以用promise
    reader.onload = () => {
        // 输出字符串
        console.log(reader.result)
        let data = reader.result
        let input = $(`#${id}`)
        input.val(input.val() + data)
        onTextInput(input)
    }
}

function Check(){
    if (!Scanned){
        return [false, "请先扫描一张问卷"]
    }
    if (Submitting){
        return [false, "正在启动填写服务中, 请终止后重试"]
    }
    if (!parseInt(NumInput.value)){
        return [false, "填写份数不是一个有效的数字"]
    }
    if (SelectType.value === ""){
        return [false, "请选择一种问卷类型"]
    }

    for (let index in Data){
        let item = Data[index]
        //
        if (item["type"] === "多选" || item["type"] === "单选"){
            let rateTotal = 0.0
            let rate;
            for (let option of item["options"]){
                rate = parseFloat(option["rate"])
                if (isNaN(rate)){
                    // 有错误值
                    return [false, `问题 ${item["index"]} 选项 ${option["index"]} 的参数不是一个有效数字`]
                }
                rateTotal += rate
            }
            if (item["required"]){
                if (rateTotal === 0){
                    layer.msg(`注意: 问题 ${item["index"]} 为必填项, 但问题选取率之和为 0`)
                }
            }
        }
    }

    return [true, ""]
}
// 点击运行
function ClickSubmitForm(){
    if (Submitting){
        SubmitClose()
        return;
    }
    let [res, data] = Check()
    if (!res){
        layer.msg(`启动失败! ${data} `)
        return
    }
   Submit()
}

// 开启服务
function Submit(){


    let data = {
        "url": UrlScanned,
        "data": Data,
        "num": parseInt(NumInput.value),
        "type": SelectType.value
    }

    console.log("启动")
    console.log("json")
    $.ajax({
        url:API + "/submit",
        data: JSON.stringify(data),
        async:true,
        method:"POST",
        dataType:"json",
        headers:{'Content-Type':'application/json;charset=utf8'},
        success(res){
            if (res && res["data"]){
                console.log(res)

                Submitting = true
                TaskID = res.data["id"]
                TaskSecret = res["secret"]

                StartSuccess()
                GetStatus()
            }else{
                // 失败
                layer.msg(`启动失败: ${res["msg"]}`)
            }
        },
        error(res){
            layer.msg(`服务启动失败: ${res["msg"]}`)
            Submitting = false

        }
    })
}

// 启动成功
function StartSuccess(){
    $(ButtonSubmit).text("终止服务")
        .removeClass("layui-bg-green")
        .removeClass("layui-btn-primary")
        .addClass("layui-bg-red")

    let div = document.createElement("div")
    // 删除原先的
    $("#progress").remove()
    $(div).html(`<div class="layui-progress" lay-filter="progress" id="progress" lay-showpercent="true"><div id="progress-bar" class="layui-progress-bar layui-bg-blue" lay-percent="0 / ${NumInput.value}"></div></div>`)
    document.getElementById("top-panel")
        .appendChild(div)
    layui.element.render('progress', 'progress');
}

function SubmitClose(){

    Close()
    layer.msg("已终止填写")
    // 未完成关闭
    $("#progress-bar")
        .removeClass("layui-bg-blue")
        .addClass("layui-bg-red")

}
function GetStatus(){
    IntervalID = setInterval( ()=>{
        if (Submitting){
            Status().then(()=>{
                console.log("接收到信息")
            })
        }
    }, 3000)
}

// 终止服务
function Close(){
    Submitting = false
    layer.msg("终止服务中...")
    $(ButtonSubmit).text("一键填写")
        .removeClass("layui-bg-red")
        .addClass("layui-bg-green")
    $.ajax({
        url: `${API}/close`,
        async:true,
        dataType: "json",
        headers:{'Content-Type':'application/json;charset=utf8'},
        data:JSON.stringify({"task_id":TaskID, "task_secret":TaskSecret}),
        method:"POST",
        success(res){
            Submitting = false
            clearInterval(IntervalID)
        },
        error(res){
            if (res && res["msg"]){
                layer.msg(`服务器错误 请尝试重启服务器`)
                Submitting = false
            }
        }
    })
}


// 获取服务信息
async function Status(){
    $.ajax({
        url: `${API}/status?task_id=${TaskID}&task_secret=${TaskSecret}`,
        async:true,
        method:"GET",
        success(res){
            if (res && res["data"]){
                let info = res["data"]
                StatusInfo = info
                console.log(info)
                $("#progress-bar").attr("lay-percent", `${info["current_num"]} / ${info["num"]}`)
                layui.element.render('progress', 'progress');

                if (info["num"]===info["current_num"]){
                    // 完成自动关闭
                    Close()
                    // 完成
                    layer.msg("已完成填写")
                    $("#progress-bar")
                        .removeClass("layui-bg-blue")
                        .addClass("layui-bg-green")
                }
            }
        },
        error(res){
            StatusInfo = null
            layer.msg("无法连接到服务, 已终止")
            console.log(res)
            $("#progress-bar")
                .removeClass("layui-bg-blue")
                .addClass("layui-bg-red")
            close()
        }
    })

}
