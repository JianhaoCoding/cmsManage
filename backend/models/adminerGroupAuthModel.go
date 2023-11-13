package models

import (
	"cms/helpers"
	"time"
)

type AdminerGroupAuth struct {
	AdminerGroupId uint   `gorm:"primary_key" json:"adminer_group_id"`
	AuthId         uint   `json:"auth_id"`
	CreateTime     uint64 `json:"create_time"`
	CreateAdminer  uint64 `json:"-"`
}

func (AdminerGroupAuth) TableName() string {
	return "adminer_group_auth"
}

// GetGroupAuthsByGroupId 通过groupId获取组权限列表
func GetGroupAuthsByGroupId(groupId uint) ([]*AdminerGroupAuth, error) {
	var adminerGroupAuths []*AdminerGroupAuth
	err := CmsDb.Model(&AdminerGroupAuth{}).Where("adminer_group_id = ?", groupId).Offset(0).Limit(-1).Find(&adminerGroupAuths).Error
	if err != nil {
		return adminerGroupAuths, err
	}

	return adminerGroupAuths, err
}

// AddGroupAuths 添加组权限
func AddGroupAuths(authIDs []uint, groupId uint, createAdminer uint64) error {
	// 查询权限列表
	newAuths, nErr := GetAdminerAuthByIds(authIDs)
	if nErr != nil {
		return nErr
	}

	// 循环出父级列表
	var prentAuthIds []uint
	seen := make(map[uint]bool)
	for _, v := range *newAuths {
		if v.PrentId > 0 && !seen[v.PrentId] {
			seen[v.PrentId] = true
			prentAuthIds = append(prentAuthIds, v.PrentId)
		}
	}

	// 合并切片
	authIDs = append(authIDs, prentAuthIds...)
	authIDs = helpers.SliceUnique(authIDs)

	// 查询用户组历史权限
	var oldAuths []AdminerGroupAuth
	var oldAuthIds []uint
	oldSeen := make(map[uint]bool)
	err := CmsDb.Where("adminer_group_id = ?", groupId).Find(&oldAuths).Error
	if err != nil {
		return err
	}
	for _, v := range oldAuths {
		if !oldSeen[v.AuthId] {
			oldSeen[v.AuthId] = true
			oldAuthIds = append(oldAuthIds, v.AuthId)
		}
	}

	// 分析新增权限和删除权限
	var addGroupAuths []AdminerGroupAuth
	if len(oldAuths) == 0 {
		// 批量添加权限
		for _, v := range authIDs {
			addGroupAuths = append(addGroupAuths, AdminerGroupAuth{
				AdminerGroupId: groupId,
				AuthId:         v,
				CreateTime:     uint64(time.Now().Unix()),
				CreateAdminer:  createAdminer,
			})
		}
	} else {
		for _, v := range authIDs {
			if !oldSeen[v] {
				addGroupAuths = append(addGroupAuths, AdminerGroupAuth{
					AdminerGroupId: groupId,
					AuthId:         v,
					CreateTime:     uint64(time.Now().Unix()),
					CreateAdminer:  createAdminer,
				})
			} else {
				oldAuthIds = removeAuthId(v, oldAuthIds)
			}
		}
	}

	// 删除权限
	if len(oldAuthIds) > 0 {
		CmsDb.Unscoped().Where("adminer_group_id = ", groupId).Where("auth_id IN ?", oldAuthIds).Delete(&AdminerGroupAuth{})
	}

	// 新增权限
	if len(addGroupAuths) > 0 {
		CmsDb.Create(&addGroupAuths)
	}

	return nil
}

func removeAuthId(authID uint, delAuths []uint) []uint {
	for i, id := range delAuths {
		if id == authID {
			return append(delAuths[:i], delAuths[i+1:]...)
		}
	}
	return delAuths
}
