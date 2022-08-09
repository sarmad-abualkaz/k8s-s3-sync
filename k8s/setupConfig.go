package callk8s

import (
	"flag"
	"path/filepath"
	
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/rest"
	log "github.com/sirupsen/logrus"

)

func SetupK8sClient(location string) (*kubernetes.Clientset, error) {

	var configError error
	var config *rest.Config

	log.WithFields(log.Fields{
		"config_location": location,
	  }).Info("Creating kubeconfig ...")

	if location == "local" {

		var kubeconfig *string

		if home := homedir.HomeDir(); home != "" {

			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

		} else {

			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

		}

		flag.Parse()

		// use the current context in kubeconfig
		config, configError = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		

	} else if location == "in-cluster" {

		config, configError = rest.InClusterConfig()

	}

	if configError != nil {
		return nil, configError
	}

	// create the clientset
	clientset, clientsetErr := kubernetes.NewForConfig(config)
	if clientsetErr != nil {
		return nil, clientsetErr
	}

	return clientset, nil
}
