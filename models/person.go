package models

type Person struct {
	Id        int
	Name      string
	ResumeUrl string
}

func (Person) TableName() string {
	return "person"
}
