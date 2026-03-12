package repositories

import (
	"backend-commons/models"
	"context"

	"github.com/google/uuid"
)

type CommonRepository interface {
	ReadCategoryByID(context context.Context, categoryID uuid.UUID) (models.Category, error)
	ReadCreditCardByID(context context.Context, creditCardId uuid.UUID) (models.CreditCard, error)
	ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error)
	ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error)
	ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error)
}
