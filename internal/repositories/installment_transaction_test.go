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

// ============= STORE MOCK FOR INSTALLMENT =============

type InstallmentStoreMock struct {
	Error                           error
	InstallmentTransactionResult    pgstore.InstallmentTransaction
	InstallmentTransactionRowResult pgstore.GetInstallmentTransactionByIDRow
	ShortInstallmentRowResult       pgstore.InstallmentTransaction
	InstallmentTransactionsResult   []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow
}

// Installment Transaction methods
func (m *InstallmentStoreMock) CreateInstallmentTransaction(ctx context.Context, arg pgstore.CreateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	if m.Error != nil {
		return pgstore.InstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionResult, nil
}

func (m *InstallmentStoreMock) ListInstallmentTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListInstallmentTransactionsByUserIDPaginatedParams) ([]pgstore.ListInstallmentTransactionsByUserIDPaginatedRow, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.InstallmentTransactionsResult, nil
}

func (m *InstallmentStoreMock) GetInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetInstallmentTransactionByIDRow, error) {
	if m.Error != nil {
		return pgstore.GetInstallmentTransactionByIDRow{}, m.Error
	}
	return m.InstallmentTransactionRowResult, nil
}

func (m *InstallmentStoreMock) UpdateInstallmentTransaction(ctx context.Context, arg pgstore.UpdateInstallmentTransactionParams) (pgstore.InstallmentTransaction, error) {
	if m.Error != nil {
		return pgstore.InstallmentTransaction{}, m.Error
	}
	return m.InstallmentTransactionResult, nil
}

func (m *InstallmentStoreMock) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return m.Error
}

// Other required interface methods (stubs)
func (m *InstallmentStoreMock) CountCategoriesByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *InstallmentStoreMock) CountCreditCardsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *InstallmentStoreMock) CountMonthlyTransactionsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	return 0, nil
}
func (m *InstallmentStoreMock) CreateCategory(ctx context.Context, arg pgstore.CreateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *InstallmentStoreMock) CreateCreditCard(ctx context.Context, arg pgstore.CreateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *InstallmentStoreMock) CreateAnnualTransaction(ctx context.Context, arg pgstore.CreateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *InstallmentStoreMock) CreateMonthlyTransaction(ctx context.Context, arg pgstore.CreateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *InstallmentStoreMock) CreateTransaction(ctx context.Context, arg pgstore.CreateTransactionParams) (pgstore.CreateTransactionRow, error) {
	return pgstore.CreateTransactionRow{}, nil
}
func (m *InstallmentStoreMock) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *InstallmentStoreMock) DeleteCreditCard(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *InstallmentStoreMock) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *InstallmentStoreMock) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *InstallmentStoreMock) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *InstallmentStoreMock) GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]pgstore.Category, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) GetCategoryByID(ctx context.Context, id uuid.UUID) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *InstallmentStoreMock) GetCreditCardByID(ctx context.Context, id uuid.UUID) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *InstallmentStoreMock) GetAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetAnnualTransactionByIDRow, error) {
	return pgstore.GetAnnualTransactionByIDRow{}, nil
}
func (m *InstallmentStoreMock) GetMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetMonthlyTransactionByIDRow, error) {
	return pgstore.GetMonthlyTransactionByIDRow{}, nil
}
func (m *InstallmentStoreMock) GetTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.GetTransactionByIDRow, error) {
	return pgstore.GetTransactionByIDRow{}, nil
}
func (m *InstallmentStoreMock) HasTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	return false, nil
}
func (m *InstallmentStoreMock) HasTransactionsByCreditCard(ctx context.Context, creditCardID pgtype.UUID) (bool, error) {
	return false, nil
}
func (m *InstallmentStoreMock) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]pgstore.CreditCard, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) ListAnnualTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListAnnualTransactionsByUserIDPaginatedParams) ([]pgstore.ListAnnualTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) ListMonthlyTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListMonthlyTransactionsByUserIDPaginatedParams) ([]pgstore.ListMonthlyTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) ListTransactionsByUserAndDate(ctx context.Context, arg pgstore.ListTransactionsByUserAndDateParams) ([]pgstore.ListTransactionsByUserAndDateRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) ListTransactionsByUserAndMonthYearPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserAndMonthYearPaginatedParams) ([]pgstore.ListTransactionsByUserAndMonthYearPaginatedRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) ListTransactionsByUserIDPaginated(ctx context.Context, arg pgstore.ListTransactionsByUserIDPaginatedParams) ([]pgstore.ListTransactionsByUserIDPaginatedRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) PayTransaction(ctx context.Context, arg pgstore.PayTransactionParams) error {
	return nil
}
func (m *InstallmentStoreMock) UpdateCategory(ctx context.Context, arg pgstore.UpdateCategoryParams) (pgstore.Category, error) {
	return pgstore.Category{}, nil
}
func (m *InstallmentStoreMock) UpdateCreditCard(ctx context.Context, arg pgstore.UpdateCreditCardParams) (pgstore.CreditCard, error) {
	return pgstore.CreditCard{}, nil
}
func (m *InstallmentStoreMock) UpdateAnnualTransaction(ctx context.Context, arg pgstore.UpdateAnnualTransactionParams) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *InstallmentStoreMock) UpdateMonthlyTransaction(ctx context.Context, arg pgstore.UpdateMonthlyTransactionParams) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *InstallmentStoreMock) UpdateTransaction(ctx context.Context, arg pgstore.UpdateTransactionParams) (pgstore.Transaction, error) {
	return pgstore.Transaction{}, nil
}
func (m *InstallmentStoreMock) GetMonthlyBalance(ctx context.Context, arg pgstore.GetMonthlyBalanceParams) (pgstore.GetMonthlyBalanceRow, error) {
	return pgstore.GetMonthlyBalanceRow{}, nil
}
func (m *InstallmentStoreMock) GetCategoriesSpending(ctx context.Context, arg pgstore.GetCategoriesSpendingParams) ([]pgstore.GetCategoriesSpendingRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) GetCreditCardsSpending(ctx context.Context, arg pgstore.GetCreditCardsSpendingParams) ([]pgstore.GetCreditCardsSpendingRow, error) {
	return nil, nil
}
func (m *InstallmentStoreMock) GetShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.MonthlyTransaction, error) {
	return pgstore.MonthlyTransaction{}, nil
}
func (m *InstallmentStoreMock) GetShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.AnnualTransaction, error) {
	return pgstore.AnnualTransaction{}, nil
}
func (m *InstallmentStoreMock) GetShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (pgstore.InstallmentTransaction, error) {
	if m.Error != nil {
		return pgstore.InstallmentTransaction{}, m.Error
	}
	return m.ShortInstallmentRowResult, nil
}
func (m *InstallmentStoreMock) WithTx(tx pgx.Tx) *pgstore.Queries {
	return nil
}

// ============= HELPER FUNCTIONS =============

func createTimestamptzInstallment(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func createPgUUIDInstallment(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func createPgTextInstallment(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}

func createPgInt4Installment(i int32) pgtype.Int4 {
	return pgtype.Int4{Int32: i, Valid: true}
}

func createPgFloat8Installment(f float64) pgtype.Float8 {
	return pgtype.Float8{Float64: f, Valid: true}
}

func newInstallmentStoreMock() *InstallmentStoreMock {
	return &InstallmentStoreMock{}
}

// ============= CREATE TESTS =============

func TestCreateInstallmentTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionResult = pgstore.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Test Installment",
		Value:       utils.Float64ToNumeric(100.50),
		InitialDate: createTimestamptzInstallment(initialDate),
		FinalDate:   createTimestamptzInstallment(finalDate),
		CategoryID:  categoryID,
		CreatedAt:   createTimestamptzInstallment(now),
		UpdatedAt:   createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	request := models.CreateInstallmentTransaction{
		UserID:      userID,
		Name:        "Test Installment",
		Value:       100.50,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}

	result, err := repo.CreateInstallmentTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Value != 100.50 {
		t.Errorf("Expected Value 100.50, got %v", result.Value)
	}
}

func TestCreateInstallmentTransactionWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionResult = pgstore.InstallmentTransaction{
		ID:           transactionID,
		UserID:       userID,
		Name:         "Test with CC",
		Value:        utils.Float64ToNumeric(200.00),
		InitialDate:  createTimestamptzInstallment(initialDate),
		FinalDate:    createTimestamptzInstallment(finalDate),
		CategoryID:   categoryID,
		CreditCardID: createPgUUIDInstallment(creditCardID),
		CreatedAt:    createTimestamptzInstallment(now),
		UpdatedAt:    createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	request := models.CreateInstallmentTransaction{
		UserID:       userID,
		Name:         "Test with CC",
		Value:        200.00,
		InitialDate:  initialDate,
		FinalDate:    finalDate,
		CategoryID:   categoryID,
		CreditCardID: &creditCardID,
	}

	result, err := repo.CreateInstallmentTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}
}

func TestCreateInstallmentTransactionError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	request := models.CreateInstallmentTransaction{
		UserID:      uuid.New(),
		Name:        "Test",
		Value:       100.00,
		InitialDate: time.Now(),
		FinalDate:   time.Now().AddDate(0, 6, 0),
		CategoryID:  uuid.New(),
	}

	_, err := repo.CreateInstallmentTransaction(ctx, request)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= READ PAGINATED TESTS =============

func TestReadInstallmentTransactionsPaginatedSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionsResult = []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Test Installment",
			Value:                   utils.Float64ToNumeric(100.00),
			InitialDate:             createTimestamptzInstallment(initialDate),
			FinalDate:               createTimestamptzInstallment(finalDate),
			CategoryID:              createPgUUIDInstallment(categoryID),
			CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
			CategoryName:            createPgTextInstallment("Category"),
			CategoryIcon:            createPgTextInstallment("icon"),
			CreditcardID:            pgtype.UUID{Valid: false},
			CreatedAt:               createTimestamptzInstallment(now),
			UpdatedAt:               createTimestamptzInstallment(now),
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %v", count)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %v", len(result))
	}

	if result[0].Name != "Test Installment" {
		t.Errorf("Expected name 'Test Installment', got %v", result[0].Name)
	}
}

func TestReadInstallmentTransactionsPaginatedWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionsResult = []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow{
		{
			ID:                         transactionID,
			UserID:                     userID,
			Name:                       "Test with CC",
			Value:                      utils.Float64ToNumeric(200.00),
			InitialDate:                createTimestamptzInstallment(initialDate),
			FinalDate:                  createTimestamptzInstallment(finalDate),
			CategoryID:                 createPgUUIDInstallment(categoryID),
			CategoryTransactionType:    createPgInt4Installment(int32(models.Credit)),
			CategoryName:               createPgTextInstallment("Credit Category"),
			CategoryIcon:               createPgTextInstallment("credit-icon"),
			CreditcardID:               createPgUUIDInstallment(creditCardID),
			CreditcardName:             createPgTextInstallment("My Card"),
			CreditcardFirstFourNumbers: createPgTextInstallment("1234"),
			CreditcardCreditLimit:      createPgFloat8Installment(5000.00),
			CreditcardCloseDay:         createPgInt4Installment(10),
			CreditcardExpireDay:        createPgInt4Installment(20),
			CreditcardBackgroundColor:  createPgTextInstallment("#000000"),
			CreditcardTextColor:        createPgTextInstallment("#FFFFFF"),
			CreatedAt:                  createTimestamptzInstallment(now),
			UpdatedAt:                  createTimestamptzInstallment(now),
			TotalCount:                 1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, _, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result[0].Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result[0].Creditcard.Name != "My Card" {
		t.Errorf("Expected creditcard name 'My Card', got %v", result[0].Creditcard.Name)
	}
}

func TestReadInstallmentTransactionsPaginatedMultipleResults(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionsResult = []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow{
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Installment 1",
			Value:                   utils.Float64ToNumeric(100.00),
			InitialDate:             createTimestamptzInstallment(initialDate),
			FinalDate:               createTimestamptzInstallment(finalDate),
			CategoryID:              createPgUUIDInstallment(categoryID),
			CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
			CategoryName:            createPgTextInstallment("Category"),
			CategoryIcon:            createPgTextInstallment("icon"),
			CreatedAt:               createTimestamptzInstallment(now),
			UpdatedAt:               createTimestamptzInstallment(now),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Installment 2",
			Value:                   utils.Float64ToNumeric(200.00),
			InitialDate:             createTimestamptzInstallment(initialDate),
			FinalDate:               createTimestamptzInstallment(finalDate),
			CategoryID:              createPgUUIDInstallment(categoryID),
			CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
			CategoryName:            createPgTextInstallment("Category"),
			CategoryIcon:            createPgTextInstallment("icon"),
			CreatedAt:               createTimestamptzInstallment(now),
			UpdatedAt:               createTimestamptzInstallment(now),
			TotalCount:              3,
		},
		{
			ID:                      uuid.New(),
			UserID:                  userID,
			Name:                    "Installment 3",
			Value:                   utils.Float64ToNumeric(300.00),
			InitialDate:             createTimestamptzInstallment(initialDate),
			FinalDate:               createTimestamptzInstallment(finalDate),
			CategoryID:              createPgUUIDInstallment(categoryID),
			CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
			CategoryName:            createPgTextInstallment("Category"),
			CategoryIcon:            createPgTextInstallment("icon"),
			CreatedAt:               createTimestamptzInstallment(now),
			UpdatedAt:               createTimestamptzInstallment(now),
			TotalCount:              3,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

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

func TestReadInstallmentTransactionsPaginatedError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: uuid.New(),
		Limit:  10,
		Offset: 0,
	}

	_, _, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReadInstallmentTransactionsPaginatedEmpty(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionsResult = []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow{}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: uuid.New(),
		Limit:  10,
		Offset: 0,
	}

	result, count, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

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

func TestReadInstallmentTransactionsPaginatedWithNoCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionsResult = []pgstore.ListInstallmentTransactionsByUserIDPaginatedRow{
		{
			ID:                      transactionID,
			UserID:                  userID,
			Name:                    "Installment without CC",
			Value:                   utils.Float64ToNumeric(150.00),
			InitialDate:             createTimestamptzInstallment(initialDate),
			FinalDate:               createTimestamptzInstallment(finalDate),
			CategoryID:              createPgUUIDInstallment(categoryID),
			CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
			CategoryName:            createPgTextInstallment("Category"),
			CategoryIcon:            createPgTextInstallment("icon"),
			CreditcardID:            pgtype.UUID{Valid: false},
			CreatedAt:               createTimestamptzInstallment(now),
			UpdatedAt:               createTimestamptzInstallment(now),
			TotalCount:              1,
		},
	}

	repo := Repository{store: mock}

	params := commonsmodels.PaginatedParams{
		UserID: userID,
		Limit:  10,
		Offset: 0,
	}

	result, _, err := repo.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result[0].Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}
}

// ============= READ BY ID TESTS =============

func TestReadInstallmentTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionRowResult = pgstore.GetInstallmentTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Test Installment",
		Value:                   utils.Float64ToNumeric(100.00),
		InitialDate:             createTimestamptzInstallment(initialDate),
		FinalDate:               createTimestamptzInstallment(finalDate),
		CategoryID:              createPgUUIDInstallment(categoryID),
		CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
		CategoryName:            createPgTextInstallment("Category"),
		CategoryIcon:            createPgTextInstallment("icon"),
		CreatedAt:               createTimestamptzInstallment(now),
		UpdatedAt:               createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadInstallmentTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Name != "Test Installment" {
		t.Errorf("Expected name 'Test Installment', got %v", result.Name)
	}
}

func TestReadInstallmentTransactionByIDWithCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionRowResult = pgstore.GetInstallmentTransactionByIDRow{
		ID:                         transactionID,
		UserID:                     userID,
		Name:                       "Test with CC",
		Value:                      utils.Float64ToNumeric(200.00),
		InitialDate:                createTimestamptzInstallment(initialDate),
		FinalDate:                  createTimestamptzInstallment(finalDate),
		CategoryID:                 createPgUUIDInstallment(categoryID),
		CategoryTransactionType:    createPgInt4Installment(int32(models.Credit)),
		CategoryName:               createPgTextInstallment("Credit Category"),
		CategoryIcon:               createPgTextInstallment("credit-icon"),
		CreditcardID:               createPgUUIDInstallment(creditCardID),
		CreditcardName:             createPgTextInstallment("My Card"),
		CreditcardFirstFourNumbers: createPgTextInstallment("1234"),
		CreditcardCreditLimit:      createPgFloat8Installment(5000.00),
		CreditcardCloseDay:         createPgInt4Installment(10),
		CreditcardExpireDay:        createPgInt4Installment(20),
		CreditcardBackgroundColor:  createPgTextInstallment("#000000"),
		CreditcardTextColor:        createPgTextInstallment("#FFFFFF"),
		CreatedAt:                  createTimestamptzInstallment(now),
		UpdatedAt:                  createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadInstallmentTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard == nil {
		t.Error("Expected Creditcard to be present")
	}

	if result.Creditcard.Name != "My Card" {
		t.Errorf("Expected creditcard name 'My Card', got %v", result.Creditcard.Name)
	}
}

func TestReadInstallmentTransactionByIDError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	_, err := repo.ReadInstallmentTransactionByID(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestReadInstallmentTransactionByIDWithoutCreditCard(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionRowResult = pgstore.GetInstallmentTransactionByIDRow{
		ID:                      transactionID,
		UserID:                  userID,
		Name:                    "Test without CC",
		Value:                   utils.Float64ToNumeric(150.00),
		InitialDate:             createTimestamptzInstallment(initialDate),
		FinalDate:               createTimestamptzInstallment(finalDate),
		CategoryID:              createPgUUIDInstallment(categoryID),
		CategoryTransactionType: createPgInt4Installment(int32(models.Debit)),
		CategoryName:            createPgTextInstallment("Category"),
		CategoryIcon:            createPgTextInstallment("icon"),
		CreditcardID:            pgtype.UUID{Valid: false},
		CreatedAt:               createTimestamptzInstallment(now),
		UpdatedAt:               createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadInstallmentTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Creditcard != nil {
		t.Error("Expected Creditcard to be nil")
	}
}

func TestReadShortInstallmentTransactionByIDSuccess(t *testing.T) {
	ctx := context.Background()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 5, 0)

	mock := newInstallmentStoreMock()
	mock.ShortInstallmentRowResult = pgstore.InstallmentTransaction{
		ID:          transactionID,
		Name:        "Short Installment",
		Value:       utils.Float64ToNumeric(199.90),
		InitialDate: createTimestamptzInstallment(initialDate),
		FinalDate:   createTimestamptzInstallment(finalDate),
		CreatedAt:   createTimestamptzInstallment(now),
		UpdatedAt:   createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	result, err := repo.ReadShortInstallmentTransactionByID(ctx, transactionID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, result.ID)
	}

	if result.Name != "Short Installment" {
		t.Errorf("Expected Name 'Short Installment', got %v", result.Name)
	}
}

func TestReadShortInstallmentTransactionByIDError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("not found")

	repo := Repository{store: mock}

	_, err := repo.ReadShortInstallmentTransactionByID(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= UPDATE TESTS =============

func TestUpdateInstallmentTransactionSuccess(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionResult = pgstore.InstallmentTransaction{
		ID:           transactionID,
		UserID:       userID,
		Name:         "Updated Installment",
		Value:        utils.Float64ToNumeric(150.00),
		InitialDate:  createTimestamptzInstallment(initialDate),
		FinalDate:    createTimestamptzInstallment(finalDate),
		CategoryID:   categoryID,
		CreditCardID: createPgUUIDInstallment(creditCardID),
		CreatedAt:    createTimestamptzInstallment(now),
		UpdatedAt:    createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	model := models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID: categoryID,
		},
		Creditcard: &models.CreditCard{
			ID: creditCardID,
		},
	}

	result, err := repo.UpdateInstallmentTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 150.00 {
		t.Errorf("Expected Value 150.00, got %v", result.Value)
	}
}

func TestUpdateInstallmentTransactionError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	model := models.InstallmentTransaction{
		ID:          uuid.New(),
		Name:        "Test",
		Value:       100.00,
		InitialDate: time.Now(),
		FinalDate:   time.Now().AddDate(0, 6, 0),
		Category: models.Category{
			ID: uuid.New(),
		},
	}

	_, err := repo.UpdateInstallmentTransaction(ctx, model)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUpdateInstallmentTransactionWithDifferentDates(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now.AddDate(0, 1, 0)
	finalDate := now.AddDate(0, 12, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionResult = pgstore.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Updated Installment",
		Value:       utils.Float64ToNumeric(250.00),
		InitialDate: createTimestamptzInstallment(initialDate),
		FinalDate:   createTimestamptzInstallment(finalDate),
		CategoryID:  categoryID,
		CreatedAt:   createTimestamptzInstallment(now),
		UpdatedAt:   createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	model := models.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Updated Installment",
		Value:       250.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: models.Category{
			ID: categoryID,
		},
	}

	result, err := repo.UpdateInstallmentTransaction(ctx, model)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 250.00 {
		t.Errorf("Expected Value 250.00, got %v", result.Value)
	}
}

// ============= DELETE TESTS =============

func TestDeleteInstallmentTransactionSuccess(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()

	repo := Repository{store: mock}

	err := repo.DeleteInstallmentTransaction(ctx, uuid.New())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteInstallmentTransactionError(t *testing.T) {
	ctx := context.Background()

	mock := newInstallmentStoreMock()
	mock.Error = errors.New("database error")

	repo := Repository{store: mock}

	err := repo.DeleteInstallmentTransaction(ctx, uuid.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// ============= EDGE CASES =============

func TestCreateInstallmentTransactionWithZeroValue(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := newInstallmentStoreMock()
	mock.InstallmentTransactionResult = pgstore.InstallmentTransaction{
		ID:          transactionID,
		UserID:      userID,
		Name:        "Zero Value Installment",
		Value:       utils.Float64ToNumeric(0),
		InitialDate: createTimestamptzInstallment(initialDate),
		FinalDate:   createTimestamptzInstallment(finalDate),
		CategoryID:  categoryID,
		CreatedAt:   createTimestamptzInstallment(now),
		UpdatedAt:   createTimestamptzInstallment(now),
	}

	repo := Repository{store: mock}

	request := models.CreateInstallmentTransaction{
		UserID:      userID,
		Name:        "Zero Value Installment",
		Value:       0,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}

	result, err := repo.CreateInstallmentTransaction(ctx, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Value != 0 {
		t.Errorf("Expected Value 0, got %v", result.Value)
	}
}
