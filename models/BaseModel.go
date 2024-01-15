package models

type BaseModel struct {
	InvestorId int
	Mid        string
	TimeStamp  string
	Checksum   string
}

type BaseModelFunc interface {
	GenerateCheckum(secret string)
}
