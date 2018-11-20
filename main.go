package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/HAL-Future-Creation-Exhibition/s3-upload-server/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	Env_load()
	fmt.Println(os.Getenv("HOGE"))
	r := gin.Default()
	// web
	r.Use(cors)
	r.GET("/alive", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ヤッホー!!",
		})
	})
	r.POST("/upload", func(c *gin.Context) {
		myS3, err := util.NewS3(os.Getenv("S3ACCESSKEY"), os.Getenv("S3SECKEY"), os.Getenv("REGION"), os.Getenv("BUCKETNAME"))
		if err != nil {
			panic(err)
		}
		fmt.Println(myS3)

		path := c.DefaultQuery("path", "")

		file, _, err := c.Request.FormFile("file")
		fmt.Println(file)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		params := &s3.PutObjectInput{
			Bucket: aws.String(myS3.BucketName),  // Required
			Key:    aws.String("nellow/" + path), // Required
			ACL:    aws.String("public-read"),
			Body:   file,
		}

		resp, err := myS3.Svc.PutObject(params)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
		fmt.Println(params)

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func cors(c *gin.Context) {
	headers := c.Request.Header.Get("Access-Control-Request-Headers")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,HEAD,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", headers)
	if c.Request.Method == "OPTIONS" {
		c.Status(200)
		c.Abort()
	}
	c.Set("start_time", time.Now())
	c.Next()

}
