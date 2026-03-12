package services

import (
	"context"
	stdErrors "errors"
	"backend-commons/constants"
	"backend-commons/models"
	"testing"

	"github.com/google/uuid"
)

type commonRepositoryMock struct {
	category       models.Category
	categoryErr    error
	creditcard     models.CreditCard
	creditcardErr  error
	annual         models.ShortAnnualTransaction
	annualErr      error
	monthly        models.ShortMonthlyTransaction
	monthlyErr     error
	installment    models.ShortInstallmentTransaction
	installmentErr error
}

func (m *commonRepositoryMock) ReadCategoryByID(ctx context.Context, categoryID uuid.UUID) (models.Category, error) {
	if m.categoryErr != nil {
		return models.Category{}, m.categoryErr
	}
	return m.category, nil
}

func (m *commonRepositoryMock) ReadCreditCardByID(ctx context.Context, creditCardId uuid.UUID) (models.CreditCard, error) {
	if m.creditcardErr != nil {
		return models.CreditCard{}, m.creditcardErr
	}
	return m.creditcard, nil
}

func (m *commonRepositoryMock) ReadShortAnnualTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortAnnualTransaction, error) {
	if m.annualErr != nil {
		return models.ShortAnnualTransaction{}, m.annualErr
	}
	return m.annual, nil
}

func (m *commonRepositoryMock) ReadShortMonthlyTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortMonthlyTransaction, error) {
	if m.monthlyErr != nil {
		return models.ShortMonthlyTransaction{}, m.monthlyErr
	}
	return m.monthly, nil
}

func (m *commonRepositoryMock) ReadShortInstallmentTransactionByID(ctx context.Context, id uuid.UUID) (models.ShortInstallmentTransaction, error) {
	if m.installmentErr != nil {
		return models.ShortInstallmentTransaction{}, m.installmentErr
	}
	return m.installment, nil
}

func TestReadAnnualTransactionBranches(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock := &commonRepositoryMock{annual: models.ShortAnnualTransaction{ID: id, UserID: userID}}
		_, apiErr := readAnnualTransaction(mock, ctx, userID, id)
		if apiErr != nil {
			t.Fatalf("expected nil error")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &commonRepositoryMock{annualErr: stdErrors.New(constants.StoreErrorNoRowsMsg)}
		_, apiErr := readAnnualTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mock := &commonRepositoryMock{annualErr: stdErrors.New("db")}
		_, apiErr := readAnnualTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 500 {
			t.Fatalf("expected 500")
		}
	})

	t.Run("user mismatch", func(t *testing.T) {
		mock := &commonRepositoryMock{annual: models.ShortAnnualTransaction{ID: id, UserID: uuid.New()}}
		_, apiErr := readAnnualTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})
}

func TestReadMonthlyTransactionBranches(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock := &commonRepositoryMock{monthly: models.ShortMonthlyTransaction{ID: id, UserID: userID}}
		_, apiErr := readMonthlyTransaction(mock, ctx, userID, id)
		if apiErr != nil {
			t.Fatalf("expected nil error")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &commonRepositoryMock{monthlyErr: stdErrors.New(constants.StoreErrorNoRowsMsg)}
		_, apiErr := readMonthlyTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mock := &commonRepositoryMock{monthlyErr: stdErrors.New("db")}
		_, apiErr := readMonthlyTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 500 {
			t.Fatalf("expected 500")
		}
	})

	t.Run("user mismatch", func(t *testing.T) {
		mock := &commonRepositoryMock{monthly: models.ShortMonthlyTransaction{ID: id, UserID: uuid.New()}}
		_, apiErr := readMonthlyTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})
}

func TestReadInstallmentTransactionBranches(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mock := &commonRepositoryMock{installment: models.ShortInstallmentTransaction{ID: id, UserID: userID}}
		_, apiErr := readInstallmentTransaction(mock, ctx, userID, id)
		if apiErr != nil {
			t.Fatalf("expected nil error")
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &commonRepositoryMock{installmentErr: stdErrors.New(constants.StoreErrorNoRowsMsg)}
		_, apiErr := readInstallmentTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mock := &commonRepositoryMock{installmentErr: stdErrors.New("db")}
		_, apiErr := readInstallmentTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 500 {
			t.Fatalf("expected 500")
		}
	})

	t.Run("user mismatch", func(t *testing.T) {
		mock := &commonRepositoryMock{installment: models.ShortInstallmentTransaction{ID: id, UserID: uuid.New()}}
		_, apiErr := readInstallmentTransaction(mock, ctx, userID, id)
		if apiErr == nil || apiErr.GetStatus() != 404 {
			t.Fatalf("expected 404")
		}
	})
}
