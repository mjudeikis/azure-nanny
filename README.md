# Testing

```
export AZURE_SUBSCRIPTION_ID=zzzzzzz
export AZURE_CLIENT_ID=zzzzzzzzz
export AZURE_CLIENT_SECRET=zzzzzzzzzzzzz
export AZURE_TENANT_ID=zzzzzzzzzzz
export RESOURCE_GROUP=mjudeikis-nanny 
export CONTAINER_NAME=etcdbackup 
export LOCATION=eastus 
export STORAGE_ACCOUNT=etcdbackup 
export MY_POD_NAMESPACE=mjudeikis-nanny

go run cmd/azure-nanny/main.go --kubeconfig /home/mjudeiki/.aks/admin.kubeconfig
```
