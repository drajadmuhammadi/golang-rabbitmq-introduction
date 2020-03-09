package models

type BasicResponse struct {
	Status bool        `json:"status"`
	Value  interface{} `json:"value"`
}
