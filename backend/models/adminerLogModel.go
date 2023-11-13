package models

type AdminerLog struct {
	LogId          uint   `gorm:"primary_key" json:"log_id"`
	ObjId          uint   `json:"obj_id"`
	ObjType        string `json:"obj_type"`
	OperateType    string `json:"operate_type"`
	Sql            string `json:"sql"`
	OperateTime    uint64 `json:"operate_time"`
	OperateIp      string `json:"operate_ip"`
	OperateAdminer uint   `json:"operate_adminer"`
}

type AdminerLogQueryResult struct {
	AdminerLogs []*AdminerLog
	Total       int64
}

type AdminerLogForm struct {
	ObjId          uint   `form:"obj_id"`
	OperateType    string `form:"operate_type"`
	ObjType        string `form:"obj_type"`
	OperateAdminer string `form:"operate_adminer"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

func (AdminerLog) TableName() string {
	return "adminer_log"
}

// AdminerLogList 管理员日志列表
func AdminerLogList(form AdminerLogForm, sort string) (AdminerLogQueryResult, error) {
	var result AdminerLogQueryResult

	// 计算分页起始值
	offset := (form.Page - 1) * form.PageSize
	// 设置默认排序
	if sort == "" {
		sort = "log_id desc"
	}
	// 初始化模型
	query := CmsDb.Debug().Model(&AdminerLog{})

	// 组装条件
	if form.ObjId != 0 {
		query = query.Where("obj_id = ?", form.ObjId)
	}

	if form.OperateType != "" {
		query = query.Where("operate_type = ?", form.OperateType)
	}

	if form.ObjType != "" {
		query = query.Where("obj_type = ?", form.ObjType)
	}

	if form.OperateAdminer != "" {
		// 查询管理员
		adminer, err := GetAdminerByName(form.OperateAdminer)
		if err != nil {
			return result, err
		}
		query = query.Where("operate_adminer = ?", adminer.AdminerId)
	}

	// 计算总数
	totalErr := query.Count(&result.Total).Error
	if totalErr != nil {
		return result, totalErr
	}

	// 处理分页
	query = query.Offset(offset).Limit(form.PageSize)
	// 处理排序
	query = query.Order(sort)

	// 执行查询
	err := query.Find(&result.AdminerLogs).Error
	if err != nil {
		return result, err
	}

	return result, err
}

// addOperateLog 添加操作日志
func addOperateLog(objId uint, objType string, opreateType string, opreateSql string) {

}
