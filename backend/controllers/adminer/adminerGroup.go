package adminer

import (
	"cms/helpers"
	"cms/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AdminerGroupShow 管理员组详情
func AdminerGroupShow(c *gin.Context) {
	// 获取并转换groupId为uint
	groupId, err := strconv.ParseUint(c.Param("group_id"), 10, 32)
	if err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "Invalid or missing groupId")
		return
	}

	// 根据groupId查询分组信息
	groupInfo, gErr := models.GetAdminerGroupById(uint(groupId))
	if gErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, gErr.Error())
		return
	}

	// 获取全部权限
	allAuths, authErr := models.GetAdminerAuths()
	if authErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, authErr.Error())
		return
	}

	// 获取组权限列表
	groupAuths, gErr := models.GetGroupAuthsByGroupId(uint(groupId))
	if gErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, gErr.Error())
		return
	}
	var filterGroupAuths []*models.AdminerGroupAuth
	if groupAuths != nil {
		for _, auth := range allAuths {
			for _, groupAuth := range groupAuths {
				if auth.PrentId != 0 && auth.AuthId == groupAuth.AuthId {
					filterGroupAuths = append(filterGroupAuths, groupAuth)
				}
			}
		}

	}
	authTree := models.BuildAuthTree(allAuths, filterGroupAuths)

	result := make(map[string]interface{})
	result["group_id"] = uint(groupId)
	result["name"] = groupInfo.Name
	result["remark"] = groupInfo.Remark
	result["auth_tree"] = authTree
	helpers.EndRequest(c, 200, result, "获取成功")
}

// AdminerGroupList 管理员组列表
func AdminerGroupList(c *gin.Context) {
	var searchParam models.AdminerGroupSearchParam
	if err := c.ShouldBindQuery(&searchParam); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 初始化分页相关参数
	if searchParam.Page <= 0 {
		searchParam.Page = 1
	}
	if searchParam.PageSize <= 0 {
		searchParam.PageSize = 10
	}

	// 管理员组列表
	groupList, lErr := models.AdminerGroupList(searchParam.Page, searchParam.PageSize, "group_id desc")
	if lErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, lErr.Error())
		return
	}

	// 获取权限列表
	allAuths, authErr := models.GetAdminerAuths()
	if authErr != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, authErr.Error())
		return
	}
	authTree := models.BuildAuthTree(allAuths, nil)

	res := make(map[string]interface{})
	res["total"] = groupList.Total
	res["list"] = groupList.AdminerGroups
	res["auth_tree"] = authTree
	helpers.EndRequest(c, 200, res, "获取成功")
}

// SaveAdminerGroup 保存管理员分组
func SaveAdminerGroup(c *gin.Context) {
	// 绑定参数
	var formData models.AdminerGroupForm
	if err := c.ShouldBindJSON(&formData); err != nil {
		helpers.EndRequest(c, http.StatusBadRequest, resData, err.Error())
		return
	}

	// 验证参数
	if formData.Name == "" {
		helpers.EndRequest(c, http.StatusBadRequest, resData, "必填项不能为空")
		return
	}
	// 添加操作人员ID
	formData.CreateAdminer = uint64(helpers.GetCacheAdminerId(c))

	// 保存
	var resStr string
	if formData.GroupId != nil && *formData.GroupId > 0 {
		resStr = "更新成功"
	} else {
		resStr = "添加成功"
	}
	saveRes, addErr := models.SaveGroupForm(formData)
	if addErr != nil {
		helpers.EndRequest(c, -1, resData, addErr.Error())
	} else {
		helpers.EndRequest(c, 200, saveRes, resStr)
	}
}
