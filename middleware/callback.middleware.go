package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira-commerce/messaging"
	"github.com/ivohutasoit/alira/util"
)

type Callback struct{}

func (m *Callback) WebLoginToken(args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		callback := c.Query("callback")
		session := sessions.Default(c)

		if callback != "" {
			data := map[string]string{
				"token": callback,
			}
			// https://tutorialedge.net/golang/consuming-restful-api-with-go/
			payload, _ := json.Marshal(data)
			resp, err := http.Post(fmt.Sprintf("%s%s", os.Getenv("URL_ACCOUNT"), os.Getenv("API_TOKENCALLBACK")),
				"application/json", bytes.NewBuffer(payload))
			if err != nil {
				session.Clear()
				session.AddFlash("error", err.Error())
				session.Save()
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				session.Clear()
				session.AddFlash("error", err.Error())
				session.Save()
				c.Redirect(http.StatusMovedPermanently, "/")
				return
			}
			var loggedUser messaging.LoggedUser
			parser := &util.Parser{}
			_, err = parser.UnmarshalResponse(respData, http.StatusOK, &loggedUser)
			if err != nil {
				fmt.Println(err.Error())
			}

			if loggedUser.AccessToken != "" {
				session.Set("access_token", loggedUser.AccessToken)
				session.Set("refresh_token", loggedUser.RefreshToken)
				session.Save()
			}
		}
		c.Next()
	}
}
