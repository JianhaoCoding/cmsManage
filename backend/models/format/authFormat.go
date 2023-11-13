package format

// AuthItem 表示原始数据项
type AuthItem struct {
	AuthId      uint   `json:"auth_id"`
	AuthName    string `json:"auth_name"`
	AuthConAct  string `json:"auth_con_act"`
	AuthGroupId uint   `json:"auth_group_id"`
	Status      uint   `json:"status"`
	Remark      string `json:"remark"`
	PrentId     uint   `json:"prent_id"`
	IsMenuShow  uint   `json:"is_menu_show"`
}

// FormattedAuthItem 表示格式化后的数据项，包括子项
type FormattedAuthItem struct {
	AuthId     uint                `json:"auth_id"`
	AuthName   string              `json:"auth_name"`
	Path       string              `json:"path"`
	Status     uint                `json:"status"`
	IsMenuShow uint                `json:"is_menu_show"`
	Children   []FormattedAuthItem `json:"children,omitempty"` //omitempty 如果children为空则在json中不显示
}

// BuildAuthTree 构建权限树
func BuildAuthTree(items []AuthItem) []FormattedAuthItem {
	childrenMap := make(map[uint][]AuthItem)
	for _, item := range items {
		childrenMap[item.PrentId] = append(childrenMap[item.PrentId], item)
	}

	var buildTree func(parentID uint) []FormattedAuthItem
	buildTree = func(parentID uint) []FormattedAuthItem {
		var tree []FormattedAuthItem
		for _, item := range childrenMap[parentID] {
			tree = append(tree, FormattedAuthItem{
				AuthId:     item.AuthId,
				AuthName:   item.AuthName,
				Path:       item.AuthConAct,
				Status:     item.Status,
				IsMenuShow: item.IsMenuShow,
				Children:   buildTree(item.AuthId), // 递归构建子树
			})
		}
		return tree
	}
	// 假设顶级节点的parent_id为0
	return buildTree(0)
}
