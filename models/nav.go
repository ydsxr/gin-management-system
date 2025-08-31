package models

type Nav struct {
	Id int
	Title string
	Link string
	Position int
	IsOpennew int
	Relation string
	Sort int
	Status int
	AddTime int
}

func (Nav) TableName() string{
	return "nav"
}