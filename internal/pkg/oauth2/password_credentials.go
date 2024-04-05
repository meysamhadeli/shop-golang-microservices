package oauth2

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	srv        *server.Server
	once       sync.Once
	manager    *manage.Manager
	privateKey = []byte(`secret`)
	clients    = []*models.Client{{ID: "clientId", Secret: "clientSecret"}, {ID: "clientId2", Secret: "clientSecret2"}}
)

// User model
type User struct {
	UserId    uuid.UUID `json:"userId" gorm:"primaryKey"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserName  string    `json:"userName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func init() {
	manager = manage.NewDefaultManager()
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", privateKey, jwt.SigningMethodHS512))
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	once.Do(func() {
		srv = server.NewDefaultServer(manager)
	})
}

func clientStore(clients ...*models.Client) oauth2.ClientStore {
	clientStore := store.NewClientStore()

	for _, client := range clients {
		if client != nil {
			err := clientStore.Set(client.ID, &models.Client{
				ID:     client.ID,
				Secret: client.Secret,
				Domain: client.Domain,
			})
			if err != nil {
				return nil
			}
		}
	}
	return clientStore
}

// ref: https://github.com/go-oauth2/oauth2
func RunOauthServer(e *echo.Echo, gorm *gorm.DB) {

	manager.MapClientStorage(clientStore(clients...))

	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {

		u := User{}
		gorm.Where("user_name = ?", username).First(&u)

		// now use p
		isMatch, err := utils.ComparePasswords(u.Password, password)

		if err != nil {
			return "", err
		}

		if isMatch {
			return u.UserId.String(), nil
		}
		return
	})

	srv.SetClientScopeHandler(func(tgr *oauth2.TokenGenerateRequest) (allowed bool, err error) {
		if tgr.Scope == "all" {
			allowed = true
		}
		return
	})
	// for using querystring
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	e.GET("connect/token", token)
	e.GET("validate-token", validateBearerToken)
}

func validateBearerToken(c echo.Context) error {
	token, err := srv.ValidationBearerToken(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, token)
}

func token(c echo.Context) error {
	err := srv.HandleTokenRequest(c.Response().Writer, c.Request())
	if err != nil {
		return err
	}
	return nil
}
