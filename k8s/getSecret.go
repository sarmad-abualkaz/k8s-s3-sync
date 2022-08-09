package callk8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ReturnTlsCert(k8sClient *kubernetes.Clientset, secret string, namespace string) ([]byte, []byte, error) {

	secrcontent, err := k8sClient.CoreV1().Secrets(namespace).Get(context.TODO(), secret, metav1.GetOptions{})


	if err != nil {
		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"namespace": namespace,
		"kind": "secret",
		"name": secret,
	}).Info("secret found ...")


	cert := secrcontent.Data["tls.crt"]
	key := secrcontent.Data["tls.key"]
	
	return cert, key, nil
}
