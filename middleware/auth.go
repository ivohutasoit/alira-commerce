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
			c.Redirect(http.StatusMovedPermanently, redirect)
			c.Abort()
			return
		}
		if accessToken != nil {
			claims := &account.AccessTokenClaims{}
			token, err := jwt.ParseWithClaims(accessToken.(string), claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})
			if err != nil || !token.Valid {
				c.Redirect(http.StatusMovedPermanently, redirect)
				c.Abort()
				return
			}

			data := map[string]string{
				"type":  "Bearer",
				"token": accessToken.(string),
			}
			// https://tutorialedge.net/golang/consuming-restful-api-with-go/
			payload, _ := json.Marshal(data)
			resp, err := http.Post(os.Getenv("URL_AUTH"), "application/json", bytes.NewBuffer(payload))
			if err != nil && !opt {
				c.Redirect(http.StatusMovedPermanently, redirect)
				c.Abort()
				return
			}
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil && !opt {
				c.Redirect(http.StatusMovedPermanently, redirect)
				c.Abort()
				return
			}

			var userProfile messaging.UserProfile
			parser := &util.Parser{}
			parser.UnmarshalResponse(respData, http.StatusOK, &userProfile)

			fmt.Println(userProfile)

			c.Set("user_id", userProfile.ID)
			alira.ViewData = gin.H{
				"user_id":    userProfile.ID,
				"username":   userProfile.Username,
				"url_logout": fmt.Sprintf("%s?redirect=%s", os.Getenv("URL_LOGOUT"), url),
			}
		}
		c.Next()
	}
}

func (m *Auth) TokenRequired(args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
