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

type AnnualTransactionServiceMock struct {
	CreateResult   dtos.AnnualTransactionResponse
	CreateError    apierrors.ApiError
	ReadResult     commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]
	ReadError      apierrors.ApiError
	ReadByIdResult dtos.AnnualTransactionResponse
	ReadByIdError  apierrors.ApiError
	UpdateResult   dtos.AnnualTransactionResponse
	UpdateError    apierrors.ApiError
	DeleteError    apierrors.ApiError
}

func NewAnnualTransactionServiceMock() *AnnualTransactionServiceMock {
	return &AnnualTransactionServiceMock{}
}

func (m *AnnualTransactionServiceMock) Create(ctx context.Context, userID uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *AnnualTransactionServiceMock) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *AnnualTransactionServiceMock) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.AnnualTransactionResponse, apierrors.ApiError) {
	return m.ReadByIdResult, m.ReadByIdError
}

func (m *AnnualTransactionServiceMock) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *AnnualTransactionServiceMock) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

// ============= HELPER FUNCTIONS =============

func setupAnnualRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setAnnualUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestAnnualTransactionHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionServiceMock()
	mock.CreateResult = dtos.AnnualTransactionResponse{
		ID:    transactionID,
		Value: 100.00,
		Day:   15,
		Month: 6,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.POST("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/annual-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response dtos.AnnualTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, response.ID)
	}

	if response.Month != 6 {
		t.Errorf("Expected Month 6, got %v", response.Month)
	}
}

func TestAnnualTransactionHandlerCreateNoUserID(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.POST("/annual-transactions", handler.Create())

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/annual-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAnnualTransactionHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.POST("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Create()(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/annual-transactions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.POST("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/annual-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.POST("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Test Annual",
		Value:      100.00,
		Day:        15,
		Month:      6,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/annual-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestAnnualTransactionHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]{
		Items: []dtos.AnnualTransactionResponse{
			{
				ID:    transactionID,
				Value: 100.00,
				Day:   5,
				Month: 3,
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

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(response.Items))
	}
}

func TestAnnualTransactionHandlerReadNoUserID(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions", handler.Read())

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAnnualTransactionHandlerReadInvalidPage(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions?page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestAnnualTransactionHandlerReadWithPagination(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]{
		Items: []dtos.AnnualTransactionResponse{
			{
				ID:    uuid.New(),
				Value: 100.00,
				Day:   5,
				Month: 7,
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

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// ============= READ BY ID TESTS =============

func TestAnnualTransactionHandlerReadByIdSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionServiceMock()
	mock.ReadByIdResult = dtos.AnnualTransactionResponse{
		ID:    transactionID,
		Value: 150.00,
		Day:   15,
		Month: 9,
		Category: dtos.CategoryResponse{
			ID:              uuid.New(),
			TransactionType: models.Debit,
			Name:            "Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.AnnualTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, response.ID)
	}

	if response.Month != 9 {
		t.Errorf("Expected Month 9, got %v", response.Month)
	}
}

func TestAnnualTransactionHandlerReadByIdNoUserID(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions/:id", handler.ReadById())

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAnnualTransactionHandlerReadByIdInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerReadByIdNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.ReadByIdError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("annual transaction"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.GET("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/annual-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestAnnualTransactionHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewAnnualTransactionServiceMock()
	mock.UpdateResult = dtos.AnnualTransactionResponse{
		ID:    transactionID,
		Value: 200.00,
		Day:   20,
		Month: 10,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Updated Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.AnnualTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.Value != 200.00 {
		t.Errorf("Expected Value 200.00, got %v", response.Value)
	}

	if response.Month != 10 {
		t.Errorf("Expected Month 10, got %v", response.Month)
	}
}

func TestAnnualTransactionHandlerUpdateNoUserID(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", handler.Update())

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/"+uuid.New().String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAnnualTransactionHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
		CategoryID: uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Update()(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/"+transactionID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerUpdateValidationError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("annual transaction"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.PUT("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.AnnualTransactionRequest{
		Name:       "Updated Annual",
		Value:      200.00,
		Day:        20,
		Month:      10,
		CategoryID: categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/annual-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestAnnualTransactionHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.DELETE("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/annual-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestAnnualTransactionHandlerDeleteNoUserID(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.DELETE("/annual-transactions/:id", handler.Delete())

	req, _ := http.NewRequest(http.MethodDelete, "/annual-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAnnualTransactionHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.DELETE("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/annual-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAnnualTransactionHandlerDeleteNotFound(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("annual transaction"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.DELETE("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/annual-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestAnnualTransactionHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewAnnualTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("database error"))

	handler := NewAnnualTransactionsHandler(mock)
	router := setupAnnualRouter()
	router.DELETE("/annual-transactions/:id", func(c *gin.Context) {
		setAnnualUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/annual-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= CONSTRUCTOR TEST =============

func TestNewAnnualTransactionsHandler(t *testing.T) {
	mock := NewAnnualTransactionServiceMock()
	handler := NewAnnualTransactionsHandler(mock)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}
