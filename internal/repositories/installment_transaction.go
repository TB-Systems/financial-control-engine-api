package repositories

import (
	"backend-commons/models"
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/google/uuid"
)

type InstallmentTransaction interface {
	ReadCategoryByID(context context.Context, categoryID uuid.UUID) (models.Category, error)
	ReadCreditCardByID(context context.Context, creditCardId uuid.UUID) (models.CreditCard, error)
	CreateInstallmentTransaction(ctx context.Context, request models.CreateInstallmentTransaction) (models.ShortInstallmentTransaction, error)
	ReadInstallmentTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.InstallmentTransaction, int64, error)
	ReadInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.InstallmentTransaction, error)
	ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error)
	ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error)
	ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error)
	UpdateInstallmentTransaction(ctx context.Context, model models.InstallmentTransaction) (models.ShortInstallmentTransaction, error)
	DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error
}

func (r Repository) CreateInstallmentTransaction(ctx context.Context, request models.CreateInstallmentTransaction) (models.ShortInstallmentTransaction, error) {
	param := pgstore.CreateInstallmentTransactionParams{
		UserID:       request.UserID,
		Name:         request.Name,
		Value:        utils.Float64ToNumeric(request.Value),
		InitialDate:  utils.TimeToPgTimestamptz(request.InitialDate),
		FinalDate:    utils.TimeToPgTimestamptz(request.FinalDate),
		CategoryID:   request.CategoryID,
		CreditCardID: utils.UUIDToPgTypeUUID(request.CreditCardID),
	}

	response, err := r.store.CreateInstallmentTransaction(ctx, param)

	if err != nil {
		return models.ShortInstallmentTransaction{}, err
	}

	return models.ShortInstallmentTransaction{
		ID:          response.ID,
		Name:        response.Name,
		Value:       utils.NumericToFloat64(response.Value),
		InitialDate: response.InitialDate.Time,
		FinalDate:   response.FinalDate.Time,
		CreatedAt:   response.CreatedAt.Time,
		UpdatedAt:   response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadInstallmentTransactionsByUserIDPaginated(ctx context.Context, params commonsmodels.PaginatedParams) ([]models.InstallmentTransaction, int64, error) {
	param := pgstore.ListInstallmentTransactionsByUserIDPaginatedParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	responses, err := r.store.ListInstallmentTransactionsByUserIDPaginated(ctx, param)

	if err != nil {
		return nil, 0, err
	}

	if len(responses) == 0 {
		return []models.InstallmentTransaction{}, 0, nil
	}

	var transactions []models.InstallmentTransaction
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

		transactions = append(transactions, models.InstallmentTransaction{
			ID:          response.ID,
			Name:        response.Name,
			Value:       utils.NumericToFloat64(response.Value),
			InitialDate: response.InitialDate.Time,
			FinalDate:   response.FinalDate.Time,
			Category:    category,
			Creditcard:  creditcard,
			CreatedAt:   response.CreatedAt.Time,
			UpdatedAt:   response.UpdatedAt.Time,
		})
	}

	return transactions, count, nil
}

func (r Repository) ReadInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.InstallmentTransaction, error) {
	response, err := r.store.GetInstallmentTransactionByID(ctx, id)

	if err != nil {
		return models.InstallmentTransaction{}, err
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

	return models.InstallmentTransaction{
		ID:          response.ID,
		UserID:      response.UserID,
		Name:        response.Name,
		Value:       utils.NumericToFloat64(response.Value),
		InitialDate: response.InitialDate.Time,
		FinalDate:   response.FinalDate.Time,
		Category:    category,
		Creditcard:  creditcard,
		CreatedAt:   response.CreatedAt.Time,
		UpdatedAt:   response.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	response, err := r.store.GetShortInstallmentTransactionByID(ctx, id)

	if err != nil {
		return models.ShortInstallmentTransaction{}, err
	}

	return models.ShortInstallmentTransaction{
		ID:           response.ID,
		UserID:       response.UserID,
		Name:         response.Name,
		Value:        utils.NumericToFloat64(response.Value),
		InitialDate:  response.InitialDate.Time,
		FinalDate:    response.FinalDate.Time,
		CategoryID:   response.CategoryID,
		CreditCardID: utils.PgTypeUUIDToUUID(response.CreditCardID),
		CreatedAt:    response.CreatedAt.Time,
		UpdatedAt:    response.UpdatedAt.Time,
	}, nil
}

func (r Repository) UpdateInstallmentTransaction(ctx context.Context, model models.InstallmentTransaction) (models.ShortInstallmentTransaction, error) {
	var creditCardID *uuid.UUID
	if model.Creditcard != nil {
		creditCardID = &model.Creditcard.ID
	}

	param := pgstore.UpdateInstallmentTransactionParams{
		ID:           model.ID,
		Name:         model.Name,
		Value:        utils.Float64ToNumeric(model.Value),
		InitialDate:  utils.TimeToPgTimestamptz(model.InitialDate),
		FinalDate:    utils.TimeToPgTimestamptz(model.FinalDate),
		CategoryID:   model.Category.ID,
		CreditCardID: utils.UUIDToPgTypeUUID(creditCardID),
	}

	response, err := r.store.UpdateInstallmentTransaction(ctx, param)

	if err != nil {
		return models.ShortInstallmentTransaction{}, err
	}

	return models.ShortInstallmentTransaction{
		ID:          response.ID,
		Name:        response.Name,
		Value:       utils.NumericToFloat64(response.Value),
		InitialDate: response.InitialDate.Time,
		FinalDate:   response.FinalDate.Time,
		CreatedAt:   response.CreatedAt.Time,
		UpdatedAt:   response.UpdatedAt.Time,
	}, nil
}

func (r Repository) DeleteInstallmentTransaction(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteInstallmentTransaction(ctx, id)
}
