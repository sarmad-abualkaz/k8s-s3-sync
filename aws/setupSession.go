package calls3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"

)

func SetupSession(profile string, region string) (*session.Session, error) {

		log.WithFields(log.Fields{
			"profile": profile,
			"region": region,
		}).Info("Creating aws session ...")

		// Set AWS session
		sess, err := session.NewSessionWithOptions(session.Options{
			// Specify profile to load for the session's config
			Profile: profile,
			SharedConfigState: session.SharedConfigEnable,
		
			// Provide SDK Config options, such as Region.
			Config: aws.Config{
				Region: aws.String(region),
			},
		})

		if err != nil {
			return nil, err
		}

		return sess, nil
}
