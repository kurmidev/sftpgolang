package models

type ApiResponse struct {
	Status   bool        `json:"status"`
	Data     interface{} `json:"data"`
	Code     int         `json:"code"`
	CheckSum string      `json:"checksum"`
}
