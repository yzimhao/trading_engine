package types

type Status int

const (
	StatusEnabled Status = 0
	StatusDisable Status = 1
)

func ParseStatusString(v string) Status {
	switch v {
	case "on", "ON", "0":
		return StatusEnabled
	default:
		return StatusDisable
	}
}
