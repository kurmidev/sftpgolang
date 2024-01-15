package models

import (
	"fmt"
	"log"
	"time"

	"github.com/liquiloans/sftp/config"
	"github.com/liquiloans/sftp/connection"
	"gorm.io/gorm"
)

type SftpLogging struct {
	Id           int       `json:"id" gorm:"primaryKey"`
	IfaId        int       `json:"ifa_id"`
	FileName     string    `json:"file_name"`
	ReportDate   time.Time `json:"report_date"`
	TotalCount   int       `json:"total_count"`
	ErrorCount   int       `json:"error_count"`
	SuccessCount int       `json:"success_count"`
	Status       string    `json:"status"`
	IsDeleted    string    `json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (sl *SftpLogging) TableName() string {
	return "sftp_logging"
}

func init() {
	connection.Connect()
	db = connection.GetDb()
}

func (sl *SftpLogging) Create() *gorm.DB {
	return db.Save(&sl)
}

func (sl *SftpLogging) Update() *gorm.DB {
	return db.Model(&sl).Updates(sl)
}

func Find(FileName string, ReportDate time.Time, IfaId int) SftpLogging {
	var result SftpLogging
	now := time.Now()
	db.Raw("SELECT * FROM sftp_logging WHERE file_name = ? and report_date=? and ifa_id=?", FileName, now.Format(config.YYYY_MM_DD), IfaId).Scan(&result)
	return result
}

func GetDataForEmailOld(IfaId int) []map[string]string {
	var LoggingDetails []map[string]string
	now := time.Now()
	for _, url := range config.Sftpconfig[IfaId] {
		key := fmt.Sprintf("%s_Initial_%d_%s", url.FileName, IfaId, now.Format(config.YYYYMMDD))
		vals, err := redisdb.HGetAll(ctx, key).Result()
		if err != nil {
			log.Fatal(err)
		}
		LoggingDetails = append(LoggingDetails, map[string]string{
			"FileType": url.FileName,
			"Success":  vals["Success"],
			"Error":    vals["Error"],
			"Total":    vals["Total"],
		})
	}
	return LoggingDetails
}

func GetDataForEmail(IfaId int) []SftpLogging {
	var results []SftpLogging
	now := time.Now()
	db.Raw("SELECT file_name,success_count,error_count,total_count FROM sftp_logging WHERE  report_date=? and ifa_id=?", now.Format(config.YYYY_MM_DD), IfaId).Scan(&results)
	return results
}
