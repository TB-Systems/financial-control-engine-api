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

type InstallmentTransaction interface {
	Create(ctx context.Context, userID uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, errors.ApiError)
	Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse], errors.ApiError)
	ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.InstallmentTransactionResponse, errors.ApiError)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, errors.ApiError)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError
}

type installmentTransaction struct {
	repository repositories.InstallmentTransaction
}

func NewInstallmentTransactionService(repository repositories.InstallmentTransaction) InstallmentTransaction {
	return &installmentTransaction{
		repository: repository,
	}
}

func (s *installmentTransaction) Create(ctx context.Context, userID uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(s.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.InstallmentTransactionResponse{}, apiErr
	}

	createModel := modelsdto.CreateInstallmentTransactionFromRequest(request, userID)

	model, err := s.repository.CreateInstallmentTransaction(ctx, createModel)

	if err != nil {
		return dtos.InstallmentTransactionResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	response := modelsdto.InstallmentTransactionResponseFromShortModel(model, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (s *installmentTransaction) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse], errors.ApiError) {
	responses, count, err := s.repository.ReadInstallmentTransactionsByUserIDPaginated(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.InstallmentTransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.InstallmentTransactionResponseFromModel(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.InstallmentTransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.Limit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (s *installmentTransaction) ReadById(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.InstallmentTransactionResponse, errors.ApiError) {
	transaction, apiErr := s.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.InstallmentTransactionResponse{}, apiErr
	}

	response := modelsdto.InstallmentTransactionResponseFromModel(transaction)

	return response, nil
}

func (s *installmentTransaction) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.InstallmentTransactionRequest) (dtos.InstallmentTransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(s.repository, ctx, userID, request.CreditCardID, request.CategoryID)

	if apiErr != nil {
		return dtos.InstallmentTransactionResponse{}, apiErr
	}

	transaction, apiErr := s.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.InstallmentTransactionResponse{}, apiErr
	}

	transaction.Name = request.Name
	transaction.Value = request.Value
	transaction.InitialDate = request.InitialDate
	transaction.FinalDate = request.FinalDate
	transaction.Category = relations.CategoryModel
	transaction.Creditcard = relations.CreditcardModel

	transactionUpdated, err := s.repository.UpdateInstallmentTransaction(ctx, transaction)

	if err != nil {
		return dtos.InstallmentTransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	response := modelsdto.InstallmentTransactionResponseFromShortModel(transactionUpdated, relations.CategoryResponse, relations.CreditcardResponse)

	return response, nil
}

func (s *installmentTransaction) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError {
	transaction, apiErr := s.read(ctx, userID, id)

	if apiErr != nil {
		return apiErr
	}

	err := s.repository.DeleteInstallmentTransaction(ctx, transaction.ID)

	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	return nil
}

func (s installmentTransaction) read(ctx context.Context, userID uuid.UUID, id uuid.UUID) (models.InstallmentTransaction, errors.ApiError) {
	transaction, err := s.repository.ReadInstallmentTransactionByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.InstallmentTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.InstallmentTransactionNotFoundMsg))
		}
		return models.InstallmentTransaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if transaction.UserID != userID {
		return models.InstallmentTransaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.InstallmentTransactionNotFoundMsg))
	}

	return transaction, nil
}
