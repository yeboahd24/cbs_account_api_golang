package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
	"strconv"
)

func CreateAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data struct {
			AccountID string `json:"coa_id"`
			Name      string `json:"name"`
		}
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		accountID, err := uuid.Parse(data.AccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coa_id"})
			return
		}
		if AccountExists(db, data.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Account with name=%s already exists", data.Name)})
			return
		}
		account, err := ReturnAccount(db, accountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Account with chart_of_account=%s does not exist", data.AccountID)})
			return
		}
		accountNumber, err := generateAccountNumber(db, account)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
		newAccount := Account{
			COA:           account,
			Name:          data.Name,
			AccountNumber: accountNumber,
		}
		err = db.Create(&newAccount).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		// Create a new account balance record
		newAccountBalance := AccountBalance{
			AccountID: newAccount.AccountID,
			Balance:   0, // Set the initial balance to 0
		}
		err = db.Create(&newAccountBalance).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		response := SuccessResponse{Message: "Account created successfully", StatusCode: http.StatusOK}
		c.JSON(http.StatusOK, response)
	}
}

func CreateAccountTypeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			StartRange  int    `json:"start_range"`
			EndRange    int    `json:"end_range"`
		}

		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if data.Name == "" || data.Description == "" || data.StartRange == 0 || data.EndRange == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		existingAccountType := AccountType{}
		err = db.Where("name = ?", data.Name).First(&existingAccountType).Error
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Account type with this name: %s already exists", data.Name)})
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		newAccountType := AccountType{
			Name:        data.Name,
			Description: data.Description,
			StartRange:  data.StartRange,
			EndRange:    data.EndRange,
		}

		err = db.Create(&newAccountType).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		response := SuccessResponse{Message: "Account type created successfully",
			StatusCode: http.StatusOK}

		c.JSON(http.StatusOK, response)
	}
}

// ListAccountTypeResponse
type ListAccountTypeResponse struct {
	AccountID   uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartRange  int       `json:"start_range"`
	EndRange    int       `json:"end_range"`
}

func ListAccountTypeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accountTypes []AccountType
		err := db.Find(&accountTypes).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		var response []ListAccountTypeResponse
		for _, accountType := range accountTypes {
			response = append(response, ListAccountTypeResponse{
				AccountID:   accountType.AccountID,
				Name:        accountType.Name,
				Description: accountType.Description,
				StartRange:  accountType.StartRange,
				EndRange:    accountType.EndRange,
			})
		}

		c.JSON(http.StatusOK, SuccessResponse{Message:"Account Type List", StatusCode: http.StatusOK, Data: response})
	}
}

func CreateChartOfAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data struct {
			AccountTypeID string `json:"account_type_id"`
			AccountNumber int    `json:"account_number"`
			Name          string `json:"name"`
		}

		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		accountTypeID, err := uuid.Parse(data.AccountTypeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account type ID"})
			return
		}

		if accountTypeID == uuid.Nil || data.AccountNumber == 0 || data.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}

		if !AccountTypeExist(db, accountTypeID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account type does not exist"})
			return
		}

		coa := ChartOfAccount{
			AccountTypeID: accountTypeID,
			AccountNumber: data.AccountNumber,
			Name:          data.Name,
		}

		err = db.Create(&coa).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}

		response := SuccessResponse{Message: "Chart of account created successfully",
			StatusCode: http.StatusOK}

		c.JSON(http.StatusOK, response)
	}
}

type ListAccountResponse struct {
	AccountID     uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	AccountNumber int       `json:"account_number"`
}

func ListAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []Account
		err := db.Find(&accounts).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		var response []ListAccountResponse
		for _, account := range accounts {
			response = append(response, ListAccountResponse{
				AccountID:     account.AccountID,
				Name:          account.Name,
				AccountNumber: account.AccountNumber,
			})
		}

		c.JSON(http.StatusOK, SuccessResponse{Message:"Account List", StatusCode: http.StatusOK, Data: response})
	}
}

type CharOfAccountResponse struct {
	AccountID     uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	AccountNumber int       `json:"account_number"`
}

func ListChartOfAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []ChartOfAccount
		err := db.Find(&accounts).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		var response []CharOfAccountResponse
		for _, account := range accounts {
			response = append(response, CharOfAccountResponse{
				AccountID:     account.AccountID,
				Name:          account.Name,
				AccountNumber: account.AccountNumber,
			})
		}

		c.JSON(http.StatusOK, SuccessResponse{Message:"Chart of Account List", StatusCode: http.StatusOK, Data: response})
	}
}

func GetAccountIdByAccountNumber(db *gorm.DB, accountNumber uint) (string, error) {
	var account AccountBalance
	result := db.Where("accountnumber =?", accountNumber).First(&account).Error
	if result != nil {
		if result == gorm.ErrRecordNotFound {
			// Handle case where account does not exist
			// Optionally, create a new account here
			return "", fmt.Errorf("account with number %d not found", accountNumber)
		}
		return "", result
	}
	return account.AccountID.String(), nil
}

func CreateJournalEntryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data struct {
			AccountCredit int    `json:"credit_account"`
			AccountDebit  int    `json:"debit_account"`
			Amount        int    `json:"amount"`
			Description   string `json:"description"`
		}

		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		accountDebit, debitErr := getAccountNumber(db, data.AccountDebit)

		if accountDebit == 0 || debitErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid debit account number"})
			return
		}

		accountCredit, creditErr := getAccountNumber(db, data.AccountCredit)
		if accountCredit == 0 || creditErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credit account number"})
			return
		}

		processTransaction(db, accountDebit, accountCredit, data.Amount)

		// Fetch current balances
		entry := JournalEntry{
			AccountCreditNumber: accountCredit,
			AccountDebitNumber:  accountDebit,
			Amount:              data.Amount,
			Description:         data.Description,
			Date:                time.Now(),
		}

		err = db.Create(&entry).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		response := SuccessResponse{Message: "Journal entry created successfully",
			StatusCode: http.StatusOK}

		c.JSON(http.StatusOK, response)
	}
}

type JournalEntryResponse struct {
	TransactionID       uuid.UUID `json:"transactionid"`
	AccountDebitNumber  int       `json:"accounttodebitid"`
	AccountCreditNumber int       `json:"accounttocreditid"`
	Date                time.Time `json:"date"`
	Amount              int       `json:"amount"`
	Description         string    `json:"description"`
}

func ListJournalEntryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entries []JournalEntry
		err := db.Find(&entries).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		var response []JournalEntryResponse
		for _, entry := range entries {
			response = append(response, JournalEntryResponse{
				TransactionID:       entry.TransactionID,
				AccountDebitNumber:  entry.AccountDebitNumber,
				AccountCreditNumber: entry.AccountCreditNumber,
				Date:                entry.Date,
				Amount:              entry.Amount,
				Description:         entry.Description,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

func ProfitAndLossHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Query("date")
		var date *time.Time
		if dateStr != "" {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date format", StatusCode: http.StatusBadRequest})
				return
			}
			date = &parsedDate
		}

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		

		profitAndLossData, err := profitAndLost(db, date, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Profit and loss data retrieved successfully", StatusCode: http.StatusOK, Data: profitAndLossData})
	}
}

func BalanceSheetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []Account
		err := db.Find(&accounts).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		balanceSheetData := make(map[string]interface{})
		assets := make([]map[string]interface{}, 0)
		liabilities := make([]map[string]interface{}, 0)
		equity := make([]map[string]interface{}, 0)

		for _, account := range accounts {
			var balance AccountBalance
			err = db.Where("accountid =?", account.AccountID).First(&balance).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// Handle accounts with no balance
					continue
				}
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
				return
			}

			accountData := map[string]interface{}{
				"account_name":   account.Name,
				"account_number": account.AccountNumber,
				"balance":        balance.Balance,
			}

			// Categorize accounts based on their number ranges
			if account.AccountNumber >= 1000 && account.AccountNumber <= 1999 {
				assets = append(assets, accountData)
			} else if account.AccountNumber >= 2000 && account.AccountNumber <= 2999 {
				liabilities = append(liabilities, accountData)
			} else if account.AccountNumber >= 3000 && account.AccountNumber <= 3999 {
				equity = append(equity, accountData)
			}
		}

		balanceSheetData["assets"] = assets
		balanceSheetData["liabilities"] = liabilities
		balanceSheetData["equity"] = equity

		c.JSON(http.StatusOK, SuccessResponse{Message: "Balance sheet data retrieved successfully", StatusCode: http.StatusOK, Data: balanceSheetData})

	}
}
