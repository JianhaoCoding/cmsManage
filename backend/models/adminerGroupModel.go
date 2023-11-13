package models

import (
	"cms/helpers"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// AdminerGroup
// `gorm:"primaryKey;autoIncrement;column:group_id" json:"group_id"`
// 如果数据库中该列的名称与 Go 结构体中的字段名称不匹配，你可以使用 column 标签来指定数据库中实际的列名称
type AdminerGroup struct {
	GroupId       uint   `gorm:"primaryKey;autoIncrement" json:"group_id"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Remark        string `json:"remark"`
	CreateTime    uint64 `json:"create_time"`
	CreateAdminer uint64 `json:"-"`
}

type AdminerGroupQueryResult struct {
	AdminerGroups []*AdminerGroup
	Total         int64
}

type AdminerGroupSearchParam struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type UpdateGroupData struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Remark string `json:"remark"`
}

type AdminerGroupForm struct {
	GroupId       *uint  `json:"group_id"` // 使用指针允许空值
	Name          string `json:"name" binding:"required"`
	Remark        string `json:"remark"`
	AuthIDs       []uint `json:"auth_ids"`
	CreateAdminer uint64 `json:"-"`
}

type GroupOption struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

func (AdminerGroup) TableName() string {
	return "adminer_group"
}

func (AdminerGroup) PageSize() int {
	return 20
}

// GetAdminerGroupById 通过管理员组获取管理组
func GetAdminerGroupById(groupId uint) (*AdminerGroup, error) {
	var adminerGroup AdminerGroup
	searchRes := CmsDb.First(&adminerGroup, groupId)
	if searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("管理组不存在")
		} else {
			return nil, searchRes.Error
		}
	}
	return &adminerGroup, nil
}

// AdminerGroupList 管理组列表
func AdminerGroupList(page int, pageSize int, sort string) (AdminerGroupQueryResult, error) {
	// 初始化值
	var result AdminerGroupQueryResult // 设置返回值
	offset := (page - 1) * pageSize    // 计算分页起始值
	if sort == "" {                    // 设置默认排序
		sort = "group_id asc"
	}
	query := CmsDb.Debug().Model(&AdminerGroup{}) // 初始化模型

	// 获取总记录数
	totalErr := query.Count(&result.Total).Error
	if totalErr != nil {
		return result, totalErr
	}

	// 处理分页
	query = query.Offset(offset).Limit(pageSize)
	// 处理排序
	query = query.Order(sort)

	// 执行查询
	err := query.Find(&result.AdminerGroups).Error
	if err != nil {
		return result, err
	}

	return result, err
}

// UpdateAdminerGroup 更新用户组
func UpdateAdminerGroup(groupId uint, data UpdateGroupData) (bool, error) {
	// 管理组是否存在
	group, resErr := GetAdminerGroupById(groupId)
	if resErr != nil || group == nil {
		return false, fmt.Errorf("更新出错: %w", resErr)
	}

	// 创建映射对象
	updates := make(map[string]interface{})
	updates["Name"] = data.Name
	updates["Role"] = data.Role
	updates["Remark"] = data.Remark
	result := CmsDb.Model(&AdminerGroup{}).Where("group_id = ?", groupId).Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("更新管理组时出错: %w", result.Error)
	}
	return true, nil
}

// SaveAdminerGroup 保存管理组
func SaveAdminerGroup(data AdminerGroup, groupId uint) (bool, error) {
	if len(data.Name) < 1 {
		return false, fmt.Errorf("组名称不能为空")
	}

	if groupId < 1 {
		// 添加分组
		adminerGroup := AdminerGroup{
			Name:          data.Name,
			Role:          data.Role,
			Remark:        data.Remark,
			CreateTime:    uint64(time.Now().Unix()),
			CreateAdminer: data.CreateAdminer,
		}

		if err := CmsDb.Create(&adminerGroup).Error; err != nil {
			return false, fmt.Errorf("添加管理组时出错: %w", err)
		}
	} else {
		// 查询管理员组是否存在
		group, gErr := GetAdminerGroupById(groupId)
		if gErr != nil {
			return false, fmt.Errorf("添加管理组时出错: %w", gErr.Error)
		}

		// 更新分组
		updates := make(map[string]interface{})
		updates["Name"] = data.Name
		updates["Remark"] = data.Remark
		result := CmsDb.Model(&AdminerGroup{}).Where("group_id = ?", group.GroupId).Updates(updates)
		if result.Error != nil {
			return false, fmt.Errorf("更新管理组时出错: %w", result.Error)
		}
		if result.RowsAffected < 1 {
			return false, fmt.Errorf("无更新内容")
		}
	}

	return true, nil
}

// SaveGroupForm 保存分组数据
func SaveGroupForm(formData AdminerGroupForm) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var groupId uint

	// 添加保存用户组信息
	if formData.GroupId != nil && *formData.GroupId > 0 {
		// 更新操作
		updates := make(map[string]interface{})
		updates["Name"] = formData.Name
		updates["Remark"] = formData.Remark
		upRes := CmsDb.Model(&AdminerGroup{}).Where("group_id = ?", *formData.GroupId).Updates(updates)
		if upRes.Error != nil {
			return result, fmt.Errorf("更新管理组时出错: %w", upRes.Error)
		}
		groupId = *formData.GroupId
	} else {
		// 添加操作
		adminerGroup := AdminerGroup{
			Name:          formData.Name,
			Role:          "",
			Remark:        formData.Remark,
			CreateTime:    uint64(time.Now().Unix()),
			CreateAdminer: formData.CreateAdminer,
		}
		addRes := CmsDb.Create(&adminerGroup)
		if addRes.Error != nil {
			return result, fmt.Errorf("添加管理组时出错: %w", addRes.Error)
		}
		groupId = adminerGroup.GroupId
	}
	result["group_id"] = groupId

	// 去除权限列表的空值
	authids := helpers.SliceUnique(formData.AuthIDs)
	fmt.Println(authids)
	if len(authids) > 0 {
		if authErr := AddGroupAuths(authids, groupId, formData.CreateAdminer); authErr != nil {
			return result, authErr
		}
	}

	return result, nil
}

// GetGroupOptions 获取分组选项
func GetGroupOptions() ([]*GroupOption, error) {
	var groups []AdminerGroup
	var groupOptions []*GroupOption

	// 从数据库获取所有管理员组
	if err := CmsDb.Order("group_id asc").Find(&groups).Error; err != nil {
		return nil, err
	}

	// 将管理员组转换为 GroupOption 列表
	for _, group := range groups {
		option := &GroupOption{
			Value: group.GroupId,
			Label: group.Name,
		}
		groupOptions = append(groupOptions, option)
	}
	return groupOptions, nil
}
