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

type AnnualTransactionRepositoryMock struct {
	Error                       error
	CategoryError               error
	CreditcardError             error
	CreateError                 error
	UpdateError                 error
	DeleteError                 error
	AnnualTransactionResult     models.ShortAnnualTransaction
	AnnualTransactionFullResult models.AnnualTransaction
	AnnualTransactionsResult    []models.AnnualTransaction
	AnnualTransactionsCount     int64
	CategoryResult              models.Category
	CreditcardResult            models.CreditCard
}

func NewAnnualTransactionRepositoryMock() *AnnualTransactionRepositoryMock {
	return &AnnualTransactionRepositoryMock{}
}

func (m *AnnualTransactionRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if m.CategoryError != nil {
		return models.Category{}, m.CategoryError
	}
	return m.CategoryResult, nil
}

func (m *AnnualTransactionRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardID uuid.UUID) (models.CreditCard, error) {
	if m.CreditcardError != nil {
		return models.CreditCard{}, m.CreditcardError
	}
	return m.CreditcardResult, nil
}

func (m *AnnualTransactionRepositoryMock) CreateAnnualTransaction(ctx context.Context, request models.CreateAnnualTransaction) (models.ShortAnnualTransaction, error) {
	if m.CreateError != nil {
		return models.ShortAnnualTransaction{}, m.CreateError
	}
	if m.Error != nil {
		return models.ShortAnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionResult, nil
}

func (m *AnnualTransactionRepositoryMock) ReadAnnualTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.AnnualTransaction, int64, error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.AnnualTransactionsResult, m.AnnualTransactionsCount, nil
}

func (m *AnnualTransactionRepositoryMock) ReadAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.AnnualTransaction, error) {
	if m.Error != nil {
		return models.AnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionFullResult, nil
}

func (m *AnnualTransactionRepositoryMock) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	if m.Error != nil {
		return models.ShortAnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionResult, nil
}

func (m *AnnualTransactionRepositoryMock) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	if m.Error != nil {
		return models.ShortMonthlyTransaction{}, m.Error
	}
	return models.ShortMonthlyTransaction{}, nil
}

func (m *AnnualTransactionRepositoryMock) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	if m.Error != nil {
		return models.ShortInstallmentTransaction{}, m.Error
	}
	return models.ShortInstallmentTransaction{}, nil
}

func (m *AnnualTransactionRepositoryMock) UpdateAnnualTransaction(ctx context.Context, model models.AnnualTransaction) (models.ShortAnnualTransaction, error) {
	if m.UpdateError != nil {
		return models.ShortAnnualTransaction{}, m.UpdateError
	}
	if m.Error != nil {
		return models.ShortAnnualTransaction{}, m.Error
	}
	return m.AnnualTransactionResult, nil
}

func (m *AnnualTransactionRepositoryMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

// ============= CREATE TESTS =============

func TestAnnualTransactionCreateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.AnnualTransactionResult = models.ShortAnnualTransaction{
		ID:        transactionID,
		Name:      "Test Annual",
		Value:     100.00,
		Day:       15,
		Month:     6,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
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

	if result.Month != 6 {
		t.Errorf("Expected Month 6, got %v", result.Month)
	}
}

func TestAnnualTransactionCreateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
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
	mock.AnnualTransactionResult = models.ShortAnnualTransaction{
		ID:        transactionID,
		Name:      "Credit Annual",
		Value:     500.00,
		Day:       10,
		Month:     12,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Credit Annual",
		Value:        500.00,
		Day:          10,
		Month:        12,
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

	if result.Month != 12 {
		t.Errorf("Expected Month 12, got %v", result.Month)
	}
}

func TestAnnualTransactionCreateCategoryNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
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

func TestAnnualTransactionCreateCategoryBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          otherUserID,
		TransactionType: models.Debit,
		Name:            "Other User Category",
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
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

func TestAnnualTransactionCreateCreditCardNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Credit Annual",
		Value:        500.00,
		Day:          10,
		Month:        12,
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

func TestAnnualTransactionCreateCreditCardBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
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

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Credit Annual",
		Value:        500.00,
		Day:          10,
		Month:        12,
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

func TestAnnualTransactionCreateCreditWithoutCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit, // Credit type without creditcard
		Name:            "Credit Category",
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Credit Annual",
		Value:      500.00,
		Day:        10,
		Month:      12,
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

func TestAnnualTransactionCreateDebitWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
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

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Debit Annual",
		Value:        500.00,
		Day:          10,
		Month:        12,
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

func TestAnnualTransactionCreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
	}
	mock.CreateError = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
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

func TestAnnualTransactionCreateCategoryInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryError = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
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

func TestAnnualTransactionCreateCreditcardInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
	}
	mock.CreditcardError = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Credit Annual",
		Value:        500.00,
		Day:          10,
		Month:        12,
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

func TestAnnualTransactionReadSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionsResult = []models.AnnualTransaction{
		{
			ID:     transactionID,
			UserID: userID,
			Name:   "Annual 1",
			Value:  100.00,
			Day:    5,
			Month:  3,
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
	mock.AnnualTransactionsCount = 1

	service := NewAnnualTransactionService(mock)

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

	if result.Items[0].Month != 3 {
		t.Errorf("Expected Month 3, got %v", result.Items[0].Month)
	}
}

func TestAnnualTransactionReadMultiplePages(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	transactions := make([]models.AnnualTransaction, 25)
	for i := 0; i < 25; i++ {
		transactions[i] = models.AnnualTransaction{
			ID:     uuid.New(),
			UserID: userID,
			Name:   "Annual",
			Value:  100.00,
			Day:    int32((i % 28) + 1),
			Month:  int32((i % 12) + 1),
			Category: models.Category{
				ID:   uuid.New(),
				Name: "Category",
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionsResult = transactions[:10]
	mock.AnnualTransactionsCount = 25

	service := NewAnnualTransactionService(mock)

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

func TestAnnualTransactionReadEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionsResult = []models.AnnualTransaction{}
	mock.AnnualTransactionsCount = 0

	service := NewAnnualTransactionService(mock)

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

func TestAnnualTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewAnnualTransactionService(mock)

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

func TestAnnualTransactionReadByIdSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Annual Test",
		Value:  150.00,
		Day:    15,
		Month:  7,
		Category: models.Category{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

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

	if result.Month != 7 {
		t.Errorf("Expected Month 7, got %v", result.Month)
	}
}

func TestAnnualTransactionReadByIdNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestAnnualTransactionReadByIdBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Annual Test",
		Value:     150.00,
		Day:       15,
		Month:     7,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestAnnualTransactionReadByIdInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= UPDATE TESTS =============

func TestAnnualTransactionUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Updated Category",
		Icon:            "new-icon",
	}
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Annual",
		Value:  100.00,
		Day:    10,
		Month:  5,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.AnnualTransactionResult = models.ShortAnnualTransaction{
		ID:        transactionID,
		Name:      "Updated Annual",
		Value:     200.00,
		Day:       20,
		Month:     10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
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

	if result.Month != 10 {
		t.Errorf("Expected Month 10, got %v", result.Month)
	}
}

func TestAnnualTransactionUpdateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
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
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Annual",
		Value:  100.00,
		Day:    10,
		Month:  5,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.AnnualTransactionResult = models.ShortAnnualTransaction{
		ID:        transactionID,
		Name:      "Updated Credit Annual",
		Value:     300.00,
		Day:       15,
		Month:     8,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:         "Updated Credit Annual",
		Value:        300.00,
		Day:          15,
		Month:        8,
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

	if result.Month != 8 {
		t.Errorf("Expected Month 8, got %v", result.Month)
	}
}

func TestAnnualTransactionUpdateRelationsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
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

func TestAnnualTransactionUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Category",
	}
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
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

func TestAnnualTransactionUpdateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Category",
	}
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Old Annual",
		Value:  100.00,
		Day:    10,
		Month:  5,
		Category: models.Category{
			ID: categoryID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.UpdateError = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	request := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
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

func TestAnnualTransactionDeleteSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Annual to Delete",
		Value:     100.00,
		Day:       10,
		Month:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}
}

func TestAnnualTransactionDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewAnnualTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestAnnualTransactionDeleteBelongsToOtherUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Annual to Delete",
		Value:     100.00,
		Day:       10,
		Month:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewAnnualTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %v", apiErr.GetStatus())
	}
}

func TestAnnualTransactionDeleteRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Annual to Delete",
		Value:     100.00,
		Day:       10,
		Month:     5,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.DeleteError = errors.New("database error")

	service := NewAnnualTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %v", apiErr.GetStatus())
	}
}

// ============= EDGE CASES =============

func TestAnnualTransactionReadWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionsResult = []models.AnnualTransaction{
		{
			ID:     transactionID,
			UserID: userID,
			Name:   "Annual with CC",
			Value:  200.00,
			Day:    15,
			Month:  9,
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
	mock.AnnualTransactionsCount = 1

	service := NewAnnualTransactionService(mock)

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

	if result.Items[0].Month != 9 {
		t.Errorf("Expected Month 9, got %v", result.Items[0].Month)
	}
}

func TestAnnualTransactionReadByIdWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionRepositoryMock()
	mock.AnnualTransactionFullResult = models.AnnualTransaction{
		ID:     transactionID,
		UserID: userID,
		Name:   "Annual with CC",
		Value:  200.00,
		Day:    15,
		Month:  11,
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

	service := NewAnnualTransactionService(mock)

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

	if result.Month != 11 {
		t.Errorf("Expected Month 11, got %v", result.Month)
	}
}
