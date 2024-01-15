package models

import (
	"time"

	"github.com/liquiloans/sftp/connection"
	"gorm.io/gorm"
)

var db *gorm.DB

type InvInvestorUser struct {
	Id                   int       `json:"id"`
	IfaId                int       `json:"ifa_id"`
	CategoryId           int64     `json:"category_id"`
	Name                 string    `json:"name"`
	FolioMasterId        int64     `json:"folio_master_id"`
	NameFetchedBy        string    `json:"name_fetched_by"`
	InvestorType         string    `json:"investor_type"`
	DashboardView        string    `json:"dashboard_view"`
	EntityType           string    `json:"entity_type"`
	HoldingType          string    `json:"holding_type"`
	Gender               string    `json:"gender"`
	Dob                  time.Time `json:"dob"`
	Email                string    `json:"email"`
	EmailHash            string    `json:"email_hash"`
	ContactNumber        string    `json:"contact_number"`
	ContactNumberHash    string    `json:"contact_number_hash"`
	Pan                  string    `json:"pan"`
	PanHash              string    `json:"pan_hash"`
	AadhaarNumber        string    `json:"aadhaar_number"`
	AadhaarNumberHash    string    `json:"aadhaar_number_hash"`
	GstNumber            string    `json:"gst_number"`
	ProfileRoi           float64   `jsno:"profile_roi"`
	ApprovalStatus       string    `json:"approval_status"`
	EAgreementMode       string    `json:"e_agreement_mode"`
	EAgreementStatus     string    `json:"e_agreement_status"`
	EAgreementDate       time.Time `json:"e_agreement_date"`
	EAgreementSignedDate time.Time
	EAgreementDocumentId string    `json:"e_agreement_document_id"`
	CkycNumber           string    `json:"ckyc_number"`
	MaxInvestmentBucket  int64     `json:"max_investment_bucket"`
	AuthenticatedBy      int64     `json:"authenticated_by"`
	AuthenticatedAt      time.Time `json:"authenticated_at"`
	KycVerificationMode  string    `json:"kyc_verification_mode"`
	KycStatus            string    `json:"kyc_status"`
	KycDate              time.Time `json:"kyc_date"`
	CaCertificateId      int64     `json:"ca_certificate_id"`
	CaCertificateStatus  string    `json:"ca_certificate_status"`
	OnboardingSource     string    `json:"onboarding_source"`
	Status               string    `json:"status"`
	IsDeleted            string    `json:"is_deleted"`
	CreatedBy            int64     `json:"created_by"`
	UpdatedBy            int64     `json:"updated_by"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Mid                  string    `json:"mid"`
	Secret               string    `json:"secret"`
}

func init() {
	connection.Connect()
	db = connection.GetDb()
}

func GetInvestorByIfaQuery(IfaId int) *gorm.DB {
	return db.Table("inv_investor_user").Select("inv_investor_user.*,ifa_api_intergration.mid,ifa_api_intergration.key as secret").Joins("inner join ifa_api_intergration on ifa_api_intergration.ifa_id = inv_investor_user.ifa_id").Where("inv_investor_user.ifa_id= ?", IfaId).Where("inv_investor_user.status = ? and inv_investor_user.is_deleted=? ", "Active", "False").Limit(100)
}

func GetInvestorByIfaCount(IfaId int) int64 {
	var count int64
	db.Table("inv_investor_user").Select("inv_investor_user.*,ifa_api_intergration.mid,ifa_api_intergration.key as secret").Where("inv_investor_user.ifa_id= ?", IfaId).Where("inv_investor_user.status = ? and inv_investor_user.is_deleted=? ", "Active", "False").Count(&count)
	return count
}
