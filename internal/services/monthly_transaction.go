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

type MonthlyTransaction interface {
	Create(ctx context.Context, userID uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, errors.ApiError)
	Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse], errors.ApiError)
	ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.MonthlyTransactionResponse, errors.ApiError)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, errors.ApiError)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError
}

type monthlyTransaction struct {
	repository repositories.MonthlyTransaction
}

func NewMonthlyTransactionService(repository repositories.MonthlyTransaction) MonthlyTransaction {
	return &monthlyTransaction{
		repository: repository,
	}
}

func (m *monthlyTransaction) Create(ctx context.Context, userID uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(m.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.MonthlyTransactionResponse{}, apiErr
	}

	createModel := modelsdto.CreateMonthlyTransactionFromRequest(request, userID)

	model, err := m.repository.CreateMonthlyTransaction(ctx, createModel)

	if err != nil {
		return dtos.MonthlyTransactionResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	response := modelsdto.MonthlyTransactionResponseFromShortModel(model, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (m *monthlyTransaction) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse], errors.ApiError) {
	responses, count, err := m.repository.ReadMonthlyTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.MonthlyTransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.MonthlyTransactionResponseFromModel(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.MonthlyTransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.Limit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (m *monthlyTransaction) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.MonthlyTransactionResponse, errors.ApiError) {
	transaction, apiErr := m.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.MonthlyTransactionResponse{}, apiErr
	}

	response := modelsdto.MonthlyTransactionResponseFromModel(transaction)

	return response, nil
}

func (m *monthlyTransaction) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.MonthlyTransactionRequest) (dtos.MonthlyTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(m.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.MonthlyTransactionResponse{}, apiErr
	}

	transaction, apiErr := m.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.MonthlyTransactionResponse{}, apiErr
	}

	transaction.Name = request.Name
	transaction.Value = request.Value
	transaction.Day = request.Day
	transaction.Category = relations.CategoryModel
	transaction.Creditcard = relations.CreditcardModel

	transactionUpdated, err := m.repository.UpdateMonthlyTransaction(ctx, transaction)

	if err != nil {
		return dtos.MonthlyTransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	response := modelsdto.MonthlyTransactionResponseFromShortModel(transactionUpdated, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (m *monthlyTransaction) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError {
	transaction, apiErr := m.read(ctx, userID, id)

	if apiErr != nil {
		return apiErr
	}

	err := m.repository.DeleteMonthlyTransaction(ctx, transaction.ID)

	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	return nil
}

func (m monthlyTransaction) read(ctx context.Context, userID uuid.UUID, id uuid.UUID) (models.MonthlyTransaction, errors.ApiError) {
	transaction, err := m.repository.ReadMonthlyTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.MonthlyTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.MonthlyTransactionNotFoundMsg))
		}
		return models.MonthlyTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if transaction.UserID != userID {
		return models.MonthlyTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.MonthlyTransactionNotFoundMsg))
	}

	return transaction, nil
}
