package util

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Svc        *s3.S3
	BucketName string
}

func NewS3(AccessKeyId, SecretAccessKey, Region, BucketName string) (*S3, error) {
	creds := credentials.NewStaticCredentials(AccessKeyId, SecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
		return nil, err
	}
	cfg := aws.NewConfig().WithRegion(Region).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	return &S3{Svc: svc, BucketName: BucketName}, nil
}
