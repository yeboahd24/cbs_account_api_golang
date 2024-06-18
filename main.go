package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "ledger_api/docs" // Import your docs package for Swagger
)

// @title LedgerAPI
// @version 1.0
// @description API for Ledger management
// @host localhost:8000
// @securityDefinitions.apikey JWT
// @in header
// @name token

// @scheme http
func main() {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(JWTAuthMiddleware(), JSONMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1/cbs/accounting"
	docs.SwaggerInfo.Schemes = []string{"http"}
	v1 := router.Group("/api/v1/cbs/accounting")

	v1.POST("/account", CreateAccountHandler(db))

	v1.POST("/accounttype", CreateAccountTypeHandler(db))

	v1.POST("/coa", CreateChartOfAccountHandler(db))

	v1.POST("/journalentry", CreateJournalEntryHandler(db))

	v1.GET("/account", ListAccountHandler(db))
	v1.GET("/coa", ListChartOfAccountHandler(db))
	v1.GET("/accounttype", ListAccountTypeHandler(db))
	v1.GET("/journalentry", ListJournalEntryHandler(db))
	v1.GET("/profitandloss", ProfitAndLossHandler(db))
	v1.GET("/balancesheet", BalanceSheetHandler(db))

	// Swagger endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))

	router.Run(":8000")
}

