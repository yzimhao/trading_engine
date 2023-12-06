package models

type AdminuserStatus string

const (
	AdminuserStatusNormal  AdminuserStatus = "normal"
	AdminuserStatusDisable AdminuserStatus = "disable"
	AdminuserStatusLocked  AdminuserStatus = "locked"
)

func (s AdminuserStatus) List() []string {
	return []string{
		AdminuserStatusNormal.String(),
		AdminuserStatusDisable.String(),
		AdminuserStatusLocked.String(),
	}
}

func (s AdminuserStatus) String() string {
	return string(s)
}
