package controller

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")

	dst := "./" + file.Filename
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

}

func UploadtoS3(c *gin.Context) {
	file, _ := c.FormFile("file")
	println(file)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	println("client connected")
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("ymc024"),
		Key:    aws.String(file.Filename),
		Body:   c.Request.Body,
	})
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success!": "file uploaded",
	})
}
