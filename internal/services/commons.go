package services

import (
	"context"
	"backend-commons/constants"
	"backend-commons/dtos"
	"backend-commons/models"
	"backend-commons/modelsdto"
	"financialcontrol/internal/repositories"
	"net/http"

	"github.com/TB-Systems/go-commons/errors"
	"github.com/google/uuid"
)

func getRelations(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, creditCardID *uuid.UUID, categoryID uuid.UUID) (dtos.TransactionRelations, errors.ApiError) {
	category, apiErr := readCategory(r, ctx, userID, categoryID)
	categoryResponse := modelsdto.CategoryResponseFromModel(category)

	if apiErr != nil {
		return dtos.TransactionRelations{}, apiErr
	}

	var creditcardResponse *dtos.CreditCardResponse
	var creditcard *models.CreditCard
	if creditCardID != nil {
		creditcard, apiErr = readCreditcard(r, ctx, userID, *creditCardID)

		if apiErr != nil {
			return dtos.TransactionRelations{}, apiErr
		}

		resp := modelsdto.CreditCardResponseFromCreditCard(*creditcard)
		creditcardResponse = &resp
	}

	if creditCardID == nil && category.TransactionType == models.Credit {
		return dtos.TransactionRelations{}, errors.NewApiError(http.StatusBadRequest, errors.BadRequestError(constants.CreditWithoutCreditcardMsg))
	}

	if creditCardID != nil && category.TransactionType != models.Credit {
		return dtos.TransactionRelations{}, errors.NewApiError(http.StatusBadRequest, errors.BadRequestError(constants.DebitOrIncomeWithCreditcardMsg))
	}

	return dtos.TransactionRelations{
		CategoryModel:      category,
		CreditcardModel:    creditcard,
		CategoryResponse:   categoryResponse,
		CreditcardResponse: creditcardResponse,
	}, nil
}

func readCategory(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, id uuid.UUID) (models.Category, errors.ApiError) {
	category, err := r.ReadCategoryByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.Category{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CategoryNotFoundMsg))
		}
		return models.Category{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if category.UserID != userID {
		return models.Category{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CategoryNotFoundMsg))
	}

	return category, nil
}

func readCreditcard(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, id uuid.UUID) (*models.CreditCard, errors.ApiError) {
	creditcard, err := r.ReadCreditCardByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return &models.CreditCard{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CreditcardNotFoundMsg))
		}
		return &models.CreditCard{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if creditcard.UserID != userID {
		return &models.CreditCard{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CreditcardNotFoundMsg))
	}

	return &creditcard, nil
}

func readAnnualTransaction(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, id uuid.UUID) (*models.ShortAnnualTransaction, errors.ApiError) {
	annualTransaction, err := r.ReadShortAnnualTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return &models.ShortAnnualTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.AnnualTransactionNotFoundMsg))
		}
		return &models.ShortAnnualTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if annualTransaction.UserID != userID {
		return &models.ShortAnnualTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.AnnualTransactionNotFoundMsg))
	}

	return &annualTransaction, nil
}

func readMonthlyTransaction(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, id uuid.UUID) (*models.ShortMonthlyTransaction, errors.ApiError) {
	monthlyTransaction, err := r.ReadShortMonthlyTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return &models.ShortMonthlyTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.MonthlyTransactionNotFoundMsg))
		}
		return &models.ShortMonthlyTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if monthlyTransaction.UserID != userID {
		return &models.ShortMonthlyTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.MonthlyTransactionNotFoundMsg))
	}

	return &monthlyTransaction, nil
}

func readInstallmentTransaction(r repositories.CommonRepository, ctx context.Context, userID uuid.UUID, id uuid.UUID) (*models.ShortInstallmentTransaction, errors.ApiError) {
	installmentTransaction, err := r.ReadShortInstallmentTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return &models.ShortInstallmentTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.InstallmentTransactionNotFoundMsg))
		}
		return &models.ShortInstallmentTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if installmentTransaction.UserID != userID {
		return &models.ShortInstallmentTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.InstallmentTransactionNotFoundMsg))
	}

	return &installmentTransaction, nil
}
