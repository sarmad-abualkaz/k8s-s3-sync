package main

import (
	"flag"
	"github.com/sarmad-abualkaz/k8s-s3-sync/cmd"

	log "github.com/sirupsen/logrus"
)

func main() {

	awsProfile := flag.String("aws-profile", "tools", "name of aws profile")
	awsRegion := flag.String("aws-region", "us-east-1", "aws region")
	certObject :=  flag.String("cert-object-name", "cert.text", "name of object for cert data")
	certKeyObject := flag.String("cert-key-object-name", "key.text", "name of object for cert key data")
	k8sConfig := flag.String("kube-config", "in-cluster", "kubeconfig setup")
	kind := flag.String("kind", "secret", "kubernetes kind type - only secret works atm")
	namespace := flag.String("namespace", "local", "kubernetes namespace where secret exists")
	s3buecktName := flag.String("s3-bucket-name", "", "Name of s3 bucket where objects are synced")
	secretName := flag.String("secret-name", "local-certificate-secret", "kubernetes secret name to sync from/to")
	sleeping := flag.Int("sleeping", 20, "sleep time between syncs")
	syncS3 := flag.Bool("sync-s3", false, "should sync s3 from secret")
	
	flag.Parse()

	log.WithFields(log.Fields{
		"aws-profile": *awsProfile,
		"aws-region": *awsRegion,
		"cert-object-name": *certObject,
		"cert-key-object-name": *certKeyObject,
		"kube-config": *k8sConfig,
		"kind": *kind,
		"namespace": *namespace,
		"s3-bucket-name": *s3buecktName,
		"secret-name": *secretName,
		"sleeping": *sleeping,
		"sync-s3": *syncS3,
	  }).Info("program started ...")

	if *syncS3 {

		log.WithFields(log.Fields{
			"namespace": *namespace,
			"kind": *kind,
			"name": *secretName,
		  }).Info("program will aim to sync s3 bucket objects from a kubernetes secret ...")

		cmd.SyncS3Objects(*awsProfile, *awsRegion, *certObject, *certKeyObject, *k8sConfig, *kind, *namespace, *s3buecktName, *secretName, *sleeping)

	} else {

		log.WithFields(log.Fields{
			"s3-bucket-name": *s3buecktName,
			"cert-object-name": *certObject,
			"cert-key-object-name": *certKeyObject,
		}).Info("program will aim sync a kubernetes secret from s3 bucket object ...")
		
		cmd.SyncSecret(*awsProfile, *awsRegion, *certObject, *certKeyObject, *k8sConfig, *kind, *namespace, *s3buecktName, *secretName, *sleeping)
	}

}
