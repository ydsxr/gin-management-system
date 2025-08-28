package models

type GoodsCate struct {
	Id          int
	Title       string
	CateImg     string
	Link        string
	Template    string
	Pid         int
	SubTitle    string
	Keywords    string
	Description string
	Status      int
	Sort        int
	AddTime     int
	GoodsCateItems []GoodsCate `gorm:"foreignKey:Pid;references:Id"`
}

func (GoodsCate) TableName() string {
	return "goods_cate"
}
