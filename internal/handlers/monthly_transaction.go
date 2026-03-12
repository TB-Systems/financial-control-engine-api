package handlers

import (
	"backend-commons/dtos"
	"financialcontrol/internal/services"
	"net/http"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/gin-gonic/gin"
)

type MonthlyTransaction struct {
	service services.MonthlyTransaction
}

func NewMonthlyTransactionsHandler(service services.MonthlyTransaction) *MonthlyTransaction {
	return &MonthlyTransaction{
		service: service,
	}
}

func (h *MonthlyTransaction) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		request, apiErr := utils.DecodeValidJson[dtos.MonthlyTransactionRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := h.service.Create(ctx.Request.Context(), userID, request)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (h *MonthlyTransaction) Read() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		page, apiErr := utils.GetQueryPage(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		limit := utils.GetQueryLimit(ctx)
		offset := utils.CalculateOffset(page, limit)

		params := commonsmodels.PaginatedParams{
			UserID: userID,
			Limit:  limit,
			Offset: offset,
			Page:   page,
		}

		data, apiErr := h.service.Read(ctx.Request.Context(), params)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (h *MonthlyTransaction) ReadById() gin.HandlerFunc {
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

		data, apiErr := h.service.ReadById(ctx.Request.Context(), userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (h *MonthlyTransaction) Update() gin.HandlerFunc {
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

		request, apiErr := utils.DecodeValidJson[dtos.MonthlyTransactionRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := h.service.Update(ctx.Request.Context(), userID, id, request)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (h *MonthlyTransaction) Delete() gin.HandlerFunc {
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

		apiErr = h.service.Delete(ctx.Request.Context(), userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, commonsmodels.NewResponseSuccess(), http.StatusNoContent)
	}
}
