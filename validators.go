package main

import (
	// "errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"

	"github.com/google/uuid"

	"time"
)

// AccountExists
func AccountExists(db *gorm.DB, name string) bool {
	existingAccount := Account{}

	err := db.Where("name = ?", name).First(&existingAccount).Error
	if err != nil {
		return false
	}

	return true
}

// AccountTypeExist
func AccountTypeExist(db *gorm.DB, accountTypeID uuid.UUID) bool {
	existingAccountType := AccountType{}
	err := db.Where("accountid = ?", accountTypeID).First(&existingAccountType).Error
	if err != nil {
		return false
	}

	return true
}

func ReturnAccount(db *gorm.DB, accountID uuid.UUID) (ChartOfAccount, error) {
	var account ChartOfAccount

	err := db.Table("chartofaccount").Where("accountid = ?", accountID).First(&account).Error
	if err != nil {
		return account, err
	}

	fmt.Println("account", account)

	return account, nil
}

// function to get the next number of Account if nothing then use the chartofaccount + 1
func generateAccountNumber(db *gorm.DB, coa ChartOfAccount) (int, error) {
	var maxAccountNumber sql.NullInt64
	err := db.Raw("SELECT COALESCE(MAX(accountnumber), ?) FROM account WHERE coaid = ?", coa.AccountNumber, coa.AccountID).Scan(&maxAccountNumber).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch max account number: %v", err)
	}

	if !maxAccountNumber.Valid {
		// If maxAccountNumber is NULL (no existing accounts), start from the base chart of account number
		return int(coa.AccountNumber) + 1, nil
	}

	// If maxAccountNumber is valid, increment it by 1 to get the next available account number
	return int(maxAccountNumber.Int64) + 1, nil
}

func getAccountNumber(db *gorm.DB, accountNumber int) (int, error) {

	accountExist := Account{}

	err := db.Where("accountnumber = ?", accountNumber).First(&accountExist).Error
	if err != nil {
		return 0, fmt.Errorf("failed to fetch max account number: %v", err)
	}
	return accountExist.AccountNumber, nil

}

func getAccountID(db *gorm.DB, accountNumber int) (string, error) {

	accountExist := Account{}

	err := db.Where("accountnumber = ?", accountNumber).First(&accountExist).Error
	if err != nil {
		return "", fmt.Errorf("failed to fetch max account number: %v", err)
	}
	return accountExist.AccountID.String(), nil
}

func processTransaction(db *gorm.DB, debitAccountNumber int, creditAccountNumber int, amount int) error {
	// Retrieve account IDs
	debitAccountID, err := getAccountUUID(db, debitAccountNumber)
	if err != nil {
		return fmt.Errorf("failed to get debit account ID: %v", err)
	}
	creditAccountID, err := getAccountUUID(db, creditAccountNumber)
	if err != nil {
		return fmt.Errorf("failed to get credit account ID: %v", err)
	}

	// Retrieve or create debit account balance
	debitBalance := AccountBalance{AccountID: debitAccountID}
	err = db.FirstOrCreate(&debitBalance, AccountBalance{AccountID: debitAccountID}).Error
	if err != nil {
		return fmt.Errorf("failed to fetch or create debit account balance: %v", err)
	}

	// Retrieve or create credit account balance
	creditBalance := AccountBalance{AccountID: creditAccountID}
	err = db.FirstOrCreate(&creditBalance, AccountBalance{AccountID: creditAccountID}).Error
	if err != nil {
		return fmt.Errorf("failed to fetch or create credit account balance: %v", err)
	}

	// Update balances
	debitBalance.Balance -= amount
	creditBalance.Balance += amount

	// Save updated balances
	err = db.Save(&debitBalance).Error
	if err != nil {
		return fmt.Errorf("failed to update debit account balance: %v", err)
	}

	err = db.Save(&creditBalance).Error
	if err != nil {
		return fmt.Errorf("failed to update credit account balance: %v", err)
	}

	return nil
}

func getAccountUUID(db *gorm.DB, accountNumber int) (uuid.UUID, error) {
	accountExist := Account{}
	err := db.Where("accountnumber = ?", accountNumber).First(&accountExist).Error
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to fetch account UUID: %v", err)
	}
	return accountExist.AccountID, nil
}

func updateAccountBalance(db *gorm.DB, accountID uuid.UUID, amount int) error {
	var accountBalance AccountBalance
	err := db.Where("accountid = ?", accountID).FirstOrCreate(&accountBalance).Error
	if err != nil {
		return fmt.Errorf("failed to fetch or create account balance: %v", err)
	}

	accountBalance.Balance += amount

	err = db.Save(&accountBalance).Error
	if err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	return nil
}

func profitAndLost(db *gorm.DB, date *time.Time, limit, offset int) (map[string]interface{}, error) {
	// Fetch paginated "Expenses" and "Incomes" accounts
	var matchingAccounts []Account
	err := db.Where("name IN ?", []string{"Expenses", "Incomes"}).Limit(limit).Offset(offset).Find(&matchingAccounts).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("Matching Accounts:")
	for _, account := range matchingAccounts {
		fmt.Printf("AccountID: %s, Name: %s, AccountNumber: %d, COAID: %s\n", account.AccountID, account.Name, account.AccountNumber, account.COAID)
	}

	if len(matchingAccounts) == 0 {
		return nil, nil
	}

	profitAndLossData := make(map[string]interface{})
	expenses := make([]map[string]interface{}, 0)
	incomes := make([]map[string]interface{}, 0)
	totalExpenses := 0
	totalIncomes := 0

	for _, account := range matchingAccounts {
		var debitTransactions []JournalEntry
		var creditTransactions []JournalEntry

		debitQuery := db.Where("AccountDebitNumber = ?", account.AccountNumber)
		if date != nil {
			debitQuery = debitQuery.Where("Date BETWEEN ? AND ?", date.Truncate(24*time.Hour), date.AddDate(0, 0, 1).Add(-time.Second))
		}
		err = debitQuery.Find(&debitTransactions).Error
		if err != nil {
			return nil, err
		}

		creditQuery := db.Where("AccountCreditNumber = ?", account.AccountNumber)
		if date != nil {
			creditQuery = creditQuery.Where("Date BETWEEN ? AND ?", date.Truncate(24*time.Hour), date.AddDate(0, 0, 1).Add(-time.Second))
		}
		err = creditQuery.Find(&creditTransactions).Error
		if err != nil {
			return nil, err
		}

		debitBalance := sumTransactionAmounts(debitTransactions)
		creditBalance := sumTransactionAmounts(creditTransactions)
		netBalance := abs(debitBalance - creditBalance)

		fmt.Printf("Account: %s\n", account.Name)
		fmt.Println("Debit Transactions:")
		for _, transaction := range debitTransactions {
			fmt.Printf("TransactionID: %s, AccountDebitNumber: %d, AccountCreditNumber: %d, Date: %s, Amount: %d, Description: %s\n", transaction.TransactionID, transaction.AccountDebitNumber, transaction.AccountCreditNumber, transaction.Date, transaction.Amount, transaction.Description)
		}
		fmt.Println("Credit Transactions:")
		for _, transaction := range creditTransactions {
			fmt.Printf("TransactionID: %s, AccountDebitNumber: %d, AccountCreditNumber: %d, Date: %s, Amount: %d, Description: %s\n", transaction.TransactionID, transaction.AccountDebitNumber, transaction.AccountCreditNumber, transaction.Date, transaction.Amount, transaction.Description)
		}

		if account.Name == "Expenses" {
			expenses = append(expenses, map[string]interface{}{
				"account_name": account.Name,
				"balance":      netBalance,
			})
			totalExpenses += netBalance
		} else if account.Name == "Incomes" {
			incomes = append(incomes, map[string]interface{}{
				"account_name": account.Name,
				"balance":      netBalance,
			})
			totalIncomes += netBalance
		}
	}

	profitAndLossData["expenses"] = expenses
	profitAndLossData["incomes"] = incomes
	profitAndLossData["total_expenses"] = totalExpenses
	profitAndLossData["total_incomes"] = totalIncomes

	return profitAndLossData, nil
}

func sumTransactionAmounts(transactions []JournalEntry) int {
	total := 0
	for _, transaction := range transactions {
		total += transaction.Amount
	}
	return total
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
