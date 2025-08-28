package models

type GoodsColor struct {
	Id         int
	ColorName  string
	ColorValue string
	Status     int
}

func (GoodsColor) TableName() string {
	return "goods_color"
}
