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

type CategoryServiceMock struct {
	CreateResult   dtos.CategoryResponse
	CreateError    apierrors.ApiError
	ReadResult     commonsmodels.ResponseList[dtos.CategoryResponse]
	ReadError      apierrors.ApiError
	ReadByIdResult dtos.CategoryResponse
	ReadByIdError  apierrors.ApiError
	UpdateResult   dtos.CategoryResponse
	UpdateError    apierrors.ApiError
	DeleteError    apierrors.ApiError
}

func NewCategoryServiceMock() *CategoryServiceMock {
	return &CategoryServiceMock{}
}

func (m *CategoryServiceMock) Create(ctx context.Context, userID uuid.UUID, request dtos.CategoryRequest) (dtos.CategoryResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *CategoryServiceMock) Read(ctx context.Context, userID uuid.UUID) (commonsmodels.ResponseList[dtos.CategoryResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *CategoryServiceMock) ReadByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.CategoryResponse, apierrors.ApiError) {
	return m.ReadByIdResult, m.ReadByIdError
}

func (m *CategoryServiceMock) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.CategoryRequest) (dtos.CategoryResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *CategoryServiceMock) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

// ============= HELPER FUNCTIONS =============

func setupCategoryRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setCategoryUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestCategoryHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	transactionType := models.Debit

	mock := NewCategoryServiceMock()
	mock.CreateResult = dtos.CategoryResponse{
		ID:              categoryID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.POST("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCategoryHandlerCreateNoUserID(t *testing.T) {
	transactionType := models.Debit
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.POST("/categories", handler.Create())

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCategoryHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.POST("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Create()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCategoryHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.POST("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: nil,
		Name:            "",
		Icon:            "",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCategoryHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()
	transactionType := models.Debit

	mock := NewCategoryServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.POST("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Test Category",
		Icon:            "icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestCategoryHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewCategoryServiceMock()
	mock.ReadResult = commonsmodels.ResponseList[dtos.CategoryResponse]{
		Items: []dtos.CategoryResponse{
			{
				ID:              uuid.New(),
				TransactionType: models.Debit,
				Name:            "Category 1",
				Icon:            "icon1.png",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
		Total: 1,
	}

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryHandlerReadNoUserID(t *testing.T) {
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories", handler.Read())

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCategoryHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewCategoryServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ BY ID TESTS =============

func TestCategoryHandlerReadByIdSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()

	mock := NewCategoryServiceMock()
	mock.ReadByIdResult = dtos.CategoryResponse{
		ID:              categoryID,
		TransactionType: models.Debit,
		Name:            "Test Category",
		Icon:            "icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.ReadByID()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryHandlerReadByIdNoUserID(t *testing.T) {
	categoryID := uuid.New()
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories/:id", handler.ReadByID())

	req := httptest.NewRequest(http.MethodGet, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCategoryHandlerReadByIdInvalidID(t *testing.T) {
	userID := uuid.New()
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.ReadByID()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/categories/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCategoryHandlerReadByIdNotFound(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	mock.ReadByIdError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("category not found"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.GET("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.ReadByID()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestCategoryHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	now := time.Now()
	transactionType := models.Income

	mock := NewCategoryServiceMock()
	mock.UpdateResult = dtos.CategoryResponse{
		ID:              categoryID,
		TransactionType: models.Income,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryHandlerUpdateNoUserID(t *testing.T) {
	categoryID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", handler.Update())

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCategoryHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCategoryHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Update()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.String(), bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCategoryHandlerUpdateValidationError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: nil,
		Name:            "",
		Icon:            "",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCategoryHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()
	transactionType := models.Income

	mock := NewCategoryServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.PUT("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CategoryRequest{
		TransactionType: &transactionType,
		Name:            "Updated Category",
		Icon:            "new_icon.png",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/categories/"+categoryID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestCategoryHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCategoryHandlerDeleteNoUserID(t *testing.T) {
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.DELETE("/categories/:id", handler.Delete())

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCategoryHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/categories/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCategoryHandlerDeleteNotFound(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("category not found"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCategoryHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := NewCategoryServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCategoriesHandler(mock)
	router := setupCategoryRouter()
	router.DELETE("/categories/:id", func(c *gin.Context) {
		setCategoryUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+categoryID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= CONSTRUCTOR TEST =============

func TestNewCategoriesHandler(t *testing.T) {
	mock := NewCategoryServiceMock()
	handler := NewCategoriesHandler(mock)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}
