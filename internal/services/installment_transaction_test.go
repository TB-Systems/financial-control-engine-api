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

type InstallmentTransactionRepositoryMock struct {
	Error                            error
	CategoryError                    error
	CreditcardError                  error
	CreateError                      error
	UpdateError                      error
	DeleteError                      error
	InstallmentTransactionResult     models.ShortInstallmentTransaction
	InstallmentTransactionFullResult models.InstallmentTransaction
	InstallmentTransactionsResult    []models.InstallmentTransaction
	InstallmentTransactionsCount     int64
	CategoryResult                   models.Category
	CreditcardResult                 models.CreditCard
}

func NewInstallmentTransactionRepositoryMock() *InstallmentTransactionRepositoryMock {
	return &InstallmentTransactionRepositoryMock{}
}

func (m *InstallmentTransactionRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if m.CategoryError != nil {
		return models.Category{}, m.CategoryError
	}
	return m.CategoryResult, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardID uuid.UUID) (models.CreditCard, error) {
	if m.CreditcardError != nil {
		return models.CreditCard{}, m.CreditcardError
	}
	return m.CreditcardResult, nil
}

func (m *InstallmentTransactionRepositoryMock) CreateInstallmentTransaction(ctx context.Context, request models.CreateInstallmentTransaction) (models.ShortInstallmentTransaction, error) {
	if m.CreateError != nil {
		return models.ShortInstallmentTransaction{}, m.CreateError
	}
	if m.Error != nil {
		return models.ShortInstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionResult, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadInstallmentTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.InstallmentTransaction, int64, error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.InstallmentTransactionsResult, m.InstallmentTransactionsCount, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.InstallmentTransaction, error) {
	if m.Error != nil {
		return models.InstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionFullResult, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	if m.Error != nil {
		return models.ShortInstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionResult, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	if m.Error != nil {
		return models.ShortAnnualTransaction{}, m.Error
	}
	return models.ShortAnnualTransaction{}, nil
}

func (m *InstallmentTransactionRepositoryMock) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	if m.Error != nil {
		return models.ShortMonthlyTransaction{}, m.Error
	}
	return models.ShortMonthlyTransaction{}, nil
}

func (m *InstallmentTransactionRepositoryMock) UpdateInstallmentTransaction(ctx context.Context, model models.InstallmentTransaction) (models.ShortInstallmentTransaction, error) {
	if m.UpdateError != nil {
		return models.ShortInstallmentTransaction{}, m.UpdateError
	}
	if m.Error != nil {
		return models.ShortInstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionResult, nil
}

func (m *InstallmentTransactionRepositoryMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

// ============= CREATE TESTS =============

func TestInstallmentTransactionCreateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.InstallmentTransactionResult = models.ShortInstallmentTransaction{
		ID:          transactionID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
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

func TestInstallmentTransactionCreateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
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
	mock.InstallmentTransactionResult = models.ShortInstallmentTransaction{
		ID:          transactionID,
		Name:        "Credit Installment",
		Value:       500.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:         "Credit Installment",
		Value:        500.00,
		InitialDate:  initialDate,
		FinalDate:    finalDate,
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
}

func TestInstallmentTransactionCreateCategoryNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateCategoryWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          uuid.New(), // Different user
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateCreditcardNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
	}
	mock.CreditcardError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:         "Test Installment",
		Value:        100.00,
		InitialDate:  now,
		FinalDate:    now.AddDate(0, 6, 0),
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateCreditcardWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: uuid.New(), // Different user
		Name:   "My Card",
		Limit:  5000.00,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:         "Test Installment",
		Value:        100.00,
		InitialDate:  now,
		FinalDate:    now.AddDate(0, 6, 0),
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateCreditWithoutCreditcard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit, // Credit type
		Name:            "Credit Category",
		Icon:            "icon",
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
		// No CreditCardID - should fail
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateDebitWithCreditcard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit, // Debit type
		Name:            "Debit Category",
		Icon:            "icon",
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditCardID,
		UserID: userID,
		Name:   "My Card",
		Limit:  5000.00,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:         "Test Installment",
		Value:        100.00,
		InitialDate:  now,
		FinalDate:    now.AddDate(0, 6, 0),
		CategoryID:   categoryID,
		CreditCardID: &creditCardID, // Should fail for debit
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionCreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.CreateError = errors.New("database error")

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ TESTS =============

func TestInstallmentTransactionReadSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionsResult = []models.InstallmentTransaction{
		{
			ID:          transactionID,
			UserID:      userID,
			Name:        "Test Installment",
			Value:       100.00,
			InitialDate: initialDate,
			FinalDate:   finalDate,
			Category: models.Category{
				ID:              categoryID,
				TransactionType: models.Debit,
				Name:            "Category",
				Icon:            "icon",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	mock.InstallmentTransactionsCount = 1

	service := NewInstallmentTransactionService(mock)

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
}

func TestInstallmentTransactionReadEmpty(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionsResult = []models.InstallmentTransaction{}
	mock.InstallmentTransactionsCount = 0

	service := NewInstallmentTransactionService(mock)

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

func TestInstallmentTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewInstallmentTransactionService(mock)

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
}

func TestInstallmentTransactionReadMultiplePages(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionsResult = []models.InstallmentTransaction{
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Installment 1",
			Value:       100.00,
			InitialDate: initialDate,
			FinalDate:   finalDate,
			Category: models.Category{
				ID:   uuid.New(),
				Name: "Category",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	mock.InstallmentTransactionsCount = 25 // 3 pages with limit 10

	service := NewInstallmentTransactionService(mock)

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
}

// ============= READ BY ID TESTS =============

func TestInstallmentTransactionReadByIdSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewInstallmentTransactionService(mock)

	result, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

func TestInstallmentTransactionReadByIdNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionReadByIdWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      uuid.New(), // Different user
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		Category: models.Category{
			ID:   uuid.New(),
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewInstallmentTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionReadByIdError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewInstallmentTransactionService(mock)

	_, apiErr := service.ReadById(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestInstallmentTransactionUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Original Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID:   categoryID,
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.InstallmentTransactionResult = models.ShortInstallmentTransaction{
		ID:          transactionID,
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}

	result, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Value != 150.00 {
		t.Errorf("Expected Value 150.00, got %v", result.Value)
	}
}

func TestInstallmentTransactionUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionUpdateCategoryNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryError = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionUpdateError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
	}
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Original Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID:   categoryID,
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.UpdateError = errors.New("database error")

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}

	_, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionUpdateWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionRepositoryMock()
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
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Original Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID:   categoryID,
			Name: "Category",
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.InstallmentTransactionResult = models.ShortInstallmentTransaction{
		ID:          transactionID,
		Name:        "Updated Installment",
		Value:       200.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	service := NewInstallmentTransactionService(mock)

	request := dtos.InstallmentTransactionRequest{
		Name:         "Updated Installment",
		Value:        200.00,
		InitialDate:  initialDate,
		FinalDate:    finalDate,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, apiErr := service.Update(ctx, userID, transactionID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Value != 200.00 {
		t.Errorf("Expected Value 200.00, got %v", result.Value)
	}
}

// ============= DELETE TESTS =============

func TestInstallmentTransactionDeleteSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		Category: models.Category{
			ID:   categoryID,
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewInstallmentTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}
}

func TestInstallmentTransactionDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewInstallmentTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionDeleteError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		Category: models.Category{
			ID:   categoryID,
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mock.DeleteError = errors.New("database error")

	service := NewInstallmentTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInstallmentTransactionDeleteWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionRepositoryMock()
	mock.InstallmentTransactionFullResult = models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      uuid.New(), // Different user
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		Category: models.Category{
			ID:   uuid.New(),
			Name: "Category",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	service := NewInstallmentTransactionService(mock)

	apiErr := service.Delete(ctx, userID, transactionID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}
}
