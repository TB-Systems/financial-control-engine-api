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

type TransactionsRepositoryMock struct {
	Error                  error
	CreateError            error
	CategoryError          error
	CreditcardError        error
	ShortAnnualError       error
	ShortMonthlyError      error
	ShortInstallmentError  error
	DeleteError            error
	PayError               error
	UpdateError            error
	TransactionResult      models.ShortTransaction
	TransactionFullResult  models.Transaction
	ShortAnnualResult      models.ShortAnnualTransaction
	ShortMonthlyResult     models.ShortMonthlyTransaction
	ShortInstallmentResult models.ShortInstallmentTransaction
	TransactionsResult     []models.Transaction
	TransactionsCount      int64
	CategoryResult         models.Category
	CreditcardResult       models.CreditCard
	LastCreatedTransaction models.CreateTransaction
}

func NewTransactionsRepositoryMock() *TransactionsRepositoryMock {
	return &TransactionsRepositoryMock{
		Error:                  nil,
		CreateError:            nil,
		CategoryError:          nil,
		CreditcardError:        nil,
		ShortAnnualError:       nil,
		ShortMonthlyError:      nil,
		ShortInstallmentError:  nil,
		UpdateError:            nil,
		TransactionResult:      models.ShortTransaction{},
		TransactionFullResult:  models.Transaction{},
		ShortAnnualResult:      models.ShortAnnualTransaction{},
		ShortMonthlyResult:     models.ShortMonthlyTransaction{},
		ShortInstallmentResult: models.ShortInstallmentTransaction{},
		TransactionsResult:     []models.Transaction{},
		TransactionsCount:      0,
		CategoryResult:         models.Category{},
		CreditcardResult:       models.CreditCard{},
	}
}

func (t *TransactionsRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if t.CategoryError != nil {
		return models.Category{}, t.CategoryError
	}
	return t.CategoryResult, nil
}

func (t *TransactionsRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardId uuid.UUID) (models.CreditCard, error) {
	if t.CreditcardError != nil {
		return models.CreditCard{}, t.CreditcardError
	}
	return t.CreditcardResult, nil
}

func (t *TransactionsRepositoryMock) CreateTransaction(ctx context.Context, transaction models.CreateTransaction) (models.ShortTransaction, error) {
	t.LastCreatedTransaction = transaction

	if t.CreateError != nil {
		return models.ShortTransaction{}, t.CreateError
	}

	if t.Error != nil {
		return models.ShortTransaction{}, t.Error
	}
	return t.TransactionResult, nil
}

func (t *TransactionsRepositoryMock) ReadTransactions(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.Transaction, int64, error) {
	if t.Error != nil {
		return []models.Transaction{}, 0, t.Error
	}
	return t.TransactionsResult, t.TransactionsCount, nil
}

func (t *TransactionsRepositoryMock) ReadTransactionsInToDates(ctx context.Context, params commonsmodels.PaginatedParamsWithDateRange) ([]models.Transaction, int64, error) {
	if t.Error != nil {
		return []models.Transaction{}, 0, t.Error
	}
	return t.TransactionsResult, t.TransactionsCount, nil
}

func (t *TransactionsRepositoryMock) ReadTransactionsByMonthYear(ctx context.Context, params commonsmodels.PaginatedParamsWithMonthYear) ([]models.Transaction, int64, error) {
	if t.Error != nil {
		return []models.Transaction{}, 0, t.Error
	}
	return t.TransactionsResult, t.TransactionsCount, nil
}

func (t *TransactionsRepositoryMock) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	if t.ShortAnnualError != nil {
		return models.ShortAnnualTransaction{}, t.ShortAnnualError
	}
	if t.Error != nil {
		return models.ShortAnnualTransaction{}, t.Error
	}
	return t.ShortAnnualResult, nil
}

func (t *TransactionsRepositoryMock) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	if t.ShortMonthlyError != nil {
		return models.ShortMonthlyTransaction{}, t.ShortMonthlyError
	}
	if t.Error != nil {
		return models.ShortMonthlyTransaction{}, t.Error
	}
	return t.ShortMonthlyResult, nil
}

func (t *TransactionsRepositoryMock) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	if t.ShortInstallmentError != nil {
		return models.ShortInstallmentTransaction{}, t.ShortInstallmentError
	}
	if t.Error != nil {
		return models.ShortInstallmentTransaction{}, t.Error
	}
	return t.ShortInstallmentResult, nil
}

func (t *TransactionsRepositoryMock) ReadTransactionById(ctx context.Context, id uuid.UUID) (models.Transaction, error) {
	if t.Error != nil {
		return models.Transaction{}, t.Error
	}
	return t.TransactionFullResult, nil
}

func (t *TransactionsRepositoryMock) UpdateTransaction(ctx context.Context, transaction models.Transaction) (models.ShortTransaction, error) {
	if t.UpdateError != nil {
		return models.ShortTransaction{}, t.UpdateError
	}
	if t.Error != nil {
		return models.ShortTransaction{}, t.Error
	}
	return t.TransactionResult, nil
}

func (t *TransactionsRepositoryMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	if t.DeleteError != nil {
		return t.DeleteError
	}
	if t.Error != nil {
		return t.Error
	}
	return nil
}

func (t *TransactionsRepositoryMock) PayTransaction(ctx context.Context, id uuid.UUID, paid bool) error {
	if t.PayError != nil {
		return t.PayError
	}
	if t.Error != nil {
		return t.Error
	}
	return nil
}

// ============= CREATE TESTS =============

func TestCreateTransactionSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.TransactionResult = models.ShortTransaction{
		ID:        transactionID,
		Name:      request.Name,
		Date:      request.Date,
		Value:     request.Value,
		Paid:      request.Paid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service := NewTransactionsService(mock)

	response, apiErr := service.Create(context.Background(), request, userID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestCreateTransactionCategoryNotFound(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryError = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionCreditRequiresCreditcard(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionDebitCannotHaveCreditcard(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test Transaction",
		Date:         time.Now(),
		Value:        100.00,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Debit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditcardID,
		UserID: userID,
		Name:   "Test Card",
	}

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionWithValidCreditcard(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()
	transactionID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test Transaction",
		Date:         time.Now(),
		Value:        100.00,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditcardID,
		UserID: userID,
		Name:   "Test Card",
		Limit:  10000,
	}
	mock.TransactionResult = models.ShortTransaction{
		ID:        transactionID,
		Name:      request.Name,
		Date:      request.Date,
		Value:     request.Value,
		Paid:      request.Paid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service := NewTransactionsService(mock)

	response, apiErr := service.Create(context.Background(), request, userID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestCreateTransactionRepositoryError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.Error = errors.New("database error")

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionCategoryRepositoryError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryError = errors.New("database error")

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionCreditcardRepositoryError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test Transaction",
		Date:         time.Now(),
		Value:        100.00,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardError = errors.New("database error")

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionWrongCategoryUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       time.Now(),
		Value:      100.00,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          otherUserID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionWrongCreditcardUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test Transaction",
		Date:         time.Now(),
		Value:        100.00,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditcardID,
		UserID: otherUserID,
		Name:   "Test Card",
	}

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestCreateTransactionReadCreditcardNotFound(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test Transaction",
		Date:         time.Now(),
		Value:        100.00,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardError = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	_, apiErr := service.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

// ============= READ TESTS =============

func TestReadTransactionsSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID1 := uuid.New()
	transactionID2 := uuid.New()

	transactions := []models.Transaction{
		{
			ID:        transactionID1,
			UserID:    userID,
			Name:      "Transaction 1",
			Date:      time.Now(),
			Value:     100.00,
			Paid:      false,
			Category:  models.Category{ID: categoryID, Name: "Category"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        transactionID2,
			UserID:    userID,
			Name:      "Transaction 2",
			Date:      time.Now(),
			Value:     200.00,
			Paid:      true,
			Category:  models.Category{ID: categoryID, Name: "Category"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionsResult = transactions
	mock.TransactionsCount = 2

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	response, apiErr := service.Read(context.Background(), params)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if len(response.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(response.Items))
	}
}

func TestReadTransactionsRepositoryError(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	_, apiErr := service.Read(context.Background(), params)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestReadTransactionsEmpty(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.TransactionsResult = []models.Transaction{}
	mock.TransactionsCount = 0

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
		Page:   1,
	}

	response, apiErr := service.Read(context.Background(), params)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if len(response.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(response.Items))
	}
}

// ============= READ BY ID TESTS =============

func TestReadByIdTransactionSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Test Transaction",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	response, apiErr := service.ReadById(context.Background(), userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestReadByIdTransactionNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	_, apiErr := service.ReadById(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestReadByIdTransactionRepositoryError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewTransactionsService(mock)

	_, apiErr := service.ReadById(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestReadByIdTransactionWrongUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Test Transaction",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	_, apiErr := service.ReadById(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error for wrong user, got none")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateTransactionSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       time.Now(),
		Value:      200.00,
		Paid:       true,
		CategoryID: categoryID,
	}

	originalTransaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Original Transaction",
		Date:      time.Now().AddDate(0, 0, -1),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.TransactionFullResult = originalTransaction
	mock.TransactionResult = models.ShortTransaction{
		ID:        transactionID,
		Name:      request.Name,
		Date:      request.Date,
		Value:     request.Value,
		Paid:      request.Paid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service := NewTransactionsService(mock)

	response, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestUpdateTransactionNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       time.Now(),
		Value:      200.00,
		Paid:       true,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.Error = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	_, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestUpdateTransactionRepositoryError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       time.Now(),
		Value:      200.00,
		Paid:       true,
		CategoryID: categoryID,
	}

	originalTransaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Original",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.TransactionFullResult = originalTransaction
	mock.UpdateError = errors.New("database error")

	service := NewTransactionsService(mock)

	_, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestUpdateTransactionWithCreditcard(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()
	transactionID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Updated",
		Date:         time.Now(),
		Value:        123,
		Paid:         true,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Credit}
	mock.CreditcardResult = models.CreditCard{ID: creditcardID, UserID: userID, Limit: 10000}
	mock.TransactionFullResult = models.Transaction{ID: transactionID, UserID: userID}
	mock.TransactionResult = models.ShortTransaction{ID: transactionID}

	service := NewTransactionsService(mock)

	response, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestUpdateTransactionWrongUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Updated",
		Date:       time.Now(),
		Value:      200.00,
		Paid:       true,
		CategoryID: categoryID,
	}

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Original",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	_, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

// ============= DELETE TESTS =============

func TestDeleteTransactionSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Test Transaction",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	apiErr := service.Delete(context.Background(), userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}
}

func TestDeleteTransactionNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	apiErr := service.Delete(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestDeleteTransactionRepositoryError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = models.Transaction{
		ID:     transactionID,
		UserID: userID,
	}
	mock.DeleteError = errors.New("delete failed")

	service := NewTransactionsService(mock)

	apiErr := service.Delete(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestDeleteTransactionWrongUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Test",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	apiErr := service.Delete(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

// ============= PAY TESTS =============

func TestPayTransactionSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Test Transaction",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	apiErr := service.Pay(context.Background(), userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}
}

func TestPayTransactionNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New(string(constants.StoreErrorNoRowsMsg))

	service := NewTransactionsService(mock)

	apiErr := service.Pay(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestPayTransactionToggleStatus(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Test Transaction",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      true,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	apiErr := service.Pay(context.Background(), userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}
}

func TestPayTransactionRepositoryError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = models.Transaction{
		ID:     transactionID,
		UserID: userID,
	}
	mock.PayError = errors.New("pay failed")

	service := NewTransactionsService(mock)

	apiErr := service.Pay(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestPayTransactionWrongUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	transaction := models.Transaction{
		ID:        transactionID,
		UserID:    otherUserID,
		Name:      "Test",
		Date:      time.Now(),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionFullResult = transaction

	service := NewTransactionsService(mock)

	apiErr := service.Pay(context.Background(), userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

// ============= GETRELATIONS TESTS =============

func TestGetRelationsCategoryInternalError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Test",
		Date:       time.Now(),
		Value:      10,
		Paid:       false,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryError = errors.New("db error")

	svc := NewTransactionsService(mock)

	_, apiErr := svc.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetRelationsDebitWithCreditcard(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Test",
		Date:         time.Now(),
		Value:        100,
		Paid:         false,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Debit,
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditcardID,
		UserID: userID,
		Limit:  10000,
	}

	svc := NewTransactionsService(mock)

	_, apiErr := svc.Create(context.Background(), request, userID)

	if apiErr == nil {
		t.Fatalf("expected error, got none")
	}
}

// ============= READINTODATES TESTS =============

func TestReadInToDatesSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	transactions := []models.Transaction{
		{
			ID:        transactionID,
			UserID:    userID,
			Name:      "Transaction 1",
			Date:      now,
			Value:     100.00,
			Paid:      false,
			Category:  models.Category{ID: categoryID, Name: "Category"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mock := NewTransactionsRepositoryMock()
	mock.TransactionsResult = transactions
	mock.TransactionsCount = 1

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		Page:      1,
		StartDate: now.AddDate(0, 0, -7),
		EndDate:   now,
	}

	svcImpl := service.(*transaction)
	response, apiErr := svcImpl.ReadInToDates(context.Background(), params)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if len(response.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(response.Items))
	}
}

func TestReadInToDatesRepositoryError(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		Page:      1,
		StartDate: now.AddDate(0, 0, -7),
		EndDate:   now,
	}

	svcImpl := service.(*transaction)
	_, apiErr := svcImpl.ReadInToDates(context.Background(), params)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestReadInToDatesEmpty(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.TransactionsResult = []models.Transaction{}
	mock.TransactionsCount = 0

	service := NewTransactionsService(mock)

	params := commonsmodels.PaginatedParamsWithDateRange{
		UserID:    userID,
		Limit:     10,
		Offset:    0,
		Page:      1,
		StartDate: now.AddDate(0, 0, -7),
		EndDate:   now,
	}

	svcImpl := service.(*transaction)
	response, apiErr := svcImpl.ReadInToDates(context.Background(), params)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if len(response.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(response.Items))
	}
}

// ============= UPDATE WITH CREDITCARD BRANCH TEST =============

func TestUpdateTransactionWithCreditcardFullBranch(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	creditcardID := uuid.New()
	transactionID := uuid.New()

	request := dtos.TransactionRequest{
		Name:         "Updated with Creditcard",
		Date:         time.Now(),
		Value:        150.00,
		Paid:         true,
		CategoryID:   categoryID,
		CreditcardID: &creditcardID,
	}

	originalTransaction := models.Transaction{
		ID:        transactionID,
		UserID:    userID,
		Name:      "Original Transaction",
		Date:      time.Now().AddDate(0, 0, -1),
		Value:     100.00,
		Paid:      false,
		Category:  models.Category{ID: categoryID, Name: "Category"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{
		ID:              categoryID,
		UserID:          userID,
		TransactionType: models.Credit,
		Name:            "Credit Category",
		Icon:            "icon",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	mock.CreditcardResult = models.CreditCard{
		ID:     creditcardID,
		UserID: userID,
		Name:   "Test Card",
		Limit:  10000,
	}
	mock.TransactionFullResult = originalTransaction
	mock.TransactionResult = models.ShortTransaction{
		ID:        transactionID,
		Name:      request.Name,
		Date:      request.Date,
		Value:     request.Value,
		Paid:      request.Paid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service := NewTransactionsService(mock)

	response, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if response.ID != transactionID {
		t.Errorf("expected ID %v, got %v", transactionID, response.ID)
	}

	if response.Creditcard == nil {
		t.Errorf("expected creditcard to be set, got nil")
	}
}

func TestUpdateTransactionGetRelationsError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	request := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       time.Now(),
		Value:      200.00,
		Paid:       true,
		CategoryID: categoryID,
	}

	mock := NewTransactionsRepositoryMock()
	mock.CategoryError = errors.New("category not found")

	service := NewTransactionsService(mock)

	_, apiErr := service.Update(context.Background(), request, userID, transactionID)

	if apiErr == nil {
		t.Errorf("expected error from getRelations, got none")
	}
}

func TestReadAtMonthAndYearSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.TransactionsResult = []models.Transaction{
		{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      "Transaction month/year",
			Date:      now,
			Value:     250.00,
			Paid:      false,
			Category:  models.Category{ID: categoryID, Name: "Category"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	mock.TransactionsCount = 1

	service := NewTransactionsService(mock)
	params := commonsmodels.PaginatedParamsWithMonthYear{
		UserID:     userID,
		Year:       2025,
		Month:      9,
		Page:       1,
		PageLimit:  10,
		PageOffset: 0,
	}

	response, apiErr := service.(*transaction).ReadAtMonthAndYear(context.Background(), params)

	if apiErr != nil {
		t.Errorf("expected no error, got %v", apiErr)
	}

	if len(response.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(response.Items))
	}
}

func TestReadAtMonthAndYearError(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.Error = errors.New("database error")

	service := NewTransactionsService(mock)
	params := commonsmodels.PaginatedParamsWithMonthYear{
		UserID:     userID,
		Year:       2025,
		Month:      9,
		Page:       1,
		PageLimit:  10,
		PageOffset: 0,
	}

	_, apiErr := service.(*transaction).ReadAtMonthAndYear(context.Background(), params)

	if apiErr == nil {
		t.Errorf("expected error, got none")
	}
}

func TestGetRelationsWithRelatedTransactionIDs(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}

	relations, annualErr := getRelations(mock, ctx, userID, nil, categoryID)
	if annualErr != nil {
		t.Fatalf("expected no error for annual relation id, got %v", annualErr)
	}
	if relations.CategoryResponse.ID != categoryID {
		t.Fatalf("expected category id %v, got %v", categoryID, relations.CategoryResponse.ID)
	}

	_, monthlyErr := getRelations(mock, ctx, userID, nil, categoryID)
	if monthlyErr != nil {
		t.Fatalf("expected no error for monthly relation id, got %v", monthlyErr)
	}

	_, installmentErr := getRelations(mock, ctx, userID, nil, categoryID)
	if installmentErr != nil {
		t.Fatalf("expected no error for installment relation id, got %v", installmentErr)
	}

	_, combinedErr := getRelations(mock, ctx, userID, nil, categoryID)
	if combinedErr != nil {
		t.Fatalf("expected no error for combined relation ids, got %v", combinedErr)
	}
}

func TestCreateFromMonthlyTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	monthlyID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortMonthlyResult = models.ShortMonthlyTransaction{
		ID:         monthlyID,
		UserID:     userID,
		Name:       "Internet",
		Day:        7,
		Value:      120.5,
		CategoryID: categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.TransactionResult = models.ShortTransaction{ID: uuid.New(), Name: "Internet", Date: now, Value: 120.5}

	svc := NewTransactionsService(mock).(*transaction)
	response, apiErr := svc.CreateFromMonthlyTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: monthlyID}, userID)

	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}

	if response.MonthlyTransaction == nil || response.MonthlyTransaction.ID != monthlyID {
		t.Fatalf("expected monthly transaction relation in response")
	}

	if mock.LastCreatedTransaction.Date.Day() != 7 || mock.LastCreatedTransaction.Date.Month() != now.Month() || mock.LastCreatedTransaction.Date.Year() != now.Year() {
		t.Fatalf("expected created date to use current month/year and monthly day")
	}
}

func TestCreateFromMonthlyTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	monthlyID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.ShortMonthlyError = errors.New(string(constants.StoreErrorNoRowsMsg))

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromMonthlyTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: monthlyID}, userID)

	if apiErr == nil {
		t.Fatalf("expected read monthly error")
	}
}

func TestCreateFromMonthlyTransactionGetRelationsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	monthlyID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.ShortMonthlyResult = models.ShortMonthlyTransaction{ID: monthlyID, UserID: userID, Name: "Internet", Day: 7, Value: 120.5, CategoryID: categoryID}
	mock.CategoryError = errors.New("category error")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromMonthlyTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: monthlyID}, userID)

	if apiErr == nil {
		t.Fatalf("expected getRelations error")
	}
}

func TestCreateFromMonthlyTransactionCreateError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	monthlyID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.ShortMonthlyResult = models.ShortMonthlyTransaction{ID: monthlyID, UserID: userID, Name: "Internet", Day: 7, Value: 120.5, CategoryID: categoryID}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.CreateError = errors.New("create failed")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromMonthlyTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: monthlyID}, userID)

	if apiErr == nil {
		t.Fatalf("expected create error")
	}
}

func TestCreateFromAnnualTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	annualID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortAnnualResult = models.ShortAnnualTransaction{
		ID:         annualID,
		UserID:     userID,
		Name:       "IPVA",
		Day:        10,
		Month:      int32(now.Month()),
		Value:      999.99,
		CategoryID: categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.TransactionResult = models.ShortTransaction{ID: uuid.New(), Name: "IPVA", Date: now, Value: 999.99}

	svc := NewTransactionsService(mock).(*transaction)
	response, apiErr := svc.CreateFromAnnualTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: annualID}, userID)

	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}

	if response.AnnualTransaction == nil || response.AnnualTransaction.ID != annualID {
		t.Fatalf("expected annual transaction relation in response")
	}

	if mock.LastCreatedTransaction.Date.Day() != 10 || mock.LastCreatedTransaction.Date.Month() != now.Month() || mock.LastCreatedTransaction.Date.Year() != now.Year() {
		t.Fatalf("expected created date to use current year with annual month/day")
	}
}

func TestCreateFromAnnualTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	annualID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.ShortAnnualError = errors.New(string(constants.StoreErrorNoRowsMsg))

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromAnnualTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: annualID}, userID)

	if apiErr == nil {
		t.Fatalf("expected read annual error")
	}
}

func TestCreateFromAnnualTransactionGetRelationsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	annualID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortAnnualResult = models.ShortAnnualTransaction{ID: annualID, UserID: userID, Name: "IPVA", Day: 10, Month: int32(now.Month()), Value: 999.99, CategoryID: categoryID}
	mock.CategoryError = errors.New("category error")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromAnnualTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: annualID}, userID)

	if apiErr == nil {
		t.Fatalf("expected getRelations error")
	}
}

func TestCreateFromAnnualTransactionCreateError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	annualID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortAnnualResult = models.ShortAnnualTransaction{ID: annualID, UserID: userID, Name: "IPVA", Day: 10, Month: int32(now.Month()), Value: 999.99, CategoryID: categoryID}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.CreateError = errors.New("create failed")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromAnnualTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: annualID}, userID)

	if apiErr == nil {
		t.Fatalf("expected create error")
	}
}

func TestCreateFromInstallmentTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()
	initial := time.Date(now.Year(), now.Month()-1, 5, 0, 0, 0, 0, now.Location())
	final := time.Date(now.Year(), now.Month()+1, 5, 0, 0, 0, 0, now.Location())

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentResult = models.ShortInstallmentTransaction{
		ID:          installmentID,
		UserID:      userID,
		Name:        "Notebook",
		InitialDate: initial,
		FinalDate:   final,
		Value:       300.0,
		CategoryID:  categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.TransactionResult = models.ShortTransaction{ID: uuid.New(), Name: "Notebook 2/3", Date: now, Value: 300.0}

	svc := NewTransactionsService(mock).(*transaction)
	response, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}

	if response.InstallmentTransaction == nil || response.InstallmentTransaction.ID != installmentID {
		t.Fatalf("expected installment relation in response")
	}

	if mock.LastCreatedTransaction.Name != "Notebook 2/3" {
		t.Fatalf("expected installment name with current/total, got %s", mock.LastCreatedTransaction.Name)
	}

	if mock.LastCreatedTransaction.Date.Day() != 5 || mock.LastCreatedTransaction.Date.Month() != now.Month() || mock.LastCreatedTransaction.Date.Year() != now.Year() {
		t.Fatalf("expected created date to use current month/year and installment initial day")
	}
}

func TestCreateFromInstallmentTransactionReadError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	installmentID := uuid.New()

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentError = errors.New(string(constants.StoreErrorNoRowsMsg))

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr == nil {
		t.Fatalf("expected read installment error")
	}
}

func TestCreateFromInstallmentTransactionGetRelationsError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentResult = models.ShortInstallmentTransaction{
		ID:          installmentID,
		UserID:      userID,
		Name:        "Notebook",
		InitialDate: time.Date(now.Year(), now.Month(), 5, 0, 0, 0, 0, now.Location()),
		FinalDate:   time.Date(now.Year(), now.Month()+2, 5, 0, 0, 0, 0, now.Location()),
		Value:       300.0,
		CategoryID:  categoryID,
	}
	mock.CategoryError = errors.New("category error")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr == nil {
		t.Fatalf("expected getRelations error")
	}
}

func TestCreateFromInstallmentTransactionCreateError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentResult = models.ShortInstallmentTransaction{
		ID:          installmentID,
		UserID:      userID,
		Name:        "Notebook",
		InitialDate: time.Date(now.Year(), now.Month(), 5, 0, 0, 0, 0, now.Location()),
		FinalDate:   time.Date(now.Year(), now.Month()+2, 5, 0, 0, 0, 0, now.Location()),
		Value:       300.0,
		CategoryID:  categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.CreateError = errors.New("create failed")

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr == nil {
		t.Fatalf("expected create error")
	}
}

func TestCreateFromInstallmentTransactionClampsToOneWhenInvalidRangeAndFutureInitial(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()

	initial := time.Date(now.Year(), now.Month()+2, 5, 0, 0, 0, 0, now.Location())
	final := time.Date(now.Year(), now.Month(), 5, 0, 0, 0, 0, now.Location())

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentResult = models.ShortInstallmentTransaction{
		ID:          installmentID,
		UserID:      userID,
		Name:        "Curso",
		InitialDate: initial,
		FinalDate:   final,
		Value:       150.0,
		CategoryID:  categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.TransactionResult = models.ShortTransaction{ID: uuid.New(), Name: "Curso 1/1", Date: now, Value: 150.0}

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}

	if mock.LastCreatedTransaction.Name != "Curso 1/1" {
		t.Fatalf("expected clamped installment name Curso 1/1, got %s", mock.LastCreatedTransaction.Name)
	}
}

func TestCreateFromInstallmentTransactionClampsCurrentToTotalWhenPastFinal(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	installmentID := uuid.New()
	now := time.Now()

	initial := time.Date(now.Year(), now.Month()-5, 5, 0, 0, 0, 0, now.Location())
	final := time.Date(now.Year(), now.Month()-3, 5, 0, 0, 0, 0, now.Location())

	mock := NewTransactionsRepositoryMock()
	mock.ShortInstallmentResult = models.ShortInstallmentTransaction{
		ID:          installmentID,
		UserID:      userID,
		Name:        "TV",
		InitialDate: initial,
		FinalDate:   final,
		Value:       200.0,
		CategoryID:  categoryID,
	}
	mock.CategoryResult = models.Category{ID: categoryID, UserID: userID, TransactionType: models.Debit}
	mock.TransactionResult = models.ShortTransaction{ID: uuid.New(), Name: "TV 3/3", Date: now, Value: 200.0}

	svc := NewTransactionsService(mock).(*transaction)
	_, apiErr := svc.CreateFromInstallmentTransaction(ctx, dtos.TransactionRequestFromRecurrentTransaction{ID: installmentID}, userID)

	if apiErr != nil {
		t.Fatalf("expected no error, got %v", apiErr)
	}

	if mock.LastCreatedTransaction.Name != "TV 3/3" {
		t.Fatalf("expected clamped installment name TV 3/3, got %s", mock.LastCreatedTransaction.Name)
	}
}
