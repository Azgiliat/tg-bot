package aws

import (
	"awesomeProject/internal/config"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/http"
)

type S3Bucket struct {
	buketPublicURL string
	bucketName     string
	s3Client       *s3.Client
}

func newClient(bucketName string, bucketPublicURL string) *S3Bucket {
	awsConfig := config.GetAWSConfig()
	creds := credentials.NewStaticCredentialsProvider(awsConfig.Key, awsConfig.Secret, "")
	conf := aws.Config{
		Region:      awsConfig.Region,
		Credentials: creds,
	}
	client := s3.NewFromConfig(conf)

	return &S3Bucket{bucketPublicURL, bucketName, client}
}

func (s3Bucket *S3Bucket) UploadImage(fileName string, file io.ReadSeeker) error {
	contentTypeBuff := make([]byte, 512)
	bytesRead, err := file.Read(contentTypeBuff)

	if err != nil {
		return err
	}

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	contentTypeBuff = contentTypeBuff[:bytesRead]
	_, err = s3Bucket.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s3Bucket.bucketName),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(http.DetectContentType(contentTypeBuff)),
	})

	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Image uploaded")
		return nil
	}
}

func (s3Bucket *S3Bucket) GenerateImageURL(imageName string) string {
	return fmt.Sprintf("%s/%s", s3Bucket.buketPublicURL, imageName)
}
