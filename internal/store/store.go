package store

import (
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store interface {
	CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error)
	CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error)
	CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error)
	CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error)
	CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error)
	CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error)
	DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	DeleteCreditCard(ctx context.Context, id uuid.UUID) error
	DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error
	DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error
	GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error)
	GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error)
	GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error)
	GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error)
	GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error)
	GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error)
	HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error)
	HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error)
	ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error)
	ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error)
	ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error)
	ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error)
	ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error)
	ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error)
	PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error
	UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error)
	UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error)
	UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error)
	UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error)
	UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error)
	UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error)
	GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error)
	GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error)
	GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error)
	GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error)
	GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error)
	GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error)
	ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error)
	WithTx(tx pgx.Tx) *pgstore.Queries
}
