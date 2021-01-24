package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"strings"
)

const UserSession = "user_session"

// Credentials stores the ClientID and ClientSecret
type Credentials struct {
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"secret"`
}

type OpenIDConnect struct {
	conf  Config
	store cookie.Store
}

type Config struct {
	BaseUrl          string
	Credentials      Credentials
	Oauth            *oauth2.Config
	UserInfoEndpoint string
	LoginURL         string
	Logger           *logrus.Logger
	SaveFunc         func(user *User)
}

// User is a retrieved and authenticated user.
type User struct {
	Id             string              `json:"user_id"`
	Name           string              `json:"name"`
	UserName       string              `json:"user_name"`
	GivenName      string              `json:"given_name"`
	FamilyName     string              `json:"family_name"`
	PhoneNumber    string              `json:"phone_number"`
	Email          string              `json:"email"`
	UserAttributes map[string][]string `json:"user_attributes"`
	Roles          []string            `json:"roles"`
}

func randToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// Setup the authorization path
func New(oidcConfig Config, secret []byte) OpenIDConnect {
	return OpenIDConnect{
		conf:  oidcConfig,
		store: cookie.NewStore(secret),
	}
}

func (o *OpenIDConnect) Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, o.store)
}

func (o *OpenIDConnect) LoginHandler(ctx *gin.Context) {
	state := randToken()
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()
	ctx.Writer.Write([]byte("<html><title>Golang OpenIDConnect</title> <body> <a href='" + o.GetLoginURL(state, o.conf.Oauth) + "'><button>Login with OpenIDConnect!</button> </a> </body></html>"))
}

func (o *OpenIDConnect) GetLoginURL(state string, oauth *oauth2.Config) string {
	return oauth.AuthCodeURL(state)
}

func (o *OpenIDConnect) EndHandler(ctx *gin.Context) {
	state := ctx.Request.URL.Query().Get("state")
	session := sessions.Default(ctx)
	stateCookie := session.Get("state")
	redirectUrlCookie := session.Get("redirect_url")

	if stateCookie == nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("unable to parse request, state cookie missing"))
		return
	}

	if state != stateCookie {
		o.conf.Logger.Warnf("%v != %v", state, stateCookie)
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("provided state differs from the one saved in cookies"))
		return
	}

	if redirectUrlCookie == nil {
		ctx.JSON(http.StatusOK, gin.H{"success": "true"})
		return
	}

	if !strings.HasPrefix(redirectUrlCookie.(string), o.conf.BaseUrl) {
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid redirect URL"))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrlCookie.(string))
}

// Auth is the OpenIDConnect authorization middleware. You can use it to protect a Router Group.
// Example:
//
//        private.Use(oidc.Auth())
//        private.GET("/", UserInfoHandler)
//        private.GET("/api", func(ctx *gin.Context) {
//            ctx.JSON(200, gin.H{"message": "Hello from private for groups"})
//        })
//    func UserInfoHandler(ctx *gin.Context) {
//        ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": ctx.MustGet("user").(oidc.User)})
//    }
func (o *OpenIDConnect) Auth(ctx *gin.Context) {
	// Check if user has logged in
	session := sessions.Default(ctx)
	userCookie := session.Get("user")
	if userCookie != nil {
		ctx.Set("user", userCookie)
		return
	}

	// Handle the exchange code to initiate a transport.
	retrievedState := session.Get("state")

	oauthCopy := *o.conf.Oauth
	oauthCopy.RedirectURL = o.conf.BaseUrl + ctx.Request.URL.String()

	if retrievedState == nil {
		state := randToken()
		session.Set("state", state)
		session.Set("redirect_url", oauthCopy.RedirectURL)
		session.Save()
		ctx.Redirect(http.StatusTemporaryRedirect, o.GetLoginURL(state, &oauthCopy))
		return
	}

	if retrievedState != ctx.Query("state") {
		state := randToken()
		session.Set("state", state)
		session.Set("redirect_url", oauthCopy.RedirectURL)
		session.Save()
		ctx.Redirect(http.StatusTemporaryRedirect, o.GetLoginURL(state, &oauthCopy))
		return
	}

	redirectUrl := session.Get("redirect_url")
	if redirectUrl == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing redirect_url"})
		return
	}
	oauthCopy.RedirectURL = redirectUrl.(string)

	tok, err := oauthCopy.Exchange(context.Background(), ctx.Query("code"))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := oauthCopy.Client(context.Background(), tok)
	email, err := client.Get(o.conf.UserInfoEndpoint)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer email.Body.Close()
	data, err := ioutil.ReadAll(email.Body)
	if err != nil {
		glog.Errorf("[Gin-OAuth] Could not read Body: %s", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		glog.Errorf("[Gin-OAuth] Unmarshal userinfo failed: %s", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Save user in DB
	o.conf.SaveFunc(&user)
	
	session.Set("user", user.Id)
	session.Save()
	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl.(string))
	ctx.Abort()
	return
}

func (o *OpenIDConnect) LogoutHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("user")
	session.Save()
}
