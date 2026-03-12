package repositories

import (
	"backend-commons/models"
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/google/uuid"
)

type AnnualTransaction interface {
	ReadCategoryByID(context context.Context, categoryID uuid.UUID) (models.Category, error)
	ReadCreditCardByID(context context.Context, creditCardId uuid.UUID) (models.CreditCard, error)
	CreateAnnualTransaction(ctx context.Context, request models.CreateAnnualTransaction) (models.ShortAnnualTransaction, error)
	ReadAnnualTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.AnnualTransaction, int64, error)
	ReadAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.AnnualTransaction, error)
	ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error)
	ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error)
	ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error)
	UpdateAnnualTransaction(ctx context.Context, model models.AnnualTransaction) (models.ShortAnnualTransaction, error)
	DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error
}

func (r Repository) CreateAnnualTransaction(ctx context.Context, request models.CreateAnnualTransaction) (models.ShortAnnualTransaction, error) {
	param := pgstore.CreateAnnualTransactionParams{
		UserID:       request.UserID,
		Name:         request.Name,
		Value:        utils.Float64ToNumeric(request.Value),
		Day:          request.Day,
		Month:        request.Month,
		CategoryID:   request.CategoryID,
		CreditCardID: utils.UUIDToPgTypeUUID(request.CreditCardID),
	}

	response, err := r.store.CreateAnnualTransaction(ctx, param)

	if err != nil {
		return models.ShortAnnualTransaction{}, err
	}

	return models.ShortAnnualTransaction{
		ID:        response.ID,
		Value:     utils.NumericToFloat64(response.Value),
		Day:       response.Day,
		Month:     response.Month,
		CreatedAt: response.CreatedAt.Time,
		UpdatedAt: response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadAnnualTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.AnnualTransaction, int64, error) {
	param := pgstore.ListAnnualTransactionsByUserIDPaginatedParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	responses, err := r.store.ListAnnualTransactionsByUserIDPaginated(ctx, param)

	if err != nil {
		return nil, 0, err
	}

	if len(responses) == 0 {
		return []models.AnnualTransaction{}, 0, nil
	}

	var transactions []models.AnnualTransaction
	count := responses[len(responses)-1].TotalCount

	for _, response := range responses {
		category := models.Category{
			ID:              *utils.PgTypeUUIDToUUID(response.CategoryID),
			TransactionType: models.TransactionType(response.CategoryTransactionType.Int32),
			Name:            response.CategoryName.String,
			Icon:            response.CategoryIcon.String,
			CreatedAt:       response.CategoryCreatedAt.Time,
			UpdatedAt:       response.CategoryUpdatedAt.Time,
		}

		var creditcard *models.CreditCard
		if response.CreditcardID.Valid {
			creditcardValue := models.CreditCard{
				ID:               *utils.PgTypeUUIDToUUID(response.CreditcardID),
				Name:             response.CreditcardName.String,
				FirstFourNumbers: response.CreditcardFirstFourNumbers.String,
				Limit:            response.CreditcardCreditLimit.Float64,
				CloseDay:         response.CreditcardCloseDay.Int32,
				ExpireDay:        response.CreditcardExpireDay.Int32,
				BackgroundColor:  response.CreditcardBackgroundColor.String,
				TextColor:        response.CreditcardTextColor.String,
				CreatedAt:        response.CreditcardCreatedAt.Time,
				UpdatedAt:        response.CreditcardUpdatedAt.Time,
			}

			creditcard = &creditcardValue
		}

		transactions = append(transactions, models.AnnualTransaction{
			ID:         response.ID,
			Name:       response.Name,
			Value:      utils.NumericToFloat64(response.Value),
			Day:        response.Day,
			Month:      response.Month,
			Category:   category,
			Creditcard: creditcard,
			CreatedAt:  response.CreatedAt.Time,
			UpdatedAt:  response.UpdatedAt.Time,
		})
	}

	return transactions, count, nil
}

func (r Repository) ReadAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.AnnualTransaction, error) {
	response, err := r.store.GetAnnualTransactionByID(ctx, id)

	if err != nil {
		return models.AnnualTransaction{}, err
	}

	category := models.Category{
		ID:              *utils.PgTypeUUIDToUUID(response.CategoryID),
		TransactionType: models.TransactionType(response.CategoryTransactionType.Int32),
		Name:            response.CategoryName.String,
		Icon:            response.CategoryIcon.String,
		CreatedAt:       response.CategoryCreatedAt.Time,
		UpdatedAt:       response.CategoryUpdatedAt.Time,
	}

	var creditcard *models.CreditCard
	if response.CreditcardID.Valid {
		creditcardValue := models.CreditCard{
			ID:               *utils.PgTypeUUIDToUUID(response.CreditcardID),
			Name:             response.CreditcardName.String,
			FirstFourNumbers: response.CreditcardFirstFourNumbers.String,
			Limit:            response.CreditcardCreditLimit.Float64,
			CloseDay:         response.CreditcardCloseDay.Int32,
			ExpireDay:        response.CreditcardExpireDay.Int32,
			BackgroundColor:  response.CreditcardBackgroundColor.String,
			TextColor:        response.CreditcardTextColor.String,
			CreatedAt:        response.CreditcardCreatedAt.Time,
			UpdatedAt:        response.CreditcardUpdatedAt.Time,
		}

		creditcard = &creditcardValue
	}

	return models.AnnualTransaction{
		ID:         response.ID,
		UserID:     response.UserID,
		Name:       response.Name,
		Value:      utils.NumericToFloat64(response.Value),
		Day:        response.Day,
		Month:      response.Month,
		Category:   category,
		Creditcard: creditcard,
		CreatedAt:  response.CreatedAt.Time,
		UpdatedAt:  response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	response, err := r.store.GetShortAnnualTransactionByID(ctx, id)

	if err != nil {
		return models.ShortAnnualTransaction{}, err
	}

	return models.ShortAnnualTransaction{
		ID:           response.ID,
		UserID:       response.UserID,
		Name:         response.Name,
		Value:        utils.NumericToFloat64(response.Value),
		Day:          response.Day,
		Month:        response.Month,
		CategoryID:   response.CategoryID,
		CreditCardID: utils.PgTypeUUIDToUUID(response.CreditCardID),
		CreatedAt:    response.CreatedAt.Time,
		UpdatedAt:    response.UpdatedAt.Time,
	}, nil
}

func (r Repository) UpdateAnnualTransaction(ctx context.Context, model models.AnnualTransaction) (models.ShortAnnualTransaction, error) {
	param := pgstore.UpdateAnnualTransactionParams{
		ID:           model.ID,
		Name:         model.Name,
		Value:        utils.Float64ToNumeric(model.Value),
		Day:          model.Day,
		Month:        model.Month,
		CategoryID:   model.Category.ID,
		CreditCardID: utils.UUIDToPgTypeUUID(&model.Creditcard.ID),
	}

	response, err := r.store.UpdateAnnualTransaction(ctx, param)

	if err != nil {
		return models.ShortAnnualTransaction{}, err
	}

	return models.ShortAnnualTransaction{
		ID:        response.ID,
		Name:      response.Name,
		Value:     utils.NumericToFloat64(response.Value),
		Day:       response.Day,
		Month:     response.Month,
		CreatedAt: response.CreatedAt.Time,
		UpdatedAt: response.UpdatedAt.Time,
	}, nil
}

func (r Repository) DeleteAnnualTransaction(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteAnnualTransaction(ctx, id)
}
