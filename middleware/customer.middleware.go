package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/ivohutasoit/alira"
	"github.com/ivohutasoit/alira/database/commerce"

	"github.com/gin-gonic/gin"
)

type Customer struct{}

func (m *Customer) OwnerRequired(args ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		api := strings.Contains(c.Request.URL.Path, os.Getenv("URL_API"))
		userId, exist := c.Get("user_id")
		if !exist {
			if api {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":   http.StatusUnauthorized,
					"status": http.StatusText(http.StatusUnauthorized),
					"error":  "access denied",
				})
				return
			}
		}
		custUser := &commerce.CustomerUser{}
		alira.GetConnection().Where("user_id = ? AND ACTIVE = ? AND role = 'OWNER'",
			userId, true).First(&custUser)

		if custUser.Model.ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": http.StatusText(http.StatusUnauthorized),
				"error":  "access denied",
			})
			return
		}
		c.Set("customer_id", custUser.CustomerID)
		c.Next()
	}
}
