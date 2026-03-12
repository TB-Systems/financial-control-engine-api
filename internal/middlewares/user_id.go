package middlewares

import (
	"backend-commons/constants"
	"net/http"

	"github.com/TB-Systems/go-commons/errors"
	"github.com/TB-Systems/go-commons/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDString, err := ctx.Cookie(constants.UserID)

		if err != nil {
			utils.SendErrorResponse(
				ctx,
				errors.NewApiError(http.StatusUnauthorized, errors.UserNotFound(err.Error())),
			)
			return
		}

		userID, err := uuid.Parse(userIDString)

		if err != nil {
			utils.SendErrorResponse(
				ctx,
				errors.NewApiError(http.StatusUnauthorized, errors.UserNotFound(err.Error())),
			)
			return
		}

		ctx.Set(constants.UserID, userID)
		ctx.Next()
	}
}
