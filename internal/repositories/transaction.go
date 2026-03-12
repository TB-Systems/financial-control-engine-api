package repositories

import (
	"backend-commons/models"
	"context"
	"financialcontrol/internal/store/pgstore"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction interface {
	ReadCategoryByID(context context.Context, categoryID uuid.UUID) (models.Category, error)
	ReadCreditCardByID(context context.Context, creditCardId uuid.UUID) (models.CreditCard, error)
	CreateTransaction(context context.Context, transaction models.CreateTransaction) (models.ShortTransaction, error)
	ReadTransactions(context context.Context, params commonsmodels.PaginatedParams) ([]models.Transaction, int64, error)
	ReadTransactionsInToDates(context context.Context, params commonsmodels.PaginatedParamsWithDateRange) ([]models.Transaction, int64, error)
	ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error)
	ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error)
	ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error)
	ReadTransactionsByMonthYear(context context.Context, params commonsmodels.PaginatedParamsWithMonthYear) ([]models.Transaction, int64, error)
	ReadTransactionById(context context.Context, id uuid.UUID) (models.Transaction, error)
	UpdateTransaction(context context.Context, transaction models.Transaction) (models.ShortTransaction, error)
	DeleteTransaction(context context.Context, id uuid.UUID) error
	PayTransaction(context context.Context, id uuid.UUID, paid bool) error
}

func (r Repository) CreateTransaction(context context.Context, transaction models.CreateTransaction) (models.ShortTransaction, error) {
	value := utils.Float64ToNumeric(transaction.Value)

	param := pgstore.CreateTransactionParams{
		UserID:                    transaction.UserID,
		Name:                      transaction.Name,
		Date:                      pgtype.Timestamptz{Time: transaction.Date, Valid: true},
		Value:                     value,
		Paid:                      transaction.Paid,
		CategoryID:                transaction.CategoryID,
		CreditCardID:              utils.UUIDToPgTypeUUID(transaction.CreditcardID),
		MonthlyTransactionsID:     utils.UUIDToPgTypeUUID(transaction.MonthlyTransactionID),
		AnnualTransactionsID:      utils.UUIDToPgTypeUUID(transaction.AnnualTransactionID),
		InstallmentTransactionsID: utils.UUIDToPgTypeUUID(transaction.InstallmentTransactionID),
	}

	createdTransaction, err := r.store.CreateTransaction(context, param)

	if err != nil {
		return models.ShortTransaction{}, err
	}

	return models.ShortTransaction{
		ID:        createdTransaction.ID,
		Name:      createdTransaction.Name,
		Date:      createdTransaction.Date.Time,
		Value:     utils.NumericToFloat64(createdTransaction.Value),
		Paid:      createdTransaction.Paid,
		CreatedAt: createdTransaction.CreatedAt.Time,
		UpdatedAt: createdTransaction.UpdatedAt.Time,
	}, nil
}

func (r Repository) ReadTransactions(context context.Context, params commonsmodels.PaginatedParams) ([]models.Transaction, int64, error) {
	args := pgstore.ListTransactionsByUserIDPaginatedParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	transactions, err := r.store.ListTransactionsByUserIDPaginated(context, args)

	if err != nil {
		return []models.Transaction{}, 0, err
	}

	if len(transactions) == 0 {
		return []models.Transaction{}, 0, nil
	}

	var transactionModels []models.Transaction
	count := transactions[len(transactions)-1].TotalCount

	for _, transaction := range transactions {
		transactionModels = append(transactionModels, storeTransactionToTransaction(storeTransactionPaginatedToStoreTransaction(transaction)))
	}

	return transactionModels, count, nil
}

func (r Repository) ReadTransactionsInToDates(context context.Context, params commonsmodels.PaginatedParamsWithDateRange) ([]models.Transaction, int64, error) {
	args := pgstore.ListTransactionsByUserAndDateParams{
		UserID: params.UserID,
		Limit:  params.Limit,
		Offset: params.Offset,
		Date:   pgtype.Timestamptz{Time: params.StartDate, Valid: true},
		Date_2: pgtype.Timestamptz{Time: params.EndDate, Valid: true},
	}

	transactions, err := r.store.ListTransactionsByUserAndDate(context, args)

	if err != nil {
		return []models.Transaction{}, 0, err
	}

	if len(transactions) == 0 {
		return []models.Transaction{}, 0, nil
	}

	var transactionModels []models.Transaction
	count := transactions[len(transactions)-1].TotalCount

	for _, transaction := range transactions {
		transactionModels = append(transactionModels, storeTransactionToTransaction(storeTransactionListToStoreTransaction(transaction)))
	}

	return transactionModels, count, nil
}

func (r Repository) ReadTransactionsByMonthYear(context context.Context, params commonsmodels.PaginatedParamsWithMonthYear) ([]models.Transaction, int64, error) {
	args := pgstore.ListTransactionsByUserAndMonthYearPaginatedParams{
		UserID:     params.UserID,
		Year:       params.Year,
		Month:      params.Month,
		PageLimit:  params.PageLimit,
		PageOffset: params.PageOffset,
	}

	transactions, err := r.store.ListTransactionsByUserAndMonthYearPaginated(context, args)

	if err != nil {
		return []models.Transaction{}, 0, err
	}

	if len(transactions) == 0 {
		return []models.Transaction{}, 0, nil
	}

	var transactionModels []models.Transaction
	count := transactions[len(transactions)-1].TotalCount

	for _, transaction := range transactions {
		transactionModels = append(transactionModels, storeTransactionToTransaction(storeTransactionListToStoreTransaction(pgstore.ListTransactionsByUserAndDateRow(transaction))))
	}

	return transactionModels, count, nil
}

func (r Repository) ReadTransactionById(context context.Context, id uuid.UUID) (models.Transaction, error) {
	transaction, err := r.store.GetTransactionByID(context, id)

	if err != nil {
		return models.Transaction{}, err
	}

	return storeTransactionToTransaction(transaction), nil
}

func (r Repository) UpdateTransaction(context context.Context, transaction models.Transaction) (models.ShortTransaction, error) {
	value := utils.Float64ToNumeric(transaction.Value)
	var creditCardID pgtype.UUID
	var monthlyTransactionID *uuid.UUID
	var annualTransactionID *uuid.UUID
	var installmentTransactionID *uuid.UUID

	if transaction.Creditcard != nil {
		creditCardID = utils.UUIDToPgTypeUUID(&transaction.Creditcard.ID)
	}

	if transaction.MonthlyTransaction != nil {
		monthlyTransactionID = &transaction.MonthlyTransaction.ID
	}

	if transaction.AnnualTransaction != nil {
		annualTransactionID = &transaction.AnnualTransaction.ID
	}

	if transaction.InstallmentTransaction != nil {
		installmentTransactionID = &transaction.InstallmentTransaction.ID
	}

	params := pgstore.UpdateTransactionParams{
		ID:                        transaction.ID,
		Name:                      transaction.Name,
		Date:                      pgtype.Timestamptz{Time: transaction.Date, Valid: true},
		Value:                     value,
		Paid:                      transaction.Paid,
		CategoryID:                transaction.Category.ID,
		CreditCardID:              creditCardID,
		MonthlyTransactionsID:     utils.UUIDToPgTypeUUID(monthlyTransactionID),
		AnnualTransactionsID:      utils.UUIDToPgTypeUUID(annualTransactionID),
		InstallmentTransactionsID: utils.UUIDToPgTypeUUID(installmentTransactionID),
	}

	transactionUpdated, err := r.store.UpdateTransaction(context, params)

	if err != nil {
		return models.ShortTransaction{}, err
	}

	return models.ShortTransaction{
		ID:        transactionUpdated.ID,
		Name:      transactionUpdated.Name,
		Date:      transactionUpdated.Date.Time,
		Value:     utils.NumericToFloat64(transactionUpdated.Value),
		Paid:      transactionUpdated.Paid,
		CreatedAt: transactionUpdated.CreatedAt.Time,
		UpdatedAt: transactionUpdated.UpdatedAt.Time,
	}, nil
}

func (r Repository) DeleteTransaction(context context.Context, id uuid.UUID) error {
	err := r.store.DeleteTransaction(context, id)

	if err != nil {
		return err
	}

	return nil
}

func (r Repository) PayTransaction(context context.Context, id uuid.UUID, paid bool) error {
	params := pgstore.PayTransactionParams{
		ID:   id,
		Paid: paid,
	}

	err := r.store.PayTransaction(context, params)

	if err != nil {
		return err
	}

	return nil
}

func storeTransactionListToStoreTransaction(transaction pgstore.ListTransactionsByUserAndDateRow) pgstore.GetTransactionByIDRow {
	return storeTransactionPaginatedToStoreTransaction(pgstore.ListTransactionsByUserIDPaginatedRow(transaction))
}

func storeTransactionPaginatedToStoreTransaction(transaction pgstore.ListTransactionsByUserIDPaginatedRow) pgstore.GetTransactionByIDRow {
	return pgstore.GetTransactionByIDRow{
		ID:                                 transaction.ID,
		UserID:                             transaction.UserID,
		Name:                               transaction.Name,
		Date:                               transaction.Date,
		Value:                              transaction.Value,
		Paid:                               transaction.Paid,
		CreatedAt:                          transaction.CreatedAt,
		UpdatedAt:                          transaction.UpdatedAt,
		CategoryID:                         transaction.CategoryID,
		CategoryTransactionType:            transaction.CategoryTransactionType,
		CategoryName:                       transaction.CategoryName,
		CategoryIcon:                       transaction.CategoryIcon,
		CategoryCreatedAt:                  transaction.CategoryCreatedAt,
		CategoryUpdatedAt:                  transaction.CategoryUpdatedAt,
		CreditcardID:                       transaction.CreditcardID,
		CreditcardName:                     transaction.CreditcardName,
		CreditcardFirstFourNumbers:         transaction.CreditcardFirstFourNumbers,
		CreditcardCreditLimit:              transaction.CreditcardCreditLimit,
		CreditcardCloseDay:                 transaction.CreditcardCloseDay,
		CreditcardExpireDay:                transaction.CreditcardExpireDay,
		CreditcardBackgroundColor:          transaction.CreditcardBackgroundColor,
		CreditcardTextColor:                transaction.CreditcardTextColor,
		CreditcardCreatedAt:                transaction.CreditcardCreatedAt,
		CreditcardUpdatedAt:                transaction.CreditcardUpdatedAt,
		MonthlyTransactionsID:              transaction.MonthlyTransactionsID,
		MonthlyTransactionsDay:             transaction.MonthlyTransactionsDay,
		MonthlyTransactionsName:            transaction.MonthlyTransactionsName,
		MonthlyTransactionsValue:           transaction.MonthlyTransactionsValue,
		MonthlyTransactionsCreatedAt:       transaction.MonthlyTransactionsCreatedAt,
		MonthlyTransactionsUpdatedAt:       transaction.MonthlyTransactionsUpdatedAt,
		AnnualTransactionsID:               transaction.AnnualTransactionsID,
		AnnualTransactionsMonth:            transaction.AnnualTransactionsMonth,
		AnnualTransactionsDay:              transaction.AnnualTransactionsDay,
		AnnualTransactionsName:             transaction.AnnualTransactionsName,
		AnnualTransactionsValue:            transaction.AnnualTransactionsValue,
		AnnualTransactionsCreatedAt:        transaction.AnnualTransactionsCreatedAt,
		AnnualTransactionsUpdatedAt:        transaction.AnnualTransactionsUpdatedAt,
		InstallmentTransactionsID:          transaction.InstallmentTransactionsID,
		InstallmentTransactionsName:        transaction.InstallmentTransactionsName,
		InstallmentTransactionsValue:       transaction.InstallmentTransactionsValue,
		InstallmentTransactionsInitialDate: transaction.InstallmentTransactionsInitialDate,
		InstallmentTransactionsFinalDate:   transaction.InstallmentTransactionsFinalDate,
		InstallmentTransactionsCreatedAt:   transaction.InstallmentTransactionsCreatedAt,
		InstallmentTransactionsUpdatedAt:   transaction.InstallmentTransactionsUpdatedAt,
	}
}

func storeTransactionToTransaction(transaction pgstore.GetTransactionByIDRow) models.Transaction {
	category := models.Category{
		ID:              *utils.PgTypeUUIDToUUID(transaction.CategoryID),
		TransactionType: models.TransactionType(transaction.CategoryTransactionType.Int32),
		Name:            transaction.CategoryName.String,
		Icon:            transaction.CategoryIcon.String,
		CreatedAt:       transaction.CategoryCreatedAt.Time,
		UpdatedAt:       transaction.CategoryUpdatedAt.Time,
	}

	var creditcard *models.CreditCard
	creditCardID := utils.PgTypeUUIDToUUID(transaction.CreditcardID)
	if transaction.CreditcardID.Valid {
		creditcardValue := models.CreditCard{
			ID:               *creditCardID,
			Name:             transaction.CreditcardName.String,
			FirstFourNumbers: transaction.CreditcardFirstFourNumbers.String,
			Limit:            transaction.CreditcardCreditLimit.Float64,
			CloseDay:         transaction.CreditcardCloseDay.Int32,
			ExpireDay:        transaction.CreditcardExpireDay.Int32,
			BackgroundColor:  transaction.CreditcardBackgroundColor.String,
			TextColor:        transaction.CreditcardTextColor.String,
			CreatedAt:        transaction.CreditcardCreatedAt.Time,
			UpdatedAt:        transaction.CreditcardUpdatedAt.Time,
		}

		creditcard = &creditcardValue
	}

	var monthlyTransaction *models.ShortMonthlyTransaction
	if transaction.MonthlyTransactionsID.Valid {
		monthlyTransactionValue := models.ShortMonthlyTransaction{
			ID:           *utils.PgTypeUUIDToUUID(transaction.MonthlyTransactionsID),
			UserID:       transaction.UserID,
			Day:          transaction.MonthlyTransactionsDay.Int32,
			Name:         transaction.MonthlyTransactionsName.String,
			Value:        utils.NumericToFloat64(transaction.MonthlyTransactionsValue),
			CategoryID:   category.ID,
			CreditCardID: creditCardID,
			CreatedAt:    transaction.MonthlyTransactionsCreatedAt.Time,
			UpdatedAt:    transaction.MonthlyTransactionsUpdatedAt.Time,
		}

		monthlyTransaction = &monthlyTransactionValue
	}

	var annualTransaction *models.ShortAnnualTransaction
	if transaction.AnnualTransactionsID.Valid {
		annualTransactionValue := models.ShortAnnualTransaction{
			ID:           *utils.PgTypeUUIDToUUID(transaction.AnnualTransactionsID),
			UserID:       transaction.UserID,
			Month:        transaction.AnnualTransactionsMonth.Int32,
			Day:          transaction.AnnualTransactionsDay.Int32,
			Name:         transaction.AnnualTransactionsName.String,
			Value:        utils.NumericToFloat64(transaction.AnnualTransactionsValue),
			CategoryID:   category.ID,
			CreditCardID: creditCardID,
			CreatedAt:    transaction.AnnualTransactionsCreatedAt.Time,
			UpdatedAt:    transaction.AnnualTransactionsUpdatedAt.Time,
		}

		annualTransaction = &annualTransactionValue
	}

	var installmentTransaction *models.ShortInstallmentTransaction
	if transaction.InstallmentTransactionsID.Valid {
		installmentTransactionValue := models.ShortInstallmentTransaction{
			ID:           *utils.PgTypeUUIDToUUID(transaction.InstallmentTransactionsID),
			UserID:       transaction.UserID,
			Name:         transaction.InstallmentTransactionsName.String,
			Value:        utils.NumericToFloat64(transaction.InstallmentTransactionsValue),
			CategoryID:   category.ID,
			CreditCardID: creditCardID,
			InitialDate:  transaction.InstallmentTransactionsInitialDate.Time,
			FinalDate:    transaction.InstallmentTransactionsFinalDate.Time,
			CreatedAt:    transaction.InstallmentTransactionsCreatedAt.Time,
			UpdatedAt:    transaction.InstallmentTransactionsUpdatedAt.Time,
		}

		installmentTransaction = &installmentTransactionValue
	}

	return models.Transaction{
		ID:                     transaction.ID,
		UserID:                 transaction.UserID,
		Name:                   transaction.Name,
		Date:                   transaction.Date.Time,
		Value:                  utils.NumericToFloat64(transaction.Value),
		Paid:                   transaction.Paid,
		Category:               category,
		Creditcard:             creditcard,
		MonthlyTransaction:     monthlyTransaction,
		AnnualTransaction:      annualTransaction,
		InstallmentTransaction: installmentTransaction,
		CreatedAt:              transaction.CreatedAt.Time,
		UpdatedAt:              transaction.UpdatedAt.Time,
	}
}
