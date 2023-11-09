package middleware

import (
	"awesomeProject1/initializers"
	"awesomeProject1/modules"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	fmt.Printf("log", "in middleware")
	tokenstring, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepected method", token.Header)
		}
		return []byte(os.Getenv("SECRET")), nil

	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["expr"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user modules.Users
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		c.Next()

	}
}
