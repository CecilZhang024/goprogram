package controller

import (
	"awesomeProject1/initializers"
	"awesomeProject1/modules"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	var body struct {
		Email    string
		Username string
		Password string
		Type     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to hash password",
		})
		return
	}

	user := modules.Users{Email: body.Email, Password: string(hash), Username: body.Username, Type: "editor"}
	fmt.Printf("bodyjson", body)
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "create user successfully"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Username string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}
	//var user modules.Users
	user := modules.Users{Email: body.Email, Password: body.Password, Username: body.Username}
	initializers.DB.First(&user, "email = ?", user.Email)
	fmt.Printf("bodyjson", user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invaild email or username",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invaild password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"expr": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  200,
			"message": "Login Succeed",
			"data":    map[string]string{"token": tokenString},
		},
	)

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})

}

func ResetPassword(c *gin.Context) {
	var body struct {
		Email string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse body",
		})
		return
	}
	user := modules.Users{Email: body.Email}
	initializers.SendEmail(user.Email)
	c.JSON(http.StatusOK, gin.H{
		"status": "Please check your mail and reset the password",
	})
}

func GetUserInfo(c *gin.Context) {
	tokenstring := c.Request.FormValue("token")
	fmt.Printf("tokenstring", tokenstring)

	token, _ := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepected method", token.Header)
		}
		return []byte(os.Getenv("SECRET")), nil

	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["expr"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invaid token",
			})
			return
		}

		var user modules.Users
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invaid token",
			})
			return
		}
		var roles = [1]string{user.Type}
		//	var answer = modules.Userinfo{Roles: roles, Introduction: "", Avatar: "", Name: ""}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string][1]string{
				"roles": roles,
			},
		})

	}
}

func Logout(c *gin.Context) {
	tokenstring, err := c.Cookie("Authorization")
	if err != nil {
		c.SetCookie("token", tokenstring, -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}

}
