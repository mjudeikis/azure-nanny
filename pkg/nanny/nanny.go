package nanny

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var dryRun = flag.Bool("n", false, "dry-run")
var nannyMode = flag.String("mode", "etcd-operator", "Type to nanny")
var kubeConfig = flag.String("kubeconfig", "~/.kube/config", "absolute path to kubeconfig")

func Run() error {

	flag.Parse()

	switch *nannyMode {
	case "etcd-operator":
		return runEtcdOperatorNanny()
	default:
		return fmt.Errorf("no nanny mode selected")
	}
	return nil
}

func runEtcdOperatorNanny() error {
	log.Println("start etcd-operator nanny")

	err := checkEnvConfig()
	if err != nil {
		return err
	}

	nanny := Nanny{
		EtcdOperatorConfig: &EtcdOperatorNanny{
			BackupStorageAccount: os.Getenv("STORAGE_ACCOUNT"),
			ResourceGroup:        os.Getenv("RESOURCE_GROUP"),
			Location:             os.Getenv("LOCATION"),
			ContainerName:        os.Getenv("CONTAINER_NAME"),
			ABSSecretName:        "etcd-backup-abs-credentials",
		},
		Clients: &NannyClients{},
	}
	log.Println("configure azure")
	if err := nanny.configureAzure(); err != nil {
		return err
	}
	log.Println("configure HCP")
	if err := nanny.configureKube(); err != nil {
		return err
	}

	return nil

}
