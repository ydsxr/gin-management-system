package models

type RoleAccess struct {
	RoleId   int
	AccessId int
}

func (RoleAccess) TableName() string {
	return "role_access"
}
