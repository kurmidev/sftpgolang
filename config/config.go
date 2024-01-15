package config

const PrudentIfaId int = 1799
const CredIfaId int = 573
const BharatPeIfaId2 int = 255
const BharatPeIfaId3 int = 890
const EpifiIfaId int = 1827
const UniIfaId int = 2300
const AbchlorIfaId int = 1596
const FincartIfaId int = 436
const LKPIfaId int = 307

const YYYYMMDD = "20060102"
const YYYYMMDD_FMT = "02-01-2006"
const YYYY_MM_DD = "2006-01-02"

type UrlDetails struct {
	Url      string
	FileName string
}

var SFTP_PARTNERS = map[int]string{
	CredIfaId:      "Cred",
	BharatPeIfaId2: "BharatPe2",
	BharatPeIfaId3: "BharatPe3",
	PrudentIfaId:   "Prudent",
	EpifiIfaId:     "Epifi",
	UniIfaId:       "Uni",
	AbchlorIfaId:   "Abchlor",
	FincartIfaId:   "Fincart",
	LKPIfaId:       "LKP",
}

const SFTP_BUCKET = "ll-sftp-storage"

const InvestmentSummary string = "GetInvestmentSummary"
const CashLedger string = "GetCashLedger"
const DashboardSummary string = "GetInvestorDashboard"
const TransactionReport string = "GetTransactionReport"

const FromEmails string = "chandraprakash.patel@liquiloans.com"

var ToEmails []string = []string{"chandraprakash.patel@liquiloans.com"}

var GetInvestmentSummary UrlDetails = UrlDetails{
	Url:      "https://supply-integration.liquiloans.com/api/GetInvestmentSummary",
	FileName: "GetInvestmentSummary",
}

var GetCashLedger UrlDetails = UrlDetails{
	Url:      "https://supply-integration.liquiloans.com/api/GetCashLedger",
	FileName: "GetCashLedger",
}

var GetDashboardSummary UrlDetails = UrlDetails{
	Url:      "https://supply-integration.liquiloans.com/api/GetInvestorDashboard",
	FileName: "GetInvestorDashboard",
}

var GetTransactionReport UrlDetails = UrlDetails{
	Url:      "https://supply-integration.liquiloans.com/api/GetTransactionReport",
	FileName: "GetTransactionReport",
}

var UrlList []UrlDetails = []UrlDetails{
	GetInvestmentSummary,
	GetDashboardSummary,
	GetCashLedger,
}

var Sftpconfig = map[int][]UrlDetails{
	PrudentIfaId:   UrlList,
	BharatPeIfaId3: UrlList,
}

var BatchSize int = 1000

var FileSize int = 15000

const EMAIL_USERNAME string = "aaf01c4efaa8c2"
const EMAIL_PASSWORD string = "22d3e5079e03fe"
const EMAIL_HOST string = "sandbox.smtp.mailtrap.io"
