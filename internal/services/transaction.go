package services

import (
	"context"
	"backend-commons/constants"
	"backend-commons/dtos"
	"backend-commons/models"
	"backend-commons/modelsdto"
	"financialcontrol/internal/repositories"
	"fmt"
	"net/http"
	"time"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/errors"
	"github.com/google/uuid"
)

type Transaction interface {
	Create(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	CreateFromMonthlyTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	CreateFromAnnualTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	CreateFromInstallmentTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError)
	ReadInToDates(ctx context.Context, params commonsmodels.PaginatedParamsWithDateRange) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError)
	ReadAtMonthAndYear(ctx context.Context, params commonsmodels.PaginatedParamsWithMonthYear) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError)
	ReadById(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	Update(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, errors.ApiError)
	Delete(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) errors.ApiError
	Pay(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) errors.ApiError
}

type transaction struct {
	repository repositories.Transaction
}

func NewTransactionsService(repository repositories.Transaction) Transaction {
	return &transaction{
		repository: repository,
	}
}

func (t transaction) Create(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(t.repository, ctx, userID, request.CreditcardID, request.CategoryID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	data := modelsdto.CreateTransactionFromTransactionRequest(request, userID)

	transaction, err := t.repository.CreateTransaction(ctx, data)

	if err != nil {
		return dtos.TransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	response := modelsdto.TransactionResponseFromShortTransaction(
		transaction,
		relations.CategoryResponse,
		relations.CreditcardResponse,
		nil,
		nil,
		nil,
	)

	return response, nil
}

func (t transaction) CreateFromMonthlyTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	monthlyTransaction, apiErr := readMonthlyTransaction(t.repository, ctx, userID, request.ID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), int(monthlyTransaction.Day), 0, 0, 0, 0, now.Location())

	createRequest := dtos.TransactionRequest{
		Name:                 monthlyTransaction.Name,
		Date:                 date,
		Value:                monthlyTransaction.Value,
		Paid:                 false,
		CategoryID:           monthlyTransaction.CategoryID,
		CreditcardID:         monthlyTransaction.CreditCardID,
		MonthlyTransactionID: &monthlyTransaction.ID,
	}

	createModel := modelsdto.CreateTransactionFromTransactionRequest(createRequest, userID)

	relations, apiErr := getRelations(t.repository, ctx, userID, createRequest.CreditcardID, createRequest.CategoryID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	transaction, err := t.repository.CreateTransaction(ctx, createModel)

	if err != nil {
		return dtos.TransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	shortMonthltTransactionResponse := modelsdto.ShortMonthlyTransactionResponseFromShortModel(*monthlyTransaction)

	response := modelsdto.TransactionResponseFromShortTransaction(
		transaction,
		relations.CategoryResponse,
		relations.CreditcardResponse,
		&shortMonthltTransactionResponse,
		nil,
		nil,
	)

	return response, nil
}

func (t transaction) CreateFromAnnualTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	annualTransaction, apiErr := readAnnualTransaction(t.repository, ctx, userID, request.ID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	now := time.Now()
	date := time.Date(now.Year(), time.Month(annualTransaction.Month), int(annualTransaction.Day), 0, 0, 0, 0, now.Location())

	createRequest := dtos.TransactionRequest{
		Name:                annualTransaction.Name,
		Date:                date,
		Value:               annualTransaction.Value,
		Paid:                false,
		CategoryID:          annualTransaction.CategoryID,
		CreditcardID:        annualTransaction.CreditCardID,
		AnnualTransactionID: &annualTransaction.ID,
	}

	createModel := modelsdto.CreateTransactionFromTransactionRequest(createRequest, userID)

	relations, apiErr := getRelations(t.repository, ctx, userID, createRequest.CreditcardID, createRequest.CategoryID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	transaction, err := t.repository.CreateTransaction(ctx, createModel)

	if err != nil {
		return dtos.TransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	annualTransactionResponse := modelsdto.ShortAnnualTransactionResponseFromShortModel(*annualTransaction)

	response := modelsdto.TransactionResponseFromShortTransaction(
		transaction,
		relations.CategoryResponse,
		relations.CreditcardResponse,
		nil,
		&annualTransactionResponse,
		nil,
	)

	return response, nil
}

func (t transaction) CreateFromInstallmentTransaction(ctx context.Context, request dtos.TransactionRequestFromRecurrentTransaction, userID uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	installmentTransaction, apiErr := readInstallmentTransaction(t.repository, ctx, userID, request.ID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), installmentTransaction.InitialDate.Day(), 0, 0, 0, 0, now.Location())

	totalInstallments := ((installmentTransaction.FinalDate.Year()-installmentTransaction.InitialDate.Year())*12 +
		int(installmentTransaction.FinalDate.Month()-installmentTransaction.InitialDate.Month())) + 1

	if totalInstallments < 1 {
		totalInstallments = 1
	}

	currentInstallment := ((now.Year()-installmentTransaction.InitialDate.Year())*12 +
		int(now.Month()-installmentTransaction.InitialDate.Month())) + 1

	if currentInstallment < 1 {
		currentInstallment = 1
	}

	if currentInstallment > totalInstallments {
		currentInstallment = totalInstallments
	}

	installmentName := fmt.Sprintf("%s %d/%d", installmentTransaction.Name, currentInstallment, totalInstallments)

	createRequest := dtos.TransactionRequest{
		Name:                     installmentName,
		Date:                     date,
		Value:                    installmentTransaction.Value,
		Paid:                     false,
		CategoryID:               installmentTransaction.CategoryID,
		CreditcardID:             installmentTransaction.CreditCardID,
		InstallmentTransactionID: &installmentTransaction.ID,
	}

	createModel := modelsdto.CreateTransactionFromTransactionRequest(createRequest, userID)

	relations, apiErr := getRelations(t.repository, ctx, userID, createRequest.CreditcardID, createRequest.CategoryID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	transaction, err := t.repository.CreateTransaction(ctx, createModel)

	if err != nil {
		return dtos.TransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	shortInstallmentTransactionResponse := modelsdto.ShortInstallmentTransactionResponseFromShortModel(*installmentTransaction)

	response := modelsdto.TransactionResponseFromShortTransaction(
		transaction,
		relations.CategoryResponse,
		relations.CreditcardResponse,
		nil,
		nil,
		&shortInstallmentTransactionResponse,
	)

	return response, nil
}

func (t transaction) Read(ctx context.Context, params commonsmodels.PaginatedParams) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError) {
	responses, count, err := t.repository.ReadTransactions(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.TransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.TransactionResponseFromTransaction(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.Limit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (t transaction) ReadInToDates(ctx context.Context, params commonsmodels.PaginatedParamsWithDateRange) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError) {
	responses, count, err := t.repository.ReadTransactionsInToDates(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.TransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.TransactionResponseFromTransaction(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.Limit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (t transaction) ReadAtMonthAndYear(ctx context.Context, params commonsmodels.PaginatedParamsWithMonthYear) (commonsmodels.PaginatedResponse[dtos.TransactionResponse], errors.ApiError) {
	responses, count, err := t.repository.ReadTransactionsByMonthYear(ctx, params)

	if err != nil {
		return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	transactionsResponse := make([]dtos.TransactionResponse, 0, len(responses))

	for _, transaction := range responses {
		transactionsResponse = append(transactionsResponse, modelsdto.TransactionResponseFromTransaction(transaction))
	}

	return commonsmodels.PaginatedResponse[dtos.TransactionResponse]{
		Items:     transactionsResponse,
		PageCount: (count / int64(params.PageLimit)) + 1,
		Page:      int64(params.Page),
	}, nil
}

func (t transaction) ReadById(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	transaction, apiErr := t.read(ctx, userID, transactionId)
	return modelsdto.TransactionResponseFromTransaction(transaction), apiErr
}

func (t transaction) Update(ctx context.Context, request dtos.TransactionRequest, userID uuid.UUID, transactionId uuid.UUID) (dtos.TransactionResponse, errors.ApiError) {
	relations, apiErr := getRelations(t.repository, ctx, userID, request.CreditcardID, request.CategoryID)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	transaction, apiErr := t.read(ctx, userID, transactionId)

	if apiErr != nil {
		return dtos.TransactionResponse{}, apiErr
	}

	transaction.Name = request.Name
	transaction.Date = request.Date
	transaction.Paid = request.Paid
	transaction.Category = relations.CategoryModel
	transaction.Creditcard = relations.CreditcardModel

	transactionUpdated, err := t.repository.UpdateTransaction(ctx, transaction)

	if err != nil {
		return dtos.TransactionResponse{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	response := modelsdto.TransactionResponseFromShortTransaction(
		transactionUpdated,
		relations.CategoryResponse,
		relations.CreditcardResponse,
		nil,
		nil,
		nil,
	)

	return response, nil
}

func (t transaction) Delete(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) errors.ApiError {
	transaction, apiErr := t.read(ctx, userID, transactionId)

	if apiErr != nil {
		return apiErr
	}

	err := t.repository.DeleteTransaction(ctx, transaction.ID)

	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	return nil
}

func (t transaction) Pay(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) errors.ApiError {
	transaction, apiErr := t.read(ctx, userID, transactionId)

	if apiErr != nil {
		return apiErr
	}

	err := t.repository.PayTransaction(ctx, transaction.ID, !transaction.Paid)

	if err != nil {
		return errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	return nil
}

func (t transaction) read(ctx context.Context, userID uuid.UUID, transactionId uuid.UUID) (models.Transaction, errors.ApiError) {
	transaction, err := t.repository.ReadTransactionById(ctx, transactionId)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.Transaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.TransactionNotFoundMsg))
		}
		return models.Transaction{}, errors.NewApiError(http.StatusInternalServerError, errors.InternalServerError(err.Error()))
	}

	if transaction.UserID != userID {
		return models.Transaction{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.TransactionNotFoundMsg))
	}

	return transaction, nil
}
