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

type CategoryRepositoryMock struct {
	Error            error
	CreateError      error
	UpdateError      error
	DeleteError      error
	HasTransError    error
	CountError       error
	CategoryResult   models.Category
	CategoriesResult []models.Category
	CategoryCount    int64
	HasTransactions  bool
}

func NewCategoryRepositoryMock() *CategoryRepositoryMock {
	return &CategoryRepositoryMock{}
}

func (m *CategoryRepositoryMock) CreateCategory(ctx context.Context, data models.CreateCategory) (models.Category, error) {
	if m.CreateError != nil {
		return models.Category{}, m.CreateError
	}
	if m.Error != nil {
		return models.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryRepositoryMock) ReadCategories(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.CategoriesResult, nil
}

func (m *CategoryRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if m.Error != nil {
		return models.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryRepositoryMock) GetCategoryCountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	if m.CountError != nil {
		return 0, m.CountError
	}
	return m.CategoryCount, nil
}

func (m *CategoryRepositoryMock) UpdateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	if m.UpdateError != nil {
		return models.Category{}, m.UpdateError
	}
	if m.Error != nil {
		return models.Category{}, m.Error
	}
	return m.CategoryResult, nil
}

func (m *CategoryRepositoryMock) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	return m.Error
}

func (m *CategoryRepositoryMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	if m.HasTransError != nil {
		return false, m.HasTransError
	}
	return m.HasTransactions, nil
}

// ============= CREATE TESTS =============

func TestCategoryCreateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	transactionType := models.Income

	mock := NewCategoryRepositoryMock()
	mock.CategoryCount = 5
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}

	result, apiErr := service.Create(ctx, userID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}

	if result.Name != "Test Category" {
		t.Errorf("Expected name 'Test Category', got '%s'", result.Name)
	}
}

func TestCategoryCreateCountError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryRepositoryMock()
	mock.CountError = errors.New("database error")

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestCategoryCreateLimitReached(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryRepositoryMock()
	mock.CategoryCount = 10

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}

	_, apiErr := service.Create(ctx, userID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 403 {
		t.Errorf("Expected status 403, got %d", apiErr.GetStatus())
	}
}

func TestCategoryCreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryRepositoryMock()
	mock.CategoryCount = 5
	mock.CreateError = errors.New("database error")

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
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

func TestCategoryReadSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoriesResult = []models.Category{
		{
			ID:              uuid.New(),
			UserID:          userID,
			TransactionType: models.Income,
			Name:            "Category 1",
			Icon:            "icon1.png",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		{
			ID:              uuid.New(),
			UserID:          userID,
			TransactionType: models.Debit,
			Name:            "Category 2",
			Icon:            "icon2.png",
			CreatedAt:       now,
			UpdatedAt:       now,
		},
	}

	service := NewCategoriesService(mock)

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

func TestCategoryReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mock := NewCategoryRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewCategoriesService(mock)

	_, apiErr := service.Read(ctx, userID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= READ BY ID TESTS =============

func TestCategoryReadByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	service := NewCategoriesService(mock)

	result, apiErr := service.ReadByID(ctx, userID, categoryID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}
}

func TestCategoryReadByIDNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCategoriesService(mock)

	_, apiErr := service.ReadByID(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCategoryReadByIDWrongUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          otherUserID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	service := NewCategoriesService(mock)

	_, apiErr := service.ReadByID(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCategoryReadByIDInternalError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryRepositoryMock()
	mock.Error = errors.New("some database error")

	service := NewCategoriesService(mock)

	_, apiErr := service.ReadByID(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= UPDATE TESTS =============

func TestCategoryUpdateSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	transactionType := models.Debit

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Original Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}

	result, apiErr := service.Update(ctx, userID, categoryID, request)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}

	if result.ID != categoryID {
		t.Errorf("Expected ID %v, got %v", categoryID, result.ID)
	}
}

func TestCategoryUpdateNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionType := models.Debit

	mock := NewCategoryRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}

	_, apiErr := service.Update(ctx, userID, categoryID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCategoryUpdateRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	transactionType := models.Debit

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Original Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mock.UpdateError = errors.New("database error")

	service := NewCategoriesService(mock)

	request := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}

	_, apiErr := service.Update(ctx, userID, categoryID, request)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

// ============= DELETE TESTS =============

func TestCategoryDeleteSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mock.HasTransactions = false

	service := NewCategoriesService(mock)

	apiErr := service.Delete(ctx, userID, categoryID)

	if apiErr != nil {
		t.Errorf("Expected no error, got %v", apiErr)
	}
}

func TestCategoryDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryRepositoryMock()
	mock.Error = errors.New(constants.StoreErrorNoRowsMsg)

	service := NewCategoriesService(mock)

	apiErr := service.Delete(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 404 {
		t.Errorf("Expected status 404, got %d", apiErr.GetStatus())
	}
}

func TestCategoryDeleteHasTransactionsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mock.HasTransError = errors.New("database error")

	service := NewCategoriesService(mock)

	apiErr := service.Delete(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}

func TestCategoryDeleteHasTransactions(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mock.HasTransactions = true

	service := NewCategoriesService(mock)

	apiErr := service.Delete(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 400 {
		t.Errorf("Expected status 400, got %d", apiErr.GetStatus())
	}
}

func TestCategoryDeleteRepositoryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Income,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	mock.HasTransactions = false
	mock.DeleteError = errors.New("database error")

	service := NewCategoriesService(mock)

	apiErr := service.Delete(ctx, userID, categoryID)

	if apiErr == nil {
		t.Error("Expected error, got nil")
	}

	if apiErr.GetStatus() != 500 {
		t.Errorf("Expected status 500, got %d", apiErr.GetStatus())
	}
}
