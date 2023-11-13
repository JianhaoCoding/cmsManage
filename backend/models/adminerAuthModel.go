package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AdminerAuth struct {
	AuthId        uint   `gorm:"primary_key" json:"auth_id"`
	AuthName      string `json:"auth_name"`
	AuthConAct    string `json:"auth_con_act"`
	AuthRoles     string `json:"auth_roles"`
	AuthGroupId   uint   `json:"auth_group_id"`
	Status        uint   `json:"status"`
	Remark        string `json:"remark"`
	PrentId       uint   `json:"prent_id"`
	IsMenuShow    uint   `json:"is_menu_show"`
	Sort          uint   `json:"sort"`
	CreateTime    uint64 `json:"create_time"`
	CreateAdminer uint   `json:"create_adminer"`
	LastTime      uint64 `json:"last_time"`
	LastAdminer   uint   `json:"last_adminer"`
}

// AdminerAuthSearchParams 权限列表查询搜索参数
type AdminerAuthSearchParams struct {
	AuthName    string `form:"auth_name"`
	AuthConAct  string `form:"auth_con_act"`
	Status      uint   `form:"status"`
	AuthGroupId uint   `form:"auth_group_id"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

type AdminerAuthQueryResult struct {
	AdminerAuths []*AdminerAuth
	Total        int64
}

// AuthWithChildren 页面权限选项结构体
type AuthWithChildren struct {
	AuthId   uint                `json:"auth_id"`
	AuthName string              `json:"auth_name"`
	IsCheck  uint                `json:"is_check"`
	Children []*AuthWithChildren `json:"children,omitempty"`
}

// GroupWithAuths 页面权限选项结构体
type GroupWithAuths struct {
	AuthGroupId   uint                `json:"auth_group_id"`
	AuthGroupName string              `json:"auth_group_name"`
	AuthName      string              `json:"auth_name"`
	Children      []*AuthWithChildren `json:"children"`
}

// AdminerAuthStatusOptions 状态开关结构
type AdminerAuthStatusOptions struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// AdminerAuthGroupOptions  权限组选项结构
type AdminerAuthGroupOptions struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// AdminerAuthIsMenuOptions 是否菜单展示开关结构
type AdminerAuthIsMenuOptions struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// AdminerPrentAuthOptions 父级菜单选择结构
type AdminerPrentAuthOptions struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// AuthSwitchStatusForm 权限状态更新表单
type AuthSwitchStatusForm struct {
	AuthId      uint `json:"auth_id"`
	Status      uint `json:"status"`
	LastAdminer uint `json:"last_adminer"`
}

// AuthSwitchMenuForm 权限是否菜单展示表单
type AuthSwitchMenuForm struct {
	AuthId      uint `json:"auth_id"`
	IsMenuShow  uint `json:"is_menu_show"`
	LastAdminer uint `json:"last_adminer"`
}

func (AdminerAuth) TableName() string {
	return "adminer_auth"
}

// AuthGroups 定义权限组
var AuthGroups = map[uint]string{
	1: "无所不能",
	2: "管理员",
	3: "内容管理",
	4: "用户管理",
	5: "其他管理",
	6: "APP管理",
	7: "活动管理",
	8: "视频管理",
}

func (AdminerAuth) PageSize() int {
	return 20
}

// authListFields 权限列表查询字段
var authListFields string = "auth_id, auth_name, auth_con_act, auth_roles, auth_group_id, remark, prent_id, is_menu_show,create_time"

// GetAuthGroupOptions 获取分组选项
func GetAuthGroupOptions() ([]*AdminerAuthGroupOptions, error) {
	var authGroupOptions []*AdminerAuthGroupOptions

	// 先获取所有的key，并排序
	var keys []int
	for k := range AuthGroups {
		keys = append(keys, int(k)) // 转换为int以便排序
	}
	sort.Ints(keys) // 按照key的int值排序

	// 按照排序好的key的顺序来构造authGroupOptions
	for _, k := range keys {
		option := &AdminerAuthGroupOptions{
			Value: uint(k), // 再将int转换回uint
			Label: AuthGroups[uint(k)],
		}
		authGroupOptions = append(authGroupOptions, option)
	}
	return authGroupOptions, nil
}

// GetAdminerAuthStatusOptions 权限状态选项切片的指针
func GetAdminerAuthStatusOptions() []*AdminerAuthStatusOptions {
	return []*AdminerAuthStatusOptions{
		{Value: 1, Label: "禁用"},
		{Value: 2, Label: "启用"},
		{Value: 99, Label: "已删除"},
	}
}

// GetAdminerAuthIsMenuOptions  权限是否菜单选项切片的指针
func GetAdminerAuthIsMenuOptions() []*AdminerAuthIsMenuOptions {
	return []*AdminerAuthIsMenuOptions{
		{Value: 1, Label: "否"},
		{Value: 2, Label: "是"},
	}
}

// GetAdminerAuthPrentOptions  权限父级选项切片的指针
func GetAdminerAuthPrentOptions() ([]*AdminerPrentAuthOptions, error) {
	var authList []*AdminerAuth
	var prentOptions []*AdminerPrentAuthOptions

	// 查询所有父级权限
	err := CmsDb.Where("status != ? ", 99).Where("prent_id = ?", 0).Find(&authList).Error
	if err != nil {
		return nil, err
	}

	// 循环处理选项
	if len(authList) > 0 {
		for _, v := range authList {
			option := &AdminerPrentAuthOptions{
				Value: v.AuthId,
				Label: v.AuthName,
			}
			prentOptions = append(prentOptions, option)
		}
	}

	return prentOptions, nil
}

// GetAdminerAuthById 通过权限ID获取权限
func GetAdminerAuthById(authId uint) (*AdminerAuth, error) {
	var adminerAuth AdminerAuth
	searchRes := CmsDb.Where("auth_id = ?", authId).First(&adminerAuth, authId)
	if searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("权限不存在")
		} else {
			return nil, searchRes.Error
		}
	}
	return &adminerAuth, nil
}

// GetAdminerAuthByIds 通过权限ID列表获取权限
func GetAdminerAuthByIds(authIds []uint) (*[]AdminerAuth, error) {
	var auths []AdminerAuth
	if len(authIds) == 0 {
		return &auths, nil
	}

	result := CmsDb.Where("auth_id IN ?", authIds).Where("status = ?", 2).Find(&auths)
	if result.Error != nil {
		return nil, result.Error
	}

	return &auths, nil
}

// GetAdminerAuthByAuthGroupId 根据分组id查询权限列表
func GetAdminerAuthByAuthGroupId(groupId uint) ([]*AdminerAuth, error) {
	var adminerAuths []*AdminerAuth
	err := CmsDb.Model(&AdminerAuth{}).
		Select(authListFields).Where("auth_group_id = ?", groupId).Where("status = ?", 2).Offset(0).Limit(-1).Order("sort asc").Find(&adminerAuths).Error
	if err != nil {
		return adminerAuths, err
	}

	return adminerAuths, err
}

// GetAdminerAuthsByGroupId 根据分组id查询权限列表
func GetAdminerAuthsByGroupId(groupId uint) ([]*AdminerAuth, error) {
	var adminerAuths []*AdminerAuth

	// 获取分组权限列表
	adminerGroupAuths, err := GetGroupAuthsByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var authIds []uint
	for _, groupAuth := range adminerGroupAuths {
		authIds = append(authIds, groupAuth.AuthId)
	}
	if len(authIds) > 0 {
		result := CmsDb.Select(authListFields).Where("auth_id IN ?", authIds).Find(&adminerAuths)
		return adminerAuths, result.Error
	}

	return adminerAuths, nil
}

// GetAdminerAuths 查询所有权限列表
func GetAdminerAuths() ([]*AdminerAuth, error) {
	var adminerAuths []*AdminerAuth
	err := CmsDb.Model(&AdminerAuth{}).
		Select(authListFields).Where("status = ?", 2).Offset(0).Limit(-1).Order("sort asc").Find(&adminerAuths).Error
	if err != nil {
		return adminerAuths, err
	}

	return adminerAuths, err
}

// AdminerAuthList 权限列表
func AdminerAuthList(params AdminerAuthSearchParams, sort string) (AdminerAuthQueryResult, error) {
	// 初始化值
	var result AdminerAuthQueryResult             // 设置返回值
	offset := (params.Page - 1) * params.PageSize // 计算分页起始值
	if sort == "" {                               // 设置默认排序
		sort = "auth_id desc"
	}
	query := CmsDb.Debug().Model(&AdminerAuth{}) // 初始化模型

	// 处理条件
	if params.AuthName != "" {
		query = query.Where("auth_name LIKE ?", "%"+params.AuthName+"%")
	}
	if params.AuthConAct != "" {
		query = query.Where("auth_con_act = ?", params.AuthConAct)
	}
	if params.AuthGroupId != 0 {
		query = query.Where("auth_group_id = ?", params.AuthGroupId)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	} else {
		query = query.Where("status != ?", 99)
	}

	// 计算总数
	totalErr := query.Count(&result.Total).Error
	if totalErr != nil {
		return result, totalErr
	}

	// 处理分页
	query = query.Offset(offset).Limit(params.PageSize)
	// 处理排序
	query = query.Order(sort)

	// 执行查询
	err := query.Find(&result.AdminerAuths).Error
	if err != nil {
		return result, err
	}

	return result, err
}

// GetAdminerAuthsByRoles 根据角色查询权限列表
func GetAdminerAuthsByRoles(roles string) ([]*AdminerAuth, error) {
	// 分割 roles 字符串得到ID列表
	roleIDs := strings.Split(roles, ",")

	// 将字符串ID列表转换为uint切片
	var roleIDsUint []uint
	for _, idStr := range roleIDs {
		idUint, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			// 或者你可以选择忽略错误的ID
			// return nil, err
			continue
		}
		roleIDsUint = append(roleIDsUint, uint(idUint))
	}

	// 查询这些ID的 adminer_auth 记录
	var adminerAuths []*AdminerAuth
	result := CmsDb.Select(authListFields).Where("auth_id IN ?", roleIDsUint).Find(&adminerAuths)
	return adminerAuths, result.Error
}

// SaveAdminerAuth 保存权限
func SaveAdminerAuth(data AdminerAuth, authId uint) (bool, error) {
	if len(data.AuthName) < 1 {
		return false, fmt.Errorf("权限名称不能为空")
	}
	if data.PrentId > 0 && len(data.AuthConAct) < 1 {
		return false, fmt.Errorf("权限不能为空")
	}
	if data.AuthGroupId < 1 {
		return false, fmt.Errorf("请选择权限组")
	}
	var authRoles string
	if data.AuthConAct != "" {
		authRoles = strings.Replace(data.AuthRoles, "/", "", -1)
	}

	if authId < 1 {
		// 查询排序值
		sortNum := 0
		sortRes := CmsDb.Model(&AdminerAuth{}).Debug().Select("sort").Where("prent_id = ?", data.PrentId).Order("sort desc").Limit(1).Find(&sortNum)
		if sortRes.Error != nil {
			return false, fmt.Errorf("排序值发生未知错误：: %w", sortRes.Error)
		}

		if data.PrentId > 0 {
			sortNum = sortNum + 1
		} else {
			sortNum = sortNum + 100
		}

		// 添加分组
		adminerAuth := AdminerAuth{
			AuthName:      data.AuthName,
			AuthConAct:    data.AuthConAct,
			AuthRoles:     authRoles,
			AuthGroupId:   data.AuthGroupId,
			Status:        data.Status,
			PrentId:       data.PrentId,
			IsMenuShow:    data.IsMenuShow,
			Remark:        data.Remark,
			Sort:          uint(sortNum),
			CreateTime:    uint64(time.Now().Unix()),
			CreateAdminer: data.CreateAdminer,
			LastTime:      uint64(time.Now().Unix()),
			LastAdminer:   data.CreateAdminer,
		}

		if err := CmsDb.Debug().Create(&adminerAuth).Error; err != nil {
			return false, fmt.Errorf("添加权限时出错: %w", err)
		}
	} else {
		// 更新权限
		updates := make(map[string]interface{})
		updates["Status"] = data.Status
		updates["Remark"] = data.Remark
		updates["PrentId"] = data.PrentId
		updates["AuthName"] = data.AuthName
		updates["AuthRoles"] = authRoles
		updates["AuthConAct"] = data.AuthConAct
		updates["IsMenuShow"] = data.IsMenuShow
		updates["AuthGroupId"] = data.AuthGroupId
		updates["LastTime"] = uint64(time.Now().Unix())
		updates["LastAdminer"] = data.LastAdminer
		result := CmsDb.Model(&AdminerAuth{}).Where("auth_id = ?", authId).Updates(updates)
		if result.Error != nil {
			return false, fmt.Errorf("更新权限时出错: %w", result.Error)
		}
	}

	return true, nil
}

// SwitchAuthStatus 开关状态更新
func SwitchAuthStatus(statusForm AuthSwitchStatusForm) (bool, error) {
	if statusForm.AuthId < 1 {
		return false, fmt.Errorf("权限ID错误")
	}

	// 允许status的值
	invalidStatuses := map[uint]bool{
		1:  true,
		2:  true,
		99: true,
	}
	if _, exists := invalidStatuses[statusForm.Status]; !exists {
		return false, fmt.Errorf("非法请求：状态值不正")
	}

	// 查询权限是否存在
	auth, _ := GetAdminerAuthById(statusForm.AuthId)
	if auth == nil || auth.AuthId < 1 {
		return false, fmt.Errorf("权限不存在")
	}

	// 更新分组
	updates := make(map[string]interface{})
	updates["Status"] = statusForm.Status
	updates["LastTime"] = time.Now().Unix()
	updates["LastAdminer"] = statusForm.LastAdminer
	result := CmsDb.Model(&AdminerAuth{}).Where("auth_id = ?", statusForm.AuthId).Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("状态更新出错: %w", result.Error)
	}

	return true, nil
}

// SwitchAuthMenuShow 菜单显示状态更新
func SwitchAuthMenuShow(menuForm AuthSwitchMenuForm) (bool, error) {
	if menuForm.AuthId < 1 {
		return false, fmt.Errorf("权限ID错误")
	}

	// 允许status的值
	invalidShows := map[uint]bool{
		1: true,
		2: true,
	}
	if _, exists := invalidShows[menuForm.IsMenuShow]; !exists {
		return false, fmt.Errorf("非法请求：状态值不正")
	}

	// 查询权限是否存在
	auth, _ := GetAdminerAuthById(menuForm.AuthId)
	if auth == nil || auth.AuthId < 1 {
		return false, fmt.Errorf("权限不存在")
	}

	// 更新分组
	updates := make(map[string]interface{})
	updates["IsMenuShow"] = menuForm.IsMenuShow
	updates["LastTime"] = time.Now().Unix()
	updates["LastAdminer"] = menuForm.LastAdminer
	result := CmsDb.Model(&AdminerAuth{}).Where("auth_id = ?", menuForm.AuthId).Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("状态更新出错: %w", result.Error)
	}

	return true, nil
}

// DeleteAdminerAuth 删除管理员权限
func DeleteAdminerAuth(authId uint, adminerId uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	// 判断参数是否合法
	if authId < 1 {
		return res, fmt.Errorf("权限ID不能为空")
	}

	// 查询权限是否存在
	auth, dErr := GetAdminerAuthById(authId)
	if dErr != nil {
		return res, fmt.Errorf("删除出错: %w", dErr.Error)
	}
	if auth.Status == 99 {
		return res, fmt.Errorf("删除出错: %w", "当前状态已是删除状态无需再次删除")
	}

	// 更新分组
	updates := make(map[string]interface{})
	updates["Status"] = 99
	updates["LastAdminer"] = adminerId
	updates["LastTime"] = time.Now().Unix()
	result := CmsDb.Model(&AdminerAuth{}).Where("auth_id = ?", authId).Updates(updates)
	if result.Error != nil {
		return res, fmt.Errorf("删除权限出错: %w", result.Error)
	}

	res["auth_id"] = authId
	res["affected_rows"] = result.RowsAffected
	return res, nil
}

// BuildAuthTree 接收所有权限和用户拥有的权限列表
func BuildAuthTree(allAuths []*AdminerAuth, userAuths []*AdminerGroupAuth) []*GroupWithAuths {
	// 1. 创建一个映射来跟踪用户拥有的权限
	userAuthsMap := make(map[uint]bool)
	for _, ua := range userAuths {
		userAuthsMap[ua.AuthId] = true
	}

	// 2. 创建权限ID到AuthWithChildren对象的映射
	authMap := make(map[uint]*AuthWithChildren)
	for _, aa := range allAuths {
		isCheck := uint(1)
		if _, found := userAuthsMap[aa.AuthId]; found {
			isCheck = 2
		}
		authMap[aa.AuthId] = &AuthWithChildren{
			AuthId:   aa.AuthId,
			AuthName: aa.AuthName,
			IsCheck:  isCheck,
			Children: nil,
		}
	}

	// 3. 构建权限的树形结构
	for _, aa := range allAuths {
		child, _ := authMap[aa.AuthId]
		if aa.PrentId == 0 {
			// 对于顶级权限，我们只构建对应权限组的树结构
			if aa.AuthGroupId != 1 { // 跳过“无所不能”
				continue
			}
		} else {
			// 对于非顶级权限，添加到父权限的children中
			parent, found := authMap[aa.PrentId]
			if found {
				parent.Children = append(parent.Children, child)
			}
		}
	}

	// 4. 将权限按组组织
	groupWithAuthsList := make([]*GroupWithAuths, 0, len(AuthGroups))
	for groupId, groupName := range AuthGroups {
		if groupId == 1 { // 跳过“无所不能”
			continue
		}
		groupWithAuths := &GroupWithAuths{
			AuthGroupId:   groupId,
			AuthGroupName: groupName,
			AuthName:      groupName,
			Children:      make([]*AuthWithChildren, 0),
		}
		for _, aa := range allAuths {
			if aa.AuthGroupId == groupId && aa.PrentId == 0 {
				groupWithAuths.Children = append(groupWithAuths.Children, authMap[aa.AuthId])
			}
		}
		groupWithAuthsList = append(groupWithAuthsList, groupWithAuths)
	}

	// 5. 更新每个组的子权限的is_check状态
	for _, groupWithAuths := range groupWithAuthsList {
		for _, auth := range groupWithAuths.Children {
			updateIsCheckStatus(auth)
		}
	}

	return groupWithAuthsList
}

// 更新权限的is_check状态，递归到第三层
func updateIsCheckStatus(auth *AuthWithChildren) {
	// 第三层权限的状态直接由用户的权限决定
	for _, child := range auth.Children {
		// 如果存在子权限，则为第二层，需要根据子权限状态更新is_check
		if len(child.Children) > 0 {
			allChecked := true
			for _, grandchild := range child.Children {
				if grandchild.IsCheck != 2 {
					allChecked = false
					break
				}
			}
			if allChecked {
				child.IsCheck = 2
			} else {
				child.IsCheck = 3
			}
		}
	}
}
