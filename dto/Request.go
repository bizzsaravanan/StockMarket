package dto

import "net/url"

type Request struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Cookie   string  `json:"cookie"`
	Origin   string  `json:"origin"`
	FormData url.Values
}
