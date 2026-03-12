package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend-commons/constants"
	"backend-commons/dtos"
	apierrors "github.com/TB-Systems/go-commons/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MonthlyReportServiceMock struct {
	response dtos.MonthlyReportResponse
	apiErr   apierrors.ApiError
}

func (m *MonthlyReportServiceMock) GenerateMonthlyReport(ctx context.Context, userID uuid.UUID, year int32, month int32) (dtos.MonthlyReportResponse, apierrors.ApiError) {
	return m.response, m.apiErr
}

func setupMonthlyReportRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setMonthlyReportUserIDContext(ctx *gin.Context, userID uuid.UUID) {
	ctx.Set(constants.UserID, userID)
}

func TestMonthlyReportHandlerGetMonthlyBalanceSuccess(t *testing.T) {
	userID := uuid.New()
	categoryID := uuid.New()

	mock := &MonthlyReportServiceMock{
		response: dtos.MonthlyReportResponse{
			TotalIncome: 1000,
			TotalDebit:  300,
			TotalCredit: 100,
			Balance:     600,
			Categories: []dtos.CategoriesSpendingResponse{{
				ID:              categoryID,
				Name:            "Food",
				Icon:            "restaurant",
				TransactionType: 1,
				Value:           250,
			}},
		},
	}

	handler := NewMonthlyReportHandler(mock)
	router := setupMonthlyReportRouter()
	router.GET("/monthly_report", func(c *gin.Context) {
		setMonthlyReportUserIDContext(c, userID)
		handler.GetMonthlyBalance()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly_report?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dtos.MonthlyReportResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Expected valid json response, got %v", err)
	}

	if response.Balance != 600 {
		t.Errorf("Expected balance 600, got %v", response.Balance)
	}
}

func TestMonthlyReportHandlerGetMonthlyBalanceNoUserID(t *testing.T) {
	mock := &MonthlyReportServiceMock{}
	handler := NewMonthlyReportHandler(mock)
	router := setupMonthlyReportRouter()
	router.GET("/monthly_report", handler.GetMonthlyBalance())

	req, _ := http.NewRequest(http.MethodGet, "/monthly_report?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMonthlyReportHandlerGetMonthlyBalanceInvalidQuery(t *testing.T) {
	userID := uuid.New()
	mock := &MonthlyReportServiceMock{}
	handler := NewMonthlyReportHandler(mock)
	router := setupMonthlyReportRouter()
	router.GET("/monthly_report", func(c *gin.Context) {
		setMonthlyReportUserIDContext(c, userID)
		handler.GetMonthlyBalance()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly_report?month=99&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMonthlyReportHandlerGetMonthlyBalanceServiceError(t *testing.T) {
	userID := uuid.New()
	mock := &MonthlyReportServiceMock{
		apiErr: apierrors.NewApiError(http.StatusInternalServerError, apierrors.InternalServerError("internal error")),
	}

	handler := NewMonthlyReportHandler(mock)
	router := setupMonthlyReportRouter()
	router.GET("/monthly_report", func(c *gin.Context) {
		setMonthlyReportUserIDContext(c, userID)
		handler.GetMonthlyBalance()(c)
	})

	req, _ := http.NewRequest(http.MethodGet, "/monthly_report?month=9&year=2025", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
