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

type registerUserEndpoint struct {
	*contracts.IdentityEndpointBase[contracts.InfrastructureConfiguration]
}

func NewCreteUserEndpoint(endpointBase *contracts.IdentityEndpointBase[contracts.InfrastructureConfiguration]) *registerUserEndpoint {
	return &registerUserEndpoint{endpointBase}
}

func (ep *registerUserEndpoint) MapRoute() {
	ep.ProductsGroup.POST("", ep.createUser(), middleware.ValidateBearerToken())
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
func (ep *registerUserEndpoint) createUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx := c.Request().Context()
		request := &dtos.RegisterUserRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerUserEndpoint_handler.Bind] error in the binding request")
			ep.Configuration.Log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := v1.NewRegisterUser(request.FirstName, request.LastName, request.UserName, request.Email, request.Password)

		if err := ep.Configuration.Validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerUserEndpoint_handler.StructCtx] command validation failed")
			ep.Configuration.Log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*v1.RegisterUser, *dtos.RegisterUserResponseDto](ctx, command)

		if err != nil {
			ep.Configuration.Log.Errorf("(RegisterUser.Handle) id: {%s}, err: {%v}", command.UserId, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		ep.Configuration.Log.Infof("(user registered) id: {%s}", command.UserId)
		return c.JSON(http.StatusCreated, result)
	}
}
