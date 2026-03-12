package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend-commons/constants"
	"backend-commons/dtos"
	"backend-commons/models"

	"github.com/google/uuid"
)

// ============= MOCK IMPLEMENTATION =============

type CreditCardRepositoryMock struct {
	Error             error
	CreateError       error
	UpdateError       error
	DeleteError       error
	HasTransError     error
	CountError        error
	CreditCardResult  models.CreditCard
	CreditCardsResult []models.CreditCard
	CreditCardCount   int
	HasTransactions   bool
}

func NewCreditCardRepositoryMock() *CreditCardRepositoryMock {
	return &CreditCardRepositoryMock{}
}

func (m *CreditCardRepositoryMock) CreateCreditCard(ctx context.Context, creditCard models.CreateCreditCard) (models.CreditCard, error) {
	if m.CreateError != nil {
		return models.CreditCard{}, m.CreateError
	}
	if m.Error != nil {
		return models.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardRepositoryMock) ReadCreditCards(ctx context.Context, userID uuid.UUID) ([]models.CreditCard, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.CreditCardsResult, nil
}

func (m *CreditCardRepositoryMock) ReadCountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	if m.CountError != nil {
		return 0, m.CountError
	}
	return m.CreditCardCount, nil
}

func (m *CreditCardRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardID uuid.UUID) (models.CreditCard, error) {
	if m.Error != nil {
		return models.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardRepositoryMock) UpdateCreditCard(ctx context.Context, creditCard models.CreditCard) (models.CreditCard, error) {
	if m.UpdateError != nil {
		return models.CreditCard{}, m.UpdateError
	}
	if m.Error != nil {
		return models.CreditCard{}, m.Error
	}
	return m.CreditCardResult, nil
}

func (m *CreditCardRepositoryMock) DeleteCreditCard(ctx context.Context, creditCardID uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

func (m *CreditCardRepositoryMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID uuid.UUID) (bool, error) {
	if m.HasTransError != nil {
		return false, m.HasTransError
	}
	return m.HasTransactions, nil
}

// ============= CREATE TESTS =============

func TestCreditCardCreateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardCount = 5
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	result, apiErr := service.Create(ctx, userID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}

	if result.Name != "Test Card" {
		t.Errorf("Expected name 'Test Card', got '%s'", result.Name)
	}
}

func TestCreditCardCreateCountError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.CountError = errors.New("database error")

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardCreateLimitReached(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardCount = 10

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 403 {
		t.Errorf("Expected status 403, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardCreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardCount = 5
	mock.CreateError = errors.New("database error")

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= READ TESTS =============

func TestCreditCardReadSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardsResult = []models.CreditCard{
		{
			ID:               uuid.New(),
			UserID:           userID,
			Name:             "Card 1",
			FirstFourNumbers: "1234",
			Limit:            5000.00,
			CloseDay:         15,
			ExpireDay:        25,
			BackgroundColor:  "#000000",
			TextColor:        "#FFFFFF",
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               uuid.New(),
			UserID:           userID,
			Name:             "Card 2",
			FirstFourNumbers: "5678",
			Limit:            10000.00,
			CloseDay:         10,
			ExpireDay:        20,
			BackgroundColor:  "#FFFFFF",
			TextColor:        "#000000",
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	service := NewCreditCardsService(mock)

	result, apiErr := service.Read(ctx, userID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}

	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}
}

func TestCreditCardReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewCreditCardsService(mock)

	_, apiErr := service.Read(ctx, userID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= READ AT TESTS =============

func TestCreditCardReadAtSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	service := NewCreditCardsService(mock)

	result, apiErr := service.ReadAt(ctx, userID, creditCardID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}
}

func TestCreditCardReadAtNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCreditCardsService(mock)

	_, apiErr := service.ReadAt(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardReadAtWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           otherUserID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	service := NewCreditCardsService(mock)

	_, apiErr := service.ReadAt(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardReadAtInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.Error = errors.New("some database error")

	service := NewCreditCardsService(mock)

	_, apiErr := service.ReadAt(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= UPDATE TESTS =============

func TestCreditCardUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Original Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9876",
		Limit:            10000.00,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#FFFFFF",
		TextColor:        "#000000",
	}

	result, apiErr := service.Update(ctx, userID, creditCardID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != creditCardID {
		t.Errorf("Expected ID %v, got %v", creditCardID, result.ID)
	}
}

func TestCreditCardUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9876",
		Limit:            10000.00,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#FFFFFF",
		TextColor:        "#000000",
	}

	_, apiErr := service.Update(ctx, userID, creditCardID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardUpdateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Original Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	mock.UpdateError = errors.New("database error")

	service := NewCreditCardsService(mock)

	request := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9876",
		Limit:            10000.00,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#FFFFFF",
		TextColor:        "#000000",
	}

	_, apiErr := service.Update(ctx, userID, creditCardID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= DELETE TESTS =============

func TestCreditCardDeleteSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	mock.HasTransactions = false

	service := NewCreditCardsService(mock)

	apiErr := service.Delete(ctx, userID, creditCardID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}
}

func TestCreditCardDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCreditCardsService(mock)

	apiErr := service.Delete(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardDeleteHasTransactionsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	mock.HasTransError = errors.New("database error")

	service := NewCreditCardsService(mock)

	apiErr := service.Delete(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardDeleteHasTransactions(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	mock.HasTransactions = true

	service := NewCreditCardsService(mock)

	apiErr := service.Delete(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 400 {
		t.Errorf("Expected status 400, got %d", apiErr.GetStatus())
	}
}

func TestCreditCardDeleteRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardRepositoryMock()
	mock.CreditCardResult = models.CreditCard{
		ID:               creditCardID,
		UserID:           userID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.00,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	mock.HasTransactions = false
	mock.DeleteError = errors.New("database error")

	service := NewCreditCardsService(mock)

	apiErr := service.Delete(ctx, userID, creditCardID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}
