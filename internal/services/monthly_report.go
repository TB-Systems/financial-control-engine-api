package services

import (
	"context"
	"financialcontrol/internal/repositories"
	"net/http"
	"sync"

	"financialcontrol/internal/constants"

	"github.com/TB-Systems/financial-control-backend-commons/dtos"
	"github.com/TB-Systems/financial-control-backend-commons/models"
	"github.com/TB-Systems/financial-control-backend-commons/modelsdto"

	"github.com/TB-Systems/go-commons/errors"
	"github.com/google/uuid"
)

type MonthlyReport interface {
	GenerateMonthlyReport(ctx context.Context, userID uuid.UUID, year int32, month int32) (dtos.MonthlyReportResponse, errors.ApiError)
}

type monthlyReport struct {
	repository repositories.MonthlyReport
}

func NewMonthlyReportService(repository repositories.MonthlyReport) MonthlyReport {
	return &monthlyReport{
		repository: repository,
	}
}

func (m *monthlyReport) GenerateMonthlyReport(ctx context.Context, userID uuid.UUID, year int32, month int32) (dtos.MonthlyReportResponse, errors.ApiError) {
	var wg sync.WaitGroup

	var balance models.MonthlyReport
	var spendingCategories []models.CategoriesSpending
	var spendingCreditCards []models.CreditCardsSpending

	var balanceErr error
	var spendingCategoriesErr error
	var spendingCreditCardsErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
		balance, balanceErr = m.repository.GetMonthlyBalance(ctx, userID, year, month)
	}()

	go func() {
		defer wg.Done()
		spendingCategories, spendingCategoriesErr = m.repository.GetCategoriesSpending(ctx, userID, year, month)
	}()

	go func() {
		defer wg.Done()
		spendingCreditCards, spendingCreditCardsErr = m.repository.GetCreditCardsSpending(ctx, userID, year, month)
	}()

	wg.Wait()

	if balanceErr != nil {
		return dtos.MonthlyReportResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(balanceErr.Error()),
		)
	}

	if spendingCategoriesErr != nil && spendingCategoriesErr.Error() != constants.StoreErrorNoRowsMsg {
		return dtos.MonthlyReportResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(spendingCategoriesErr.Error()),
		)
	}

	if spendingCreditCardsErr != nil && spendingCreditCardsErr.Error() != constants.StoreErrorNoRowsMsg {
		return dtos.MonthlyReportResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(spendingCreditCardsErr.Error()),
		)
	}

	response := modelsdto.MonthlyReportResponseFromModels(balance, spendingCategories, spendingCreditCards)

	return response, nil
}
