package nanny

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	azstorage "github.com/Azure/azure-sdk-for-go/storage"
	blob "github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	envList = []string{"AZURE_SUBSCRIPTION_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET", "AZURE_TENANT_ID", "RESOURCE_GROUP", "CONTAINER_NAME", "LOCATION", "STORAGE_ACCOUNT"}
)

func (n *Nanny) configureAzure() error {

	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	n.Clients.Accounts = storage.NewAccountsClient(subscriptionID)
	n.Clients.Accounts.Authorizer = authorizer
	n.Clients.Groups = resources.NewGroupsClient(subscriptionID)
	n.Clients.Groups.Authorizer = authorizer
	n.Clients.Images = compute.NewImagesClient(subscriptionID)
	n.Clients.Images.Authorizer = authorizer

	// I should not be here
	if err := n.createResourceGroup(); err != nil {
		return err
	}
	if _, err := n.createStorageAccount(); err != nil {
		return err
	}

	keys, err := n.Clients.Accounts.ListKeys(context.Background(), n.EtcdOperatorConfig.ResourceGroup, n.EtcdOperatorConfig.BackupStorageAccount)
	if err != nil {
		return err
	}

	n.EtcdOperatorConfig.StorageKey = *(*keys.Keys)[0].Value

	n.Clients.Storage, err = azstorage.NewClient(n.EtcdOperatorConfig.BackupStorageAccount, *(*keys.Keys)[0].Value, azstorage.DefaultBaseURL, azstorage.DefaultAPIVersion, true)
	if err != nil {
		return err
	}

	if err := n.createContainer(context.Background(), *(*keys.Keys)[0].Value); err != nil {
		return err
	}

	return nil
}

func (n Nanny) createResourceGroup() error {

	rg := resources.Group{
		Location: to.StringPtr(n.EtcdOperatorConfig.Location),
	}

	_, err := n.Clients.Groups.CreateOrUpdate(context.Background(), n.EtcdOperatorConfig.ResourceGroup, rg)
	if err != nil {
		return err
	}
	return nil
}

func (n Nanny) createStorageAccount() (s storage.Account, err error) {

	storageAccountsClient := n.Clients.Accounts

	result, err := n.Clients.Accounts.CheckNameAvailability(
		context.Background(),
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(n.EtcdOperatorConfig.BackupStorageAccount),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		log.Printf("%s: %v", "storage account creation failed", err)
	}
	if *result.NameAvailable != true {
		log.Printf("%s [%s]: %v: %v", "storage account name not available", n.EtcdOperatorConfig.BackupStorageAccount, err, *result.Message)
		return s, nil
	}

	future, err := n.Clients.Accounts.Create(
		context.Background(),
		n.EtcdOperatorConfig.ResourceGroup,
		n.EtcdOperatorConfig.BackupStorageAccount,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Kind:     storage.Storage,
			Location: to.StringPtr(n.EtcdOperatorConfig.Location),
			AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{},
		})
	if err != nil {
		return s, fmt.Errorf("cannot create storage account: %v", err)
	}

	err = future.WaitForCompletion(context.Background(), n.Clients.Accounts.Client)
	if err != nil {
		return s, fmt.Errorf("cannot get the storage account create future response: %v", err)
	}

	return future.Result(storageAccountsClient)
}

func (n Nanny) createContainer(ctx context.Context, key string) error {

	c := blob.NewSharedKeyCredential(n.EtcdOperatorConfig.BackupStorageAccount, key)
	p := blob.NewPipeline(c, blob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", n.EtcdOperatorConfig.BackupStorageAccount, n.EtcdOperatorConfig.ContainerName))

	containerURL := blob.NewContainerURL(*u, p)

	// Create the container
	fmt.Printf("Creating a container named %s\n", n.EtcdOperatorConfig.ContainerName)
	_, err := containerURL.Create(context.Background(), blob.Metadata{}, blob.PublicAccessNone)

	if strings.Contains(err.Error(), "The specified container already exists.") {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func checkEnvConfig() error {
	var errstrings []string
	for _, env := range envList {
		_, ok := os.LookupEnv(env)
		if ok == false {
			errstrings = append(errstrings, fmt.Sprintf("%s not set", env))
		}
	}
	if len(errstrings) > 0 {
		return fmt.Errorf(strings.Join(errstrings, "\n"))
	}
	return nil
}
