package models

type Assets struct {
	Base
	UserId    string `json:"user_id,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Total     string `json:"total,omitempty"`
	Freeze    string `json:"freeze,omitempty"`
	Available string `json:"avail,omitempty"`
}

type CreateAssets struct {
	UserId    string  `json:"user_id,omitempty"`
	Symbol    string  `json:"symbol,omitempty"`
	Total     *string `json:"total,omitempty"`
	Freeze    *string `json:"freeze,omitempty"`
	Available *string `json:"avail,omitempty"`
}

type UpdateAssets struct {
	UserId    *string `json:"user_id,omitempty"`
	Symbol    *string `json:"symbol,omitempty"`
	Total     *string `json:"total,omitempty"`
	Freeze    *string `json:"freeze,omitempty"`
	Available *string `json:"avail,omitempty"`
}
