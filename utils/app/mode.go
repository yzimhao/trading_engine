package app

type ModeType string

const (
	ModeProd  ModeType = "prod"
	ModeDev   ModeType = "dev"
	ModeDebug ModeType = "debug"
	ModeDemo  ModeType = "demo"
)

func (m ModeType) String() string {
	return string(m)
}
