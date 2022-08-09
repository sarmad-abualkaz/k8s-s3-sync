package callk8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func WriteTlsCert(k8sClient *kubernetes.Clientset, secret string, namespace string, cert []byte, key []byte, exists bool) error {

	var err error

	secertsClient := k8sClient.CoreV1().Secrets(namespace)

	secreManifest := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secret,
			Namespace: namespace,
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.crt": cert,
			"tls.key": key,
		},
	}

	if exists {
		
		log.WithFields(log.Fields{
			"namespace": namespace,
			"kind": "secret",
			"name": secret,
		}).Info("Updating secret ...")

		_, err = secertsClient.Update(context.TODO(), secreManifest, metav1.UpdateOptions{})
		
	
	} else {

		log.WithFields(log.Fields{
			"namespace": namespace,
			"kind": "secret",
			"name": secret,
		}).Info("Creating secret ...")

		_, err = secertsClient.Create(context.TODO(), secreManifest, metav1.CreateOptions{})

	}

	if err != nil {

		return err
	}
	
	return nil
}
