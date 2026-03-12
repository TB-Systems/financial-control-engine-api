package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"backend-commons/constants"
	"backend-commons/dtos"
	apierrors "github.com/TB-Systems/go-commons/errors"
	"backend-commons/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============= MOCK IMPLEMENTATION =============

type TransactionServiceMock struct {
	CreateResult                dtos.TransactionResponse
	CreateError                 apierrors.ApiError
	CreateFromMonthlyResult     dtos.TransactionResponse
	CreateFromMonthlyError      apierrors.ApiError
	CreateFromAnnualResult      dtos.TransactionResponse
	CreateFromAnnualError       apierrors.ApiError
	CreateFromInstallmentResult dtos.TransactionResponse
	CreateFromInstallmentError  apierrors.ApiError
	ReadResult                  commonsmodels.PaginatedResponse[dtos.TransactionResponse]
	ReadError                   apierrors.ApiError
	// ============= READ BY MONTH/YEAR TESTS =============
	ReadInToDatesResult      commonsmodels.PaginatedResponse[dtos.TransactionResponse]
	ReadInToDatesError       apierrors.ApiError
	ReadAtMonthAndYearResult commonsmodels.PaginatedResponse[dtos.TransactionResponse]
	ReadAtMonthAndYearError  apierrors.ApiError
	ReadByIdResult           dtos.TransactionResponse
	ReadByIdError            apierrors.ApiError
	UpdateResult             dtos.TransactionResponse
	UpdateError              apierrors.ApiError
	DeleteError              apierrors.ApiError
	PayError                 apierrors.ApiError
}

func NewTransactionServiceMock() *TransactionServiceMock {
	return &TransactionServiceMock{}
}

func (m *TransactionServiceMock) Create(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *TransactionServiceMock) CreateFromMonthlyTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.CreateFromMonthlyResult, m.CreateFromMonthlyError
}

func (m *TransactionServiceMock) CreateFromAnnualTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.CreateFromAnnualResult, m.CreateFromAnnualError
}

func (m *TransactionServiceMock) CreateFromInstallmentTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.CreateFromInstallmentResult, m.CreateFromInstallmentError
}

func (m *TransactionServiceMock) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *TransactionServiceMock) ReadInToDates(ctx context.Context, params commonsmodels.PaginatedParamsWithDateRange) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], apierrors.ApiError) {
	return m.ReadInToDatesResult, m.ReadInToDatesError
}

func (m *TransactionServiceMock) ReadAtMonthAndYear(ctx context.Context, params commonsmodels.PaginatedParamsWithMonthYear) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], apierrors.ApiError) {
	return m.ReadAtMonthAndYearResult, m.ReadAtMonthAndYearError
}

func (m *TransactionServiceMock) ReadById(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.ReadByIdResult, m.ReadByIdError
}

func (m *TransactionServiceMock) Update(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *TransactionServiceMock) Delete(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

func (m *TransactionServiceMock) Pay(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) apierrors.ApiError {
	return m.PayError
}

// ============= HELPER FUNCTIONS =============

func setupTransactionRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setTransactionUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestTransactionHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.CreateResult = dtos.TransactionResponse{
		ID:    transactionID,
		Name:  "Test Transaction",
		Date:  now,
		Value: 100.00,
		Paid:  true,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Category",
			Icon:            "icon.png",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       now,
		Value:      100.00,
		Paid:       true,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestTransactionHandlerCreateNoUserID(t *testing.T) {
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions", handler.Create())

	requestBody := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       now,
		Value:      100.00,
		Paid:       true,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Create()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:  "",
		Value: -1,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestTransactionHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:       "Test Transaction",
		Date:       now,
		Value:      100.00,
		Paid:       true,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionHandlerCreateFromMonthlySuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	recurrentID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.CreateFromMonthlyResult = dtos.TransactionResponse{ID: transactionID, Name: "Internet", Date: now, Value: 100, Paid: false, CreatedAt: now, UpdatedAt: now}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromMonthlyTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/monthly", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestTransactionHandlerCreateFromMonthlyNoUserID(t *testing.T) {
	recurrentID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/monthly", handler.CreateFromMonthlyTransaction())

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/monthly", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerCreateFromMonthlyInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromMonthlyTransaction()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/monthly", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerCreateFromMonthlyServiceError(t *testing.T) {
	userID := uuid.New()
	recurrentID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.CreateFromMonthlyError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromMonthlyTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/monthly", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionHandlerCreateFromAnnualSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	recurrentID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.CreateFromAnnualResult = dtos.TransactionResponse{ID: transactionID, Name: "IPVA", Date: now, Value: 500, Paid: false, CreatedAt: now, UpdatedAt: now}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/annual", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromAnnualTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/annual", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestTransactionHandlerCreateFromAnnualNoUserID(t *testing.T) {
	recurrentID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/annual", handler.CreateFromAnnualTransaction())

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/annual", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerCreateFromAnnualInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/annual", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromAnnualTransaction()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/annual", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerCreateFromAnnualServiceError(t *testing.T) {
	userID := uuid.New()
	recurrentID := uuid.New()
	mock := NewTransactionServiceMock()
	mock.CreateFromAnnualError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/annual", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromAnnualTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/annual", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionHandlerCreateFromInstallmentSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	recurrentID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.CreateFromInstallmentResult = dtos.TransactionResponse{ID: transactionID, Name: "Notebook 2/12", Date: now, Value: 350, Paid: false, CreatedAt: now, UpdatedAt: now}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/installment", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromInstallmentTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/installment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestTransactionHandlerCreateFromInstallmentNoUserID(t *testing.T) {
	recurrentID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/installment", handler.CreateFromInstallmentTransaction())

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/installment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerCreateFromInstallmentInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/installment", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromInstallmentTransaction()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/installment", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerCreateFromInstallmentServiceError(t *testing.T) {
	userID := uuid.New()
	recurrentID := uuid.New()
	mock := NewTransactionServiceMock()
	mock.CreateFromInstallmentError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.POST("/transactions/recurrent/installment", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.CreateFromInstallmentTransaction()(c)
	})

	requestBody := dtos.TransactionRequestFromRecurrentTransaction{ID: recurrentID}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/transactions/recurrent/installment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestTransactionHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items: []dtos.TransactionResponse{
			{
				ID:    uuid.New(),
				Name:  "Transaction 1",
				Date:  now,
				Value: 100.00,
				Paid:  true,
				Category: dtos.CategoryResponse{
					ID:              categoryID,
					TransactionType: models.Debit,
					Name:            "Category",
					Icon:            "icon.png",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Page:      1,
		PageCount: 1,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerReadNoUserID(t *testing.T) {
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", handler.Read())

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerReadInvalidPage(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions?page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionHandlerReadWithDateRange(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.ReadInToDatesResult = commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items: []dtos.TransactionResponse{
			{
				ID:    uuid.New(),
				Name:  "Transaction 1",
				Date:  now,
				Value: 100.00,
				Paid:  true,
				Category: dtos.CategoryResponse{
					ID:              categoryID,
					TransactionType: models.Debit,
					Name:            "Category",
					Icon:            "icon.png",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Page:      1,
		PageCount: 1,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	startDate := now.AddDate(0, -1, 0).Format(time.DateOnly)
	endDate := now.Format(time.DateOnly)
	req := httptest.NewRequest(http.MethodGet, "/transactions?start_date="+startDate+"&end_date="+endDate, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerReadWithDateRangeError(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.ReadInToDatesError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	now := time.Now()
	startDate := now.AddDate(0, -1, 0).Format(time.DateOnly)
	endDate := now.Format(time.DateOnly)
	req := httptest.NewRequest(http.MethodGet, "/transactions?start_date="+startDate+"&end_date="+endDate, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestTransactionHandlerReadWithInvalidDateRange(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions?start_date=invalid&end_date=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ============= READ BY ID TESTS =============

func TestTransactionHandlerReadByIdSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.ReadByIdResult = dtos.TransactionResponse{
		ID:    transactionID,
		Name:  "Test Transaction",
		Date:  now,
		Value: 100.00,
		Paid:  true,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Category",
			Icon:            "icon.png",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerReadByIdNoUserID(t *testing.T) {
	transactionID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/:id", handler.ReadById())

	req := httptest.NewRequest(http.MethodGet, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerReadByIdInvalidID(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerReadByIdNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.ReadByIdError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("transaction not found"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestTransactionHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.UpdateResult = dtos.TransactionResponse{
		ID:    transactionID,
		Name:  "Updated Transaction",
		Date:  now,
		Value: 200.00,
		Paid:  false,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Income,
			Name:            "Category",
			Icon:            "icon.png",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       now,
		Value:      200.00,
		Paid:       false,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerUpdateNoUserID(t *testing.T) {
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", handler.Update())

	requestBody := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       now,
		Value:      200.00,
		Paid:       false,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       now,
		Value:      200.00,
		Paid:       false,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Update()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String(), bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerUpdateValidationError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:  "",
		Value: -1,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestTransactionHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.TransactionRequest{
		Name:       "Updated Transaction",
		Date:       now,
		Value:      200.00,
		Paid:       false,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestTransactionHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.DELETE("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerDeleteNoUserID(t *testing.T) {
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.DELETE("/transactions/:id", handler.Delete())

	req := httptest.NewRequest(http.MethodDelete, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.DELETE("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerDeleteNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("transaction not found"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.DELETE("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTransactionHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.DELETE("/transactions/:id", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= PAY TESTS =============

func TestTransactionHandlerPaySuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id/pay", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Pay()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String()+"/pay", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerPayNoUserID(t *testing.T) {
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id/pay", handler.Pay())

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String()+"/pay", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerPayInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id/pay", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Pay()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/transactions/invalid-uuid/pay", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerPayServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.PayError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.PUT("/transactions/:id/pay", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.Pay()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/transactions/"+transactionID.String()+"/pay", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= CONSTRUCTOR TEST =============

func TestNewTransactionsHandler(t *testing.T) {
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

// ============= READ BY MONTH/YEAR TESTS =============

func TestTransactionHandlerReadByMonthAndYearSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewTransactionServiceMock()
	mock.ReadAtMonthAndYearResult = commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items: []dtos.TransactionResponse{
			{
				ID:    uuid.New(),
				Name:  "Transaction month/year",
				Date:  now,
				Value: 180.00,
				Paid:  true,
				Category: dtos.CategoryResponse{
					ID:              categoryID,
					TransactionType: models.Debit,
					Name:            "Category",
					Icon:            "icon.png",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Page:      1,
		PageCount: 1,
	}

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadByMonthAndYear()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/monthly?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransactionHandlerReadByMonthAndYearNoUserID(t *testing.T) {
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/monthly", handler.ReadByMonthAndYear())

	req := httptest.NewRequest(http.MethodGet, "/transactions/monthly?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestTransactionHandlerReadByMonthAndYearInvalidMonth(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadByMonthAndYear()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/monthly?month=13&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerReadByMonthAndYearInvalidPage(t *testing.T) {
	userID := uuid.New()
	mock := NewTransactionServiceMock()
	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadByMonthAndYear()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/monthly?month=9&year=2025&page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTransactionHandlerReadByMonthAndYearServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewTransactionServiceMock()
	mock.ReadAtMonthAndYearError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewTransactionsHandler(mock)
	router := setupTransactionRouter()
	router.GET("/transactions/monthly", func(c *gin.Context) {
		setTransactionUserIDContext(c, userID)
		handler.ReadByMonthAndYear()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/transactions/monthly?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
