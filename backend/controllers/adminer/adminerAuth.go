package adminer

import (
	"cms/helpers"
	"cms/models"
	"cms/models/format"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// 空接口返回Data
var resData map[string]interface{}

func MenuAuths(c *gin.Context) {
	// 获取用户ID
	tmpId, exists := c.Get("adminer_id")
	if !exists {
		helpers.EndRequest(c, -403, helpers.BuildNullRequestStruct(), "非法请求")
		return
	}

	adminerId, ok := tmpId.(uint)
	if !ok {
		helpers.EndRequest(c, -403, helpers.BuildNullRequestStruct(), "非法请求")
		return
	}

	// 查询当前用户分组
	var authItems []format.AuthItem
	var formattedItems []format.FormattedAuthItem
	var adminerAuths []*models.AdminerAuth
	var authErr error
	authErr = nil

	// 管理员信息
	adminer, aErr := models.GetAdminerById(adminerId)
	if aErr != nil {
		helpers.EndRequest(c, -403, helpers.BuildNullRequestStruct(), "管理员不存在")
		return
	}
	if adminer.GroupId == 1 {
		// 超级管理员全部权限
		adminerAuths, authErr = models.GetAdminerAuths()
	} else {
		// 普通管理员指定组权限
		adminer, adminerErr := models.GetAdminerById(adminerId)
		if adminerErr != nil {
			helpers.EndRequest(c, -1, helpers.BuildNullRequestStruct(), adminerErr.Error())
			return
		}
		adminerAuths, authErr = models.GetAdminerAuthsByGroupId(adminer.GroupId)
	}
	if authErr != nil {
		helpers.EndRequest(c, -1, helpers.BuildNullRequestStruct(), authErr.Error())
		return
	}
	// 格式化数据
	authItems = _convertAuthItems(adminerAuths)
	formattedItems = format.BuildAuthTree(authItems)
	//res := make(map[string]interface{})
	//res["roles"] = formattedItems
	helpers.EndRequest(c, 200, formattedItems, "获取成功")
}

// AuthList 权限管理列表
func AuthList(c *gin.Context) {
	// 绑定查询参数
	var searchParams models.AdminerAuthSearchParams
	if err := c.ShouldBindQuery(&searchParams); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 设置分页默认值
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}
	if searchParams.PageSize <= 0 {
		searchParams.PageSize = 10
	}

	// 过滤参数
	if searchParams.AuthName != "" {
		searchParams.AuthName = strings.TrimSpace(searchParams.AuthName)
	}
	if searchParams.AuthConAct != "" {
		searchParams.AuthConAct = strings.TrimSpace(searchParams.AuthConAct)
	}

	// 查询权限列表
	searchRes, sErr := models.AdminerAuthList(searchParams, "auth_id desc")
	if sErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, sErr.Error())
		return
	}

	// 获取权限分组选项
	authGroupOptions, agErr := models.GetAuthGroupOptions()
	if agErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, agErr.Error())
		return
	}

	// 获取父级权限选项
	prentAuthOptions, paErr := models.GetAdminerAuthPrentOptions()
	if paErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, paErr.Error())
		return
	}

	result := make(map[string]interface{})
	result["list"] = searchRes.AdminerAuths
	result["total"] = searchRes.Total
	result["statusOptions"] = models.GetAdminerAuthStatusOptions()
	result["authGroupOptions"] = authGroupOptions
	result["prentAuthOptions"] = prentAuthOptions
	result["isMenuOptions"] = models.GetAdminerAuthIsMenuOptions()
	helpers.EndRequest(c, 200, result, "获取成功")
}

// AuthShow 权限管理详情
func AuthShow(c *gin.Context) {
	// 接收参数
	authId, err := strconv.ParseUint(c.Param("auth_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, nil, "参数错误")
		return
	}

	// 查询权限信息
	auth, aErr := models.GetAdminerAuthById(uint(authId))
	if aErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, nil, aErr.Error())
		return
	}

	helpers.EndRequest(c, 200, auth, "获取成功")
}

// SwitchAuthStatus 权限状态切换
func SwitchAuthStatus(c *gin.Context) {
	// 接收参数
	var formReq struct {
		AuthId uint `json:"auth_id"`
		Status uint `json:"status"`
	}
	// 绑定查询字符串参数到 formReq 结构体
	if err := c.ShouldBindJSON(&formReq); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 调用模型查询列表
	var statusForm = models.AuthSwitchStatusForm{
		AuthId:      formReq.AuthId,
		Status:      formReq.Status,
		LastAdminer: helpers.GetCacheAdminerId(c),
	}
	_, err := models.SwitchAuthStatus(statusForm)
	if err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}

	helpers.EndRequest(c, 200, resData, "更新成功")
	return
}

// SwitchMenuShow  菜单状态切换
func SwitchMenuShow(c *gin.Context) {
	// 接收参数
	var formReq struct {
		AuthId     uint `json:"auth_id"`
		IsMenuShow uint `json:"is_menu_show"`
	}
	// 绑定查询字符串参数到 formReq 结构体
	if err := c.ShouldBindJSON(&formReq); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 调用模型查询列表
	var menuForm = models.AuthSwitchMenuForm{
		AuthId:      formReq.AuthId,
		IsMenuShow:  formReq.IsMenuShow,
		LastAdminer: helpers.GetCacheAdminerId(c),
	}
	_, err := models.SwitchAuthMenuShow(menuForm)
	if err != nil {
		helpers.EndRequest(c, -1, resData, err.Error())
		return
	}

	helpers.EndRequest(c, 200, resData, "更新成功")
	return
}

// SaveAuth 保存权限
func SaveAuth(c *gin.Context) {
	var auth models.AdminerAuth

	// 绑定提交参数到 Adminer 结构体
	if err := c.ShouldBindJSON(&auth); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 验证参数
	if auth.AuthName == "" || auth.AuthGroupId == 0 {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "必填项不能为空")
		return
	}
	if auth.PrentId > 0 && auth.AuthConAct == "" {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "请填写权限")
		return
	}

	// 保存
	message := ""
	if auth.AuthId > 0 {
		message = "更新成功"
	} else {
		message = "添加成功"
	}
	saveRes, addErr := models.SaveAdminerAuth(auth, auth.AuthId)
	if addErr != nil {
		helpers.EndRequest(c, -1, resData, addErr.Error())
	} else {
		helpers.EndRequest(c, 200, saveRes, message)
	}
}

// DeleteAuth 删除权限
func DeleteAuth(c *gin.Context) {
	// 获取并转换adminerId为uint
	authId, err := strconv.ParseUint(c.Param("auth_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "参数错误")
		return
	}

	// 删除管理员
	adminerId := helpers.GetCacheAdminerId(c)
	dRes, deleteErr := models.DeleteAdminerAuth(uint(authId), adminerId)
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

// _convertAuthItems 转换AuthItem数据类型
func _convertAuthItems(adminerAuths []*models.AdminerAuth) []format.AuthItem {
	var authItems []format.AuthItem
	for _, v := range adminerAuths {
		authItem := format.AuthItem{
			AuthId:      v.AuthId,
			AuthName:    v.AuthName,
			AuthConAct:  v.AuthConAct,
			AuthGroupId: v.AuthGroupId,
			Status:      v.Status,
			Remark:      v.Remark,
			PrentId:     v.PrentId,
			IsMenuShow:  v.IsMenuShow,
		}
		authItems = append(authItems, authItem)
	}
	return authItems
}
