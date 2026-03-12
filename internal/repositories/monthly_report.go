package repositories

import (
	"backend-commons/models"
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/TB-Systems/go-commons/utils"
	"github.com/google/uuid"
)

type MonthlyReport interface {
	GetMonthlyBalance(ctx context.Context, userID uuid.UUID, year int32, month int32) (models.MonthlyReport, error)
	GetCategoriesSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CategoriesSpending, error)
	GetCreditCardsSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CreditCardsSpending, error)
}

func (r Repository) GetMonthlyBalance(ctx context.Context, userID uuid.UUID, year int32, month int32) (models.MonthlyReport, error) {
	param := pgstore.GetMonthlyBalanceParams{
		UserID:  userID,
		Column2: year,
		Column3: month,
	}

	response, err := r.store.GetMonthlyBalance(ctx, param)

	if err != nil {
		return models.MonthlyReport{}, err
	}

	return models.MonthlyReport{
		TotalIncome: utils.NumericToFloat64(response.TotalIncome),
		TotalDebit:  utils.NumericToFloat64(response.TotalDebit),
		TotalCredit: utils.NumericToFloat64(response.TotalCredit),
		Balance:     utils.NumericToFloat64(response.Balance),
	}, nil
}

func (r Repository) GetCategoriesSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CategoriesSpending, error) {
	param := pgstore.GetCategoriesSpendingParams{
		UserID:  userID,
		Column2: year,
		Column3: month,
	}

	responses, err := r.store.GetCategoriesSpending(ctx, param)

	if err != nil {
		return []models.CategoriesSpending{}, err
	}

	if len(responses) == 0 {
		return []models.CategoriesSpending{}, nil
	}

	categories := make([]models.CategoriesSpending, len(responses))

	for i, response := range responses {
		categories[i] = models.CategoriesSpending{
			CategoryID:              response.CategoryID,
			CategoryName:            response.CategoryName,
			CategoryIcon:            response.CategoryIcon,
			CategoryTransactionType: models.TransactionType(response.CategoryTransactionType),
			TotalSpent:              utils.NumericToFloat64(response.TotalSpent),
		}
	}

	return categories, nil
}

func (r Repository) GetCreditCardsSpending(ctx context.Context, userID uuid.UUID, year int32, month int32) ([]models.CreditCardsSpending, error) {
	param := pgstore.GetCreditCardsSpendingParams{
		UserID:  userID,
		Column2: year,
		Column3: month,
	}

	responses, err := r.store.GetCreditCardsSpending(ctx, param)

	if err != nil {
		return []models.CreditCardsSpending{}, err
	}

	if len(responses) == 0 {
		return []models.CreditCardsSpending{}, nil
	}

	creditCards := make([]models.CreditCardsSpending, len(responses))

	for i, response := range responses {
		creditCards[i] = models.CreditCardsSpending{
			ID:               response.CreditCardID,
			Name:             response.CreditCardName,
			FirstFourNumbers: response.CreditCardFirstFourNumbers,
			Limit:            response.CreditCardLimit,
			CloseDay:         response.CreditCardCloseDay,
			ExpireDay:        response.CreditCardExpireDay,
			BackgroundColor:  response.CreditCardBackgroundColor,
			TextColor:        response.CreditCardTextColor,
			TotalSpent:       utils.NumericToFloat64(response.TotalSpent),
		}
	}

	return creditCards, nil
}
