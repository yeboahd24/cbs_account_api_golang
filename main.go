package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.Use(JWTAuthMiddleware(), JSONMiddleware())

	router.POST("/account", CreateAccountHandler(db))
	router.POST("/accounttype", CreateAccountTypeHandler(db))
	router.POST("/coa", CreateChartOfAccountHandler(db))
	router.POST("journalentry", CreateJournalEntryHandler(db))

	router.GET("/account", ListAccountHandler(db))
	router.GET("/coa", ListChartOfAccountHandler(db))
	router.GET("/accounttype", ListAccountTypeHandler(db))
	router.GET("/journalentry", ListJournalEntryHandler(db))
	router.GET("/profitandloss", ProfitAndLossHandler(db))
	router.GET("/balancesheet", BalanceSheetHandler(db))

	router.Run(":8000")
}
