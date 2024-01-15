package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/liquiloans/sftp/controllers"
)

func main() {
	start := time.Now()
	IfaId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("invalid input parameter ", err)
	}
	fmt.Println("String data fetching jobs.......")
	controllers.FetchInvestorDetails(IfaId)
	fmt.Println("File generation jobs.......")
	controllers.WriteDataToFile(IfaId)
	fmt.Println("Failed file generation jobs.......")
	controllers.WriteDataToFailedFile(IfaId)
	elapsed := time.Since(start)
	fmt.Printf("Time taken for jobs %s", elapsed)
}
