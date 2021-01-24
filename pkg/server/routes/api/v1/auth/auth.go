package auth

import (
	"crypto/rand"
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/denysvitali/dev-portal/pkg/server/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

func Setup(r *gin.RouterGroup, app *app.App) {
	var appOidcConfig = app.Config.Oidc
	oauth2Config := oauth2.Config{
		ClientID:     appOidcConfig.ClientID,
		ClientSecret: appOidcConfig.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   appOidcConfig.AuthURL,
			TokenURL:  appOidcConfig.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: app.Config.BaseUrl + r.BasePath() + "/oauth/end",
		Scopes:      nil,
	}

	saveUser := func(user *oidc.User) {
		app.Log.Infof("saving user: %v", user)
	}
	oidcConfig := oidc.Config{
		BaseUrl:          app.Config.BaseUrl,
		Oauth:            &oauth2Config,
		UserInfoEndpoint: appOidcConfig.UserInfoURL,
		Logger:           app.Log,
		SaveFunc:         saveUser,
	}

	var secret = make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		app.Log.Fatalf("unable to generate secret: %v", err)
	}

	openIdConnect := oidc.New(oidcConfig, secret)
	r.GET("/login", openIdConnect.LoginHandler)
	r.GET("/logout", openIdConnect.LogoutHandler)
	r.Use(openIdConnect.Auth) // Redirect if token is missing

	r.GET("/abc", func(ctx *gin.Context) {
		user, exists := ctx.Get("user")
		if !exists {
			ctx.Status(http.StatusForbidden)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
		ctx.Status(http.StatusOK)
		ctx.Done()
	})
}
