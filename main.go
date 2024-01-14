package main

import (
	"fmt"
	"os"
	"playground-s3/s3"
)

func main() {

	// Replace these with your own values
	accessKey := ""
	secretKey := ""
	region := "auto"
	bucket := "bucket-test"
	endpoint := ""

	uploader := s3.NewS3Uploader(accessKey, secretKey, region, bucket, endpoint)

	// imageData := []byte("your-image-data")
	imagePath := "frog.png"
	fileName := "your-frog.jpg"

	file, _ := os.Open(imagePath)
	defer file.Close()

	url, err := uploader.UploadImage(file, fileName)
	fmt.Printf("Uploaded image URL: %s\n", url)
	fmt.Println(err)

}
