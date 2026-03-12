package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"backend-commons/constants"
	"backend-commons/models"

	"github.com/google/uuid"
)

type MonthlyReportRepositoryMock struct {
	balanceResult    models.MonthlyReport
	balanceErr       error
	categoriesResult []models.CategoriesSpending
	categoriesErr    error
	cardsResult      []models.CreditCardsSpending
	cardsErr         error
}

func (m *MonthlyReportRepositoryMock) GetMonthlyBalance(ctx context.Context, userID uuid.UUID, year int32, month int32) (models.MonthlyReport, error) {
	if m.balanceErr != nil {
		return models.MonthlyReport{}, m.balanceErr
	}
	return m.balanceResult, nil
}

func (m *MonthlyReportRepositoryMock) GetCategoriesSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CategoriesSpending, error) {
	if m.categoriesErr != nil {
		return nil, m.categoriesErr
	}
	return m.categoriesResult, nil
}

func (m *MonthlyReportRepositoryMock) GetCreditCardsSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CreditCardsSpending, error) {
	if m.cardsErr != nil {
		return nil, m.cardsErr
	}
	return m.cardsResult, nil
}

func TestGenerateMonthlyReportSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	cardID := uuid.New()

	mock := &MonthlyReportRepositoryMock{
		balanceResult: models.MonthlyReport{TotalIncome: 1000, TotalDebit: 300, TotalCredit: 100, Balance: 600},
		categoriesResult: []models.CategoriesSpending{{
			CategoryID:              categoryID,
			CategoryName:            "Food",
			CategoryIcon:            "restaurant",
			CategoryTransactionType: models.Debit,
			TotalSpent:              250,
		}},
		cardsResult: []models.CreditCardsSpending{{
			ID:               cardID,
			Name:             "Main Card",
			FirstFourNumbers: "1234",
			Limit:            5000,
			CloseDay:         10,
			ExpireDay:        20,
			BackgroundColor:  "#000",
			TextColor:        "#fff",
			TotalSpent:       900,
		}},
	}

	service := NewMonthlyReportService(mock)

	result, apiErr := service.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr != nil {
		t.Fatalf("Expected no error, got %v", apiErr)
	}

	if result.Balance != 600 || result.TotalIncome != 1000 {
		t.Errorf("Unexpected summary: %+v", result)
	}

	if result.MostSpentCategory == nil || result.MostSpentCategory.ID != categoryID {
		t.Errorf("Expected most spent category to be present")
	}

	if result.MostSpentCreditCard == nil || result.MostSpentCreditCard.ID != cardID {
		t.Errorf("Expected most spent credit card to be present")
	}
}

func TestGenerateMonthlyReportBalanceError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := &MonthlyReportRepositoryMock{balanceErr: errors.New("balance error")}
	service := NewMonthlyReportService(mock)

	_, apiErr := service.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr == nil {
		t.Fatal("Expected api error, got nil")
	}

	if apiErr.GetStatus() != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestGenerateMonthlyReportCategoriesNoRowsIgnored(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := &MonthlyReportRepositoryMock{
		balanceResult: models.MonthlyReport{TotalIncome: 100, TotalDebit: 20, TotalCredit: 10, Balance: 70},
		categoriesErr: errors.New(constants.StoreErrorNoRowsMsg),
	}
	service := NewMonthlyReportService(mock)

	result, apiErr := service.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr != nil {
		t.Fatalf("Expected no error, got %v", apiErr)
	}

	if len(result.Categories) != 0 {
		t.Errorf("Expected empty categories when no rows")
	}
}

func TestGenerateMonthlyReportCategoriesUnexpectedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := &MonthlyReportRepositoryMock{
		balanceResult: models.MonthlyReport{TotalIncome: 100, TotalDebit: 20, TotalCredit: 10, Balance: 70},
		categoriesErr: errors.New("db error"),
	}
	service := NewMonthlyReportService(mock)

	_, apiErr := service.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr == nil {
		t.Fatal("Expected api error, got nil")
	}

	if apiErr.GetStatus() != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestGenerateMonthlyReportCreditCardsNoRowsIgnoredAndUnexpectedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	serviceNoRows := NewMonthlyReportService(&MonthlyReportRepositoryMock{
		balanceResult: models.MonthlyReport{TotalIncome: 100, TotalDebit: 20, TotalCredit: 10, Balance: 70},
		cardsErr:      errors.New(constants.StoreErrorNoRowsMsg),
	})

	result, apiErr := serviceNoRows.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr != nil {
		t.Fatalf("Expected no error, got %v", apiErr)
	}
	if len(result.CreditCards) != 0 {
		t.Errorf("Expected empty credit cards when no rows")
	}

	serviceErr := NewMonthlyReportService(&MonthlyReportRepositoryMock{
		balanceResult: models.MonthlyReport{TotalIncome: 100, TotalDebit: 20, TotalCredit: 10, Balance: 70},
		cardsErr:      errors.New("db error"),
	})

	_, apiErr = serviceErr.GenerateMonthlyReport(ctx, userID, 2025, 9)
	if apiErr == nil {
		t.Fatal("Expected api error, got nil")
	}
	if apiErr.GetStatus() != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}
