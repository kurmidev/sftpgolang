package models

import (
	"fmt"
	"time"

	"github.com/liquiloans/sftp/utils"
)

type InvestmentSummary struct {
	InvestorId int    `json:"InvestorId"`
	Mid        string `json:"MID"`
	TimeStamp  string `json:"Timestamp"`
	Checksum   string `json:"Checksum"`
}

func (c *InvestmentSummary) GenerateCheckum(secret string) {
	c.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	checksumString := fmt.Sprintf("%d%s%s", c.InvestorId, "||", c.TimeStamp)
	c.Checksum = utils.EncryptString(checksumString, secret)
}
