package models

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"time"
)

type SystemUserRole int

const (
	SystemUserRoleRoot  SystemUserRole = -1
	SystemUserRoleSuper SystemUserRole = 0
	SystemUserRoleAdmin SystemUserRole = 1
	SystemUserRoleUser  SystemUserRole = 1
)

type SystemUser struct {
	Id        int64          `xorm:"'id' autoincr pk"`
	Name      string         `xorm:"varchar(100) unique"`
	RealName  string         `xorm:"varchar(100)"`
	Password  string         `xorm:"varchar(65)"`
	Role      SystemUserRole `xorm:""`
	Status    uint8          `xorm:"default(1)"`
	Remark    string         `xorm:"varchar(255)"`
	CreatedAt time.Time      `xorm:"timestamp created"`
	UpdatedAt time.Time      `xorm:"timestamp updated"`
}

func SystemUserPassword(pwd string) string {
	sha1Hash := sha1.Sum([]byte(pwd))
	// Calculate the MD5 hash of the SHA1 hash
	md5Hash := md5.Sum(sha1Hash[:])
	// Convert the hashes to hexadecimal strings
	return hex.EncodeToString(md5Hash[:])
}

// http://layuimini.99php.cn/docs/init/sql.html
type SystemMenu struct {
	Id        int64     `xorm:"'id' autoincr pk"`
	Pid       int64     `xorm:"'pid'"`
	Title     string    `xorm:"'title' index varchar(100)"`
	Icon      string    `xorm:"'icon' varchar(100)"`
	Href      string    `xorm:"'href' index varchar(100)"`
	Target    string    `xorm:"'target' varchar(20)"`
	Sort      int       `xorm:"'sort'"`
	Status    uint8     `xorm:"'status'"`
	Remark    string    `xorm:"'remark' varchar(255)"`
	CreatedAt time.Time `xorm:"timestamp created"`
	UpdatedAt time.Time `xorm:"timestamp updated"`
	DeletedAt time.Time `xorm:""`
}

// 初始化结构体
type SystemInit struct {
	HomeInfo struct {
		Title string `json:"title"`
		Href  string `json:"href"`
	} `json:"homeInfo"`
	LogoInfo struct {
		Title string `json:"title"`
		Image string `json:"image"`
	} `json:"logoInfo"`
	MenuInfo []*MenuTreeList `json:"menuInfo"`
}

// 菜单结构体
type MenuTreeList struct {
	Id     int64           `json:"id"`
	Pid    int64           `json:"pid"`
	Title  string          `json:"title"`
	Icon   string          `json:"icon"`
	Href   string          `json:"href"`
	Target string          `json:"target"`
	Remark string          `json:"remark"`
	Child  []*MenuTreeList `json:"child"`
}

// 获取初始化数据
func (m *SystemMenu) GetSystemInit() SystemInit {
	var systemInit SystemInit

	// 首页
	systemInit.HomeInfo.Title = "首页"
	systemInit.HomeInfo.Href = "/admin/"

	// logo
	systemInit.LogoInfo.Title = "HaoTrader"
	systemInit.LogoInfo.Image = "/admin/images/logo.png"

	// 菜单
	systemInit.MenuInfo = m.GetMenuList()

	return systemInit
}

// 获取菜单列表
func (m *SystemMenu) GetMenuList() []*MenuTreeList {
	var menuList []SystemMenu

	// db.Table(new(SystemMenu)).Where("status=?", 1).OrderBy("sort asc").Find(&menuList)

	if len(menuList) == 0 {
		menuList = []SystemMenu{
			// SystemMenu{Id: 1, Pid: 0, Title: "系统", Icon: "fa fa-wrench", Href: "javascript:;", Target: "_self"},
			SystemMenu{Id: 2, Pid: 0, Title: "交易", Icon: "fa fa-usd", Href: "javascript:;", Target: "_self"},
			SystemMenu{Id: 3, Pid: 0, Title: "用户", Icon: "fa fa-user", Href: "javascript:;", Target: "_self"},
			SystemMenu{Id: 4, Pid: 0, Title: "统计", Icon: "fa fa-bar-chart", Href: "javascript:;", Target: "_self"},

			SystemMenu{Id: 10, Pid: 1, Title: "系统设置", Icon: "fa fa-wrench", Href: "/admin/system/settings", Target: "_self"},

			SystemMenu{Id: 20, Pid: 2, Title: "资产种类", Icon: "fa fa-btc", Href: "/admin/varieties/list", Target: "_self"},
			SystemMenu{Id: 21, Pid: 2, Title: "板块分类", Icon: "fa fa-file-text-o", Href: "/admin/symbols/category", Target: "_self"},
			SystemMenu{Id: 22, Pid: 2, Title: "交易列表", Icon: "fa fa-retweet", Href: "/admin/tradingvarieties/list", Target: "_self"},

			SystemMenu{Id: 30, Pid: 3, Title: "用户资产", Icon: "fa fa-user", Href: "/admin/user/assets", Target: "_self"},
			SystemMenu{Id: 31, Pid: 3, Title: "用户订单", Icon: "fa fa-reorder", Href: "/admin/user/order", Target: "_self"},
			SystemMenu{Id: 32, Pid: 3, Title: "用户挂单", Icon: "fa fa-newspaper-o", Href: "/admin/user/unfinished", Target: "_self"},
		}
	}

	return m.buildMenuChild(0, menuList)
}

func (m *SystemMenu) buildMenuChild(pid int64, menuList []SystemMenu) []*MenuTreeList {
	var treeList []*MenuTreeList
	for _, v := range menuList {
		if pid == v.Pid {
			node := &MenuTreeList{
				Id:     v.Id,
				Title:  v.Title,
				Icon:   v.Icon,
				Href:   v.Href,
				Target: v.Target,
				Pid:    v.Pid,
			}
			child := v.buildMenuChild(v.Id, menuList)
			if len(child) != 0 {
				node.Child = child
			}
			// todo 后续此处加上用户的权限判断
			treeList = append(treeList, node)
		}
	}
	return treeList
}
