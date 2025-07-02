package types

type Status int

const (
	StatusEnabled Status = iota + 1
	StatusDisabled
)
