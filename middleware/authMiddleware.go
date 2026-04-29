package middleware

import (
	"fmt"
	"net/http"
	helper "user_management_system/helpers"

	"github.com/gin-gonic/gin"
)

//It authenticates the request using a token
// If valid → allows request and sets user data
func Authenticate()gin.HandlerFunc{
	return func (c *gin.Context)  {
		clientToken:=c.Request.Header.Get("token")
		if clientToken==""{
		  c.JSON(http.StatusUnauthorized,gin.H{"error":fmt.Sprintln("No Authorization header provide")})
		  c.Abort()
		  return
		}
		claims,err:=helper.ValidateToken(clientToken)

		if err!=""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":err})
			c.Abort()
			return
		}
		c.Set("email",claims.Email)
		c.Set("first_name",claims.First_name)
		c.Set("last_name",claims.Last_name)
		c.Set("uid",claims.Uid)
		c.Set("user_type",claims.User_type)

		c.Next()
	}
}
