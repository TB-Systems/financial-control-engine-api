package handlers

import (
	"backend-commons/dtos"

	"financialcontrol/internal/services"
	"net/http"

	"github.com/TB-Systems/go-commons/commonsmodels"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/gin-gonic/gin"
)

type CreditCard struct {
	service services.CreditCard
}

func NewCreditCardsHandler(service services.CreditCard) *CreditCard {
	return &CreditCard{service: service}
}

func (c *CreditCard) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		request, apiErr := utils.DecodeValidJson[dtos.CreditCardRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.Create(ctx.Request.Context(), userID, request)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusCreated)
	}
}

func (c *CreditCard) Read() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, apiErr := utils.GetUserIDFromContext(ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.Read(ctx.Request.Context(), userID)
		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *CreditCard) ReadAt() gin.HandlerFunc {
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

		data, apiErr := c.service.ReadAt(ctx.Request.Context(), userID, id)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *CreditCard) Update() gin.HandlerFunc {
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

		request, apiErr := utils.DecodeValidJson[dtos.CreditCardRequest](ctx)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		data, apiErr := c.service.Update(ctx.Request.Context(), userID, id, request)

		if apiErr != nil {
			utils.SendErrorResponse(ctx, apiErr)
			return
		}

		utils.SendResponse(ctx, data, http.StatusOK)
	}
}

func (c *CreditCard) Delete() gin.HandlerFunc {
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
