package models

type Focus struct {
	Id         int
	Title      string
	FocusType  int
	FocusImage string
	Link       string
	Sort       int
	AddTime    int
	Status     int
}

func (Focus) TableName() string {
	return "focus"
}
