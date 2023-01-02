package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/middleware"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/dtos"
	v1 "github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/identity/features/registering_user/commands/v1"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity-service/shared/contracts"
	"github.com/pkg/errors"
	"net/http"
)

func MapRoute(infra *contracts.InfrastructureConfiguration) {
	group := infra.Echo.Group("/api/v1/users")
	group.POST("", createUser(infra), middleware.ValidateBearerToken())
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
func createUser(infra *contracts.InfrastructureConfiguration) echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()
		request := &dtos.RegisterUserRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerUserEndpoint_handler.Bind] error in the binding request")
			infra.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := v1.NewRegisterUser(request.FirstName, request.LastName, request.UserName, request.Email, request.Password)

		if err := infra.Validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerUserEndpoint_handler.StructCtx] command validation failed")
			infra.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*v1.RegisterUser, *dtos.RegisterUserResponseDto](ctx, command)

		if err != nil {
			infra.Log.Errorf("(RegisterUser.Handle) id: {%s}, err: {%v}", command.UserId, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		infra.Log.Infof("(user registered) id: {%s}", command.UserId)
		return c.JSON(http.StatusCreated, result)
	}
}
