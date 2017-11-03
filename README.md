# captcha
![](https://github.com/bsed/images/raw/master/2017/09/echo_captcha.png)

## 纯数字验证码

![](https://raw.githubusercontent.com/bsed/images/master/2017/09/captcha_digits.png)
## 大写字母和数字验证码
![](https://raw.githubusercontent.com/bsed/images/master/2017/09/captcha_digit_upper.png)

## 大小写和数字组合
![](https://raw.githubusercontent.com/bsed/images/master/2017/09/captcha_digits_upper_ligit.png)

## 小写字母加数字
![](https://raw.githubusercontent.com/bsed/images/master/2017/09/captcha_lowercase_digits.png)

## 白色背景/无线条（领导要求）
![](https://raw.githubusercontent.com/bsed/images/master/2017/09/captcha_background_white.png)

 
## 添加REDIS支持
 
### 用法:
 
配置文件 config.ini
```ini
[global]
# machine id
expire=10
# startup port
port=8666
# redis config
redisAddr=localhost:6379
```
用例:
```go
// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// example of HTTP server that uses the captcha package.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"captcha"

	"github.com/zpatrick/go-config"
)

var formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	d := struct {
		Id string
	}{
		captcha.New(),
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("id"), r.FormValue("captcha")) {
		io.WriteString(w, "Wrong captcha solution! No robots allowed!\n")
	} else {
		io.WriteString(w, "Great job, human! You solved the captcha.\n")
	}
	io.WriteString(w, "<br><a href='/'>Try another one</a>")
}

//获取验证码id, /captcha/{{id}}.png获取到验证码
func getCaptchaFormHandler(w http.ResponseWriter, r *http.Request) {
	lenStr := r.FormValue("len")
	if lenStr == "" {
		d := struct {
			Code int
			Data string
		}{
			0,
			captcha.New(),
		}
		json.NewEncoder(w).Encode(d)
	} else {
		len, err := strconv.Atoi(lenStr)
		if err != nil {
			d := struct {
				Code int
				Data string
			}{
				1,
				"获取验证码id失败",
			}
			json.NewEncoder(w).Encode(d)
			return
		}
		d := struct {
			Code int
			Data string
		}{
			0,
			captcha.NewLen(len),
		}
		json.NewEncoder(w).Encode(d)
	}
}

//验证验证码
func verifyFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("id"), r.FormValue("captcha")) {
		result := struct {
			Code int
			Data string
		}{
			1,
			"验证失败",
		}
		json.NewEncoder(w).Encode(result)
	} else {
		result := struct {
			Code int
			Data string
		}{
			0,
			"验证成功",
		}
		json.NewEncoder(w).Encode(result)
	}
}

func main() {
	iniFile := config.NewINIFile("config.ini")
	c := config.NewConfig([]config.Provider{iniFile})
	if err := c.Load(); err != nil {
		panic("load config error")
	}
	expire, err := c.Int("global.expire")
	if err != nil {
		panic("get expire config error")
	}
	port, err := c.Int("global.port")
	if err != nil {
		panic("get port config error")
	}
	redisAddr, err := c.String("global.redisAddr")
	if err != nil {
		panic("get redisaddr config error")
	}
	captcha.SetCustomStore(captcha.NewRedisStore(redisAddr, "", 0, time.Duration(expire)*time.Minute))
	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/process", processFormHandler)
	http.HandleFunc("/v1/captcha_id", getCaptchaFormHandler)
	http.HandleFunc("/v1/verify", verifyFormHandler)
	http.Handle("/v1/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	fmt.Println("Server is at :" + strconv.Itoa(port))
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}

const formTemplateSrc = `<!doctype html>
<head><title>Captcha Example</title></head>
<body>
<script>
function setSrcQuery(e, q) {
	var src  = e.src;
	var p = src.indexOf('?');
	if (p >= 0) {
		src = src.substr(0, p);
	}
	e.src = src + "?" + q
}

function playAudio() {
	var le = document.getElementById("lang");
	var lang = le.options[le.selectedIndex].value;
	var e = document.getElementById('audio')
	setSrcQuery(e, "lang=" + lang)
	e.style.display = 'block';
	e.autoplay = 'true';
	return false;
}

function changeLang() {
	var e = document.getElementById('audio')
	if (e.style.display == 'block') {
		playAudio();
	}
}

function reload() {
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	setSrcQuery(document.getElementById('audio'), (new Date()).getTime());
	return false;
}
</script>
<select id="lang" onchange="changeLang()">
<option value="en">English</option>
<option value="ru">Russian</option>
<option value="zh">Chinese</option>
</select>
<form action="/process" method=post>
<p>Type the numbers you see in the picture below:</p>
<p><img id=image src="/captcha/{{.Id}}.png" alt="Captcha image"></p>
<a href="#" onclick="reload()">Reload</a> | <a href="#" onclick="playAudio()">Play Audio</a>
<audio id=audio controls style="display:none" src="/captcha/{{.Id}}.wav" preload=none>
You browser doesn't support audio.
<a href="/captcha/download/{{.Id}}.wav">Download file</a> to play it in the external player.
</audio>
<input type=hidden name=id value="{{.Id}}"><br>
<input name=captcha>
<input type=submit value=Submit>
</form>
`


```
