package models

type Access struct {
	Id          int
	ModuleName  string //模块名称
	Type        int    // 节点类型 ： 1，表示模块 2，表示菜单 3,表示操作
	ActionName  string //操作名称
	Url         string //路由跳转地址
	ModuleId    int    //此module_id和当前模型的id关联
	Sort        int
	Description string
	AddTime     int
	Status      int
	AccessItem  []Access `gorm:"foreignKey:ModuleId;references:Id"`
	Checked     bool     `gorm:"-"` // 操作数据库时忽略此字段
}

func (Access) TableName() string {
	return "access"
}
