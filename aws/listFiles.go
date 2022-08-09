package calls3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"

)

// Find all files in bucket
func ListFiles(sess *session.Session, bucket string) (map[string]string, error) {
	
	// pass session to s3-clinet
	s3Client := s3.New(sess)
	
	// Declare empty object map
	objectMap := make(map[string]string)
	
	
	//setup a ListObjectsInput
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(2),
	}
	
	// call s3 for objects
	objectsRes, s3Err := s3Client.ListObjects(input)

	//
	if s3Err != nil {
		return objectMap, s3Err
	}
	
	for _, object := range objectsRes.Contents {

		// objectList = append(objectList, *object.Key)
		objectMap[*object.Key] = "found"
    }

	return objectMap, nil

}
