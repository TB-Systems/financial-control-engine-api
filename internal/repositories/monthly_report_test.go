package repositories

import (
	"context"
	"errors"
	"testing"

	"backend-commons/models"
	"financialcontrol/internal/store/pgstore"
	"github.com/TB-Systems/go-commons/utils"

	"github.com/google/uuid"
)

type MonthlyReportStoreMock struct {
	*SimpleStoreMock

	monthlyBalanceResult pgstore.GetMonthlyBalanceRow
	monthlyBalanceErr    error
	monthlyBalanceArg    pgstore.GetMonthlyBalanceParams

	categoriesResult []pgstore.GetCategoriesSpendingRow
	categoriesErr    error
	categoriesArg    pgstore.GetCategoriesSpendingParams

	creditCardsResult []pgstore.GetCreditCardsSpendingRow
	creditCardsErr    error
	creditCardsArg    pgstore.GetCreditCardsSpendingParams
}

func newMonthlyReportStoreMock() *MonthlyReportStoreMock {
	return &MonthlyReportStoreMock{SimpleStoreMock: &SimpleStoreMock{}}
}

func (m *MonthlyReportStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	m.monthlyBalanceArg = arg
	if m.monthlyBalanceErr != nil {
		return pgstore.GetMonthlyBalanceRow{}, m.monthlyBalanceErr
	}
	return m.monthlyBalanceResult, nil
}

func (m *MonthlyReportStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	m.categoriesArg = arg
	if m.categoriesErr != nil {
		return nil, m.categoriesErr
	}
	return m.categoriesResult, nil
}

func (m *MonthlyReportStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	m.creditCardsArg = arg
	if m.creditCardsErr != nil {
		return nil, m.creditCardsErr
	}
	return m.creditCardsResult, nil
}

func TestGetMonthlyBalanceSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newMonthlyReportStoreMock()
	mock.monthlyBalanceResult = pgstore.GetMonthlyBalanceRow{
		TotalIncome: utils.Float64ToNumeric(1000),
		TotalDebit:  utils.Float64ToNumeric(300),
		TotalCredit: utils.Float64ToNumeric(100),
		Balance:     utils.Float64ToNumeric(600),
	}

	repo := NewRepository(mock)

	result, err := repo.GetMonthlyBalance(ctx, userID, 2025, 9)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.TotalIncome != 1000 || result.TotalDebit != 300 || result.TotalCredit != 100 || result.Balance != 600 {
		t.Errorf("Unexpected balance mapping: %+v", result)
	}

	if mock.monthlyBalanceArg.UserID != userID || mock.monthlyBalanceArg.Column2 != 2025 || mock.monthlyBalanceArg.Column3 != 9 {
		t.Errorf("Unexpected params: %+v", mock.monthlyBalanceArg)
	}
}

func TestGetMonthlyBalanceError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newMonthlyReportStoreMock()
	mock.monthlyBalanceErr = errors.New("db error")

	repo := NewRepository(mock)

	_, err := repo.GetMonthlyBalance(ctx, userID, 2025, 9)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetCategoriesSpendingSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := newMonthlyReportStoreMock()
	mock.categoriesResult = []pgstore.GetCategoriesSpendingRow{
		{
			CategoryID:              categoryID,
			CategoryName:            "Food",
			CategoryIcon:            "restaurant",
			CategoryTransactionType: int32(models.Debit),
			TotalSpent:              utils.Float64ToNumeric(250),
		},
	}

	repo := NewRepository(mock)

	result, err := repo.GetCategoriesSpending(ctx, userID, 2025, 9)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].CategoryID != categoryID || result[0].TotalSpent != 250 {
		t.Errorf("Unexpected categories mapping: %+v", result)
	}

	if mock.categoriesArg.UserID != userID || mock.categoriesArg.Column2 != 2025 || mock.categoriesArg.Column3 != 9 {
		t.Errorf("Unexpected params: %+v", mock.categoriesArg)
	}
}

func TestGetCategoriesSpendingEmptyAndError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newMonthlyReportStoreMock()
	repo := NewRepository(mock)

	result, err := repo.GetCategoriesSpending(ctx, userID, 2025, 9)
	if err != nil {
		t.Fatalf("Expected no error for empty response, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %+v", result)
	}

	mock.categoriesErr = errors.New("categories error")
	_, err = repo.GetCategoriesSpending(ctx, userID, 2025, 9)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetCreditCardsSpendingSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	cardID := uuid.New()

	mock := newMonthlyReportStoreMock()
	mock.creditCardsResult = []pgstore.GetCreditCardsSpendingRow{
		{
			CreditCardID:               cardID,
			CreditCardName:             "Main Card",
			CreditCardFirstFourNumbers: "1234",
			CreditCardLimit:            5000,
			CreditCardCloseDay:         10,
			CreditCardExpireDay:        20,
			CreditCardBackgroundColor:  "#000",
			CreditCardTextColor:        "#fff",
			TotalSpent:                 utils.Float64ToNumeric(900),
		},
	}

	repo := NewRepository(mock)

	result, err := repo.GetCreditCardsSpending(ctx, userID, 2025, 9)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].ID != cardID || result[0].TotalSpent != 900 {
		t.Errorf("Unexpected credit cards mapping: %+v", result)
	}

	if mock.creditCardsArg.UserID != userID || mock.creditCardsArg.Column2 != 2025 || mock.creditCardsArg.Column3 != 9 {
		t.Errorf("Unexpected params: %+v", mock.creditCardsArg)
	}
}

func TestGetCreditCardsSpendingEmptyAndError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newMonthlyReportStoreMock()
	repo := NewRepository(mock)

	result, err := repo.GetCreditCardsSpending(ctx, userID, 2025, 9)
	if err != nil {
		t.Fatalf("Expected no error for empty response, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result, got %+v", result)
	}

	mock.creditCardsErr = errors.New("credit cards error")
	_, err = repo.GetCreditCardsSpending(ctx, userID, 2025, 9)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
