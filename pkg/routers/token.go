package routers

import (
	"net/http"
	"time"

	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "toc_user"

// AuthMiddleware AuthMiddleware
func AuthMiddleware(g *gin.Engine) *jwt.GinJWTMiddleware {
	randomkey := utils.RandomString(50)
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "trade_agent",
		Key:            []byte(randomkey),
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieName:     "toc_token",
		CookieSameSite: http.SameSiteDefaultMode,
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		IdentityKey:    "toc_key",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*dbagent.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &dbagent.User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: UserAuthenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*dbagent.User); ok && v.UserName == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: toc_token, cookie: toc_token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		return nil
	}
	return authMiddleware
}

// UserAuthenticator UserAuthenticator
func UserAuthenticator(c *gin.Context) (interface{}, error) {
	var loginVals dbagent.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	userID := loginVals.Username
	password := loginVals.Password
	if userID == "admin" && password == "asdf0000" {
		return &dbagent.User{
			UserName: userID,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}
