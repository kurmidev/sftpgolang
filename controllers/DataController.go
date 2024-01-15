package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/liquiloans/sftp/config"
	"github.com/liquiloans/sftp/models"
	"github.com/liquiloans/sftp/utils"
	"gorm.io/gorm"
)

/**
1. Log the file generations
2. Fetch investor details in batch
3. For each batch initiate go routine

**/

func FetchInvestorDetails(IfaId int) {
	QueryObject := models.GetInvestorByIfaQuery(IfaId)

	var investors []models.InvInvestorUser
	fileTogenerate := config.Sftpconfig[IfaId]
	wg := sync.WaitGroup{}

	initiateLoggin(IfaId)

	QueryObject.FindInBatches(&investors, config.BatchSize, func(tx *gorm.DB, batch int) error {
		if tx.Error != nil {
			log.Fatal(tx.Error)
		}
		for _, url := range fileTogenerate {
			wg.Add(1)
			go FetchDataByUrls(url.Url, url.FileName, investors, &wg)
		}
		fmt.Printf("Batch initiated done %d next batch after %d \n", batch, time.Minute)
		time.Sleep(time.Second)
		return nil
	})
	wg.Wait()
}

/*
*

	For each batch
	1 Create Request
	2 Generate Checksum
	3 Perform Curl request
	4 Save data in redis

*
*/
func FetchDataByUrls(url string, filename string, investors []models.InvInvestorUser, wg *sync.WaitGroup) {
	var request []byte
	defer func() {
		fmt.Println("Waitgroup finished for ", filename)
		wg.Done()
	}()
	fmt.Println("Total Investor found is ", len(investors))
	for _, investor := range investors {
		request = generateRequestData(filename, investor)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(request))
		if err != nil {
			utils.Error.Printf("Not able fetch data for %s of investor %d ", filename, investor.Id)
			utils.Error.Println("Error is  ", err)
		} else if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.Error.Println("Unable to read the data ", err)
			} else {
				var data models.ApiResponse
				if err := json.Unmarshal([]byte(string(body)), &data); err != nil {
					utils.Error.Println("Unable to un marshal the data", err)
				} else {
					models.SaveToRedis(filename, investor.IfaId, investor.Id, []byte(string(body)))
				}
			}
		}
		defer resp.Body.Close()
	}
}

/**
Generate the request data
**/

func generateRequestData(reqtype string, investor models.InvInvestorUser) []byte {
	var request []byte
	if reqtype == config.CashLedger {
		r := models.CashLedger{
			InvestorId: investor.Id,
			Mid:        investor.Mid,
			TimeStamp:  "",
			Checksum:   "",
		}
		r.GenerateCheckum(investor.Secret)
		request, _ = json.Marshal(r)
	} else if reqtype == config.DashboardSummary {
		r := models.DashboardSummary{
			InvestorId: investor.Id,
			Mid:        investor.Mid,
			TimeStamp:  "",
			Checksum:   "",
		}
		r.GenerateCheckum(investor.Secret)
		request, _ = json.Marshal(r)
	} else if reqtype == config.InvestmentSummary {
		r := models.InvestmentSummary{
			InvestorId: investor.Id,
			Mid:        investor.Mid,
			TimeStamp:  "",
			Checksum:   "",
		}
		r.GenerateCheckum(investor.Secret)
		request, _ = json.Marshal(r)

	}
	return request
}

/**
1. Fetch investor details in batch
2. For each batch initiate go routine to write data to file
**/

func WriteDataToFile(IfaId int) {
	var investors []models.InvInvestorUser
	fileTogenerate := config.Sftpconfig[IfaId]
	filewg := sync.WaitGroup{}
	QueryObject := models.GetInvestorByIfaQuery(IfaId)
	i := 1
	QueryObject.FindInBatches(&investors, config.FileSize, func(tx *gorm.DB, batch int) error {
		if tx.Error != nil {
			utils.Error.Println("Error fetching investor details", tx.Error)
		} else {
			for _, url := range fileTogenerate {
				filewg.Add(1)
				go getDataForFile(url.FileName, IfaId, investors, i, &filewg)
			}
			i++
			fmt.Printf("Batch initiated done %d \n", batch)
		}
		return nil
	})
	filewg.Wait()
}

/*
*
For each batch
 1. Get investor datat from redis
 2. Arrange the data in following format
    a. Succes data saved in response
    b. Not found data in redis

3. For the fetched data saved Create a file and upload in aws
*
*/
func getDataForFile(FileType string, IfaId int, investors []models.InvInvestorUser, fileno int, filewg *sync.WaitGroup) {

	defer func() {
		fmt.Printf("File generation finished for %s \n", FileType)
		filewg.Done()
	}()
	var response []interface{}

	for _, investor := range investors {
		data := models.GetInvestorData(FileType, IfaId, investor.Id)
		if data != nil {
			response = append(response, data)
		} else {
			models.SaveFailedData(FileType, IfaId, investor.Id)
		}
	}

	if len(response) > 0 {
		UpdateLogging(IfaId, FileType, len(response))
		fileName := generateFiles(fileno, FileType, response)
		bucket := config.SFTP_PARTNERS[IfaId]
		fmt.Print(fileName, bucket)
		uploadFileToAws(fileName, FileType, bucket)
	}
}

/*
*
For each generated data save the data in file
*
*/
func generateFiles(fileno int, fileType string, fileData []interface{}) string {
	filename := fmt.Sprintf("%s_%d.%s", fileType, fileno, "json")
	f, _ := os.Create(filename)
	defer f.Close()
	b, _ := json.MarshalIndent(fileData, "", "    ")
	f.WriteString(string(b))
	fmt.Printf("File generation finished for %s with %d \n", filename, fileno)
	return filename
}

/*
*
Upload file to AWS S3 bucket
*
*/
func uploadFileToAws(FileName string, FileType string, bucket string) {
	models.UploadFileToS3(FileName, FileType, bucket)
}

/*
*
Initiate login functions
*
*/
func initiateLoggin(IfaId int) {
	fileTogenerate := config.Sftpconfig[IfaId]
	count := models.GetInvestorByIfaCount(IfaId)
	var logging models.SftpLogging
	for _, url := range fileTogenerate {
		logging = models.SftpLogging{
			IfaId:        IfaId,
			FileName:     url.FileName,
			ReportDate:   time.Now(),
			TotalCount:   int(count),
			ErrorCount:   0,
			SuccessCount: 0,
			Status:       "Active",
			IsDeleted:    "False",
		}
		logging.Create()
	}
}

/*
*
Update success and error messages
*
*/
func UpdateLogging(IfaId int, FileType string, SuccessCount int) {
	fmt.Println("Staring update logging")
	result := models.Find(FileType, time.Now(), IfaId)
	if result != (models.SftpLogging{}) {
		result.SuccessCount = result.SuccessCount + SuccessCount
		result.ErrorCount = result.TotalCount - result.SuccessCount
		result.Update()
		fmt.Printf("%#v", result)
	}
}

/*
*
Send final data in emails
*
*/
func SendEmail(IfaId int) {
	now := time.Now()
	templateData := models.GetDataForEmail(IfaId)
	fmt.Println("Sending emails.")
	subject := fmt.Sprintf("SFT Data status for %s", now.Format(config.YYYYMMDD_FMT))
	r := models.NewRequest(subject, subject)
	err := r.ParseTemplate("html/template.html", templateData)
	if err != nil {
		utils.Error.Println("Error parsing the template file", err)
	} else {
		fmt.Print("Sending emails....")
		_, err := r.SendEmail()
		if err != nil {
			utils.Error.Println("Error sending emails", err)
		} else {
			fmt.Println("Sending emails......")
		}
	}
}

/*
*
Write data to failed files
*
*/
func WriteDataToFailedFile(IfaId int) {
	fileTogenerate := config.Sftpconfig[IfaId]
	filewg := sync.WaitGroup{}
	for _, url := range fileTogenerate {
		filewg.Add(1)
		go getFailedDataFile(url.FileName, IfaId, &filewg)
	}
	filewg.Wait()
}

func getFailedDataFile(FileName string, IfaId int, filewg *sync.WaitGroup) {
	defer func() {
		fmt.Printf("Failed File generation finished for %s \n", FileName)
		filewg.Done()
	}()
	response := models.GetFailedData(FileName, IfaId)
	var results []interface{}
	for _, investorIdData := range response {
		investorId, err := strconv.Atoi(investorIdData)
		if err != nil {
			utils.Error.Println("String to int conversion ", err)
		} else {
			result := models.GetPreviousData(FileName, IfaId, investorId)
			if result != nil {
				results = append(results, result)
			}
		}
	}

	if len(results) > 0 {
		fileType := fmt.Sprintf("%s_Failed", FileName)
		fileName := generateFiles(1, fileType, results)
		bucket := config.SFTP_PARTNERS[IfaId]
		fmt.Print(fileName, bucket)
		uploadFileToAws(fileName, fileName, bucket)
	}

	if len(response) > 0 {
		fileType := fmt.Sprintf("%s_Failed_Investor_List", FileName)
		fileName := generateTxtFiles(1, fileType, response)
		bucket := config.SFTP_PARTNERS[IfaId]
		fmt.Print(fileName, bucket)
		uploadFileToAws(fileName, fileName, bucket)
	}
}

func generateTxtFiles(fileno int, fileType string, fileData []string) string {
	filename := fmt.Sprintf("%s_%d.%s", fileType, fileno, "txt")
	f, _ := os.Create(filename)
	defer f.Close()
	buffer := bufio.NewWriter(f)
	for _, investorid := range fileData {
		_, err := buffer.WriteString(investorid + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := buffer.Flush(); err != nil {
		utils.Error.Println("Error writting failed investor ids ", err)
	} else {
		fmt.Printf("File generation finished for %s with %d \n", filename, fileno)
	}

	return filename
}
