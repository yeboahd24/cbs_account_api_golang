package main

import (
	"time"

	"github.com/google/uuid"
)

type AccountType struct {
	AccountID   uuid.UUID `json:"account_id" gorm:"column:accountid;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	StartRange  int       `json:"start_range" gorm:"column:startrange"`
	EndRange    int       `json:"end_range" gorm:"column:endrange"`
}

func (AccountType) TableName() string {
	return "accounttype"
}

type ChartOfAccount struct {
	AccountID     uuid.UUID `json:"accountid" gorm:"column:accountid;default:uuid_generate_v4()"`
	AccountTypeID uuid.UUID `json:"accounttypeid" gorm:"column:accounttypeid"`
	AccountNumber int       `json:"accountnumber" gorm:"column:accountnumber"`
	Name          string    `json:"name" gorm:"column:name"`
}

func (ChartOfAccount) TableName() string {
	return "chartofaccount"
}

type Account struct {
	AccountID     uuid.UUID      `json:"accountid" gorm:"column:accountid;default:uuid_generate_v4()"`
	Name          string         `json:"name" gorm:"column:name"`
	COAID         uuid.UUID      `json:"coa_id" gorm:"column:coaid"`
	COA           ChartOfAccount `gorm:"foreignKey:COAID;references:AccountID"`
	AccountNumber int            `json:"account_number" gorm:"column:accountnumber"`
}

func (Account) TableName() string {
	return "account"
}

type AccountBalance struct {
	BalanceID uuid.UUID `json:"balanceid" gorm:"column:balanceid;default:uuid_generate_v4();primarykey"`
	AccountID uuid.UUID `json:"accountid" gorm:"column:accountid;uniqueIndex"`
	Balance   int       `json:"balance" gorm:"column:balance"`
}

func (AccountBalance) TableName() string {
	return "accountbalance"
}

type JournalEntry struct {
	TransactionID       uuid.UUID `json:"transactionid" gorm:"column:transactionid;default:uuid_generate_v4()"`
	AccountDebitNumber  int       `json:"accounttodebitid" gorm:"column:accountdebitnumber"`
	AccountCreditNumber int       `json:"accounttocreditid" gorm:"column:accountcreditnumber"`
	AccountDebit        Account   `json:"accounttodebit" gorm:"column:accountdebitnumber;references:AccountNumber;foreignKey:AccountDebitNumber"`
	AccountCredit       Account   `json:"accounttocredit" gorm:"column:accountcreditnumber;references:AccountNumber;foreignKey:AccountCreditNumber"`
	Date                time.Time `json:"date" gorm:"column:date"`
	Amount              int       `json:"amount" gorm:"column:amount"`
	Description         string    `json:"description" gorm:"column:description"`
}

func (JournalEntry) TableName() string {
	return "journalentry"
}
