package models

type BasicResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}
