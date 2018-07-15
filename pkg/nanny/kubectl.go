package nanny

import (
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (n Nanny) configureKube() error {
	fmt.Println("configure kube for storage acc")
	var config *rest.Config
	var err error
	if len(*kubeConfig) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig)
		if err != nil {
			return err
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	log.Println("update secret ")
	secret, err := clientset.CoreV1().Secrets(os.Getenv("MY_POD_NAMESPACE")).Get(n.EtcdOperatorConfig.ABSSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	secretCopy := secret.DeepCopy()
	s := make(map[string]string)
	s["storage-key"] = n.EtcdOperatorConfig.StorageKey
	secretCopy.StringData = s
	_, err = clientset.CoreV1().Secrets(os.Getenv("MY_POD_NAMESPACE")).Update(secretCopy)

	return err
}
