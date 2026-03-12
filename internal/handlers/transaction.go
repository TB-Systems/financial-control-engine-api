package handlers

import (
	"backend-commons/dtos"
	"financialcontrol/internal/services"
	"net/http"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/errors"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Transaction struct {
	service services.Transaction
}

func NewTransactionsHandler(service services.Transaction) *Transaction {
	return &Transaction{service: service}
}

func (c *Transaction) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		request, apiErr := utils.DecodeValidJson[dtos.TransactionRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.Create(ctx.Request.Context(), request, userID)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (c *Transaction) CreateFromMonthlyTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, request, apiErr := c.getRecurrentRequest(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.CreateFromMonthlyTransaction(ctx.Request.Context(), request, *userID)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (c *Transaction) CreateFromAnnualTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, request, apiErr := c.getRecurrentRequest(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.CreateFromAnnualTransaction(ctx.Request.Context(), request, *userID)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (c *Transaction) CreateFromInstallmentTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, request, apiErr := c.getRecurrentRequest(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.CreateFromInstallmentTransaction(ctx.Request.Context(), request, *userID)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (c *Transaction) Read() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		limit := utils.GetQueryLimit(ctx)

		page, apiErr := utils.GetQueryPage(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		offset := utils.CalculateOffset(page, limit)

		startDate, endDate, hasDates, apiErr := utils.GetQueryDatesIfHas(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		var data commonsmodels.PaginatedResponse[dtos.TransactionResponse]

		if hasDates {
			params := commonsmodels.PaginatedParamsWithDateRange{
				UserID:    userID,
				Limit:     limit,
				Offset:    offset,
				Page:      page,
				StartDate: startDate,
				EndDate:   endDate,
			}

			data, apiErr = c.service.ReadInToDates(ctx.Request.Context(), params)
		} else {
			params := commonsmodels.PaginatedParams{
				UserID: userID,
				Limit:  limit,
				Offset: offset,
				Page:   page,
			}

			data, apiErr = c.service.Read(ctx.Request.Context(), params)
		}

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *Transaction) ReadByMonthAndYear() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		month, year, apiErr := utils.GetQueryMonthAndYear(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		limit := utils.GetQueryLimit(ctx)

		page, apiErr := utils.GetQueryPage(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		offset := utils.CalculateOffset(page, limit)

		params := commonsmodels.PaginatedParamsWithMonthYear{
			UserID:     userID,
			Month:      month,
			Year:       year,
			Page:       page,
			PageOffset: offset,
			PageLimit:  limit,
		}

		data, apiErr := c.service.ReadAtMonthAndYear(ctx, params)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *Transaction) ReadById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		id, apiErr := utils.IDFromURLParam(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.ReadById(ctx.Request.Context(), userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *Transaction) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		id, apiErr := utils.IDFromURLParam(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		request, apiErr := utils.DecodeValidJson[dtos.TransactionRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.Update(ctx.Request.Context(), request, userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *Transaction) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		id, apiErr := utils.IDFromURLParam(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		apiErr = c.service.Delete(ctx, userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, commonsmodels.NewResponseSuccess(), http.StatusOK)
	}
}

func (c *Transaction) Pay() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		id, apiErr := utils.IDFromURLParam(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		apiErr = c.service.Pay(ctx.Request.Context(), userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, commonsmodels.NewResponseSuccess(), http.StatusOK)
	}
}

func (c *Transaction) getRecurrentRequest(ctx *gin.Context) (*uuid.UUID, dtos.TransactionRequestFromRecurrentTransaction, errors.ApiError) {
	userID, apiErr := utils.GetUserIDFromContext(ctx)

	if apiErr != nil {
		return nil, dtos.TransactionRequestFromRecurrentTransaction{}, apiErr
	}

	request, apiErr := utils.DecodeValidJson[dtos.TransactionRequestFromRecurrentTransaction](ctx)

	if apiErr != nil {
		return nil, dtos.TransactionRequestFromRecurrentTransaction{}, apiErr
	}

	return &userID, request, nil
}
