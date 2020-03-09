package models

type OpenMethodAPI struct {
	Method string      `json:"method"`
	Value  interface{} `json:"value"`
}
