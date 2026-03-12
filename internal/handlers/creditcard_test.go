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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============= MOCK IMPLEMENTATION =============

type CreditCardServiceMock struct {
	CreateResult dtos.CreditCardResponse
	CreateError  apierrors.ApiError
	ReadResult   commonsmodels.ResponseList[dtos.CreditCardResponse]
	ReadError    apierrors.ApiError
	ReadAtResult dtos.CreditCardResponse
	ReadAtError  apierrors.ApiError
	UpdateResult dtos.CreditCardResponse
	UpdateError  apierrors.ApiError
	DeleteError  apierrors.ApiError
}

func NewCreditCardServiceMock() *CreditCardServiceMock {
	return &CreditCardServiceMock{}
}

func (m *CreditCardServiceMock) Create(ctx context.Context, userID uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, apierrors.ApiError) {
	return m.CreateResult, m.CreateError
}

func (m *CreditCardServiceMock) Read(ctx context.Context, userID uuid.UUID) (commonsmodels.ResponseList[dtos.CreditCardResponse], apierrors.ApiError) {
	return m.ReadResult, m.ReadError
}

func (m *CreditCardServiceMock) ReadAt(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.CreditCardResponse, apierrors.ApiError) {
	return m.ReadAtResult, m.ReadAtError
}

func (m *CreditCardServiceMock) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, apierrors.ApiError) {
	return m.UpdateResult, m.UpdateError
}

func (m *CreditCardServiceMock) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) apierrors.ApiError {
	return m.DeleteError
}

// ============= HELPER FUNCTIONS =============

func setupCreditCardRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setCreditCardUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

// ============= CREATE TESTS =============

func TestCreditCardHandlerCreateSuccess(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardServiceMock()
	mock.CreateResult = dtos.CreditCardResponse{
		ID:               creditCardID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.POST("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/credit-cards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreditCardHandlerCreateNoUserID(t *testing.T) {
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.POST("/credit-cards", handler.Create())

	requestBody := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/credit-cards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreditCardHandlerCreateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.POST("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Create()(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/credit-cards", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreditCardHandlerCreateValidationError(t *testing.T) {
	userID := uuid.New()
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.POST("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "",
		FirstFourNumbers: "12",
		Limit:            0,
		CloseDay:         0,
		ExpireDay:        0,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/credit-cards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCreditCardHandlerCreateServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.CreateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.POST("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Create()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/credit-cards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ TESTS =============

func TestCreditCardHandlerReadSuccess(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	mock := NewCreditCardServiceMock()
	mock.ReadResult = commonsmodels.ResponseList[dtos.CreditCardResponse]{
		Items: []dtos.CreditCardResponse{
			{
				ID:               uuid.New(),
				Name:             "Card 1",
				FirstFourNumbers: "1111",
				Limit:            3000.0,
				CloseDay:         10,
				ExpireDay:        20,
				BackgroundColor:  "#FF0000",
				TextColor:        "#FFFFFF",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
		Total: 1,
	}

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/credit-cards", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreditCardHandlerReadNoUserID(t *testing.T) {
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards", handler.Read())

	req := httptest.NewRequest(http.MethodGet, "/credit-cards", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreditCardHandlerReadServiceError(t *testing.T) {
	userID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.ReadError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Read()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/credit-cards", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= READ AT TESTS =============

func TestCreditCardHandlerReadAtSuccess(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardServiceMock()
	mock.ReadAtResult = dtos.CreditCardResponse{
		ID:               creditCardID,
		Name:             "Test Card",
		FirstFourNumbers: "1234",
		Limit:            5000.0,
		CloseDay:         15,
		ExpireDay:        25,
		BackgroundColor:  "#000000",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.ReadAt()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreditCardHandlerReadAtNoUserID(t *testing.T) {
	creditCardID := uuid.New()
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards/:id", handler.ReadAt())

	req := httptest.NewRequest(http.MethodGet, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreditCardHandlerReadAtInvalidID(t *testing.T) {
	userID := uuid.New()
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.ReadAt()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/credit-cards/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreditCardHandlerReadAtNotFound(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.ReadAtError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("credit card not found"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.GET("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.ReadAt()(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// ============= UPDATE TESTS =============

func TestCreditCardHandlerUpdateSuccess(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()
	now := time.Now()

	mock := NewCreditCardServiceMock()
	mock.UpdateResult = dtos.CreditCardResponse{
		ID:               creditCardID,
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/"+creditCardID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreditCardHandlerUpdateNoUserID(t *testing.T) {
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", handler.Update())

	requestBody := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/"+creditCardID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreditCardHandlerUpdateInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/invalid-uuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreditCardHandlerUpdateInvalidJSON(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Update()(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/"+creditCardID.String(), bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreditCardHandlerUpdateValidationError(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "",
		FirstFourNumbers: "12",
		Limit:            0,
		CloseDay:         0,
		ExpireDay:        0,
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/"+creditCardID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCreditCardHandlerUpdateServiceError(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.UpdateError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.PUT("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Update()(c)
	})

	requestBody := dtos.CreditCardRequest{
		Name:             "Updated Card",
		FirstFourNumbers: "9999",
		Limit:            10000.0,
		CloseDay:         20,
		ExpireDay:        30,
		BackgroundColor:  "#0000FF",
		TextColor:        "#FFFFFF",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPut, "/credit-cards/"+creditCardID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= DELETE TESTS =============

func TestCreditCardHandlerDeleteSuccess(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.DELETE("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreditCardHandlerDeleteNoUserID(t *testing.T) {
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.DELETE("/credit-cards/:id", handler.Delete())

	req := httptest.NewRequest(http.MethodDelete, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreditCardHandlerDeleteInvalidID(t *testing.T) {
	userID := uuid.New()

	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.DELETE("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/credit-cards/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreditCardHandlerDeleteNotFound(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusNotFound, apierrors.NotFoundError("credit card not found"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.DELETE("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreditCardHandlerDeleteServiceError(t *testing.T) {
	userID := uuid.New()
	creditCardID := uuid.New()

	mock := NewCreditCardServiceMock()
	mock.DeleteError = apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("service error"))

	handler := NewCreditCardsHandler(mock)
	router := setupCreditCardRouter()
	router.DELETE("/credit-cards/:id", func(c *gin.Context) {
		setCreditCardUserIDContext(c, userID)
		handler.Delete()(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/credit-cards/"+creditCardID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ============= CONSTRUCTOR TEST =============

func TestNewCreditCardsHandler(t *testing.T) {
	mock := NewCreditCardServiceMock()
	handler := NewCreditCardsHandler(mock)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}
