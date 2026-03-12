package repositories

import (
	"backend-commons/models"
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/google/uuid"
)

type MonthlyTransaction interface {
	ReadCategoryByID(context context.Context, categoryID uuid.UUID) (models.Category, error)
	ReadCreditCardByID(context context.Context, creditCardId uuid.UUID) (models.CreditCard, error)
	CreateMonthlyTransaction(ctx context.Context, request models.CreateMonthlyTransaction) (models.ShortMonthlyTransaction, error)
	ReadMonthlyTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.MonthlyTransaction, int64, error)
	ReadMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.MonthlyTransaction, error)
	ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error)
	ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error)
	ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error)
	UpdateMonthlyTransaction(ctx context.Context, model models.MonthlyTransaction) (models.ShortMonthlyTransaction, error)
	DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error
}

func (r Repository) CreateMonthlyTransaction(ctx context.Context, request models.CreateMonthlyTransaction) (models.ShortMonthlyTransaction, error) {
	param := pgstore.CreateMonthlyTransactionParams{
		UserID:       request.UserID,
		Name:         request.Name,
		Value:        utils.Float64ToNumeric(request.Value),
		Day:          request.Day,
		CategoryID:   request.CategoryID,
		CreditCardID: utils.UUIDToPgTypeUUID(request.CreditCardID),
	}

	response, err := r.store.CreateMonthlyTransaction(ctx, param)

	if err != nil {
		return models.ShortMonthlyTransaction{}, err
	}

	return models.ShortMonthlyTransaction{
		ID:        response.ID,
		Value:     utils.NumericToFloat64(response.Value),
		Day:       response.Day,
		CreatedAt: response.CreatedAt.Time,
		UpdatedAt: response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadMonthlyTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.MonthlyTransaction, int64, error) {
	param := pgstore.ListMonthlyTransactionsByUserIDPaginatedParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	responses, err := r.store.ListMonthlyTransactionsByUserIDPaginated(ctx, param)

	if err != nil {
		return nil, 0, err
	}

	if len(responses) == 0 {
		return []models.MonthlyTransaction{}, 0, nil
	}

	var transactions []models.MonthlyTransaction
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

		transactions = append(transactions, models.MonthlyTransaction{
			ID:         response.ID,
			Name:       response.Name,
			Value:      utils.NumericToFloat64(response.Value),
			Day:        response.Day,
			Category:   category,
			Creditcard: creditcard,
			CreatedAt:  response.CreatedAt.Time,
			UpdatedAt:  response.UpdatedAt.Time,
		})
	}

	return transactions, count, nil
}

func (r Repository) ReadMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.MonthlyTransaction, error) {
	response, err := r.store.GetMonthlyTransactionByID(ctx, id)

	if err != nil {
		return models.MonthlyTransaction{}, err
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

	return models.MonthlyTransaction{
		ID:         response.ID,
		UserID:     response.UserID,
		Name:       response.Name,
		Value:      utils.NumericToFloat64(response.Value),
		Day:        response.Day,
		Category:   category,
		Creditcard: creditcard,
		CreatedAt:  response.CreatedAt.Time,
		UpdatedAt:  response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	response, err := r.store.GetShortMonthlyTransactionByID(ctx, id)

	if err != nil {
		return models.ShortMonthlyTransaction{}, err
	}

	return models.ShortMonthlyTransaction{
		ID:           response.ID,
		UserID:       response.UserID,
		Name:         response.Name,
		Value:        utils.NumericToFloat64(response.Value),
		Day:          response.Day,
		CategoryID:   response.CategoryID,
		CreditCardID: utils.PgTypeUUIDToUUID(response.CreditCardID),
		CreatedAt:    response.CreatedAt.Time,
		UpdatedAt:    response.UpdatedAt.Time,
	}, nil
}

func (r Repository) UpdateMonthlyTransaction(ctx context.Context, model models.MonthlyTransaction) (models.ShortMonthlyTransaction, error) {
	param := pgstore.UpdateMonthlyTransactionParams{
		ID:           model.ID,
		Name:         model.Name,
		Value:        utils.Float64ToNumeric(model.Value),
		Day:          model.Day,
		CategoryID:   model.Category.ID,
		CreditCardID: utils.UUIDToPgTypeUUID(&model.Creditcard.ID),
	}

	response, err := r.store.UpdateMonthlyTransaction(ctx, param)

	if err != nil {
		return models.ShortMonthlyTransaction{}, err
	}

	return models.ShortMonthlyTransaction{
		ID:        response.ID,
		Value:     utils.NumericToFloat64(response.Value),
		Day:       response.Day,
		CreatedAt: response.CreatedAt.Time,
		UpdatedAt: response.UpdatedAt.Time,
	}, nil
}

func (r Repository) DeleteMonthlyTransaction(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteMonthlyTransaction(ctx, id)
}
