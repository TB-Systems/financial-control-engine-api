package api

import "financialcontrol/internal/middlewares"

func (a *Api) RegisterRoutes() {
	api := a.Router.Group("/engine/v1")

	api.Use(middlewares.UserIDMiddleware())

	categories := api.Group("/categories")
	{
		categories.POST("/", a.categoriesHandler.Create())
		categories.GET("/", a.categoriesHandler.Read())
		categories.GET("/:id", a.categoriesHandler.ReadByID())
		categories.PUT("/:id", a.categoriesHandler.Update())
		categories.DELETE("/:id", a.categoriesHandler.Delete())
	}

	creditcards := api.Group("/creditcards")
	{
		creditcards.POST("/", a.creditCardsHandler.Create())
		creditcards.GET("/", a.creditCardsHandler.Read())
		creditcards.GET("/:id", a.creditCardsHandler.ReadAt())
		creditcards.PUT("/:id", a.creditCardsHandler.Update())
		creditcards.DELETE("/:id", a.creditCardsHandler.Delete())
	}

	transactions := api.Group("/transactions")
	{
		transactions.POST("/", a.transactionsHandler.Create())
		transactions.POST("/monthly", a.transactionsHandler.CreateFromMonthlyTransaction())
		transactions.POST("/annual", a.transactionsHandler.CreateFromAnnualTransaction())
		transactions.POST("/installment", a.transactionsHandler.CreateFromInstallmentTransaction())
		transactions.GET("/", a.transactionsHandler.Read())
		transactions.GET("/:id", a.transactionsHandler.ReadById())
		transactions.GET("/report", a.transactionsHandler.ReadByMonthAndYear())
		transactions.PUT("/:id", a.transactionsHandler.Update())
		transactions.DELETE("/:id", a.transactionsHandler.Delete())
		transactions.PUT("/pay/:id", a.transactionsHandler.Pay())
	}

	monthlyTransactions := api.Group("/monthly_transactions")
	{
		monthlyTransactions.POST("/", a.monthlyTransactionsHandler.Create())
		monthlyTransactions.GET("/", a.monthlyTransactionsHandler.Read())
		monthlyTransactions.GET("/:id", a.monthlyTransactionsHandler.ReadById())
		monthlyTransactions.PUT("/:id", a.monthlyTransactionsHandler.Update())
		monthlyTransactions.DELETE("/:id", a.monthlyTransactionsHandler.Delete())
	}

	annualTransactions := api.Group("/annual_transactions")
	{
		annualTransactions.POST("/", a.annualTransactionsHandler.Create())
		annualTransactions.GET("/", a.annualTransactionsHandler.Read())
		annualTransactions.GET("/:id", a.annualTransactionsHandler.ReadById())
		annualTransactions.PUT("/:id", a.annualTransactionsHandler.Update())
		annualTransactions.DELETE("/:id", a.annualTransactionsHandler.Delete())
	}

	installmentTransactions := api.Group("/installment_transactions")
	{
		installmentTransactions.POST("/", a.installmentTransactionsHandler.Create())
		installmentTransactions.GET("/", a.installmentTransactionsHandler.Read())
		installmentTransactions.GET("/:id", a.installmentTransactionsHandler.ReadById())
		installmentTransactions.PUT("/:id", a.installmentTransactionsHandler.Update())
		installmentTransactions.DELETE("/:id", a.installmentTransactionsHandler.Delete())
	}

	monthlyReport := api.Group("/monthly_report")
	{
		monthlyReport.GET("/", a.monthlyReportHandler.GetMonthlyBalance())
	}
}
