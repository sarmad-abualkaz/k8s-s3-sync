package calls3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	
	log "github.com/sirupsen/logrus"
)


func WriteFile(sess *session.Session, bucket string, profile string, region string, body []byte, key string) error {

	// // Set AWS session
	// sess, err := SetupSession(profile, region)

	// // Check for errors from AWS session setup
	// if err != nil {
	// 	return err
	// }
	
	// Successful AWS session setup
	
	// pass session to s3-clinet
	s3Client := s3.New(sess)
	
	//setup a ListObjectsInput
	input := &s3.PutObjectInput{
		Bucket:  aws.String(bucket),
		Key: aws.String(key),
		Body: aws.ReadSeekCloser(bytes.NewReader(body)),
	}
	
	log.WithFields(log.Fields{
		"profile": profile,
		"bucket": bucket,
		"region": region,
		"key": key,
	}).Info("Writing object to S3 ...")

	// call s3 for objects
	_, s3Err := s3Client.PutObject(input)

	//
	if s3Err != nil {
		return s3Err
	}
	
	log.WithFields(log.Fields{
		"profile": profile,
		"bucket": bucket,
		"region": region,
		"key": key,
	}).Info("Succefully put object to S3 ...")

	return nil

}
