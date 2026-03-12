package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"backend-commons/models"
	"financialcontrol/internal/store/pgstore"
	"github.com/TB-Systems/go-commons/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============= STORE MOCK FOR ANNUAL =============

type AnnualStoreMock struct {
	Error                      error
	AnnualTransactionResult    pgstore.AnnualTransaction
	AnnualTransactionRowResult pgstore.GetAnnualTransactionByIDRow
	ShortAnnualRowResult       pgstore.AnnualTransaction
	AnnualTransactionsResult   []pgstore.ListAnnualTransactionsByUserIDPaginatedRow
}

// Annual Transaction methods
func (m *AnnualStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	if m.Error != nil {
		return pgstore.AnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionResult, nil
}

func (m *AnnualStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.AnnualTransactionsResult, nil
}

func (m *AnnualStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	if m.Error != nil {
		return pgstore.GetAnnualTransactionByIDRow{}, m.Error
	}
	return m.AnnualTransactionRowResult, nil
}

func (m *AnnualStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	if m.Error != nil {
		return pgstore.AnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionResult, nil
}

func (m *AnnualStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return m.Error
}

// Other required interface methods (stubs)
func (m *AnnualStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *AnnualStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *AnnualStoreMock) CountMonthlyTransactionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *AnnualStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *AnnualStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *AnnualStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *AnnualStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *AnnualStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *AnnualStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *AnnualStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *AnnualStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *AnnualStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *AnnualStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *AnnualStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *AnnualStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *AnnualStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *AnnualStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *AnnualStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *AnnualStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *AnnualStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *AnnualStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *AnnualStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *AnnualStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *AnnualStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *AnnualStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *AnnualStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *AnnualStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *AnnualStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *AnnualStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *AnnualStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *AnnualStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *AnnualStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	if m.Error != nil {
		return pgstore.AnnualTransaction{}, m.Error
	}
	return m.ShortAnnualRowResult, nil
}
func (m *AnnualStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *AnnualStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptzAnnual(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func createPgUUIDAnnual(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func createPgTextAnnual(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func createPgInt4Annual(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func createPgFloat8Annual(f float64) pgtype.Float8 {
	return pgtype.Float8{Float64: f, Valid: true}
}

func newAnnualStoreMock() *AnnualStoreMock {
	return &AnnualStoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateAnnualTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionResult = pgstore.AnnualTransaction{
		ID:         transactionID,
		UserID:     userID,
		Name:       "Test Annual",
		Value:      utils.Float64ToNumeric(100.00),
		Day:        15,
		Month:      6,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptzAnnual(now),
		UpdatedAt:  createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	request := models.CreateAnnualTransaction{
		UserID:     userID,
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
		CategoryID: categoryID,
	}

	result, err := repo.CreateAnnualTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Day != 15 {
		t.Errorf("Expected Day 15, got %v", result.Day)
	}

	if result.Month != 6 {
		t.Errorf("Expected Month 6, got %v", result.Month)
	}
}

func TestCreateAnnualTransactionWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionResult = pgstore.AnnualTransaction{
		ID:           transactionID,
		UserID:       userID,
		Name:         "Annual with CC",
		Value:        utils.Float64ToNumeric(500.00),
		Day:          10,
		Month:        12,
		CategoryID:   categoryID,
		CreditCardID: createPgUUIDAnnual(creditCardID),
		CreatedAt:    createTimestamptzAnnual(now),
		UpdatedAt:    createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	request := models.CreateAnnualTransaction{
		UserID:       userID,
		Name:         "Annual with CC",
		Value:        500.00,
		Day:          10,
		Month:        12,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, err := repo.CreateAnnualTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Month != 12 {
		t.Errorf("Expected Month 12, got %v", result.Month)
	}
}

func TestCreateAnnualTransactionError(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	request := models.CreateAnnualTransaction{
		UserID:     uuid.New(),
		Name:       "Test",
		Value:      100.00,
		Day:        1,
		Month:      1,
		CategoryID: uuid.New(),
	}

	_, err := repo.CreateAnnualTransaction(ctx, request)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ PAGINATED TESTS =============

func TestReadAnnualTransactionsPaginatedSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionsResult = []pgstore.ListAnnualTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Annual 1",
			Value:                   utils.Float64ToNumeric(200.00),
			Day:                     5,
			Month:                   3,
			CategoryID:              createPgUUIDAnnual(categoryID),
			CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
			CategoryName:            createPgTextAnnual("Category 1"),
			CategoryIcon:            createPgTextAnnual("icon1"),
			CreatedAt:               createTimestamptzAnnual(now),
			UpdatedAt:               createTimestamptzAnnual(now),
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %v", len(result))
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	if result[0].Month != 3 {
		t.Errorf("Expected Month 3, got %v", result[0].Month)
	}
}

func TestReadAnnualTransactionsPaginatedWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionsResult = []pgstore.ListAnnualTransactionsByUserIDPaginatedRow{
		{
			ID:                         transactionID,
			UserID:                     userID,
			Name:                       "Annual with CC",
			Value:                      utils.Float64ToNumeric(300.00),
			Day:                        15,
			Month:                      7,
			CategoryID:                 createPgUUIDAnnual(categoryID),
			CategoryTransactionType:    createPgInt4Annual(int32(models.Credit)),
			CategoryName:               createPgTextAnnual("Credit Category"),
			CategoryIcon:               createPgTextAnnual("credit-icon"),
			CreditcardID:               createPgUUIDAnnual(creditCardID),
			CreditcardName:             createPgTextAnnual("My Card"),
			CreditcardFirstFourNumbers: createPgTextAnnual("1234"),
			CreditcardCreditLimit:      createPgFloat8Annual(5000.00),
			CreditcardCloseDay:         createPgInt4Annual(10),
			CreditcardExpireDay:        createPgInt4Annual(20),
			CreditcardBackgroundColor:  createPgTextAnnual("#000000"),
			CreditcardTextColor:        createPgTextAnnual("#FFFFFF"),
			CreatedAt:                  createTimestamptzAnnual(now),
			UpdatedAt:                  createTimestamptzAnnual(now),
			TotalCount:                 1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, _, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result[0].Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result[0].Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result[0].Creditcard.ID)
	}
}

func TestReadAnnualTransactionsPaginatedMultipleResults(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionsResult = []pgstore.ListAnnualTransactionsByUserIDPaginatedRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Annual 1",
			Value:                   utils.Float64ToNumeric(100.00),
			Day:                     1,
			Month:                   1,
			CategoryID:              createPgUUIDAnnual(categoryID),
			CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
			CategoryName:            createPgTextAnnual("Category"),
			CategoryIcon:            createPgTextAnnual("icon"),
			CreatedAt:               createTimestamptzAnnual(now),
			UpdatedAt:               createTimestamptzAnnual(now),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Annual 2",
			Value:                   utils.Float64ToNumeric(200.00),
			Day:                     15,
			Month:                   6,
			CategoryID:              createPgUUIDAnnual(categoryID),
			CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
			CategoryName:            createPgTextAnnual("Category"),
			CategoryIcon:            createPgTextAnnual("icon"),
			CreatedAt:               createTimestamptzAnnual(now),
			UpdatedAt:               createTimestamptzAnnual(now),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Annual 3",
			Value:                   utils.Float64ToNumeric(300.00),
			Day:                     25,
			Month:                   12,
			CategoryID:              createPgUUIDAnnual(categoryID),
			CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
			CategoryName:            createPgTextAnnual("Category"),
			CategoryIcon:            createPgTextAnnual("icon"),
			CreatedAt:               createTimestamptzAnnual(now),
			UpdatedAt:               createTimestamptzAnnual(now),
			TotalCount:              3,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %v", len(result))
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %v", count)
	}
}

func TestReadAnnualTransactionsPaginatedError(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: uuid.New(),
		Limit:  10,
		Offset: 0,
	}

	_, _, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReadAnnualTransactionsPaginatedEmpty(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionsResult = []pgstore.ListAnnualTransactionsByUserIDPaginatedRow{}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: uuid.New(),
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestReadAnnualTransactionsPaginatedWithNoCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionsResult = []pgstore.ListAnnualTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Annual without CC",
			Value:                   utils.Float64ToNumeric(150.00),
			Day:                     20,
			Month:                   9,
			CategoryID:              createPgUUIDAnnual(categoryID),
			CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
			CategoryName:            createPgTextAnnual("Category"),
			CategoryIcon:            createPgTextAnnual("icon"),
			CreditcardID:            pgtype.UUID{Valid: false},
			CreatedAt:               createTimestamptzAnnual(now),
			UpdatedAt:               createTimestamptzAnnual(now),
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, _, err := repo.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result[0].Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadAnnualTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionRowResult = pgstore.GetAnnualTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Annual Test",
		Value:                   utils.Float64ToNumeric(250.00),
		Day:                     10,
		Month:                   5,
		CategoryID:              createPgUUIDAnnual(categoryID),
		CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
		CategoryName:            createPgTextAnnual("Test Category"),
		CategoryIcon:            createPgTextAnnual("test-icon"),
		CreatedAt:               createTimestamptzAnnual(now),
		UpdatedAt:               createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadAnnualTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, result.UserID)
	}

	if result.Month != 5 {
		t.Errorf("Expected Month 5, got %v", result.Month)
	}
}

func TestReadAnnualTransactionByIDWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionRowResult = pgstore.GetAnnualTransactionByIDRow{
		ID:                         transactionID,
		UserID:                     userID,
		Name:                       "Annual with CC",
		Value:                      utils.Float64ToNumeric(400.00),
		Day:                        5,
		Month:                      11,
		CategoryID:                 createPgUUIDAnnual(categoryID),
		CategoryTransactionType:    createPgInt4Annual(int32(models.Credit)),
		CategoryName:               createPgTextAnnual("Credit Category"),
		CategoryIcon:               createPgTextAnnual("credit-icon"),
		CreditcardID:               createPgUUIDAnnual(creditCardID),
		CreditcardName:             createPgTextAnnual("My Card"),
		CreditcardFirstFourNumbers: createPgTextAnnual("5678"),
		CreditcardCreditLimit:      createPgFloat8Annual(10000.00),
		CreditcardCloseDay:         createPgInt4Annual(15),
		CreditcardExpireDay:        createPgInt4Annual(25),
		CreditcardBackgroundColor:  createPgTextAnnual("#FF0000"),
		CreditcardTextColor:        createPgTextAnnual("#FFFFFF"),
		CreatedAt:                  createTimestamptzAnnual(now),
		UpdatedAt:                  createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadAnnualTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result.Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result.Creditcard.ID)
	}

	if result.Month != 11 {
		t.Errorf("Expected Month 11, got %v", result.Month)
	}
}

func TestReadAnnualTransactionByIDError(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("not found")

	repo := Repository{store: mock}

	_, err := repo.ReadAnnualTransactionByID(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReadAnnualTransactionByIDWithoutCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionRowResult = pgstore.GetAnnualTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Annual without CC",
		Value:                   utils.Float64ToNumeric(180.00),
		Day:                     25,
		Month:                   8,
		CategoryID:              createPgUUIDAnnual(categoryID),
		CategoryTransactionType: createPgInt4Annual(int32(models.Debit)),
		CategoryName:            createPgTextAnnual("Category"),
		CategoryIcon:            createPgTextAnnual("icon"),
		CreditcardID:            pgtype.UUID{Valid: false},
		CreatedAt:               createTimestamptzAnnual(now),
		UpdatedAt:               createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadAnnualTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}
}

func TestReadShortAnnualTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.ShortAnnualRowResult = pgstore.AnnualTransaction{
		ID:        transactionID,
		Name:      "Short Annual",
		Value:     utils.Float64ToNumeric(320.00),
		Day:       9,
		Month:     12,
		CreatedAt: createTimestamptzAnnual(now),
		UpdatedAt: createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadShortAnnualTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Month != 12 {
		t.Errorf("Expected Month 12, got %v", result.Month)
	}
}

func TestReadShortAnnualTransactionByIDError(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("not found")

	repo := Repository{store: mock}

	_, err := repo.ReadShortAnnualTransactionByID(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateAnnualTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionResult = pgstore.AnnualTransaction{
		ID:         transactionID,
		Name:       "Updated Annual",
		Value:      utils.Float64ToNumeric(350.00),
		Day:        20,
		Month:      10,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptzAnnual(now),
		UpdatedAt:  createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	model := models.AnnualTransaction{
		ID:    transactionID,
		Name:  "Updated Annual",
		Value: 350.00,
		Day:   20,
		Month: 10,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateAnnualTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Day != 20 {
		t.Errorf("Expected Day 20, got %v", result.Day)
	}

	if result.Month != 10 {
		t.Errorf("Expected Month 10, got %v", result.Month)
	}
}

func TestUpdateAnnualTransactionError(t *testing.T) {
	ctx := context.Background()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("update error")

	repo := Repository{store: mock}

	model := models.AnnualTransaction{
		ID:    uuid.New(),
		Name:  "Test",
		Value: 100.00,
		Day:   1,
		Month: 1,
		Category: models.Category{
			ID: uuid.New(),
		},
		Creditcard: &models.CreditCard{
			ID: uuid.New(),
		},
	}

	_, err := repo.UpdateAnnualTransaction(ctx, model)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUpdateAnnualTransactionWithDifferentMonth(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionResult = pgstore.AnnualTransaction{
		ID:         transactionID,
		Name:       "Annual Month Change",
		Value:      utils.Float64ToNumeric(500.00),
		Day:        15,
		Month:      4,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptzAnnual(now),
		UpdatedAt:  createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	model := models.AnnualTransaction{
		ID:    transactionID,
		Name:  "Annual Month Change",
		Value: 500.00,
		Day:   15,
		Month: 4,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateAnnualTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Month != 4 {
		t.Errorf("Expected Month 4, got %v", result.Month)
	}
}

// ============= DELETE TESTS =============

func TestDeleteAnnualTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newAnnualStoreMock()

	repo := Repository{store: mock}

	err := repo.DeleteAnnualTransaction(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteAnnualTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newAnnualStoreMock()
	mock.Error = errors.New("delete error")

	repo := Repository{store: mock}

	err := repo.DeleteAnnualTransaction(ctx, transactionID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= EDGE CASES =============

func TestCreateAnnualTransactionWithZeroValue(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newAnnualStoreMock()
	mock.AnnualTransactionResult = pgstore.AnnualTransaction{
		ID:         transactionID,
		UserID:     userID,
		Name:       "Zero Value Annual",
		Value:      utils.Float64ToNumeric(0.00),
		Day:        1,
		Month:      1,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptzAnnual(now),
		UpdatedAt:  createTimestamptzAnnual(now),
	}

	repo := Repository{store: mock}

	request := models.CreateAnnualTransaction{
		UserID:     userID,
		Name:       "Zero Value Annual",
		Value:      0.00,
		Day:        1,
		Month:      1,
		CategoryID: categoryID,
	}

	result, err := repo.CreateAnnualTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 0.00 {
		t.Errorf("Expected Value 0.00, got %v", result.Value)
	}
}
