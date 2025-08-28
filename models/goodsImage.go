package models

type GoodsImage struct {
	Id      int
	GoodsId int
	ImgUrl  string
	ColorId int
	Sort    int
	AddTime int
	Status  int
}

func (GoodsImage) TableName() string {
	return "goods_image"
}
