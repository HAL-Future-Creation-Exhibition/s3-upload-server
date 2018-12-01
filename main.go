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
	"image/png"
	"image/jpeg"
	"image/gif"
	"image"
	"os/exec"
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
	r.POST("/upload/icon", func(c *gin.Context) {
		myS3, err := util.NewS3(os.Getenv("S3ACCESSKEY"), os.Getenv("S3SECKEY"), os.Getenv("REGION"), os.Getenv("BUCKETNAME"))
		if err != nil {
			panic(err)
		}
		fmt.Println(myS3)

		path := c.DefaultQuery("path", "")

		icon, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		file, err := icon.Open()
		defer file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "fileOpenに失敗",
			})
			return
		}

		_, format, err := image.DecodeConfig(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "画像情報取得に失敗",
			})
			return
		}

		file, err = icon.Open()
		defer file.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "fileOpenに失敗",
			})
			return
		}

		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "画像デコードに失敗",
			})
			return
		}

		tmpf, err := os.Create("./tmp/" + "tmp-" + path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "致命的エラー",
			})
			return
		}
		defer tmpf.Close()
		switch format {
		case "png":
			err = png.Encode(tmpf, img)
		case "jpeg":
			err = jpeg.Encode(tmpf, img, nil)
		case "gif":
			err = gif.Encode(tmpf, img, nil)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "致命的エラー | 画像保存に失敗",
			})
			return
		}

		err = exec.Command("python3", "main.py", "./tmp/" + "tmp-" + path, "./tmp/" + path).Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "致命的エラー | 画像保存に失敗",
			})
			return
		}

		sendFile, _ := os.Open("./tmp/" + path)
		defer file.Close()

		params := &s3.PutObjectInput{
			Bucket: aws.String(myS3.BucketName),  // Required
			Key:    aws.String("nellow/" + path), // Required
			ACL:    aws.String("public-read"),
			Body:   sendFile,
		}


		resp, err := myS3.Svc.PutObject(params)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
		fmt.Println(params)

		err = exec.Command("rm", "./tmp/" + "tmp-" + path, "./tmp/" + path).Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "致命的エラー | 画像お掃除に失敗",
			})
			return
		}

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
