package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend-commons/models"
	"financialcontrol/internal/store/pgstore"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============= STORE MOCK FOR CREDITCARD =============

type CreditCardStoreMock struct {
	Error             error
	CreateError       error
	UpdateError       error
	DeleteError       error
	CreditCardResult  pgstore.CreditCard
	CreditCardsResult []pgstore.CreditCard
	CountResult       int64
	HasTransResult    bool
	HasTransError     error
}

// CreditCard methods
func (m *CreditCardStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	if m.CreateError != nil {
		return pgstore.CreditCard{}, m.CreateError
	}
	if m.Error != nil {
		return pgstore.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.CreditCardsResult, nil
}

func (m *CreditCardStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	if m.Error != nil {
		return pgstore.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	if m.Error != nil {
		return 0, m.Error
	}
	return m.CountResult, nil
}

func (m *CreditCardStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	if m.UpdateError != nil {
		return pgstore.CreditCard{}, m.UpdateError
	}
	if m.Error != nil {
		return pgstore.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

func (m *CreditCardStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	if m.HasTransError != nil {
		return false, m.HasTransError
	}
	return m.HasTransResult, nil
}

// Other required interface methods (stubs)
func (m *CreditCardStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *CreditCardStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CreditCardStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *CreditCardStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CreditCardStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CreditCardStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *CreditCardStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CreditCardStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CreditCardStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CreditCardStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CreditCardStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CreditCardStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *CreditCardStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *CreditCardStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *CreditCardStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *CreditCardStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *CreditCardStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *CreditCardStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *CreditCardStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CreditCardStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *CreditCardStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CreditCardStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CreditCardStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *CreditCardStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *CreditCardStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *CreditCardStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CreditCardStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CreditCardStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CreditCardStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptzCreditCard(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func newCreditCardStoreMock() *CreditCardStoreMock {
	return &CreditCardStoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateCreditCardSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newCreditCardStoreMock()
	mock.CreditCardResult = pgstore.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		CreditLimit:      5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        createTimestamptzCreditCard(now),
		UpdatedAt:        createTimestamptzCreditCard(now),
	}

	repo := NewRepository(mock)

	data := models.CreateCreditCard{
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	result, err := repo.CreateCreditCard(ctx, data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}

	if result.Name != "Test Card" {
		t.Errorf("Expected name 'Test Card', got '%s'", result.Name)
	}
}

func TestCreateCreditCardError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.CreateError = errors.New("database error")

	repo := NewRepository(mock)

	data := models.CreateCreditCard{
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	_, err := repo.CreateCreditCard(ctx, data)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ TESTS =============

func TestReadCreditCardsSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := newCreditCardStoreMock()
	mock.CreditCardsResult = []pgstore.CreditCard{
		{
			ID:               uuid.New(),
			UserID:           userID,
			Name:             "Card 1",
			FirstFourNumbers: "1111",
			CreditLimit:      3000.0,
			CloseDay:         10,
			ExpireDay:        20,
			BackgroundColor:  "#FF0000",
			TextColor:        "#FFFFFF",
			CreatedAt:        createTimestamptzCreditCard(now),
			UpdatedAt:        createTimestamptzCreditCard(now),
		},
		{
			ID:               uuid.New(),
			UserID:           userID,
			Name:             "Card 2",
			FirstFourNumbers: "2222",
			CreditLimit:      5000.0,
			CloseDay:         15,
			ExpireDay:        25,
			BackgroundColor:  "#00FF00",
			TextColor:        "#000000",
			CreatedAt:        createTimestamptzCreditCard(now),
			UpdatedAt:        createTimestamptzCreditCard(now),
		},
	}

	repo := NewRepository(mock)

	result, err := repo.ReadCreditCards(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 credit cards, got %d", len(result))
	}
}

func TestReadCreditCardsEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.CreditCardsResult = []pgstore.CreditCard{}

	repo := NewRepository(mock)

	result, err := repo.ReadCreditCards(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 credit cards, got %d", len(result))
	}
}

func TestReadCreditCardsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.ReadCreditCards(ctx, userID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadCreditCardByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newCreditCardStoreMock()
	mock.CreditCardResult = pgstore.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		CreditLimit:      5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        createTimestamptzCreditCard(now),
		UpdatedAt:        createTimestamptzCreditCard(now),
	}

	repo := NewRepository(mock)

	result, err := repo.ReadCreditCardByID(ctx, creditCardID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}
}

func TestReadCreditCardByIDError(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.Error = errors.New("not found")

	repo := NewRepository(mock)

	_, err := repo.ReadCreditCardByID(ctx, creditCardID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= COUNT TESTS =============

func TestReadCreditCardCountByUserSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.CountResult = 3

	repo := NewRepository(mock)

	result, err := repo.ReadCountByUser(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != 3 {
		t.Errorf("Expected count 3, got %d", result)
	}
}

func TestReadCreditCardCountByUserError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.ReadCountByUser(ctx, userID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateCreditCardSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := newCreditCardStoreMock()
	mock.CreditCardResult = pgstore.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		CreditLimit:      10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
		CreatedAt:        createTimestamptzCreditCard(now),
		UpdatedAt:        createTimestamptzCreditCard(now),
	}

	repo := NewRepository(mock)

	creditCard := models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
	}

	result, err := repo.UpdateCreditCard(ctx, creditCard)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Name != "Updated Card" {
		t.Errorf("Expected name 'Updated Card', got '%s'", result.Name)
	}
}

func TestUpdateCreditCardError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.UpdateError = errors.New("database error")

	repo := NewRepository(mock)

	creditCard := models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
	}

	_, err := repo.UpdateCreditCard(ctx, creditCard)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= DELETE TESTS =============

func TestDeleteCreditCardSuccess(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()

	repo := NewRepository(mock)

	err := repo.DeleteCreditCard(ctx, creditCardID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteCreditCardError(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.DeleteError = errors.New("database error")

	repo := NewRepository(mock)

	err := repo.DeleteCreditCard(ctx, creditCardID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= HAS TRANSACTIONS TESTS =============

func TestHasTransactionsByCreditCardSuccess(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.HasTransResult = true

	repo := NewRepository(mock)

	result, err := repo.HasTransactionsByCreditCard(ctx, creditCardID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestHasTransactionsByCreditCardFalse(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.HasTransResult = false

	repo := NewRepository(mock)

	result, err := repo.HasTransactionsByCreditCard(ctx, creditCardID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}
}

func TestHasTransactionsByCreditCardError(t *testing.T) {
	ctx := context.Background()
	creditCardID := uuid.New()

	mock := newCreditCardStoreMock()
	mock.HasTransError = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.HasTransactionsByCreditCard(ctx, creditCardID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= STORE CREDITCARD TO CREDITCARD TESTS =============

func TestStoreCreditCardToCreditCard(t *testing.T) {
	creditCardID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	storeCreditCard := pgstore.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		CreditLimit:      5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        createTimestamptzCreditCard(now),
		UpdatedAt:        createTimestamptzCreditCard(now),
	}

	result := storeCreditcardToCreditcard(storeCreditCard)

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}

	if result.Name != "Test Card" {
		t.Errorf("Expected name 'Test Card', got '%s'", result.Name)
	}

	if result.Limit != 5000.0 {
		t.Errorf("Expected limit 5000.0, got %f", result.Limit)
	}
}
