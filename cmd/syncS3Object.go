package cmd

import (
	"crypto/sha256"
	"fmt"
	"time"
	
	log "github.com/sirupsen/logrus"
	calls3 "github.com/sarmad-abualkaz/k8s-s3-sync/aws"
	callk8s "github.com/sarmad-abualkaz/k8s-s3-sync/k8s"
)

func SyncS3Objects(awsProfile string, awsRegion string, certObject string, certKeyObject string, k8sconfiglocation string, kind string, namespace string, s3buecktName string, secret string, sleeping int) {

	k8sclient, k8sclientErr := callk8s.SetupK8sClient(k8sconfiglocation)

	if k8sclientErr != nil {

		log.WithFields(log.Fields{
			"location": k8sconfiglocation,
			"Error": k8sclientErr.Error(),
		}).Fatal("Error setting up kubernetes client ...")

		panic(k8sclientErr.Error())
	}

	awsClient, awsClientErr := calls3.SetupSession(awsProfile, awsRegion)

	if awsClientErr != nil {

		log.WithFields(log.Fields{
			"profile": awsProfile,
			"region": awsRegion,
		}).Fatal("Error setting up aws client ...")

		panic(awsClientErr.Error())
	}
	
	for {
		
		tlsCert, tlsKey, tlsCertErr := callk8s.ReturnTlsCert(k8sclient, secret, namespace)

		if tlsCertErr != nil {

			log.WithFields(log.Fields{
				"namespace": namespace,
				"kind": kind,
				"name": secret,
			}).Error("secret not found ...")

		} else {
			
			hashedCert := sha256.New()
			hashedKey := sha256.New()
	
			hashedCert.Write([]byte(tlsCert))
			hashedKey.Write([]byte(tlsKey))
	
			s3Objects, s3Err:= calls3.ListFiles(awsClient, s3buecktName)
			
			if s3Err != nil {
				log.WithFields(log.Fields{
					"profile": awsProfile,
					"bucket": s3buecktName,
					"region": awsRegion,
					"Error": s3Err,
				}).Error("Error retreiving s3 bucket objects ...")
			
			} else {

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

						hashedS3Cert := sha256.New()
						hashedS3Cert.Write([]byte(s3CertVal))
			
						if fmt.Sprintf("%x",hashedS3Cert.Sum(nil)) != fmt.Sprintf("%x", hashedCert.Sum(nil)) {
							
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certObject,
								"namespace": namespace,
								"kind": kind,
								"name": secret,
							}).Warn("Mismatch of cert values - updating s3 object ...")
			
							s3UodateCertErr:= calls3.WriteFile(awsClient, s3buecktName, awsProfile, awsRegion, tlsCert, certObject)
			
							if s3UodateCertErr != nil {
			
								log.WithFields(log.Fields{
									"profile": awsProfile,
									"bucket": s3buecktName,
									"region": awsRegion,
									"key": certObject,
								}).Error("Failed to update object ...")
							}

						} else {
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certObject,
								"namespace": namespace,
								"kind": kind,
								"name": secret,
							}).Info("Cert values are a match - s3 object matches secret in kubernetes. No update required ...")
						}
					}
		
				} else {
		
					log.WithFields(log.Fields{
						"profile": awsProfile,
						"bucket": s3buecktName,
						"region": awsRegion,
						"key": certObject,
						}).Warn("Cert object Not found - sending request to create it ...")
		
					s3WriteCertErr:= calls3.WriteFile(awsClient, s3buecktName, awsProfile, awsRegion, tlsCert, certObject)
		
					if s3WriteCertErr != nil {
		
						log.WithFields(log.Fields{
							"profile": awsProfile,
							"bucket": s3buecktName,
							"region": awsRegion,
							"key": certObject,
						}).Error("Failed to write object ...")
					}
				}
				
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

						hashedS3CertKey := sha256.New()
						hashedS3CertKey.Write([]byte(s3CertKeyVal))
			
						if fmt.Sprintf("%x", hashedS3CertKey.Sum(nil)) != fmt.Sprintf("%x", hashedKey.Sum(nil)) {
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certKeyObject,
								"namespace": namespace,
								"kind": kind,
								"name": secret,
							}).Warn("Mismatch of key cert values - updating s3 object ...")
						
							s3UodateCertKeyErr := calls3.WriteFile(awsClient, s3buecktName, awsProfile, awsRegion, tlsKey, certKeyObject)
			
							if s3UodateCertKeyErr != nil {
			
								log.WithFields(log.Fields{
									"profile": awsProfile,
									"bucket": s3buecktName,
									"region": awsRegion,
									"key": certKeyObject,
								}).Error("Failed to update object ...")
							}
			
						} else {
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certKeyObject,
								"namespace": namespace,
								"kind": kind,
								"name": secret,
							}).Info("Cert key values are a match - s3 object matches secret in kubernetes. No update required ...")
						}
						
					}
		
				} else {
		
					log.WithFields(log.Fields{
						"profile": awsProfile,
						"bucket": s3buecktName,
						"region": awsRegion,
						"key": certKeyObject,
						}).Warn("Cert key object Not found - sending request to create it ...")
		
						s3WriteCertKeyErr:= calls3.WriteFile(awsClient, s3buecktName, awsProfile, awsRegion, tlsKey, certKeyObject)
		
						if s3WriteCertKeyErr != nil {
			
							log.WithFields(log.Fields{
								"profile": awsProfile,
								"bucket": s3buecktName,
								"region": awsRegion,
								"key": certKeyObject,
							}).Error("Failed to write object ...")
						}
				}

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
