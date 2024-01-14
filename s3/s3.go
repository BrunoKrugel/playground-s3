package s3

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader is a struct for uploading files to S3
type S3Uploader struct {
	S3Client *s3.S3
	Bucket   string
}

// NewS3Uploader creates a new S3Uploader
func NewS3Uploader(accessKey, secretKey, region, bucket, endpoint string) *S3Uploader {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	config := &aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
		Endpoint:    aws.String(endpoint),
	}

	S3Client, _ := session.NewSession()

	return &S3Uploader{
		S3Client: s3.New(S3Client, config),
		Bucket:   bucket,
	}
}

// DeleteImage deletes an image from S3
func (u *S3Uploader) DeleteImage(fileName string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(fileName),
	}

	_, err := u.S3Client.DeleteObject(params)
	return err
}

// GetImage gets the content of an image from S3
func (u *S3Uploader) GetImage(fileName string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(fileName),
	}

	resp, err := u.S3Client.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// UploadImageByURL uploads an image to S3 using its URL
func (u *S3Uploader) UploadImageByURL(imageURL, fileName string) (string, error) {
	response, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image from URL: %s", imageURL)
	}

	return u.UploadImage(response.Body, fileName)
}

// UploadImage uploads an image to S3 and returns the URL
func (u *S3Uploader) UploadImage(imageReader io.Reader, fileName string) (string, error) {
	imageData, err := io.ReadAll(imageReader)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(imageData)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(u.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(fileType),
	}

	_, err = u.S3Client.PutObject(params)
	if err != nil {
		return "", err
	}

	// Assuming your S3 bucket is configured for public read access
	url := fmt.Sprintf("https://pub-90fb5d16307649bf89d71e5328ddc51c.r2.dev/%s", fileName)

	return url, nil
}
