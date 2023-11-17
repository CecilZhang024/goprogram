package main

import (
	"awesomeProject1/controller"
	"awesomeProject1/initializers"
	"awesomeProject1/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDb()
	initializers.SyncDB()
}
func main() {
	r := gin.Default()

	r.POST("/user/signup", controller.Signup)

	r.POST("/user/login", controller.Login)
	r.GET("/user/validate", middleware.RequireAuth, controller.Validate)
	r.POST("/file/uploadfile", middleware.RequireAuth, controller.Upload)
	r.POST("/user/resetpwd", controller.ResetPassword)
	r.POST("/file/uploadfileS3", middleware.RequireAuth, controller.UploadtoS3)
	r.GET("/user/info", controller.GetUserInfo)
	r.POST("/gpt/askgpt", middleware.RequireAuth, controller.PostToGPT)
	r.POST("/user/logout", controller.Logout)
	r.Run()
}
