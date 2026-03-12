package repositories

import (
	"context"
	"testing"

	"financialcontrol/internal/store/pgstore"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============= SIMPLE STORE MOCK =============

type SimpleStoreMock struct{}

func (m *SimpleStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *SimpleStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *SimpleStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *SimpleStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *SimpleStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *SimpleStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *SimpleStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *SimpleStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *SimpleStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *SimpleStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *SimpleStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *SimpleStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *SimpleStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *SimpleStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *SimpleStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *SimpleStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *SimpleStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *SimpleStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *SimpleStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *SimpleStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *SimpleStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *SimpleStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *SimpleStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *SimpleStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *SimpleStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *SimpleStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *SimpleStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *SimpleStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *SimpleStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *SimpleStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *SimpleStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= TESTS =============

func TestNewRepositoryCreatesRepository(t *testing.T) {
	mock := &SimpleStoreMock{}

	repo := NewRepository(mock)

	if repo.store == nil {
		t.Error("Expected store to be set, got nil")
	}
}

func TestNewRepositoryStoreIsAssigned(t *testing.T) {
	mock := &SimpleStoreMock{}

	repo := NewRepository(mock)

	if repo.store != mock {
		t.Error("Expected store to be the mock, got different value")
	}
}
