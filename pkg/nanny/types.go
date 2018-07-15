package nanny

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	azstorage "github.com/Azure/azure-sdk-for-go/storage"
)

type NannyClients struct {
	Accounts storage.AccountsClient
	Groups   resources.GroupsClient
	Images   compute.ImagesClient
	Storage  azstorage.Client
}

type Nanny struct {
	Clients            *NannyClients
	EtcdOperatorConfig *EtcdOperatorNanny
}

type EtcdOperatorNanny struct {
	BackupStorageAccount string
	ResourceGroup        string
	ContainerName        string
	Location             string
	StorageKey           string
	ABSSecretName        string
}
