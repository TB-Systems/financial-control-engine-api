package api

import (
	"testing"

	"financialcontrol/internal/repositories"
	"financialcontrol/internal/store/pgstore"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewApi(t *testing.T) {
	router := gin.New()

	api := NewApi(router, nil)

	if api.Router == nil {
		t.Error("Expected Router to be set")
	}
	if api.categoriesHandler == nil {
		t.Error("Expected categoriesHandler to be set")
	}
	if api.creditCardsHandler == nil {
		t.Error("Expected creditCardsHandler to be set")
	}
	if api.transactionsHandler == nil {
		t.Error("Expected transactionsHandler to be set")
	}
	if api.monthlyTransactionsHandler == nil {
		t.Error("Expected monthlyTransactionsHandler to be set")
	}
	if api.annualTransactionsHandler == nil {
		t.Error("Expected annualTransactionsHandler to be set")
	}
	if api.installmentTransactionsHandler == nil {
		t.Error("Expected installmentTransactionsHandler to be set")
	}
	if api.monthlyReportHandler == nil {
		t.Error("Expected monthlyReportHandler to be set")
	}
}

func TestApiStruct(t *testing.T) {
	// Test that Api struct has the expected shape
	api := Api{}

	// Verify all handler fields exist (they will be nil)
	if api.Router != nil {
		t.Error("Expected Router to be nil on empty struct")
	}
	if api.categoriesHandler != nil {
		t.Error("Expected categoriesHandler to be nil on empty struct")
	}
	if api.creditCardsHandler != nil {
		t.Error("Expected creditCardsHandler to be nil on empty struct")
	}
	if api.transactionsHandler != nil {
		t.Error("Expected transactionsHandler to be nil on empty struct")
	}
	if api.monthlyTransactionsHandler != nil {
		t.Error("Expected monthlyTransactionsHandler to be nil on empty struct")
	}
	if api.annualTransactionsHandler != nil {
		t.Error("Expected annualTransactionsHandler to be nil on empty struct")
	}
	if api.installmentTransactionsHandler != nil {
		t.Error("Expected installmentTransactionsHandler to be nil on empty struct")
	}
	if api.monthlyReportHandler != nil {
		t.Error("Expected monthlyReportHandler to be nil on empty struct")
	}
}

func TestCreateHandlersFactories(t *testing.T) {
	repository := repositories.NewRepository(pgstore.New(nil))

	if createCategory(repository) == nil {
		t.Error("Expected createCategory to return handler")
	}
	if createCreditCard(repository) == nil {
		t.Error("Expected createCreditCard to return handler")
	}
	if createTransactions(repository) == nil {
		t.Error("Expected createTransactions to return handler")
	}
	if createMonthlyTransactions(repository) == nil {
		t.Error("Expected createMonthlyTransactions to return handler")
	}
	if createAnnualTransactions(repository) == nil {
		t.Error("Expected createAnnualTransactions to return handler")
	}
	if createInstallmentTransactions(repository) == nil {
		t.Error("Expected createInstallmentTransactions to return handler")
	}
	if createMonthlyReport(repository) == nil {
		t.Error("Expected createMonthlyReport to return handler")
	}
}
