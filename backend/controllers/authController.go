package controllers

import (
	"cms/helpers"
	"cms/models"
	"github.com/gin-gonic/gin"
	"time"
)

// 空接口返回Data
var resData map[string]interface{}

func Login(c *gin.Context) {
	// 接收参数
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginReq); err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}
	username := loginReq.Username
	password := loginReq.Password

	// 验证参数
	if len(username) < 1 || len(password) < 1 {
		helpers.EndRequest(c, -1, resData, "用户名或者密码不能为空!")
		return
	}

	// 调用模型进行登录操作
	loginRes, err := models.AdminerLogin(username, password)
	if err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}

	// 记录最新登录的ip和时间
	loginIp := helpers.GetClientIP(c)
	loginTime := time.Now().Unix()
	var adminerId uint
	if val, ok := loginRes["adminer_id"].(uint); ok {
		adminerId = val
	}
	models.UpdateAdminerLast(adminerId, loginIp, uint64(loginTime))

	loginRes["last_ip"] = loginIp
	loginRes["last_time_str"] = helpers.FormatTime(loginTime)
	helpers.EndRequest(c, 200, loginRes, "登录成功")
	return
}
