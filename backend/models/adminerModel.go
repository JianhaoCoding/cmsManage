package models

import (
	"cms/cache"
	conf "cms/config"
	"cms/helpers"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var tableName string = "adminer"
var pageSize int = 20

type Adminer struct {
	AdminerId  uint         `gorm:"primaryKey;autoIncrement" json:"adminer_id"`
	Username   string       `json:"username"`
	Password   string       `json:"-"`
	Nickname   string       `json:"nickname"`
	Email      string       `json:"email"`
	Mobile     string       `json:"mobile"`
	GroupId    uint         `json:"group_id"`
	Group      AdminerGroup `gorm:"foreignKey:GroupId;references:GroupId" json:"Group"`
	Status     uint         `json:"status"`
	Remark     string       `json:"remark"`
	LastIp     string       `json:"last_ip"`
	LastTime   uint64       `json:"last_time"`
	CreateTime uint64       `json:"create_time"`
	UpdateTime uint64       `json:"-"`
}

type AdminerSearchParams struct {
	Username string `form:"username"`
	Mobile   string `form:"mobile"`
	GroupId  uint   `form:"group_id"`
	Status   uint   `form:"status"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type AdminerQueryResult struct {
	Adminers []*Adminer
	Total    int64
}

type AdminerStatusOptions struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

func (Adminer) TableNameDefine() string {
	return tableName
}

func (Adminer) PageSizeDefine() int {
	return pageSize
}

func (Adminer) TableName() string {
	return tableName
}

func (Adminer) PageSize() int {
	return 20
}

// GetAdminerById 通过管理员ID获取管理员
func GetAdminerById(adminerId uint) (*Adminer, error) {
	var adminer Adminer
	searchRes := CmsDb.First(&adminer, adminerId)
	if searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("获取资源失败: %w", searchRes.Error)
		} else {
			return nil, searchRes.Error
		}
	}
	return &adminer, nil
}

// GetAdminerByName 通过用户名获取管理员
func GetAdminerByName(name string) (*Adminer, error) {
	var adminer Adminer
	searchRes := CmsDb.Where("username = ?", name).Where("status != ?", 99).First(&adminer)
	if searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("获取资源失败: %w", searchRes.Error)
		} else {
			return nil, searchRes.Error
		}
	}
	return &adminer, nil
}

// GetAdminerByMobile 通过手机号获取管理员
func GetAdminerByMobile(mobile uint64) (*Adminer, error) {
	var adminer Adminer
	searchRes := CmsDb.Where("mobile = ?", mobile).Where("status != ?", 99).First(&adminer)
	if searchRes.Error != nil {
		if errors.Is(searchRes.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("获取资源失败: %w", searchRes.Error)
		} else {
			return nil, searchRes.Error
		}
	}
	return &adminer, nil
}

// AdminerList 管理员列表
func AdminerList(params AdminerSearchParams, page int, pageSize int, sort string) (AdminerQueryResult, error) {
	// 初始化值
	var result AdminerQueryResult // 设置返回值
	var adObj Adminer
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = adObj.PageSize()
	}
	offset := (page - 1) * pageSize // 计算分页起始值
	if sort == "" {                 // 设置默认排序
		sort = "adminer_id asc"
	}
	query := CmsDb.Debug().Model(&Adminer{}) // 初始化模型

	// 处理条件
	if params.Username != "" {
		query = query.Where("username = ?", params.Username)
	}
	if params.Mobile != "" {
		query = query.Where("mobile = ?", params.Mobile)
	}
	if params.GroupId != 0 {
		query = query.Where("group_id = ?", params.GroupId)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	} else {
		query = query.Where("status != ?", 99)
	}

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
	err := query.Preload("Group").Find(&result.Adminers).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// AdminerLogin 管理员登录
func AdminerLogin(username string, password string) (map[string]interface{}, error) {
	// 根据用户名获取管理员
	adminer, err := GetAdminerByName(username)
	if err != nil {
		return nil, err
	}

	// 效验状态是否正常
	if adminer.Status != 2 {
		if adminer.Status == 1 {
			return nil, fmt.Errorf("抱歉您的账号状态不正常，请联系管理员")
		} else {
			return nil, fmt.Errorf("该账号不存在，请联系管理员新增管理账户")
		}
	}

	// 效验密码是否正确
	if adminer.Password != helpers.MD5Hash(password) {
		return nil, fmt.Errorf("密码错误请重写输入")
	}

	// 登录成功，创建token
	res := make(map[string]interface{})
	res["adminer_id"] = adminer.AdminerId
	res["username"] = adminer.Username
	res["nickname"] = adminer.Nickname
	res["email"] = adminer.Email
	res["mobile"] = adminer.Mobile
	res["group_id"] = adminer.GroupId
	token, _ := helpers.GenerateToken(adminer.AdminerId)
	res["token"] = token
	// Bearer token 前端要拼接Bearer

	// 登录成功存储redis，以token为键，AdminerId为值
	err = cache.RedisClient.Set(cache.Ctx, token, adminer.AdminerId, 2*conf.CacheTimeOneHour).Err()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return res, nil
}

// UpdateAdminerLast 更新用户最新登录信息
func UpdateAdminerLast(adminerId uint, loginIp string, loginTime uint64) (bool, error) {
	// 管理员是否存在
	adminer, resErr := GetAdminerById(adminerId)
	if resErr != nil || adminer == nil {
		return false, fmt.Errorf("更新出错: %w", resErr)
	}

	// 创建映射对象
	updates := make(map[string]interface{})
	updates["LastIp"] = loginIp
	updates["LastTime"] = loginTime
	result := CmsDb.Model(&Adminer{}).Where("adminer_id = ?", adminerId).Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("更新管理员时出错: %w", result.Error)
	}
	return true, nil
}

// SaveAdminer 保存管理员
func SaveAdminer(data Adminer, adminerId uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	if data.Username < "" {
		return res, fmt.Errorf("用户名称不能为空")
	}

	if data.Nickname == "" {
		data.Nickname = data.Username
	}

	if data.GroupId == 0 {
		return res, fmt.Errorf("请选择用户组")
	}

	var password string
	if adminerId < 1 {
		// 生成密码
		password = helpers.GenerateRandomPassword(6, 15)

		// 添加用户
		adminer := Adminer{
			Username:   data.Username,
			Password:   helpers.MD5Hash(password),
			Nickname:   data.Nickname,
			Email:      data.Email,
			Mobile:     data.Mobile,
			GroupId:    data.GroupId,
			Status:     data.Status,
			Remark:     data.Remark,
			CreateTime: uint64(time.Now().Unix()),
		}

		result := CmsDb.Create(&adminer)
		if result.Error != nil {
			return res, fmt.Errorf("添加管理员时出错: %w", result.Error)
		}
		res["adminer_id"] = adminer.AdminerId
		res["password"] = password
		res["affected_rows"] = result.RowsAffected
	} else {
		// 查询管理员是否存在
		adminer, aErr := GetAdminerById(adminerId)
		if aErr != nil {
			return res, fmt.Errorf("添加管理员时出错: %w", aErr.Error)
		}

		// 更新分组
		updates := make(map[string]interface{})
		updates["Nickname"] = data.Nickname
		updates["Email"] = data.Email
		updates["Mobile"] = data.Mobile
		updates["GroupId"] = data.GroupId
		updates["Status"] = data.Status
		updates["Remark"] = data.Remark
		updates["UpdateTime"] = time.Now().Unix()
		result := CmsDb.Model(&Adminer{}).Where("adminer_id = ?", adminerId).Updates(updates)
		if result.Error != nil {
			return res, fmt.Errorf("更新管理员时出错: %w", result.Error)
		}

		res["adminer_id"] = adminer.AdminerId
		res["password"] = adminer.Password
		res["affected_rows"] = result.RowsAffected
	}

	return res, nil
}

// AdminerStatusSwitch 状态开关
func AdminerStatusSwitch(adminerId uint, status uint) (bool, error) {
	// 判断参数是否合法
	if adminerId < 1 {
		return false, fmt.Errorf("管理员ID不能为空")
	}
	// 允许status的值
	invalidStatuses := map[uint]bool{
		1:  true,
		2:  true,
		99: true,
	}
	if _, exists := invalidStatuses[status]; !exists {
		return false, fmt.Errorf("非法请求：状态值不正")
	}
	// 超管状态不能更新
	if adminerId == 1 {
		return false, fmt.Errorf("超级管理员状态不允许被操作！")
	}

	// 更新分组
	updates := make(map[string]interface{})
	updates["Status"] = status
	updates["UpdateTime"] = time.Now().Unix()
	result := CmsDb.Model(&Adminer{}).Where("adminer_id = ?", adminerId).Updates(updates)
	if result.Error != nil {
		return false, fmt.Errorf("状态更新出错: %w", result.Error)
	}

	return true, nil
}

// DeleteAdminer 删除管理员
func DeleteAdminer(adminerId uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	// 判断参数是否合法
	if adminerId < 1 {
		return res, fmt.Errorf("管理员ID不能为空")
	}

	// 超管状态不能更新
	if adminerId == 1 {
		return res, fmt.Errorf("超级管理员状态不允许被操作！")
	}

	// 查询管理员是否存在
	adminer, dErr := GetAdminerById(adminerId)
	if dErr != nil {
		return res, fmt.Errorf("删除出错: %w", dErr.Error)
	}
	if adminer.Status == 99 {
		return res, fmt.Errorf("删除出错: %w", "当前状态已是删除状态无需再次删除")
	}

	// 更新分组
	updates := make(map[string]interface{})
	updates["Status"] = 99
	updates["UpdateTime"] = time.Now().Unix()
	result := CmsDb.Model(&Adminer{}).Where("adminer_id = ?", adminerId).Updates(updates)
	if result.Error != nil {
		return res, fmt.Errorf("删除管理员出错: %w", result.Error)
	}

	res["adminer_id"] = adminer.AdminerId
	res["affected_rows"] = result.RowsAffected
	return res, nil
}

// ResetPass 重置管理员密码
func ResetPass(adminerId uint, operatorId uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	// 判断参数是否合法
	if adminerId < 1 {
		return res, fmt.Errorf("管理员ID不能为空")
	}

	// 超管状态不能更新
	if adminerId == 1 && operatorId != adminerId {
		return res, fmt.Errorf("您无权操作超管信息！")
	}

	// 查询管理员是否存在
	adminer, aErr := GetAdminerById(adminerId)
	if aErr != nil {
		return res, fmt.Errorf("重置密码时出错: %w", aErr.Error)
	}

	// 更新分组
	password := helpers.GenerateRandomPassword(6, 15)
	updates := make(map[string]interface{})
	updates["Password"] = helpers.MD5Hash(password)
	updates["UpdateTime"] = time.Now().Unix()
	result := CmsDb.Model(&Adminer{}).Where("adminer_id = ?", adminerId).Updates(updates)
	if result.Error != nil {
		return res, fmt.Errorf("删除管理员出错: %w", result.Error)
	}

	res["adminer_id"] = adminer.AdminerId
	res["password"] = password
	res["affected_rows"] = result.RowsAffected
	return res, nil
}

// GetAdminerStatusOptions 返回指向管理员状态选项切片的指针
func GetAdminerStatusOptions() []*AdminerStatusOptions {
	return []*AdminerStatusOptions{
		{Value: 1, Label: "禁用"},
		{Value: 2, Label: "启用"},
	}
}
