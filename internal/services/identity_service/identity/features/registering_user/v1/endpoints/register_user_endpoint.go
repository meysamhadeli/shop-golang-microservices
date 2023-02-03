package endpoints

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	echomiddleware "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	commandsv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/features/registering_user/v1/commands"
	dtosv1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/identity/features/registering_user/v1/dtos"
	"github.com/pkg/errors"
	"net/http"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/users")
	group.POST("", createUser(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

// RegisterUser
// @Tags Users
// @Summary Register user
// @Description Create new user
// @Accept json
// @Produce json
// @Param RegisterUserRequestDto body dtos.RegisterUserRequestDto true "User data"
// @Success 201 {object} dtos.RegisterUserResponseDto
// @Security ApiKeyAuth
// @Router /api/v1/users [post]
func createUser(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := &dtosv1.RegisterUserRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerUserEndpoint_handler.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commandsv1.NewRegisterUser(request.FirstName, request.LastName, request.UserName, request.Email, request.Password)

		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerUserEndpoint_handler.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*commandsv1.RegisterUser, *dtosv1.RegisterUserResponseDto](ctx, command)

		if err != nil {
			log.Errorf("(RegisterUser.Handle) id: {%s}, err: {%v}", command.UserId, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("(user registered) id: {%s}", command.UserId)
		return c.JSON(http.StatusCreated, result)
	}
}
