package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"backend-commons/constants"
	"backend-commons/dtos"
	"backend-commons/models"

	"github.com/google/uuid"
)

// ============= MOCK IMPLEMENTATION =============

type MonthlyTransactionRepositoryMock struct {
	Error                        error
	CategoryError                error
	CreditcardError              error
	CreateError                  error
	UpdateError                  error
	DeleteError                  error
	MonthlyTransactionResult     models.ShortMonthlyTransaction
	MonthlyTransactionFullResult models.MonthlyTransaction
	MonthlyTransactionsResult    []models.MonthlyTransaction
	MonthlyTransactionsCount     int64
	CategoryResult               models.Category
	CreditcardResult             models.CreditCard
}

func NewMonthlyTransactionRepositoryMock() *MonthlyTransactionRepositoryMock {
	return &MonthlyTransactionRepositoryMock{}
}

func (m *MonthlyTransactionRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if m.CategoryError != nil {
		return models.Category{}, m.CategoryError
	}
	return m.CategoryResult, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardID uuid.UUID) (models.CreditCard, error) {
	if m.CreditcardError != nil {
		return models.CreditCard{}, m.CreditcardError
	}
	return m.CreditcardResult, nil
}

func (m *MonthlyTransactionRepositoryMock) CreateTransaction(ctx context.Context, transaction models.CreateTransaction) (models.ShortTransaction, error) {
	return models.ShortTransaction{}, nil
}

func (m *MonthlyTransactionRepositoryMock) CreateMonthlyTransaction(ctx context.Context, request models.CreateMonthlyTransaction) (models.ShortMonthlyTransaction, error) {
	if m.CreateError != nil {
		return models.ShortMonthlyTransaction{}, m.CreateError
	}
	if m.Error != nil {
		return models.ShortMonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionResult, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadMonthlyTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.MonthlyTransaction, int64, error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.MonthlyTransactionsResult, m.MonthlyTransactionsCount, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.MonthlyTransaction, error) {
	if m.Error != nil {
		return models.MonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionFullResult, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	if m.Error != nil {
		return models.ShortMonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionResult, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	if m.Error != nil {
		return models.ShortAnnualTransaction{}, m.Error
	}
	return models.ShortAnnualTransaction{}, nil
}

func (m *MonthlyTransactionRepositoryMock) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	if m.Error != nil {
		return models.ShortInstallmentTransaction{}, m.Error
	}
	return models.ShortInstallmentTransaction{}, nil
}

func (m *MonthlyTransactionRepositoryMock) UpdateMonthlyTransaction(ctx context.Context, model models.MonthlyTransaction) (models.ShortMonthlyTransaction, error) {
	if m.UpdateError != nil {
		return models.ShortMonthlyTransaction{}, m.UpdateError
	}
	if m.Error != nil {
		return models.ShortMonthlyTransaction{}, m.Error
	}
	return m.MonthlyTransactionResult, nil
}

func (m *MonthlyTransactionRepositoryMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

// ============= CREATE TESTS =============

func TestMonthlyTransactionCreateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.MonthlyTransactionResult = models.ShortMonthlyTransaction{
		ID:        transactionID,
		Name:      "Test Monthly",
		Value:     100.00,
		Day:       15,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}

	result, apiErr := service.Create(ctx, userID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Value != 100.00 {
		t.Errorf("Expected Value 100.00, got %v", result.Value)
	}
}

func TestMonthlyTransactionCreateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: userID,
		Name:   "My Card",
		Limit:  5000.00,
	}
	mock.MonthlyTransactionResult = models.ShortMonthlyTransaction{
		ID:        transactionID,
		Name:      "Credit Monthly",
		Value:     500.00,
		Day:       10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Credit Monthly",
		Value:        500.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, apiErr := service.Create(ctx, userID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}
}

func TestMonthlyTransactionCreateCategoryNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCategoryBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          otherUserID,
		TransactionType: models.Debit,
		Name:            "Other User Category",
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCreditCardNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Credit Monthly",
		Value:        500.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCreditCardBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: otherUserID,
		Name:   "Other User Card",
		Limit:  5000.00,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Credit Monthly",
		Value:        500.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCreditWithoutCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit, // Credit type without creditcard
		Name:            "Credit Category",
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Credit Monthly",
		Value:      500.00,
		Day:        10,
		CategoryID: categoryID,
		// No CreditCardID
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 400 {
		t.Errorf("Expected status 400, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateDebitWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit, // Debit type with creditcard
		Name:            "Debit Category",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: userID,
		Name:   "My Card",
		Limit:  5000.00,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Debit Monthly",
		Value:        500.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 400 {
		t.Errorf("Expected status 400, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
	}
	mock.CreateError = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCategoryInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryError = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionCreateCreditcardInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardError = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Credit Monthly",
		Value:        500.00,
		Day:          10,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= READ TESTS =============

func TestMonthlyTransactionReadSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionsResult = []models.MonthlyTransaction{
		{
			ID:     transactionID,
			UserID: userID,
			Name:   "Monthly 1",
			Value:  100.00,
			Day:    5,
			Category: models.Category{
				ID:              categoryID,
				TransactionType: models.Debit,
				Name:            "Category 1",
				Icon:            "icon1",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	mock.MonthlyTransactionsCount = 1

	service := NewMonthlyTransactionService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	result, apiErr := service.Read(ctx, params)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %v", len(result.Items))
	}

	if result.Items[0].ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.Items[0].ID)
	}

	if result.PageCount != 1 {
		t.Errorf("Expected PageCount 1, got %v", result.PageCount)
	}
}

func TestMonthlyTransactionReadMultiplePages(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	transactions := make([]models.MonthlyTransaction, 25)
	for i := 0; i < 25; i++ {
		transactions[i] = models.MonthlyTransaction{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Monthly",
			Value:  100.00,
			Day:    int32(i + 1),
			Category: models.Category{
				ID:   uuid.New(),
				Name: "Category",
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionsResult = transactions[:10]
	mock.MonthlyTransactionsCount = 25

	service := NewMonthlyTransactionService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	result, apiErr := service.Read(ctx, params)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.PageCount != 3 {
		t.Errorf("Expected PageCount 3, got %v", result.PageCount)
	}

	if result.Page != 1 {
		t.Errorf("Expected Page 1, got %v", result.Page)
	}
}

func TestMonthlyTransactionReadEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionsResult = []models.MonthlyTransaction{}
	mock.MonthlyTransactionsCount = 0

	service := NewMonthlyTransactionService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	result, apiErr := service.Read(ctx, params)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if len(result.Items) != 0 {
		t.Errorf("Expected 0 items, got %v", len(result.Items))
	}
}

func TestMonthlyTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	_, apiErr := service.Read(ctx, params)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= READ BY ID TESTS =============

func TestMonthlyTransactionReadByIdSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Monthly Test",
		Value:  150.00,
		Day:    15,
		Category: models.Category{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	result, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Value != 150.00 {
		t.Errorf("Expected Value 150.00, got %v", result.Value)
	}
}

func TestMonthlyTransactionReadByIdNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionReadByIdBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Monthly Test",
		Value:     150.00,
		Day:       15,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionReadByIdInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= UPDATE TESTS =============

func TestMonthlyTransactionUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Updated Category",
		Icon:            "new-icon",
	}
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Monthly",
		Value:  100.00,
		Day:    10,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.MonthlyTransactionResult = models.ShortMonthlyTransaction{
		ID:        transactionID,
		Name:      "Updated Monthly",
		Value:     200.00,
		Day:       20,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}

	result, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Value != 200.00 {
		t.Errorf("Expected Value 200.00, got %v", result.Value)
	}

	if result.Day != 20 {
		t.Errorf("Expected Day 20, got %v", result.Day)
	}
}

func TestMonthlyTransactionUpdateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: userID,
		Name:   "My Card",
		Limit:  5000.00,
	}
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Monthly",
		Value:  100.00,
		Day:    10,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.MonthlyTransactionResult = models.ShortMonthlyTransaction{
		ID:        transactionID,
		Name:      "Updated Credit Monthly",
		Value:     300.00,
		Day:       15,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:         "Updated Credit Monthly",
		Value:        300.00,
		Day:          15,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}
}

func TestMonthlyTransactionUpdateRelationsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Category",
	}
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionUpdateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Category",
	}
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Monthly",
		Value:  100.00,
		Day:    10,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.UpdateError = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	request := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= DELETE TESTS =============

func TestMonthlyTransactionDeleteSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Monthly to Delete",
		Value:     100.00,
		Day:       10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}
}

func TestMonthlyTransactionDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewMonthlyTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionDeleteBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Monthly to Delete",
		Value:     100.00,
		Day:       10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestMonthlyTransactionDeleteRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Monthly to Delete",
		Value:     100.00,
		Day:       10,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.DeleteError = errors.New("database error")

	service := NewMonthlyTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= EDGE CASES =============

func TestMonthlyTransactionReadWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionsResult = []models.MonthlyTransaction{
		{
			ID:     transactionID,
			UserID: userID,
			Name:   "Monthly with CC",
			Value:  200.00,
			Day:    15,
			Category: models.Category{
				ID:              categoryID,
				TransactionType: models.Credit,
				Name:            "Credit Category",
				Icon:            "icon",
			},
			Creditcard: &models.CreditCard{
				ID:   creditCardID,
				Name: "My Card",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	mock.MonthlyTransactionsCount = 1

	service := NewMonthlyTransactionService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	result, apiErr := service.Read(ctx, params)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Items[0].Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result.Items[0].Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result.Items[0].Creditcard.ID)
	}
}

func TestMonthlyTransactionReadByIdWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionRepositoryMock()
	mock.MonthlyTransactionFullResult = models.MonthlyTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Monthly with CC",
		Value:  200.00,
		Day:    15,
		Category: models.Category{
			ID:              categoryID,
			TransactionType: models.Credit,
			Name:            "Credit Category",
			Icon:            "icon",
		},
		Creditcard: &models.CreditCard{
			ID:   creditCardID,
			Name: "My Card",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewMonthlyTransactionService(mock)

	result, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result.Creditcard.ID != creditCardID {
		t.Errorf("Expected Creditcard ID %v, got %v", creditCardID, result.Creditcard.ID)
	}
}
