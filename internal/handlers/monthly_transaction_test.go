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

type MonthlyTransactionServiceMock struct {
	CreateResult   dtos.MonthlyTransactionResponse
	CreateError    apierrors.ApiError
	ReadResult     commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]
	ReadError      apierrors.ApiError
	ReadByIdResult dtos.MonthlyTransactionResponse
	ReadByIdError  apierrors.ApiError
	UpdateResult   dtos.MonthlyTransactionResponse
	UpdateError    apierrors.ApiError
	DeleteError    apierrors.ApiError
}

func NewMonthlyTransactionServiceMock() *MonthlyTransactionServiceMock {
	return &MonthlyTransactionServiceMock{}
}

func (m *MonthlyTransactionServiceMock) Create(ctx context.Context, userID uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *MonthlyTransactionServiceMock) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *MonthlyTransactionServiceMock) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.MonthlyTransactionResponse, apierrors.ApiError) {
	return m.ReadByIdResult, m.ReadByIdError
}

func (m *MonthlyTransactionServiceMock) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *MonthlyTransactionServiceMock) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

// ============= HELPER FUNCTIONS =============

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestMonthlyTransactionHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionServiceMock()
	mock.CreateResult = dtos.MonthlyTransactionResponse{
		ID:    transactionID,
		Value: 100.00,
		Day:   15,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.POST("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/monthly-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response dtos.MonthlyTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestMonthlyTransactionHandlerCreateNoUserID(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.POST("/monthly-transactions", handler.Create())

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/monthly-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyTransactionHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.POST("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Create()(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/monthly-transactions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.POST("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/monthly-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.POST("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Test Monthly",
		Value:      100.00,
		Day:        15,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/monthly-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestMonthlyTransactionHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]{
		Items: []dtos.MonthlyTransactionResponse{
			{
				ID:    transactionID,
				Value: 100.00,
				Day:   5,
				Category: dtos.CategoryResponse{
					ID:              uuid.New(),
					TransactionType: models.Debit,
					Name:            "Category",
					Icon:            "icon",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Page:      1,
		PageCount: 1,
	}

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(response.Items))
	}
}

func TestMonthlyTransactionHandlerReadNoUserID(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions", handler.Read())

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyTransactionHandlerReadInvalidPage(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions?page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMonthlyTransactionHandlerReadWithPagination(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]{
		Items: []dtos.MonthlyTransactionResponse{
			{
				ID:    uuid.New(),
				Value: 100.00,
				Day:   5,
				Category: dtos.CategoryResponse{
					ID:              uuid.New(),
					TransactionType: models.Debit,
					Name:            "Category",
					Icon:            "icon",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Page:      2,
		PageCount: 5,
	}

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// ============= READ BY ID TESTS =============

func TestMonthlyTransactionHandlerReadByIdSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionServiceMock()
	mock.ReadByIdResult = dtos.MonthlyTransactionResponse{
		ID:    transactionID,
		Value: 150.00,
		Day:   15,
		Category: dtos.CategoryResponse{
			ID:              uuid.New(),
			TransactionType: models.Debit,
			Name:            "Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.MonthlyTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestMonthlyTransactionHandlerReadByIdNoUserID(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions/:id", handler.ReadById())

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyTransactionHandlerReadByIdInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerReadByIdNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.ReadByIdError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("monthly transaction"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.GET("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestMonthlyTransactionHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewMonthlyTransactionServiceMock()
	mock.UpdateResult = dtos.MonthlyTransactionResponse{
		ID:    transactionID,
		Value: 200.00,
		Day:   20,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Updated Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.MonthlyTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Value != 200.00 {
		t.Errorf("Expected Value 200.00, got %v", response.Value)
	}
}

func TestMonthlyTransactionHandlerUpdateNoUserID(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", handler.Update())

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/"+uuid.New().String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyTransactionHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Update()(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/"+transactionID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerUpdateValidationError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("monthly transaction"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.PUT("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.MonthlyTransactionRequest{
		Name:       "Updated Monthly",
		Value:      200.00,
		Day:        20,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/monthly-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestMonthlyTransactionHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.DELETE("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/monthly-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestMonthlyTransactionHandlerDeleteNoUserID(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.DELETE("/monthly-transactions/:id", handler.Delete())

	req, _ := http.NewRequest(http.MethodDelete, "/monthly-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyTransactionHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.DELETE("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/monthly-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyTransactionHandlerDeleteNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("monthly transaction"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.DELETE("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/monthly-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestMonthlyTransactionHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewMonthlyTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewMonthlyTransactionsHandler(mock)
	router := setupRouter()
	router.DELETE("/monthly-transactions/:id", func(c *gin.Context) {
		setUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/monthly-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= CONSTRUCTOR TEST =============

func TestNewMonthlyTransactionsHandler(t *testing.T) {
	mock := NewMonthlyTransactionServiceMock()
	handler := NewMonthlyTransactionsHandler(mock)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}
