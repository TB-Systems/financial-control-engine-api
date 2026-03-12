package services

import (
	"context"
	"backend-commons/constants"
	"backend-commons/dtos"
	"backend-commons/models"
	"backend-commons/modelsdto"
	"financialcontrol/internal/repositories"
	"net/http"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/errors"
	"github.com/google/uuid"
)

type AnnualTransaction interface {
	Create(ctx context.Context, userID uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, errors.ApiError)
	Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse], errors.ApiError)
	ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.AnnualTransactionResponse, errors.ApiError)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, errors.ApiError)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError
}

type annualTransaction struct {
	repository repositories.AnnualTransaction
}

func NewAnnualTransactionService(repository repositories.AnnualTransaction) AnnualTransaction {
	return &annualTransaction{
		repository: repository,
	}
}

func (a *annualTransaction) Create(ctx context.Context, userID uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(a.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.AnnualTransactionResponse{}, apiErr
	}

	createModel := modelsdto.CreateAnnualTransactionFromRequest(request, userID)

	model, err := a.repository.CreateAnnualTransaction(ctx, createModel)

	if err != nil {
		return dtos.AnnualTransactionResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	response := modelsdto.AnnualTransactionResponseFromShortModel(model, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (a *annualTransaction) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse], errors.ApiError) {
	responses, count, err := a.repository.ReadAnnualTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.AnnualTransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.AnnualTransactionResponseFromModel(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.AnnualTransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.Limit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (a *annualTransaction) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.AnnualTransactionResponse, errors.ApiError) {
	transaction, apiErr := a.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.AnnualTransactionResponse{}, apiErr
	}

	response := modelsdto.AnnualTransactionResponseFromModel(transaction)

	return response, nil
}

func (a *annualTransaction) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.AnnualTransactionRequest) (dtos.AnnualTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(a.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.AnnualTransactionResponse{}, apiErr
	}

	transaction, apiErr := a.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.AnnualTransactionResponse{}, apiErr
	}

	transaction.Name = request.Name
	transaction.Value = request.Value
	transaction.Day = request.Day
	transaction.Month = request.Month
	transaction.Category = relations.CategoryModel
	transaction.Creditcard = relations.CreditcardModel

	transactionUpdated, err := a.repository.UpdateAnnualTransaction(ctx, transaction)

	if err != nil {
		return dtos.AnnualTransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	response := modelsdto.AnnualTransactionResponseFromShortModel(transactionUpdated, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (a *annualTransaction) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError {
	transaction, apiErr := a.read(ctx, userID, id)

	if apiErr != nil {
		return apiErr
	}

	err := a.repository.DeleteAnnualTransaction(ctx, transaction.ID)

	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	return nil
}

func (a annualTransaction) read(ctx context.Context, userID uuid.UUID, id uuid.UUID) (models.AnnualTransaction, errors.ApiError) {
	transaction, err := a.repository.ReadAnnualTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.AnnualTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.AnnualTransactionNotFoundMsg))
		}
		return models.AnnualTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if transaction.UserID != userID {
		return models.AnnualTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.AnnualTransactionNotFoundMsg))
	}

	return transaction, nil
}
