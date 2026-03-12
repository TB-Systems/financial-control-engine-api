package handlers

import (
	"financialcontrol/internal/services"
	"net/http"

	"github.com/TB-Systems/go-commons/utils"
	"github.com/gin-gonic/gin"
)

type MonthlyReport struct {
	service services.MonthlyReport
}

func NewMonthlyReportHandler(service services.MonthlyReport) *MonthlyReport {
	return &MonthlyReport{
		service: service,
	}
}

func (h *MonthlyReport) GetMonthlyBalance() gin.HandlerFunc {
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

		response, apiErr := h.service.GenerateMonthlyReport(ctx, userID, year, month)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, response, http.StatusOK)
	}
}
