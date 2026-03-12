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

type CreditCard interface {
	Create(ctx context.Context, userID uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, errors.ApiError)
	Read(ctx context.Context, userID uuid.UUID) (commonsmodels.ResponseList[dtos.CreditCardResponse], errors.ApiError)
	ReadAt(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.CreditCardResponse, errors.ApiError)
	Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, errors.ApiError)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError
}

type creditCard struct {
	repository repositories.CreditCard
}

func NewCreditCardsService(repository repositories.CreditCard) CreditCard {
	return &creditCard{repository: repository}
}

func (c creditCard) Create(ctx context.Context, userID uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, errors.ApiError) {
	count, err := c.repository.ReadCountByUser(ctx, userID)

	if err != nil {
		return dtos.CreditCardResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	if count >= 10 {
		return dtos.CreditCardResponse{}, errors.NewApiError(
			http.StatusForbidden,
			errors.BadRequestError(constants.LimitReachedMsg),
		)
	}

	model := modelsdto.CreateCreditCardFromCreditCardRequest(request, userID)

	creditCard, err := c.repository.CreateCreditCard(ctx, model)

	if err != nil {
		return dtos.CreditCardResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	return modelsdto.CreditCardResponseFromCreditCard(creditCard), nil
}

func (c creditCard) Read(ctx context.Context, userID uuid.UUID) (commonsmodels.ResponseList[dtos.CreditCardResponse], errors.ApiError) {
	creditCards, err := c.repository.ReadCreditCards(ctx, userID)

	if err != nil {
		return commonsmodels.ResponseList[dtos.CreditCardResponse]{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	creditCardsResponse := make([]dtos.CreditCardResponse, 0, len(creditCards))
	for _, creditCard := range creditCards {
		creditCardsResponse = append(creditCardsResponse, modelsdto.CreditCardResponseFromCreditCard(creditCard))
	}

	return commonsmodels.ResponseList[dtos.CreditCardResponse]{
		Items: creditCardsResponse,
		Total: len(creditCardsResponse),
	}, nil
}

func (c creditCard) ReadAt(ctx context.Context, userID uuid.UUID, id uuid.UUID) (dtos.CreditCardResponse, errors.ApiError) {
	creditcard, err := c.read(ctx, userID, id)

	if err != nil {
		return dtos.CreditCardResponse{}, err
	}

	return modelsdto.CreditCardResponseFromCreditCard(creditcard), nil
}

func (c creditCard) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, request dtos.CreditCardRequest) (dtos.CreditCardResponse, errors.ApiError) {
	creditcard, apiErr := c.read(ctx, userID, id)

	if apiErr != nil {
		return dtos.CreditCardResponse{}, apiErr
	}

	creditcard.Name = request.Name
	creditcard.FirstFourNumbers = request.FirstFourNumbers
	creditcard.Limit = request.Limit
	creditcard.CloseDay = request.CloseDay
	creditcard.ExpireDay = request.ExpireDay
	creditcard.BackgroundColor = request.BackgroundColor
	creditcard.TextColor = request.TextColor

	creditcard, err := c.repository.UpdateCreditCard(ctx, creditcard)

	if err != nil {
		return dtos.CreditCardResponse{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	return modelsdto.CreditCardResponseFromCreditCard(creditcard), nil
}

func (c creditCard) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) errors.ApiError {
	creditcard, apiErr := c.read(ctx, userID, id)

	if apiErr != nil {
		return apiErr
	}

	hasTransactions, err := c.repository.HasTransactionsByCreditCard(ctx, creditcard.ID)

	if err != nil {
		return errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	if hasTransactions {
		return errors.NewApiError(
			http.StatusBadRequest,
			errors.BadRequestError(constants.CannotBeDeletedMsg),
		)
	}

	err = c.repository.DeleteCreditCard(ctx, creditcard.ID)

	if err != nil {
		return errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	return nil
}

func (c creditCard) read(ctx context.Context, userID uuid.UUID, id uuid.UUID) (models.CreditCard, errors.ApiError) {
	creditcard, err := c.repository.ReadCreditCardByID(ctx, id)

	if err != nil {
		if err.Error() == constants.StoreErrorNoRowsMsg {
			return models.CreditCard{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CreditcardNotFoundMsg))
		}
		return models.CreditCard{}, errors.NewApiError(
			http.StatusInternalServerError,
			errors.InternalServerError(err.Error()),
		)
	}

	if creditcard.UserID != userID {
		return models.CreditCard{}, errors.NewApiError(http.StatusNotFound, errors.NotFoundError(constants.CreditcardNotFoundMsg))
	}

	return creditcard, nil
}
