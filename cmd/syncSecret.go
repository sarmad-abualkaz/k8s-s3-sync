package cmd

import (
	"crypto/sha256"
	"fmt"
	"time"
	
	log "github.com/sirupsen/logrus"
	calls3 "github.com/sarmad-abualkaz/k8s-s3-sync/aws"
	callk8s "github.com/sarmad-abualkaz/k8s-s3-sync/k8s"

	"k8s.io/apimachinery/pkg/api/errors"
)

func SyncSecret(awsProfile string, awsRegion string, certObject string, certKeyObject string, k8sconfiglocation string, kind string, namespace string, s3buecktName string, secret string, sleeping int) {

	awsClient, awsClientErr := calls3.SetupSession(awsProfile, awsRegion)

	if awsClientErr != nil {

		log.WithFields(log.Fields{
			"profile": awsProfile,
			"region": awsRegion,
		}).Fatal("Error setting up aws client ...")

		panic(awsClientErr.Error())
	}

	k8sclient, k8sclientErr := callk8s.SetupK8sClient(k8sconfiglocation)

	if k8sclientErr != nil {

		log.WithFields(log.Fields{
			"location": k8sconfiglocation,
			"Error": k8sclientErr.Error(),
		}).Fatal("Error setting up kubernetes client ...")

		panic(k8sclientErr.Error())
	}

	for {

		//return objects from bucket:
		s3Objects, s3Err:= calls3.ListFiles(awsClient, s3buecktName)
		
		if s3Err != nil {
			log.WithFields(log.Fields{
				"profile": awsProfile,
				"bucket": s3buecktName,
				"region": awsRegion,
				"Error": s3Err,
			}).Error("Error retreiving s3 bucket objects ...")
		
		} else {

					// look for cert object and store value
			if _, certFound := s3Objects[certObject]; certFound {
				
				log.WithFields(log.Fields{
					"profile": awsProfile,
					"bucket": s3buecktName,
					"region": awsRegion,
					"key": certObject,
				}).Info("Cert object found ...")

				log.WithFields(log.Fields{
					"profile": awsProfile,
					"bucket": s3buecktName,
					"region": awsRegion,
					"key": certObject,
				}).Info("Retreiving cert object value ...")
				
				s3CertVal, s3CertErr := calls3.ReadFile(awsClient, s3buecktName, awsProfile, awsRegion, certObject)

				if s3CertErr != nil {

					log.WithFields(log.Fields{
						"profile": awsProfile,
						"bucket": s3buecktName,
						"region": awsRegion,
						"key": certObject,		
					}).Error("Error retreiving s3 cert objects value ...")
				
				} else {

					// initialize hashed values for cert 
					hashedS3Cert := sha256.New()
					hashedS3Cert.Write([]byte(s3CertVal))

					// look for key object and store value
					if _, certKeyFound := s3Objects[certKeyObject]; certKeyFound {
				
						log.WithFields(log.Fields{
							"profile": awsProfile,
							"bucket": s3buecktName,
							"region": awsRegion,
							"key": certKeyObject,
						}).Info("Cert key object found ...")
			
						log.WithFields(log.Fields{
							"profile": awsProfile,
							"bucket": s3buecktName,
							"region": awsRegion,
							"key": certKeyObject,
						}).Info("Retreiving cert key object value ...")
						
						s3CertKeyVal, s3CertKeyErr := calls3.ReadFile(awsClient, s3buecktName, awsProfile, awsRegion, certKeyObject)
			
						if s3CertKeyErr != nil {
			
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certKeyObject,		
							}).Error("Error retreiving s3 cert key objects value ...")
						
						} else {

							// initialize hashed values for key 
							hashedS3CertKey := sha256.New()
							hashedS3CertKey.Write([]byte(s3CertKeyVal))

							// retreive secret from k8s and store values of cert and key

							tlsCert, tlsKey, tlsCertErr := callk8s.ReturnTlsCert(k8sclient, secret, namespace)

							if errors.IsNotFound(tlsCertErr) {

								log.WithFields(log.Fields{
									"namespace": namespace,
									"kind": kind,
									"name": secret,
								}).Warn("secret not found ...")

								// create secret if does not exist
								createSecretErr := callk8s.WriteTlsCert(k8sclient, secret , namespace , s3CertVal, s3CertKeyVal, false)

								if createSecretErr != nil {

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Error("Unable to create Kubernetes secrets ...", createSecretErr)
								}
					
							} else if tlsCertErr != nil {

								log.WithFields(log.Fields{
									"namespace": namespace,
									"kind": kind,
									"name": secret,
								}).Error("Secret cannot be retrieved: ", tlsCertErr)
								
							} else {
								
								// if secret exists check if cert matches

								hashedCert := sha256.New()
								hashedKey := sha256.New()
						
								hashedCert.Write([]byte(tlsCert))
								hashedKey.Write([]byte(tlsKey))

								// update if s3 does not match secret
								if fmt.Sprintf("%x",hashedS3Cert.Sum(nil)) != fmt.Sprintf("%x", hashedCert.Sum(nil)) {

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Warn("Mismatch between cert value in s3 and secret in k8s ...")

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Warn("Updating secret in k8s ...")
									
									updateSecretErr := callk8s.WriteTlsCert(k8sclient, secret , namespace , s3CertVal, s3CertKeyVal, true)

									if updateSecretErr != nil {
		
										log.WithFields(log.Fields{
											"namespace": namespace,
											"kind": kind,
											"name": secret,
										}).Error("Unable to create Kubernetes secrets ...", updateSecretErr)
									}
								
								// update again only if key is different
								// check if key matches
								} else if fmt.Sprintf("%x", hashedS3CertKey.Sum(nil)) != fmt.Sprintf("%x", hashedKey.Sum(nil)) {

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Warn("Mismatch between cert key value in s3 and secret in k8s ...")

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Warn("Updating secret in k8s ...")
									
									updateSecretErr := callk8s.WriteTlsCert(k8sclient, secret , namespace , s3CertVal, s3CertKeyVal, true)

									// update if not
									if updateSecretErr != nil {
		
										log.WithFields(log.Fields{
											"namespace": namespace,
											"kind": kind,
											"name": secret,
										}).Error("Unable to create Kubernetes secrets ...", updateSecretErr)
									}
								} else {

									log.WithFields(log.Fields{
										"namespace": namespace,
										"kind": kind,
										"name": secret,
									}).Info("Cert key values are a match - s3 object matches secret in kubernetes. No update required ...")
								}

							}

						}

					// error if object does not exist
					} else {
						log.WithFields(log.Fields{
							"profile": awsProfile,
							"bucket": s3buecktName,
							"region": awsRegion,
							"key": certKeyObject,		
						}).Error("Could not find required S3 bucket object ...")

					}
				}

			// error if object does not exist				
			} else {

				log.WithFields(log.Fields{
					"profile": awsProfile,
					"bucket": s3buecktName,
					"region": awsRegion,
					"key": certObject,		
				}).Error("Could not find required S3 bucket object ...")

			}
		}
		
		log.WithFields(log.Fields{
			"namespace": namespace,
			"kind": kind,
			"name": secret,
		}).Info("sleeping for ", sleeping, " seconds")

		time.Sleep(time.Duration(sleeping) * time.Second)
	}
}
