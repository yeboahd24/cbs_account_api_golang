package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// @summary Create a new account
// @description Create a new account
// tags User Authentication
// @accept  json
// @produce  json
// @param account body Account true "Account"
// @success 200 {object} Account
// @router /account [post]
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

// @summary Create a new account type
// @description Create a new account type
// @accept  json
// @produce  json
// @param accountType body AccountType true "AccountType"
// @success 200 {object} AccountType
// @router /accounttype [post]
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



// @summary List account type
// @description List account ListAccountTypeHandler
// @produce  json
// @success 200 {object} ListAccountTypeResponse
// @router /accounttype [get]
func ListAccountTypeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve pagination query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var accountTypes []AccountType
		err := db.Limit(limit).Offset(offset).Find(&accountTypes).Error
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

		c.JSON(http.StatusOK, SuccessResponse{Message: "Account Type List", StatusCode: http.StatusOK, Data: response})
	}
}

// @summary Create a new chart of account
// @description Create a new chart of account
// @accept  json
// @produce  json
// @param chartOfAccount body ChartOfAccount true "ChartOfAccount"
// @success 200 {object} ChartOfAccount
// @router /coa [post]
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


// @summary List account
// @description List account
// @produce  json
// @success 200 {object} ListAccountResponse
// @router /account [get]
func ListAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve pagination query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var accounts []Account
		err := db.Limit(limit).Offset(offset).Find(&accounts).Error
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

		c.JSON(http.StatusOK, SuccessResponse{Message: "Account List", StatusCode: http.StatusOK, Data: response})
	}
}

type CharOfAccountResponse struct {
	AccountID     uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	AccountNumber int       `json:"account_number"`
}


// @summary List chart of account
// @description List chart of account
// @produce  json
// @success 200 {object} CharOfAccountResponse
// @router /coa [get]
func ListChartOfAccountHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Retrieve pagination query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var accounts []ChartOfAccount
		err := db.Limit(limit).Offset(offset).Find(&accounts).Error
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

		c.JSON(http.StatusOK, SuccessResponse{Message: "Chart of Account List", StatusCode: http.StatusOK, Data: response})
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

// @summary Create a new journal entry
// @description Create a new journal entry
// @accept  json
// @produce  json
// @param journalEntry body JournalEntry true "JournalEntry"
// @success 200 {object} JournalEntry
// @router /journalentry [post]
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


// @summary List journal entry
// @description List journal entry
// @produce  json
// @success 200 {object} JournalEntryResponse
// @router /journalentry [get]
func ListJournalEntryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Retrieve pagination query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var entries []JournalEntry
		err := db.Limit(limit).Offset(offset).Find(&entries).Error
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

// @summary Profit and loss
// @description Profit and loss
// @accept  json
// @produce  json
// @success 200 {object} map[string]interface{}
// @router /profitandloss [get]
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

		profitAndLossData, err := profitAndLost(db, date, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Profit and loss data retrieved successfully", StatusCode: http.StatusOK, Data: profitAndLossData})
	}
}

// @summary Balance sheet
// @description Balance sheet
// @produce  json
// @success 200 {object} map[string]interface{}
// @router /balancesheet [get]
func BalanceSheetHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve pagination query parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit

		var accounts []Account
		err := db.Limit(limit).Offset(offset).Find(&accounts).Error
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
			err = db.Where("accountid = ?", account.AccountID).First(&balance).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// Handle accounts with no balance
					balance = AccountBalance{
						AccountID: account.AccountID,
						Balance:   0,
					}
				} else {
					c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error(), StatusCode: http.StatusInternalServerError})
					return
				}
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
