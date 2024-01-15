package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Result struct {
	Id     int
	IfaId  int
	Mid    string
	Secret string
}

type InvestmentSummary struct {
	InvestorId int    `json:"InvestorId"`
	Mid        string `json:"MID"`
	TimeStamp  string `json:"Timestamp"`
	Checksum   string `json:"Checksum"`
}

type ApiResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    interface{}
}

type InvestmentSummaryItems struct {
	TransactionId               int        `json:"transaction_id"`
	InvestorId                  int        `json:"investor_id"`
	TransactionDate             time.Time  `json:"transaction_date"`
	InvestedAmount              float64    `json:"invested_amount"`
	BalancePrincipal            float64    `json:"balance_principal"`
	BalanceInterest             float64    `json:"balance_interest"`
	SchemeId                    int        `json:"scheme_id"`
	SchemeName                  string     `json:"scheme_name"`
	InvestmentRoi               int        `json:"investment_roi"`
	ReturnType                  string     `json:"return_type"`
	PayoutType                  string     `json:"payout_type"`
	LockinTenure                int        `json:"lockin_tenure"`
	LockinBreak                 string     `json:"lockin_break"`
	MaturityStartDate           time.Time  `json:"maturity_start_date"`
	MaturityEndDate             time.Time  `json:"maturity_end_date"`
	TransactionSubType          string     `json:"transaction_sub_type"`
	InvestmentStatus            string     `json:"investment_status"`
	ParentInvestmentId          *int       `json:"parent_investment_id"`
	MasterParentInvestmentId    int        `json:"master_parent_investment_id"`
	LastWithdrawalAt            *time.Time `json:"last_withdrawal_at"`
	RedeemedPrincipal           float64    `json:"redeemed_principal"`
	RedeemedInterest            float64    `json:"redeemed_interest"`
	TotalRedemption             float64    `json:"total_redemption"`
	RedeemedInterestWithRequest float64    `json:"redeemed_interest_with_request"`
	Name                        string     `json:"name"`
	DisplayScheme               string     `json:"display_scheme"`
	LockinType                  string     `json:"lockin_type"`
	SchemeDetails               string     `json:"scheme_details"`
	LockinEndDate               time.Time  `json:"lockin_end_date"`
	NetPrincipalInvestment      float64    `json:"net_principal_investment"`
	InterestAmount              float64    `json:"interest_amount"`
	AccruedValue                float64    `json:"accrued_value"`
	WithdrawableBalance         float64    `json:"withdrawable_balance"`
	LockedWithdrawableBalance   float64    `json:"locked_withdrawable_balance"`
	RequestId                   int        `json:"request_id"`
	Source                      string     `json:"source"`
}

type InvestmentSummaryResponse struct {
	CurrentInvestment []InvestmentSummaryItems `json:"current_investments"`
	PastInvestments   []InvestmentSummaryItems `json:"past_investments"`
}

var investsummaryResponse []ApiResponse

func mainold() {
	startTime := time.Now()
	dsn := "chandraprakash-patel:PzV0Mk02Ou76@tcp(ll-rds-wr.liquiloans.com:3306)/liquiloans?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var results []Result
	counter := 1
	db.Table("inv_investor_user").Select("inv_investor_user.id,inv_investor_user.ifa_id ,ifa_api_intergration.mid,ifa_api_intergration.key as secret").
		Joins("left join ifa_api_intergration on ifa_api_intergration.ifa_id = inv_investor_user.ifa_id").Where("inv_investor_user.ifa_id= ?", 1799).Where("inv_investor_user.status = ? and inv_investor_user.is_deleted=? ", "Active", "False").FindInBatches(&results, 15000, func(tx *gorm.DB, batch int) error {
		for _, result := range results {
			result.fetchInvestmentSummary()
		}
		WriteDataToFile(counter, "GetInvestmentSummary")
		counter++
		fmt.Println(batch) // Batch 1, 2, 3
		if tx.Error != nil {
			fmt.Println(tx.Error)
		}
		// Returning an error will stop further batch processing
		return nil
	})

	//fmt.Printf("%+v\n", investsummaryResponse)

	endTime := time.Now()
	fmt.Println("Started at ", startTime, " and ended at ", endTime)
}

func (i Result) fetchInvestmentSummary() {
	url := "https://supply-integration.liquiloans.com/api/GetInvestmentSummary"
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	checksumString := fmt.Sprintf("%d%s%s", i.Id, "||", timestamp)
	hmac := hmac.New(sha256.New, []byte(i.Secret))
	hmac.Write([]byte(checksumString))

	invsumm := InvestmentSummary{
		InvestorId: i.Id,
		Mid:        i.Mid,
		TimeStamp:  timestamp,
		Checksum:   hex.EncodeToString(hmac.Sum(nil)),
	}
	//Ignoring error.
	body, _ := json.Marshal(invsumm)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			//Failed to read response.
			log.Fatal(err)
		}
		//Convert bytes to String and print
		var data ApiResponse
		if err := json.Unmarshal([]byte(string(body)), &data); err != nil {
			log.Fatal(err)
		}
		fmt.Print(data)
		investsummaryResponse = append(investsummaryResponse, data)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp.Status, resp.StatusCode)
	}
}

func WriteDataToFile(fileno int, filename string) {
	filename = fmt.Sprintf("%s%d.%s", filename, fileno, "json")
	f, _ := os.Create(filename)
	defer f.Close()
	b, _ := json.MarshalIndent(investsummaryResponse, "", "    ")
	f.WriteString(string(b))
}
