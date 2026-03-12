package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"backend-commons/models"
	"financialcontrol/internal/store/pgstore"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============= STORE MOCK FOR TRANSACTION =============

type TransactionStoreMock struct {
	Error                    error
	CreateError              error
	UpdateError              error
	DeleteError              error
	PayError                 error
	TransactionResult        pgstore.CreateTransactionRow
	TransactionByIdResult    pgstore.GetTransactionByIDRow
	TransactionUpdatedResult pgstore.Transaction
	PaginatedResult          []pgstore.ListTransactionsByUserIDPaginatedRow
	DateRangeResult          []pgstore.ListTransactionsByUserAndDateRow
	MonthYearResult          []pgstore.ListTransactionsByUserAndMonthYearPaginatedRow
}

// Transaction methods
func (m *TransactionStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	if m.CreateError != nil {
		return pgstore.CreateTransactionRow{}, m.CreateError
	}
	if m.Error != nil {
		return pgstore.CreateTransactionRow{}, m.Error
	}
	return m.TransactionResult, nil
}

func (m *TransactionStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.PaginatedResult, nil
}

func (m *TransactionStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.DateRangeResult, nil
}

func (m *TransactionStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	if m.Error != nil {
		return pgstore.GetTransactionByIDRow{}, m.Error
	}
	return m.TransactionByIdResult, nil
}

func (m *TransactionStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	if m.UpdateError != nil {
		return pgstore.Transaction{}, m.UpdateError
	}
	if m.Error != nil {
		return pgstore.Transaction{}, m.Error
	}
	return m.TransactionUpdatedResult, nil
}

func (m *TransactionStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

func (m *TransactionStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	if m.PayError != nil {
		return m.PayError
	}
	return m.Error
}

// Other required interface methods (stubs)
func (m *TransactionStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *TransactionStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *TransactionStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *TransactionStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *TransactionStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *TransactionStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *TransactionStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *TransactionStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *TransactionStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *TransactionStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *TransactionStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *TransactionStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *TransactionStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *TransactionStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *TransactionStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *TransactionStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *TransactionStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *TransactionStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *TransactionStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *TransactionStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *TransactionStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *TransactionStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *TransactionStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *TransactionStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *TransactionStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.MonthYearResult, nil
}
func (m *TransactionStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *TransactionStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *TransactionStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *TransactionStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *TransactionStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *TransactionStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *TransactionStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *TransactionStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *TransactionStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *TransactionStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *TransactionStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *TransactionStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptzTransaction(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func createNumeric(value float64) pgtype.Numeric {
	var num pgtype.Numeric
	num.Scan(value)
	return num
}

func createPgUUIDTransaction(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func newTransactionStoreMock() *TransactionStoreMock {
	return &TransactionStoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionResult = pgstore.CreateTransactionRow{
		ID:        transactionID,
		Name:      "Test Transaction",
		Date:      createTimestamptzTransaction(now),
		Value:     createNumeric(100.50),
		Paid:      true,
		CreatedAt: createTimestamptzTransaction(now),
		UpdatedAt: createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	data := models.CreateTransaction{
		UserID:     userID,
		Name:       "Test Transaction",
		Date:       now,
		Value:      100.50,
		Paid:       true,
		CategoryID: categoryID,
	}

	result, err := repo.CreateTransaction(ctx, data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Name != "Test Transaction" {
		t.Errorf("Expected name 'Test Transaction', got '%s'", result.Name)
	}
}

func TestCreateTransactionError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.CreateError = errors.New("database error")

	repo := NewRepository(mock)

	data := models.CreateTransaction{
		UserID:     userID,
		Name:       "Test Transaction",
		Date:       now,
		Value:      100.50,
		Paid:       true,
		CategoryID: categoryID,
	}

	_, err := repo.CreateTransaction(ctx, data)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCreateTransactionWithOptionalFields(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	monthlyID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionResult = pgstore.CreateTransactionRow{
		ID:        transactionID,
		Name:      "Test Transaction",
		Date:      createTimestamptzTransaction(now),
		Value:     createNumeric(100.50),
		Paid:      false,
		CreatedAt: createTimestamptzTransaction(now),
		UpdatedAt: createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	data := models.CreateTransaction{
		UserID:               userID,
		Name:                 "Test Transaction",
		Date:                 now,
		Value:                100.50,
		Paid:                 false,
		CategoryID:           categoryID,
		CreditcardID:         &creditCardID,
		MonthlyTransactionID: &monthlyID,
	}

	result, err := repo.CreateTransaction(ctx, data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

// ============= READ PAGINATED TESTS =============

func TestReadTransactionsSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.PaginatedResult = []pgstore.ListTransactionsByUserIDPaginatedRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Transaction 1",
			Date:                    createTimestamptzTransaction(now),
			Value:                   createNumeric(50.0),
			Paid:                    true,
			CategoryID:              createPgUUIDTransaction(categoryID),
			CategoryTransactionType: pgtype.Int4{Int32: 0, Valid: true},
			CategoryName:            pgtype.Text{String: "Category", Valid: true},
			CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
			CreatedAt:               createTimestamptzTransaction(now),
			UpdatedAt:               createTimestamptzTransaction(now),
			TotalCount:              2,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Transaction 2",
			Date:                    createTimestamptzTransaction(now),
			Value:                   createNumeric(75.0),
			Paid:                    false,
			CategoryID:              createPgUUIDTransaction(categoryID),
			CategoryTransactionType: pgtype.Int4{Int32: 1, Valid: true},
			CategoryName:            pgtype.Text{String: "Category 2", Valid: true},
			CategoryIcon:            pgtype.Text{String: "icon2.png", Valid: true},
			CreatedAt:               createTimestamptzTransaction(now),
			UpdatedAt:               createTimestamptzTransaction(now),
			TotalCount:              2,
		},
	}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadTransactions(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(result))
	}

	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestReadTransactionsEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newTransactionStoreMock()
	mock.PaginatedResult = []pgstore.ListTransactionsByUserIDPaginatedRow{}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadTransactions(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(result))
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestReadTransactionsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newTransactionStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	_, _, err := repo.ReadTransactions(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReadTransactionsWithMonthlyTransactionAndNoCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	monthlyID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.PaginatedResult = []pgstore.ListTransactionsByUserIDPaginatedRow{
		{
			ID:                           uuid.New(),
			UserID:                       userID,
			Name:                         "Monthly Transaction Instance",
			Date:                         createTimestamptzTransaction(now),
			Value:                        createNumeric(40.0),
			Paid:                         false,
			CategoryID:                   createPgUUIDTransaction(categoryID),
			CategoryTransactionType:      pgtype.Int4{Int32: 0, Valid: true},
			CategoryName:                 pgtype.Text{String: "Category", Valid: true},
			CategoryIcon:                 pgtype.Text{String: "icon.png", Valid: true},
			MonthlyTransactionsID:        createPgUUIDTransaction(monthlyID),
			MonthlyTransactionsDay:       pgtype.Int4{Int32: 10, Valid: true},
			MonthlyTransactionsName:      pgtype.Text{String: "Rent", Valid: true},
			MonthlyTransactionsValue:     createNumeric(40.0),
			MonthlyTransactionsCreatedAt: createTimestamptzTransaction(now),
			MonthlyTransactionsUpdatedAt: createTimestamptzTransaction(now),
			CreatedAt:                    createTimestamptzTransaction(now),
			UpdatedAt:                    createTimestamptzTransaction(now),
			TotalCount:                   1,
		},
	}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadTransactions(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(result))
	}

	if result[0].MonthlyTransaction == nil {
		t.Fatal("Expected monthly transaction, got nil")
	}

	if result[0].MonthlyTransaction.CreditCardID != nil {
		t.Error("Expected monthly transaction credit card ID to be nil")
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

// ============= READ BY DATE RANGE TESTS =============

func TestReadTransactionsInToDatesSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	mock := newTransactionStoreMock()
	mock.DateRangeResult = []pgstore.ListTransactionsByUserAndDateRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Transaction 1",
			Date:                    createTimestamptzTransaction(now),
			Value:                   createNumeric(100.0),
			Paid:                    true,
			CategoryID:              createPgUUIDTransaction(categoryID),
			CategoryTransactionType: pgtype.Int4{Int32: 0, Valid: true},
			CategoryName:            pgtype.Text{String: "Category", Valid: true},
			CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
			CreatedAt:               createTimestamptzTransaction(now),
			UpdatedAt:               createTimestamptzTransaction(now),
			TotalCount:              1,
		},
	}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, count, err := repo.ReadTransactionsInToDates(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result))
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestReadTransactionsInToDatesEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.DateRangeResult = []pgstore.ListTransactionsByUserAndDateRow{}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now,
	}

	result, count, err := repo.ReadTransactionsInToDates(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(result))
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestReadTransactionsInToDatesError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now,
	}

	_, _, err := repo.ReadTransactionsInToDates(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ BY MONTH/YEAR TESTS =============

func TestReadTransactionsByMonthYearSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.MonthYearResult = []pgstore.ListTransactionsByUserAndMonthYearPaginatedRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Transaction month/year",
			Date:                    createTimestamptzTransaction(now),
			Value:                   createNumeric(120.0),
			Paid:                    true,
			CategoryID:              createPgUUIDTransaction(categoryID),
			CategoryTransactionType: pgtype.Int4{Int32: 1, Valid: true},
			CategoryName:            pgtype.Text{String: "Category", Valid: true},
			CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
			CreatedAt:               createTimestamptzTransaction(now),
			UpdatedAt:               createTimestamptzTransaction(now),
			TotalCount:              1,
		},
	}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithMonthYear{
		UserID:     userID,
		Year:       2025,
		Month:      9,
		PageLimit:  10,
		PageOffset: 0,
	}

	result, count, err := repo.ReadTransactionsByMonthYear(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result))
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestReadTransactionsByMonthYearEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newTransactionStoreMock()
	mock.MonthYearResult = []pgstore.ListTransactionsByUserAndMonthYearPaginatedRow{}

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithMonthYear{
		UserID:     userID,
		Year:       2025,
		Month:      9,
		PageLimit:  10,
		PageOffset: 0,
	}

	result, count, err := repo.ReadTransactionsByMonthYear(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 transactions, got %d", len(result))
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestReadTransactionsByMonthYearError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newTransactionStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	params := commonsmodels.PaginatedParamsWithMonthYear{
		UserID:     userID,
		Year:       2025,
		Month:      9,
		PageLimit:  10,
		PageOffset: 0,
	}

	_, _, err := repo.ReadTransactionsByMonthYear(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadTransactionByIdSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionByIdResult = pgstore.GetTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Test Transaction",
		Date:                    createTimestamptzTransaction(now),
		Value:                   createNumeric(150.0),
		Paid:                    true,
		CategoryID:              createPgUUIDTransaction(categoryID),
		CategoryTransactionType: pgtype.Int4{Int32: 0, Valid: true},
		CategoryName:            pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
		CreatedAt:               createTimestamptzTransaction(now),
		UpdatedAt:               createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	result, err := repo.ReadTransactionById(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Name != "Test Transaction" {
		t.Errorf("Expected name 'Test Transaction', got '%s'", result.Name)
	}
}

func TestReadTransactionByIdWithCreditCard(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionByIdResult = pgstore.GetTransactionByIDRow{
		ID:                         transactionID,
		UserID:                     userID,
		Name:                       "Credit Card Transaction",
		Date:                       createTimestamptzTransaction(now),
		Value:                      createNumeric(200.0),
		Paid:                       false,
		CategoryID:                 createPgUUIDTransaction(categoryID),
		CategoryTransactionType:    pgtype.Int4{Int32: 2, Valid: true},
		CategoryName:               pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:               pgtype.Text{String: "icon.png", Valid: true},
		CreditcardID:               createPgUUIDTransaction(creditCardID),
		CreditcardName:             pgtype.Text{String: "My Card", Valid: true},
		CreditcardFirstFourNumbers: pgtype.Text{String: "1234", Valid: true},
		CreditcardCreditLimit:      pgtype.Float8{Float64: 5000.0, Valid: true},
		CreditcardCloseDay:         pgtype.Int4{Int32: 15, Valid: true},
		CreditcardExpireDay:        pgtype.Int4{Int32: 25, Valid: true},
		CreditcardBackgroundColor:  pgtype.Text{String: "#FF0000", Valid: true},
		CreditcardTextColor:        pgtype.Text{String: "#FFFFFF", Valid: true},
		CreatedAt:                  createTimestamptzTransaction(now),
		UpdatedAt:                  createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	result, err := repo.ReadTransactionById(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard == nil {
		t.Error("Expected credit card, got nil")
	}

	if result.Creditcard != nil && result.Creditcard.Name != "My Card" {
		t.Errorf("Expected credit card name 'My Card', got '%s'", result.Creditcard.Name)
	}
}

func TestReadTransactionByIdError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()
	mock.Error = errors.New("not found")

	repo := NewRepository(mock)

	_, err := repo.ReadTransactionById(ctx, transactionID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionUpdatedResult = pgstore.Transaction{
		ID:        transactionID,
		Name:      "Updated Transaction",
		Date:      createTimestamptzTransaction(now),
		Value:     createNumeric(200.0),
		Paid:      true,
		CreatedAt: createTimestamptzTransaction(now),
		UpdatedAt: createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	transaction := models.Transaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Updated Transaction",
		Date:   now,
		Value:  200.0,
		Paid:   true,
		Category: models.Category{
			ID: categoryID,
		},
	}

	result, err := repo.UpdateTransaction(ctx, transaction)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Name != "Updated Transaction" {
		t.Errorf("Expected name 'Updated Transaction', got '%s'", result.Name)
	}
}

func TestUpdateTransactionWithCreditCard(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionUpdatedResult = pgstore.Transaction{
		ID:        transactionID,
		Name:      "Updated Transaction",
		Date:      createTimestamptzTransaction(now),
		Value:     createNumeric(200.0),
		Paid:      false,
		CreatedAt: createTimestamptzTransaction(now),
		UpdatedAt: createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	transaction := models.Transaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Updated Transaction",
		Date:   now,
		Value:  200.0,
		Paid:   false,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateTransaction(ctx, transaction)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Name != "Updated Transaction" {
		t.Errorf("Expected name 'Updated Transaction', got '%s'", result.Name)
	}
}

func TestUpdateTransactionWithRecurringRelations(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	monthlyID := uuid.New()
	annualID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.TransactionUpdatedResult = pgstore.Transaction{
		ID:        transactionID,
		Name:      "Updated With Relations",
		Date:      createTimestamptzTransaction(now),
		Value:     createNumeric(250.0),
		Paid:      true,
		CreatedAt: createTimestamptzTransaction(now),
		UpdatedAt: createTimestamptzTransaction(now),
	}

	repo := NewRepository(mock)

	transaction := models.Transaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Updated With Relations",
		Date:   now,
		Value:  250.0,
		Paid:   true,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
		MonthlyTransaction:     &models.ShortMonthlyTransaction{ID: monthlyID},
		AnnualTransaction:      &models.ShortAnnualTransaction{ID: annualID},
		InstallmentTransaction: &models.ShortInstallmentTransaction{ID: installmentID},
	}

	result, err := repo.UpdateTransaction(ctx, transaction)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

func TestUpdateTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newTransactionStoreMock()
	mock.UpdateError = errors.New("database error")

	repo := NewRepository(mock)

	transaction := models.Transaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Updated Transaction",
		Date:   now,
		Value:  200.0,
		Paid:   true,
		Category: models.Category{
			ID: categoryID,
		},
	}

	_, err := repo.UpdateTransaction(ctx, transaction)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= DELETE TESTS =============

func TestDeleteTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()

	repo := NewRepository(mock)

	err := repo.DeleteTransaction(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()
	mock.DeleteError = errors.New("database error")

	repo := NewRepository(mock)

	err := repo.DeleteTransaction(ctx, transactionID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= PAY TESTS =============

func TestPayTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()

	repo := NewRepository(mock)

	err := repo.PayTransaction(ctx, transactionID, true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPayTransactionUnpay(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()

	repo := NewRepository(mock)

	err := repo.PayTransaction(ctx, transactionID, false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPayTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newTransactionStoreMock()
	mock.PayError = errors.New("database error")

	repo := NewRepository(mock)

	err := repo.PayTransaction(ctx, transactionID, true)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= CONVERTER TESTS =============

func TestStoreTransactionListToStoreTransaction(t *testing.T) {
	now := time.Now()
	transactionID := uuid.New()
	categoryID := uuid.New()

	row := pgstore.ListTransactionsByUserAndDateRow{
		ID:                      transactionID,
		UserID:                  uuid.New(),
		Name:                    "Test",
		Date:                    createTimestamptzTransaction(now),
		Value:                   createNumeric(100.0),
		Paid:                    true,
		CategoryID:              createPgUUIDTransaction(categoryID),
		CategoryTransactionType: pgtype.Int4{Int32: 0, Valid: true},
		CategoryName:            pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
		CreatedAt:               createTimestamptzTransaction(now),
		UpdatedAt:               createTimestamptzTransaction(now),
		TotalCount:              1,
	}

	result := storeTransactionListToStoreTransaction(row)

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

func TestStoreTransactionPaginatedToStoreTransaction(t *testing.T) {
	now := time.Now()
	transactionID := uuid.New()
	categoryID := uuid.New()

	row := pgstore.ListTransactionsByUserIDPaginatedRow{
		ID:                      transactionID,
		UserID:                  uuid.New(),
		Name:                    "Test",
		Date:                    createTimestamptzTransaction(now),
		Value:                   createNumeric(100.0),
		Paid:                    true,
		CategoryID:              createPgUUIDTransaction(categoryID),
		CategoryTransactionType: pgtype.Int4{Int32: 1, Valid: true},
		CategoryName:            pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
		CreatedAt:               createTimestamptzTransaction(now),
		UpdatedAt:               createTimestamptzTransaction(now),
		TotalCount:              1,
	}

	result := storeTransactionPaginatedToStoreTransaction(row)

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

func TestStoreTransactionToTransaction(t *testing.T) {
	now := time.Now()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()

	row := pgstore.GetTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Test Transaction",
		Date:                    createTimestamptzTransaction(now),
		Value:                   createNumeric(100.0),
		Paid:                    true,
		CategoryID:              createPgUUIDTransaction(categoryID),
		CategoryTransactionType: pgtype.Int4{Int32: 0, Valid: true},
		CategoryName:            pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:            pgtype.Text{String: "icon.png", Valid: true},
		CreatedAt:               createTimestamptzTransaction(now),
		UpdatedAt:               createTimestamptzTransaction(now),
	}

	result := storeTransactionToTransaction(row)

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Name != "Test Transaction" {
		t.Errorf("Expected name 'Test Transaction', got '%s'", result.Name)
	}

	if result.Category.ID != categoryID {
		t.Errorf("Expected category ID %v, got %v", categoryID, result.Category.ID)
	}

	if result.Creditcard != nil {
		t.Error("Expected no credit card, got one")
	}
}

func TestStoreTransactionToTransactionWithCreditCard(t *testing.T) {
	now := time.Now()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	row := pgstore.GetTransactionByIDRow{
		ID:                         transactionID,
		UserID:                     userID,
		Name:                       "Test Transaction",
		Date:                       createTimestamptzTransaction(now),
		Value:                      createNumeric(100.0),
		Paid:                       true,
		CategoryID:                 createPgUUIDTransaction(categoryID),
		CategoryTransactionType:    pgtype.Int4{Int32: 2, Valid: true},
		CategoryName:               pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:               pgtype.Text{String: "icon.png", Valid: true},
		CreditcardID:               createPgUUIDTransaction(creditCardID),
		CreditcardName:             pgtype.Text{String: "My Card", Valid: true},
		CreditcardFirstFourNumbers: pgtype.Text{String: "1234", Valid: true},
		CreditcardCreditLimit:      pgtype.Float8{Float64: 5000.0, Valid: true},
		CreditcardCloseDay:         pgtype.Int4{Int32: 15, Valid: true},
		CreditcardExpireDay:        pgtype.Int4{Int32: 25, Valid: true},
		CreditcardBackgroundColor:  pgtype.Text{String: "#000000", Valid: true},
		CreditcardTextColor:        pgtype.Text{String: "#FFFFFF", Valid: true},
		CreatedAt:                  createTimestamptzTransaction(now),
		UpdatedAt:                  createTimestamptzTransaction(now),
	}

	result := storeTransactionToTransaction(row)

	if result.Creditcard == nil {
		t.Error("Expected credit card, got nil")
	}

	if result.Creditcard != nil && result.Creditcard.ID != creditCardID {
		t.Errorf("Expected credit card ID %v, got %v", creditCardID, result.Creditcard.ID)
	}

	if result.Creditcard != nil && result.Creditcard.Name != "My Card" {
		t.Errorf("Expected credit card name 'My Card', got '%s'", result.Creditcard.Name)
	}
}

func TestStoreTransactionToTransactionWithRecurringRelations(t *testing.T) {
	now := time.Now()
	transactionID := uuid.New()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	monthlyID := uuid.New()
	annualID := uuid.New()
	installmentID := uuid.New()

	row := pgstore.GetTransactionByIDRow{
		ID:                                 transactionID,
		UserID:                             userID,
		Name:                               "Recurring Transaction",
		Date:                               createTimestamptzTransaction(now),
		Value:                              createNumeric(320.0),
		Paid:                               false,
		CategoryID:                         createPgUUIDTransaction(categoryID),
		CategoryTransactionType:            pgtype.Int4{Int32: 2, Valid: true},
		CategoryName:                       pgtype.Text{String: "Category", Valid: true},
		CategoryIcon:                       pgtype.Text{String: "icon.png", Valid: true},
		CreditcardID:                       createPgUUIDTransaction(creditCardID),
		CreditcardName:                     pgtype.Text{String: "My Card", Valid: true},
		CreditcardFirstFourNumbers:         pgtype.Text{String: "9876", Valid: true},
		CreditcardCreditLimit:              pgtype.Float8{Float64: 7000.0, Valid: true},
		CreditcardCloseDay:                 pgtype.Int4{Int32: 8, Valid: true},
		CreditcardExpireDay:                pgtype.Int4{Int32: 20, Valid: true},
		CreditcardBackgroundColor:          pgtype.Text{String: "#111111", Valid: true},
		CreditcardTextColor:                pgtype.Text{String: "#EEEEEE", Valid: true},
		MonthlyTransactionsID:              createPgUUIDTransaction(monthlyID),
		MonthlyTransactionsDay:             pgtype.Int4{Int32: 5, Valid: true},
		MonthlyTransactionsName:            pgtype.Text{String: "Monthly", Valid: true},
		MonthlyTransactionsValue:           createNumeric(100.0),
		MonthlyTransactionsCreatedAt:       createTimestamptzTransaction(now),
		MonthlyTransactionsUpdatedAt:       createTimestamptzTransaction(now),
		AnnualTransactionsID:               createPgUUIDTransaction(annualID),
		AnnualTransactionsMonth:            pgtype.Int4{Int32: 12, Valid: true},
		AnnualTransactionsDay:              pgtype.Int4{Int32: 15, Valid: true},
		AnnualTransactionsName:             pgtype.Text{String: "Annual", Valid: true},
		AnnualTransactionsValue:            createNumeric(1200.0),
		AnnualTransactionsCreatedAt:        createTimestamptzTransaction(now),
		AnnualTransactionsUpdatedAt:        createTimestamptzTransaction(now),
		InstallmentTransactionsID:          createPgUUIDTransaction(installmentID),
		InstallmentTransactionsName:        pgtype.Text{String: "Installment", Valid: true},
		InstallmentTransactionsValue:       createNumeric(80.0),
		InstallmentTransactionsInitialDate: createTimestamptzTransaction(now),
		InstallmentTransactionsFinalDate:   createTimestamptzTransaction(now.AddDate(0, 3, 0)),
		InstallmentTransactionsCreatedAt:   createTimestamptzTransaction(now),
		InstallmentTransactionsUpdatedAt:   createTimestamptzTransaction(now),
		CreatedAt:                          createTimestamptzTransaction(now),
		UpdatedAt:                          createTimestamptzTransaction(now),
	}

	result := storeTransactionToTransaction(row)

	if result.Creditcard == nil || result.Creditcard.ID != creditCardID {
		t.Errorf("Expected credit card ID %v", creditCardID)
	}

	if result.MonthlyTransaction == nil || result.MonthlyTransaction.ID != monthlyID {
		t.Errorf("Expected monthly transaction ID %v", monthlyID)
	}

	if result.AnnualTransaction == nil || result.AnnualTransaction.ID != annualID {
		t.Errorf("Expected annual transaction ID %v", annualID)
	}

	if result.InstallmentTransaction == nil || result.InstallmentTransaction.ID != installmentID {
		t.Errorf("Expected installment transaction ID %v", installmentID)
	}
}

// ============= REPOSITORY TEST =============

func TestNewRepository(t *testing.T) {
	mock := newTransactionStoreMock()

	repo := NewRepository(mock)

	if repo.store == nil {
		t.Error("Expected store to be set, got nil")
	}
}
