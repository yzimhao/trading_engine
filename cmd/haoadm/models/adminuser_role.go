package models

type AdminRole int

const (
	AdminRoleSuper AdminRole = iota
	AdminRoleAdmin
	AdminRoleUser
)

var statusText = map[AdminRole]string{
	AdminRoleSuper: "super",
	AdminRoleAdmin: "admin",
	AdminRoleUser:  "user",
}
