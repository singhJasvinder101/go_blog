package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/utils/jwt"
	"github.com/singhJasvinder101/go_blog/internal/utils/response"
)

// returns function for use as callback in router.use()
func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Authorization"){
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(fmt.Errorf("unauthorized")))
			c.Abort()
			return
		}

		strs := strings.Split(authHeader, " ")
		if len(strs) != 2 || strs[1] == ""{
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(fmt.Errorf("unauthorized")))
			c.Abort()
			return
		}

		user_data, err := jwt.VerifyToken(strs[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(fmt.Errorf("unauthorized")))
			c.Abort()
			return
		}

		// store user_id in context
		c.Set("user_data", user_data)
		c.Next()
	}
}


