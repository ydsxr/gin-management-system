package models

type GoodsTypeAttribute struct {
	Id        int    `json:"id"`
	CateId    int    `json:"cate_id"`
	Title     string `json:"title"`
	AttrType  int    `json:"attr_type"`
	AttrValue string `json:"attr_value"`
	Sort      int    `json:"sort"`
	AddTime   int    `json:"add_time"`
	Status    int    `json:"status"`
}

func (GoodsTypeAttribute) TableName() string {
	return "goods_type_attribute"
}
