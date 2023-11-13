package adminer

import (
	"cms/helpers"
	"cms/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Adminer 当前登录管理员基础信息
func Adminer(c *gin.Context) {
	// 获取用户ID
	tmpId, exists := c.Get("adminer_id")
	if !exists {
		helpers.EndRequest(c, -403, helpers.BuildNullRequestStruct(), "非法请求1")
		return
	}

	adminerId, ok := tmpId.(uint)
	if !ok {
		helpers.EndRequest(c, -403, helpers.BuildNullRequestStruct(), "非法请求2")
		return
	}

	// 查询管理员信息
	adminerInfo, aErr := models.GetAdminerById(adminerId)
	if aErr != nil {
		helpers.EndRequest(c, -1, resData, aErr.Error())
		return
	}

	helpers.EndRequest(c, 200, adminerInfo, "获取成功")
}

// AdminerShow 管理员详情
func AdminerShow(c *gin.Context) {
	// 获取并转换adminerId为uint
	adminerId, err := strconv.ParseUint(c.Param("adminer_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "参数错误")
		return
	}

	// 查询管理员信息
	adminerInfo, aErr := models.GetAdminerById(uint(adminerId))
	if aErr != nil {
		helpers.EndRequest(c, -1, resData, aErr.Error())
		return
	}

	helpers.EndRequest(c, 200, adminerInfo, "获取成功")
}

// AdminerList 管理员列表
func AdminerList(c *gin.Context) {
	var searchParams models.AdminerSearchParams

	// 绑定查询字符串参数到 searchParams 结构体
	if err := c.ShouldBindQuery(&searchParams); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 设置page、pageSize默认值
	setDefaultValues(&searchParams)

	// 管理员列表
	searchRes, err := models.AdminerList(searchParams, searchParams.Page, searchParams.PageSize, "adminer_id desc")
	if err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}
	// 管理组选项
	groupOptions, gErr := models.GetGroupOptions()
	if gErr != nil {
		helpers.EndRequest(c, -1, resData, gErr.Error())
		return
	}
	// 状态选项
	statusOptions := models.GetAdminerStatusOptions()

	res := make(map[string]interface{})
	res["adminerList"] = searchRes.Adminers
	res["adminerTotal"] = searchRes.Total
	res["groupOptions"] = groupOptions
	res["statusOptions"] = statusOptions
	helpers.EndRequest(c, 200, res, "获取成功")
	return
}

// AddAdminer 添加管理员
func AddAdminer(c *gin.Context) {
	var adminer models.Adminer

	// 绑定提交参数到 Adminer 结构体
	if err := c.ShouldBindJSON(&adminer); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 验证参数
	if adminer.Username == "" || adminer.Nickname == "" || adminer.GroupId == 0 || adminer.Status == 0 {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "必填项不能为空")
		return
	}

	// 保存
	saveRes, addErr := models.SaveAdminer(adminer, adminer.AdminerId)
	if addErr != nil {
		helpers.EndRequest(c, -1, resData, addErr.Error())
	} else {
		helpers.EndRequest(c, 200, saveRes, "添加成功")
	}
}

// UpdateAdminer 修改管理员
func UpdateAdminer(c *gin.Context) {
	var adminer models.Adminer

	// 绑定提交参数到 Adminer 结构体
	if err := c.ShouldBindJSON(&adminer); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 验证参数
	if adminer.Nickname == "" || adminer.GroupId == 0 || adminer.Status == 0 {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "必填项不能为空")
		return
	}

	// 保存
	saveRes, adderr := models.SaveAdminer(adminer, adminer.AdminerId)
	if adderr != nil {
		helpers.EndRequest(c, -1, resData, adderr.Error())
	} else {
		helpers.EndRequest(c, 200, saveRes, "修改成功")
	}
}

// UpdateAdminerStatus 更新管理员状态
func UpdateAdminerStatus(c *gin.Context) {
	// 接收参数
	var formReq struct {
		AdminerId uint `json:"adminer_id"`
		Status    uint `json:"status"`
	}
	// 绑定查询字符串参数到 formReq 结构体
	if err := c.ShouldBindJSON(&formReq); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 调用模型查询列表
	_, err := models.AdminerStatusSwitch(formReq.AdminerId, formReq.Status)
	if err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}

	helpers.EndRequest(c, 200, resData, "更新成功")
	return
}

// DeleteAdminer 删除管理员
func DeleteAdminer(c *gin.Context) {
	// 获取并转换adminerId为uint
	adminerId, err := strconv.ParseUint(c.Param("adminer_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "Invalid or missing adminerId")
		return
	}

	// 删除管理员
	dRes, deleteErr := models.DeleteAdminer(uint(adminerId))
	if deleteErr != nil {
		helpers.EndRequest(c, -1, resData, deleteErr.Error())
		return
	}

	// 检查删除操作是否影响了任何行
	if affectedRows, ok := dRes["affected_rows"].(int64); ok && affectedRows > 0 {
		helpers.EndRequest(c, 200, resData, "删除成功")
	} else {
		helpers.EndRequest(c, -1, resData, "删除失败")
	}
}

// ResetPass 重置管理员密码
func ResetPass(c *gin.Context) {
	// 获取并转换adminerId为uint
	adminerId, err := strconv.ParseUint(c.Query("adminer_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "Invalid or missing adminerId")
		return
	}

	// 获取当前用户ID
	tmpOperatorId, exists := c.Get("adminer_id")
	if !exists {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "operatorId is required")
		return
	}

	if operatorId, ok := tmpOperatorId.(uint); ok {
		if resetRes, err := models.ResetPass(uint(adminerId), operatorId); err != nil {
			helpers.EndRequest(c, -1, resData, err.Error())
		} else {
			helpers.EndRequest(c, http.StatusOK, resetRes, "重置成功")
		}
	} else {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "当前管理员id未获取到")
	}
}

// setDefaultValues 为搜索参数设置默认值
func setDefaultValues(params *models.AdminerSearchParams) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 1
	}
}
