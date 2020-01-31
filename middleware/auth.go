package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira/database/account"
	"github.com/ivohutasoit/alira/messaging"
	"github.com/ivohutasoit/alira/util"
)

type Auth struct{}

func (m *Auth) SessionRequired(args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path
		except := os.Getenv("WEB_EXCEPT")
		if except != "" {
			excepts := strings.Split(except, ";")
			for _, value := range excepts {
				if currentPath == strings.TrimSpace(value) {
					c.Next()
					return
				}
			}
		}

		opt := false
		optional := os.Getenv("WEB_OPTIONAL")
		if optional != "" {
			optionals := strings.Split(optional, ";")
			for _, value := range optionals {
				if value == "/" && (currentPath == "" || currentPath == "/") {
					opt = true
					break
				} else {
					if c.Request.Method == http.MethodGet {
						if strings.Index(currentPath, value) > 0 {
							opt = true
							return
						}
					}
				}
			}
		}

		url, err := util.GenerateUrl(c.Request.TLS, c.Request.Host, currentPath, true)
		if err != nil {
			fmt.Println(err)
		}
		redirect := fmt.Sprintf("%s?redirect=%s", os.Getenv("URL_LOGIN"), url)

		session := sessions.Default(c)
		accessToken := session.Get("access_token")
		if accessToken == nil && !opt {
			session.Clear()
			session.Save()
			c.Redirect(http.StatusMovedPermanently, redirect)
			c.Abort()
			return
		}
		if accessToken != nil {
			claims := &account.AccessTokenClaims{}
			token, err := jwt.ParseWithClaims(accessToken.(string), claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})
			if err != nil && !opt {
				session.Clear()
				session.Save()
				c.Redirect(http.StatusMovedPermanently, redirect)
				c.Abort()
				return
			}

			if !token.Valid && !opt {
				session.Clear()
				session.Save()
				c.Redirect(http.StatusMovedPermanently, redirect)
				c.Abort()
				return
			}

			if token.Valid {
				data := map[string]string{
					"type":  "Bearer",
					"token": accessToken.(string),
				}
				// https://tutorialedge.net/golang/consuming-restful-api-with-go/
				payload, _ := json.Marshal(data)
				resp, err := http.Post(os.Getenv("URL_AUTH"), "application/json", bytes.NewBuffer(payload))
				if err != nil && !opt {
					session.Clear()
					session.Save()
					c.Redirect(http.StatusMovedPermanently, redirect)
					c.Abort()
					return
				}
				respData, err := ioutil.ReadAll(resp.Body)
				if err != nil && !opt {
					session.Clear()
					session.Save()
					c.Redirect(http.StatusMovedPermanently, redirect)
					c.Abort()
					return
				}

				var userProfile messaging.UserProfile
				parser := &util.Parser{}
				parser.UnmarshalResponse(respData, http.StatusOK, &userProfile)

				c.Set("user_id", userProfile.ID)
				alira.ViewData = gin.H{
					"user_id":    userProfile.ID,
					"username":   userProfile.Username,
					"url_logout": fmt.Sprintf("%s?redirect=%s", os.Getenv("URL_LOGOUT"), url),
				}
			}
		}
		c.Next()
	}
}

func (m *Auth) TokenRequired(args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path
		except := os.Getenv("API_EXCEPT")
		if except != "" {
			excepts := strings.Split(except, ";")

			for _, value := range excepts {
				if c.Request.Method == http.MethodGet {
					if strings.Index(currentPath, value) > 0 {
						c.Next()
						return
					}
				} else if currentPath == strings.TrimSpace(value) {
					c.Next()
					return
				}
			}
		}

		authorization := c.Request.Header.Get("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "missing authorization token",
			})
			return
		}

		tokens := strings.Split(authorization, " ")
		if len(tokens) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "invalid token",
			})
			return
		}

		var claims jwt.Claims
		if tokens[0] == "Bearer" {
			claims = &account.AccessTokenClaims{}
		} else if tokens[0] == "Refresh" {
			if currentPath != "/api/alpha/auth/refresh" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":   http.StatusUnauthorized,
					"status": http.StatusText(http.StatusUnauthorized),
					"error":  "invalid refresh uri",
				})
				return
			}
			claims = &account.RefreshTokenClaims{}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "invalid token indentifier",
			})
			return
		}

		tokenString := tokens[1]
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  err.Error(),
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "invalid token",
			})
			return
		}
		if tokens[0] == "Refresh" {
			if claims.(*account.RefreshTokenClaims).Sub != 1 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":   http.StatusUnauthorized,
					"status": http.StatusText(http.StatusUnauthorized),
					"error":  "invalid refresh token",
				})
				return
			}
		}

		data := map[string]string{
			"type":  tokens[0],
			"token": tokenString,
		}
		// https://tutorialedge.net/golang/consuming-restful-api-with-go/
		payload, _ := json.Marshal(data)
		resp, err := http.Post(os.Getenv("URL_AUTH"), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "unable to verify token",
			})
			return
		}
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "unable to read detail token",
			})
			return
		}
		var userProfile messaging.UserProfile
		parser := &util.Parser{}
		parser.UnmarshalResponse(respData, http.StatusOK, &userProfile)

		c.Set("user_id", userProfile.ID)
		if tokens[0] == "Refresh" {
			c.Set("sub", claims.(*account.RefreshTokenClaims).Sub)
		}
		c.Next()
	}
}
