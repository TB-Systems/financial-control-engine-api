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

type InstallmentTransactionServiceMock struct {
	CreateResult   dtos.InstallmentTransactionResponse
	CreateError    apierrors.ApiError
	ReadResult     commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse]
	ReadError      apierrors.ApiError
	ReadByIdResult dtos.InstallmentTransactionResponse
	ReadByIdError  apierrors.ApiError
	UpdateResult   dtos.InstallmentTransactionResponse
	UpdateError    apierrors.ApiError
	DeleteError    apierrors.ApiError
}

func NewInstallmentTransactionServiceMock() *InstallmentTransactionServiceMock {
	return &InstallmentTransactionServiceMock{}
}

func (m *InstallmentTransactionServiceMock) Create(ctx context.Context, userID uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *InstallmentTransactionServiceMock) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *InstallmentTransactionServiceMock) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.InstallmentTransactionResponse, apierrors.ApiError) {
	return m.ReadByIdResult, m.ReadByIdError
}

func (m *InstallmentTransactionServiceMock) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *InstallmentTransactionServiceMock) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

// ============= HELPER FUNCTIONS =============

func setupInstallmentRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setInstallmentUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestInstallmentTransactionHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionServiceMock()
	mock.CreateResult = dtos.InstallmentTransactionResponse{
		ID:          transactionID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.POST("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/installment-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response dtos.InstallmentTransactionResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	if response.ID != transactionID {
		t.Errorf("Expected ID %v, got %v", transactionID, response.ID)
	}
}

func TestInstallmentTransactionHandlerCreateNoUserID(t *testing.T) {
	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.POST("/installment-transactions", handler.Create())

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: time.Now(),
		FinalDate:   time.Now().AddDate(0, 6, 0),
		CategoryID:  uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/installment-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInstallmentTransactionHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.POST("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Create()(c)
	})

	req, _ := http.NewRequest(http.MethodPost, "/installment-transactions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.POST("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := map[string]interface{}{
		"name": "",
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/installment-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("error"))

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.POST("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/installment-transactions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestInstallmentTransactionHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionServiceMock()
	mock.ReadResult = commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse]{
		Items: []dtos.InstallmentTransactionResponse{
			{
				ID:          transactionID,
				Name:        "Test Installment",
				Value:       100.00,
				InitialDate: initialDate,
				FinalDate:   finalDate,
				Category: dtos.CategoryResponse{
					ID:              categoryID,
					TransactionType: models.Debit,
					Name:            "Test Category",
					Icon:            "icon",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		PageCount: 1,
		Page:      1,
	}

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadNoUserID(t *testing.T) {
	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions", handler.Read())

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadInvalidPage(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions?page=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("error"))

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Read()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions?page=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ BY ID TESTS =============

func TestInstallmentTransactionHandlerReadByIdSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionServiceMock()
	mock.ReadByIdResult = dtos.InstallmentTransactionResponse{
		ID:          transactionID,
		Name:        "Test Installment",
		Value:       100.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadByIdNoUserID(t *testing.T) {
	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions/:id", handler.ReadById())

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadByIdInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerReadByIdServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	mock.ReadByIdError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("not found"))

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.GET("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.ReadById()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/installment-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestInstallmentTransactionHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	initialDate := now
	finalDate := now.AddDate(0, 6, 0)

	mock := NewInstallmentTransactionServiceMock()
	mock.UpdateResult = dtos.InstallmentTransactionResponse{
		ID:          transactionID,
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		Category: dtos.CategoryResponse{
			ID:              categoryID,
			TransactionType: models.Debit,
			Name:            "Test Category",
			Icon:            "icon",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.PUT("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		CategoryID:  categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/installment-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestInstallmentTransactionHandlerUpdateNoUserID(t *testing.T) {
	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.PUT("/installment-transactions/:id", handler.Update())

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: time.Now(),
		FinalDate:   time.Now().AddDate(0, 6, 0),
		CategoryID:  uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/installment-transactions/"+uuid.New().String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInstallmentTransactionHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.PUT("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: time.Now(),
		FinalDate:   time.Now().AddDate(0, 6, 0),
		CategoryID:  uuid.New(),
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/installment-transactions/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.PUT("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Update()(c)
	})

	req, _ := http.NewRequest(http.MethodPut, "/installment-transactions/"+transactionID.String(), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewInstallmentTransactionServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("not found"))

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.PUT("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.InstallmentTransactionRequest{
		Name:        "Updated Installment",
		Value:       150.00,
		InitialDate: now,
		FinalDate:   now.AddDate(0, 6, 0),
		CategoryID:  categoryID,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPut, "/installment-transactions/"+transactionID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestInstallmentTransactionHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.DELETE("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/installment-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestInstallmentTransactionHandlerDeleteNoUserID(t *testing.T) {
	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.DELETE("/installment-transactions/:id", handler.Delete())

	req, _ := http.NewRequest(http.MethodDelete, "/installment-transactions/"+uuid.New().String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInstallmentTransactionHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.DELETE("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/installment-transactions/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestInstallmentTransactionHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()

	mock := NewInstallmentTransactionServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("not found"))

	handler := NewInstallmentTransactionsHandler(mock)
	router := setupInstallmentRouter()
	router.DELETE("/installment-transactions/:id", func(c *gin.Context) {
		setInstallmentUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req, _ := http.NewRequest(http.MethodDelete, "/installment-transactions/"+transactionID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
