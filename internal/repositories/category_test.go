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

// ============= STORE MOCK FOR CATEGORY =============

type CategoryStoreMock struct {
	Error            error
	CreateError      error
	UpdateError      error
	DeleteError      error
	CategoryResult   pgstore.Category
	CategoriesResult []pgstore.Category
	CountResult      int64
	HasTransResult   bool
	HasTransError    error
}

// Category methods
func (m *CategoryStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	if m.CreateError != nil {
		return pgstore.Category{}, m.CreateError
	}
	if m.Error != nil {
		return pgstore.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.CategoriesResult, nil
}

func (m *CategoryStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	if m.Error != nil {
		return pgstore.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	if m.Error != nil {
		return 0, m.Error
	}
	return m.CountResult, nil
}

func (m *CategoryStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	if m.UpdateError != nil {
		return pgstore.Category{}, m.UpdateError
	}
	if m.Error != nil {
		return pgstore.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

func (m *CategoryStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	if m.HasTransError != nil {
		return false, m.HasTransError
	}
	return m.HasTransResult, nil
}

// Other required interface methods (stubs)
func (m *CategoryStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *CategoryStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CategoryStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *CategoryStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CategoryStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CategoryStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *CategoryStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CategoryStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CategoryStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CategoryStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CategoryStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *CategoryStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *CategoryStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *CategoryStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	return pgstore.GetInstallmentTransactionByIDRow{}, nil
}
func (m *CategoryStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *CategoryStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *CategoryStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *CategoryStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *CategoryStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CategoryStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *CategoryStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CategoryStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CategoryStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *CategoryStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *CategoryStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *CategoryStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *CategoryStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *CategoryStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	return pgstore.InstallmentTransaction{}, nil
}
func (m *CategoryStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptzCategory(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func newCategoryStoreMock() *CategoryStoreMock {
	return &CategoryStoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateCategorySuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newCategoryStoreMock()
	mock.CategoryResult = pgstore.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: 0,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       createTimestamptzCategory(now),
		UpdatedAt:       createTimestamptzCategory(now),
	}

	repo := NewRepository(mock)

	data := models.CreateCategory{
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
	}

	result, err := repo.CreateCategory(ctx, data)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}

	if result.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", result.Name)
	}
}

func TestCreateCategoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCategoryStoreMock()
	mock.CreateError = errors.New("database error")

	repo := NewRepository(mock)

	data := models.CreateCategory{
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
	}

	_, err := repo.CreateCategory(ctx, data)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ TESTS =============

func TestReadCategoriesSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := newCategoryStoreMock()
	mock.CategoriesResult = []pgstore.Category{
		{
			ID:              uuid.New(),
			UserID:          userID,
			TransactionType: 0,
			Name:            "Category 1",
			Icon:            "icon1.png",
			CreatedAt:       createTimestamptzCategory(now),
			UpdatedAt:       createTimestamptzCategory(now),
		},
		{
			ID:              uuid.New(),
			UserID:          userID,
			TransactionType: 1,
			Name:            "Category 2",
			Icon:            "icon2.png",
			CreatedAt:       createTimestamptzCategory(now),
			UpdatedAt:       createTimestamptzCategory(now),
		},
	}

	repo := NewRepository(mock)

	result, err := repo.ReadCategories(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(result))
	}
}

func TestReadCategoriesEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCategoryStoreMock()
	mock.CategoriesResult = []pgstore.Category{}

	repo := NewRepository(mock)

	result, err := repo.ReadCategories(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(result))
	}
}

func TestReadCategoriesError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCategoryStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.ReadCategories(ctx, userID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadCategoryByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newCategoryStoreMock()
	mock.CategoryResult = pgstore.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: 0,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       createTimestamptzCategory(now),
		UpdatedAt:       createTimestamptzCategory(now),
	}

	repo := NewRepository(mock)

	result, err := repo.ReadCategoryByID(ctx, categoryID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}
}

func TestReadCategoryByIDError(t *testing.T) {
	ctx := context.Background()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()
	mock.Error = errors.New("not found")

	repo := NewRepository(mock)

	_, err := repo.ReadCategoryByID(ctx, categoryID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= COUNT TESTS =============

func TestGetCategoryCountByUserSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCategoryStoreMock()
	mock.CountResult = 5

	repo := NewRepository(mock)

	result, err := repo.GetCategoryCountByUser(ctx, userID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != 5 {
		t.Errorf("Expected count 5, got %d", result)
	}
}

func TestGetCategoryCountByUserError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := newCategoryStoreMock()
	mock.Error = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.GetCategoryCountByUser(ctx, userID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateCategorySuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := newCategoryStoreMock()
	mock.CategoryResult = pgstore.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: 1,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
		CreatedAt:       createTimestamptzCategory(now),
		UpdatedAt:       createTimestamptzCategory(now),
	}

	repo := NewRepository(mock)

	category := models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}

	result, err := repo.UpdateCategory(ctx, category)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Name != "Updated Category" {
		t.Errorf("Expected name 'Updated Category', got '%s'", result.Name)
	}
}

func TestUpdateCategoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()
	mock.UpdateError = errors.New("database error")

	repo := NewRepository(mock)

	category := models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}

	_, err := repo.UpdateCategory(ctx, category)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= DELETE TESTS =============

func TestDeleteCategorySuccess(t *testing.T) {
	ctx := context.Background()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()

	repo := NewRepository(mock)

	err := repo.DeleteCategory(ctx, categoryID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteCategoryError(t *testing.T) {
	ctx := context.Background()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()
	mock.DeleteError = errors.New("database error")

	repo := NewRepository(mock)

	err := repo.DeleteCategory(ctx, categoryID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= HAS TRANSACTIONS TESTS =============

func TestHasTransactionsByCategorySuccess(t *testing.T) {
	ctx := context.Background()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()
	mock.HasTransResult = true

	repo := NewRepository(mock)

	result, err := repo.HasTransactionsByCategory(ctx, categoryID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestHasTransactionsByCategoryError(t *testing.T) {
	ctx := context.Background()
	categoryID := uuid.New()

	mock := newCategoryStoreMock()
	mock.HasTransError = errors.New("database error")

	repo := NewRepository(mock)

	_, err := repo.HasTransactionsByCategory(ctx, categoryID)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= STORE CATEGORY MODEL TO CATEGORY TESTS =============

func TestStoreCategoryModelToCategory(t *testing.T) {
	categoryID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	storeCategory := pgstore.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: 0,
		Name:            "Test",
		Icon:            "icon.png",
		CreatedAt:       createTimestamptzCategory(now),
		UpdatedAt:       createTimestamptzCategory(now),
	}

	result := StoreCategoryModelToCategory(storeCategory)

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}

	if result.TransactionType != models.Income {
		t.Errorf("Expected Income, got %v", result.TransactionType)
	}
}
