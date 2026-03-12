package api

import (
	"financialcontrol/internal/handlers"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterRoutes(t *testing.T) {
	router := gin.New()

	// Create a minimal Api with mock handlers
	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	// Get all registered routes
	routes := router.Routes()

	// Verify expected routes are registered
	expectedRoutes := map[string][]string{
		"/engine/v1/categories/":    {"POST", "GET"},
		"/engine/v1/categories/:id": {"GET", "PUT", "DELETE"},

		"/engine/v1/creditcards/":    {"POST", "GET"},
		"/engine/v1/creditcards/:id": {"GET", "PUT", "DELETE"},

		"/engine/v1/transactions/":        {"POST", "GET"},
		"/engine/v1/transactions/:id":     {"GET", "PUT", "DELETE"},
		"/engine/v1/transactions/pay/:id": {"PUT"},

		"/engine/v1/monthly_transactions/":    {"POST", "GET"},
		"/engine/v1/monthly_transactions/:id": {"GET", "PUT", "DELETE"},

		"/engine/v1/annual_transactions/":    {"POST", "GET"},
		"/engine/v1/annual_transactions/:id": {"GET", "PUT", "DELETE"},

		"/engine/v1/installment_transactions/":    {"POST", "GET"},
		"/engine/v1/installment_transactions/:id": {"GET", "PUT", "DELETE"},
		"/engine/v1/monthly_report/":              {"GET"},
	}

	// Create a map of registered routes for easier checking
	registeredRoutes := make(map[string]map[string]bool)
	for _, route := range routes {
		if registeredRoutes[route.Path] == nil {
			registeredRoutes[route.Path] = make(map[string]bool)
		}
		registeredRoutes[route.Path][route.Method] = true
	}

	// Verify all expected routes are registered
	for path, methods := range expectedRoutes {
		for _, method := range methods {
			if !registeredRoutes[path][method] {
				t.Errorf("Expected route %s %s to be registered", method, path)
			}
		}
	}
}

func TestRegisterRoutesCount(t *testing.T) {
	router := gin.New()

	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	routes := router.Routes()

	// Expected: 5 routes per resource (POST, GET, GET/:id, PUT/:id, DELETE/:id)
	// categories: 5, creditcards: 5, transactions: 10 (monthly, annual, installment, report, pay extras), monthly: 5, annual: 5, installment: 5, monthly_report: 1
	// Total: 36
	expectedCount := 36
	if len(routes) != expectedCount {
		t.Errorf("Expected %d routes, got %d", expectedCount, len(routes))
	}
}

func TestRegisterRoutesBaseGroup(t *testing.T) {
	router := gin.New()

	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	routes := router.Routes()

	// All routes should start with /engine/v1
	for _, route := range routes {
		if len(route.Path) < 11 || route.Path[:11] != "/engine/v1/" {
			t.Errorf("Expected route to start with /engine/v1/, got %s", route.Path)
		}
	}
}

func TestCategoriesRoutes(t *testing.T) {
	router := gin.New()

	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	routes := router.Routes()

	categoriesRoutes := 0
	for _, route := range routes {
		if strings.Contains(route.Path, "/categories") {
			categoriesRoutes++
		}
	}

	if categoriesRoutes != 5 {
		t.Errorf("Expected 5 categories routes, got %d", categoriesRoutes)
	}
}

func TestCreditcardsRoutes(t *testing.T) {
	router := gin.New()

	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	routes := router.Routes()

	creditcardsRoutes := 0
	for _, route := range routes {
		if strings.Contains(route.Path, "/creditcards") {
			creditcardsRoutes++
		}
	}

	if creditcardsRoutes != 5 {
		t.Errorf("Expected 5 creditcards routes, got %d", creditcardsRoutes)
	}
}

func TestTransactionsRoutes(t *testing.T) {
	router := gin.New()

	api := Api{
		Router:                         router,
		categoriesHandler:              &handlers.Category{},
		creditCardsHandler:             &handlers.CreditCard{},
		transactionsHandler:            &handlers.Transaction{},
		monthlyTransactionsHandler:     &handlers.MonthlyTransaction{},
		annualTransactionsHandler:      &handlers.AnnualTransaction{},
		installmentTransactionsHandler: &handlers.InstallmentTransaction{},
		monthlyReportHandler:           &handlers.MonthlyReport{},
	}

	api.RegisterRoutes()

	routes := router.Routes()

	transactionsRoutes := 0
	for _, route := range routes {
		if strings.Contains(route.Path, "/transactions") && !strings.Contains(route.Path, "monthly") && !strings.Contains(route.Path, "annual") && !strings.Contains(route.Path, "installment") {
			transactionsRoutes++
		}
	}

	// 7 routes because of the extra /pay/:id and /report routes
	if transactionsRoutes != 7 {
		t.Errorf("Expected 7 transactions routes, got %d", transactionsRoutes)
	}
}
