package calls3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	
	log "github.com/sirupsen/logrus"
)

func ReadFile(sess *session.Session, bucket string, profile string, region string, key string) ([]byte, error) {
	
	// pass session to s3-clinet
	s3Client := s3.New(sess)
	
	buffBody := new(bytes.Buffer)

	//setup a ListObjectsInput
	input := &s3.GetObjectInput{
		Bucket:  aws.String(bucket),
		Key: aws.String(key),
	}
	
	log.WithFields(log.Fields{
		"profile": profile,
		"bucket": bucket,
		"region": region,
		"key": key,
	}).Info("Getting object from S3 ...")

	// call s3 for objects
	resutl, s3Err := s3Client.GetObject(input)

	//
	if s3Err != nil {
		return nil, s3Err
	}
	
	log.WithFields(log.Fields{
		"profile": profile,
		"bucket": bucket,
		"region": region,
		"key": key,
	}).Info("Succefully read object to S3 ...")
	
	buffBody.ReadFrom(resutl.Body)

	return buffBody.Bytes(), nil

}
