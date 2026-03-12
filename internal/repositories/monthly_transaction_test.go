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

// ============= STORE MOCK =============

type StoreMock struct {
	Error                       error
	MonthlyTransactionResult    pgstore.MonthlyTransaction
	MonthlyTransactionRowResult pgstore.GetMonthlyTransactionByIDRow
	ShortMonthlyRowResult       pgstore.MonthlyTransaction
	MonthlyTransactionsResult   []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow
}

// Monthly Transaction methods
func (m *StoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	if m.Error != nil {
		return pgstore.MonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionResult, nil
}

func (m *StoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.MonthlyTransactionsResult, nil
}

func (m *StoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	if m.Error != nil {
		return pgstore.GetMonthlyTransactionByIDRow{}, m.Error
	}
	return m.MonthlyTransactionRowResult, nil
}

func (m *StoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	if m.Error != nil {
		return pgstore.MonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionResult, nil
}

func (m *StoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return m.Error
}

// Other required interface methods (stubs)
func (m *StoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *StoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *StoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *StoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *StoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *StoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *StoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *StoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *StoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *StoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *StoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *StoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *StoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *StoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *StoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *StoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *StoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *StoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *StoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *StoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *StoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *StoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *StoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *StoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *StoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *StoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *StoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *StoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *StoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *StoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *StoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *StoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *StoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *StoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *StoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *StoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	if m.Error != nil {
		return pgstore.MonthlyTransaction{}, m.Error
	}
	return m.ShortMonthlyRowResult, nil
}
func (m *StoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *StoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *StoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func createPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func createPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func createPgInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func createPgFloat8(f float64) pgtype.Float8 {
	return pgtype.Float8{Float64: f, Valid: true}
}

func newStoreMock() *StoreMock {
	return &StoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateMonthlyTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionResult = pgstore.MonthlyTransaction{
		ID:         transactionID,
		UserID:     userID,
		Name:       "Test Monthly",
		Value:      utils.Float64ToNumeric(150.50),
		Day:        15,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptz(now),
		UpdatedAt:  createTimestamptz(now),
	}

	repo := Repository{store: mock}

	request := models.CreateMonthlyTransaction{
		UserID:       userID,
		Name:         "Test Monthly",
		Value:        150.50,
		Day:          15,
		CategoryID:   categoryID,
		CreditCardID: nil,
	}

	result, err := repo.CreateMonthlyTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Day != 15 {
		t.Errorf("Expected Day 15, got %v", result.Day)
	}

	if result.Value != 150.50 {
		t.Errorf("Expected Value 150.50, got %v", result.Value)
	}
}

func TestCreateMonthlyTransactionWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionResult = pgstore.MonthlyTransaction{
		ID:           transactionID,
		UserID:       userID,
		Name:         "Test Monthly with CC",
		Value:        utils.Float64ToNumeric(200.00),
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: createPgUUID(creditCardID),
		CreatedAt:    createTimestamptz(now),
		UpdatedAt:    createTimestamptz(now),
	}

	repo := Repository{store: mock}

	request := models.CreateMonthlyTransaction{
		UserID:       userID,
		Name:         "Test Monthly with CC",
		Value:        200.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, err := repo.CreateMonthlyTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Value != 200.00 {
		t.Errorf("Expected Value 200.00, got %v", result.Value)
	}
}

func TestCreateMonthlyTransactionError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := newStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	request := models.CreateMonthlyTransaction{
		UserID:       userID,
		Name:         "Test Monthly",
		Value:        150.50,
		Day:          15,
		CategoryID:   categoryID,
		CreditCardID: nil,
	}

	result, err := repo.CreateMonthlyTransaction(ctx, request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Errorf("Expected 'database error', got %v", err.Error())
	}

	if result.ID != uuid.Nil {
		t.Errorf("Expected empty ID, got %v", result.ID)
	}
}

// ============= READ PAGINATED TESTS =============

func TestReadMonthlyTransactionsPaginatedSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionsResult = []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Monthly 1",
			Value:                   utils.Float64ToNumeric(100.00),
			Day:                     5,
			CreatedAt:               createTimestamptz(now),
			UpdatedAt:               createTimestamptz(now),
			CategoryID:              createPgUUID(categoryID),
			CategoryTransactionType: createPgInt4(1),
			CategoryName:            createPgText("Category 1"),
			CategoryIcon:            createPgText("icon1"),
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %v", len(result))
	}

	if result[0].ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result[0].ID)
	}

	if result[0].Category.Name != "Category 1" {
		t.Errorf("Expected Category Name 'Category 1', got %v", result[0].Category.Name)
	}
}

func TestReadMonthlyTransactionsPaginatedWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionsResult = []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow{
		{
			ID:                         transactionID,
			UserID:                     userID,
			Name:                       "Monthly with CC",
			Value:                      utils.Float64ToNumeric(250.00),
			Day:                        20,
			CreatedAt:                  createTimestamptz(now),
			UpdatedAt:                  createTimestamptz(now),
			CategoryID:                 createPgUUID(categoryID),
			CategoryTransactionType:    createPgInt4(1),
			CategoryName:               createPgText("Category 1"),
			CategoryIcon:               createPgText("icon1"),
			CreditcardID:               createPgUUID(creditCardID),
			CreditcardName:             createPgText("My Card"),
			CreditcardFirstFourNumbers: createPgText("1234"),
			CreditcardCreditLimit:      createPgFloat8(5000.00),
			CreditcardCloseDay:         createPgInt4(25),
			CreditcardExpireDay:        createPgInt4(5),
			CreditcardBackgroundColor:  createPgText("#000000"),
			CreditcardTextColor:        createPgText("#FFFFFF"),
			TotalCount:                 1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	if result[0].Creditcard == nil {
		t.Error("Expected Creditcard to be present")
		return
	}

	if result[0].Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result[0].Creditcard.ID)
	}

	if result[0].Creditcard.Name != "My Card" {
		t.Errorf("Expected Creditcard Name 'My Card', got %v", result[0].Creditcard.Name)
	}

	if result[0].Creditcard.Limit != 5000.00 {
		t.Errorf("Expected Creditcard Limit 5000.00, got %v", result[0].Creditcard.Limit)
	}
}

func TestReadMonthlyTransactionsPaginatedMultipleResults(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionsResult = []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Monthly 1",
			Value:                   utils.Float64ToNumeric(100.00),
			Day:                     5,
			CreatedAt:               createTimestamptz(now),
			UpdatedAt:               createTimestamptz(now),
			CategoryID:              createPgUUID(categoryID),
			CategoryTransactionType: createPgInt4(1),
			CategoryName:            createPgText("Category 1"),
			CategoryIcon:            createPgText("icon1"),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Monthly 2",
			Value:                   utils.Float64ToNumeric(200.00),
			Day:                     10,
			CreatedAt:               createTimestamptz(now),
			UpdatedAt:               createTimestamptz(now),
			CategoryID:              createPgUUID(categoryID),
			CategoryTransactionType: createPgInt4(1),
			CategoryName:            createPgText("Category 1"),
			CategoryIcon:            createPgText("icon1"),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Monthly 3",
			Value:                   utils.Float64ToNumeric(300.00),
			Day:                     15,
			CreatedAt:               createTimestamptz(now),
			UpdatedAt:               createTimestamptz(now),
			CategoryID:              createPgUUID(categoryID),
			CategoryTransactionType: createPgInt4(1),
			CategoryName:            createPgText("Category 1"),
			CategoryIcon:            createPgText("icon1"),
			TotalCount:              3,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 3 {
		t.Errorf("Expected count 3, got %v", count)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %v", len(result))
	}
}

func TestReadMonthlyTransactionsPaginatedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %v", count)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestReadMonthlyTransactionsPaginatedEmpty(t *testing.T) {
	ctx := context.Background()

	mock := newStoreMock()
	mock.MonthlyTransactionsResult = []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow{}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: uuid.New(),
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

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

func TestReadMonthlyTransactionsPaginatedWithNoCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionsResult = []pgstore.ListMonthlyTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Monthly No CC",
			Value:                   utils.Float64ToNumeric(50.00),
			Day:                     15,
			CreatedAt:               createTimestamptz(now),
			UpdatedAt:               createTimestamptz(now),
			CategoryID:              createPgUUID(categoryID),
			CategoryTransactionType: createPgInt4(1),
			CategoryName:            createPgText("Category"),
			CategoryIcon:            createPgText("icon"),
			CreditcardID:            pgtype.UUID{Valid: false},
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	if result[0].Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadMonthlyTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionRowResult = pgstore.GetMonthlyTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Monthly Test",
		Value:                   utils.Float64ToNumeric(175.00),
		Day:                     12,
		CreatedAt:               createTimestamptz(now),
		UpdatedAt:               createTimestamptz(now),
		CategoryID:              createPgUUID(categoryID),
		CategoryTransactionType: createPgInt4(1),
		CategoryName:            createPgText("Test Category"),
		CategoryIcon:            createPgText("test-icon"),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadMonthlyTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, result.UserID)
	}

	if result.Name != "Monthly Test" {
		t.Errorf("Expected Name 'Monthly Test', got %v", result.Name)
	}

	if result.Value != 175.00 {
		t.Errorf("Expected Value 175.00, got %v", result.Value)
	}

	if result.Day != 12 {
		t.Errorf("Expected Day 12, got %v", result.Day)
	}

	if result.Category.ID != categoryID {
		t.Errorf("Expected Category ID %v, got %v", categoryID, result.Category.ID)
	}
}

func TestReadMonthlyTransactionByIDWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionRowResult = pgstore.GetMonthlyTransactionByIDRow{
		ID:                         transactionID,
		UserID:                     userID,
		Name:                       "Monthly with CC",
		Value:                      utils.Float64ToNumeric(500.00),
		Day:                        1,
		CreatedAt:                  createTimestamptz(now),
		UpdatedAt:                  createTimestamptz(now),
		CategoryID:                 createPgUUID(categoryID),
		CategoryTransactionType:    createPgInt4(2),
		CategoryName:               createPgText("Bills"),
		CategoryIcon:               createPgText("bill-icon"),
		CreditcardID:               createPgUUID(creditCardID),
		CreditcardName:             createPgText("Platinum Card"),
		CreditcardFirstFourNumbers: createPgText("5678"),
		CreditcardCreditLimit:      createPgFloat8(10000.00),
		CreditcardCloseDay:         createPgInt4(20),
		CreditcardExpireDay:        createPgInt4(10),
		CreditcardBackgroundColor:  createPgText("#1A1A1A"),
		CreditcardTextColor:        createPgText("#GOLD"),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadMonthlyTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
		return
	}

	if result.Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result.Creditcard.ID)
	}

	if result.Creditcard.Name != "Platinum Card" {
		t.Errorf("Expected Creditcard Name 'Platinum Card', got %v", result.Creditcard.Name)
	}

	if result.Creditcard.FirstFourNumbers != "5678" {
		t.Errorf("Expected FirstFourNumbers '5678', got %v", result.Creditcard.FirstFourNumbers)
	}

	if result.Creditcard.Limit != 10000.00 {
		t.Errorf("Expected Limit 10000.00, got %v", result.Creditcard.Limit)
	}

	if result.Creditcard.CloseDay != 20 {
		t.Errorf("Expected CloseDay 20, got %v", result.Creditcard.CloseDay)
	}

	if result.Creditcard.ExpireDay != 10 {
		t.Errorf("Expected ExpireDay 10, got %v", result.Creditcard.ExpireDay)
	}
}

func TestReadMonthlyTransactionByIDError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newStoreMock()
	mock.Error = errors.New("not found")

	repo := Repository{store: mock}

	result, err := repo.ReadMonthlyTransactionByID(ctx, transactionID)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "not found" {
		t.Errorf("Expected 'not found', got %v", err.Error())
	}

	if result.ID != uuid.Nil {
		t.Errorf("Expected empty ID, got %v", result.ID)
	}
}

func TestReadMonthlyTransactionByIDWithoutCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionRowResult = pgstore.GetMonthlyTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Monthly without CC",
		Value:                   utils.Float64ToNumeric(99.99),
		Day:                     28,
		CreatedAt:               createTimestamptz(now),
		UpdatedAt:               createTimestamptz(now),
		CategoryID:              createPgUUID(categoryID),
		CategoryTransactionType: createPgInt4(1),
		CategoryName:            createPgText("No CC Category"),
		CategoryIcon:            createPgText("no-cc-icon"),
		CreditcardID:            pgtype.UUID{Valid: false},
	}

	repo := Repository{store: mock}

	result, err := repo.ReadMonthlyTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}

	if result.Name != "Monthly without CC" {
		t.Errorf("Expected Name 'Monthly without CC', got %v", result.Name)
	}
}

func TestReadShortMonthlyTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.ShortMonthlyRowResult = pgstore.MonthlyTransaction{
		ID:        transactionID,
		Name:      "Short Monthly",
		Value:     utils.Float64ToNumeric(88.50),
		Day:       17,
		CreatedAt: createTimestamptz(now),
		UpdatedAt: createTimestamptz(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadShortMonthlyTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Day != 17 {
		t.Errorf("Expected Day 17, got %v", result.Day)
	}
}

func TestReadShortMonthlyTransactionByIDError(t *testing.T) {
	ctx := context.Background()

	mock := newStoreMock()
	mock.Error = errors.New("not found")

	repo := Repository{store: mock}

	_, err := repo.ReadShortMonthlyTransactionByID(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateMonthlyTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionResult = pgstore.MonthlyTransaction{
		ID:           transactionID,
		UserID:       userID,
		Name:         "Updated Monthly",
		Value:        utils.Float64ToNumeric(300.00),
		Day:          25,
		CategoryID:   categoryID,
		CreditCardID: createPgUUID(creditCardID),
		CreatedAt:    createTimestamptz(now),
		UpdatedAt:    createTimestamptz(now),
	}

	repo := Repository{store: mock}

	model := models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Updated Monthly",
		Value:  300.00,
		Day:    25,
		Category: models.Category{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Updated Category",
			Icon:            "updated-icon",
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateMonthlyTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Value != 300.00 {
		t.Errorf("Expected Value 300.00, got %v", result.Value)
	}

	if result.Day != 25 {
		t.Errorf("Expected Day 25, got %v", result.Day)
	}
}

func TestUpdateMonthlyTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := newStoreMock()
	mock.Error = errors.New("update failed")

	repo := Repository{store: mock}

	model := models.MonthlyTransaction{
		ID:    transactionID,
		Name:  "Updated Monthly",
		Value: 300.00,
		Day:   25,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateMonthlyTransaction(ctx, model)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "update failed" {
		t.Errorf("Expected 'update failed', got %v", err.Error())
	}

	if result.ID != uuid.Nil {
		t.Errorf("Expected empty ID, got %v", result.ID)
	}
}

func TestUpdateMonthlyTransactionWithDifferentDay(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionResult = pgstore.MonthlyTransaction{
		ID:           transactionID,
		Name:         "End of Month",
		Value:        utils.Float64ToNumeric(1000.00),
		Day:          31,
		CategoryID:   categoryID,
		CreditCardID: createPgUUID(creditCardID),
		CreatedAt:    createTimestamptz(now),
		UpdatedAt:    createTimestamptz(now),
	}

	repo := Repository{store: mock}

	model := models.MonthlyTransaction{
		ID:    transactionID,
		Name:  "End of Month",
		Value: 1000.00,
		Day:   31,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateMonthlyTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Day != 31 {
		t.Errorf("Expected Day 31, got %v", result.Day)
	}
}

// ============= DELETE TESTS =============

func TestDeleteMonthlyTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newStoreMock()
	mock.Error = nil

	repo := Repository{store: mock}

	err := repo.DeleteMonthlyTransaction(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteMonthlyTransactionError(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()

	mock := newStoreMock()
	mock.Error = errors.New("delete failed")

	repo := Repository{store: mock}

	err := repo.DeleteMonthlyTransaction(ctx, transactionID)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "delete failed" {
		t.Errorf("Expected 'delete failed', got %v", err.Error())
	}
}

// ============= EDGE CASES =============

func TestCreateMonthlyTransactionWithZeroValue(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := newStoreMock()
	mock.MonthlyTransactionResult = pgstore.MonthlyTransaction{
		ID:         transactionID,
		UserID:     userID,
		Name:       "Zero Value Monthly",
		Value:      utils.Float64ToNumeric(0.00),
		Day:        1,
		CategoryID: categoryID,
		CreatedAt:  createTimestamptz(now),
		UpdatedAt:  createTimestamptz(now),
	}

	repo := Repository{store: mock}

	request := models.CreateMonthlyTransaction{
		UserID:       userID,
		Name:         "Zero Value Monthly",
		Value:        0.00,
		Day:          1,
		CategoryID:   categoryID,
		CreditCardID: nil,
	}

	result, err := repo.CreateMonthlyTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 0.00 {
		t.Errorf("Expected Value 0.00, got %v", result.Value)
	}
}
