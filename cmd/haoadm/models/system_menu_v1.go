package models

// 初始化结构体
type SystemV1Init struct {
	ClearInfo struct {
		ClearUrl string `json:"clearUrl"`
	} `json:"clearInfo"`
	HomeInfo struct {
		Title string `json:"title"`
		Href  string `json:"href"`
		Icon  string `json:"icon"`
	} `json:"homeInfo"`
	LogoInfo struct {
		Title string `json:"title"`
		Image string `json:"image"`
		Href  string `json:"href"`
	} `json:"logoInfo"`
	MenuInfo []*MenuTreeList `json:"menuInfo"`
}

// 获取初始化数据
func (m *SystemMenu) GetV1SystemInit() SystemV1Init {
	var init SystemV1Init

	// 首页
	init.HomeInfo.Title = "首页"
	init.HomeInfo.Href = "/admin/"
	init.HomeInfo.Icon = "fa fa-home"

	init.ClearInfo.ClearUrl = "/admin/api/clear.json"

	// logo
	init.LogoInfo.Title = "HaoTrader"
	init.LogoInfo.Image = "/admin/images/logo.png"

	// 菜单
	init.MenuInfo = m.GetV1MenuList()

	return init
}

func (m *SystemMenu) GetV1MenuList() []*MenuTreeList {
	return m.buildMenuChild(0, menuList)
}
