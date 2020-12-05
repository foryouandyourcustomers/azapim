package disasterrecovery

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// StorageAccount is used to retrieve the storage access keys and check for the defined blob storage
type StorageAccount struct {
	Ctx           context.Context
	AccountClient storage.AccountsClient
	Subscription  string
	ResourceGroup string
	AccountName   string
	BlobName      string
	Key           string
}

// InitializeClient inializes the storage account client
func (sc *StorageAccount) InitializeClient(s string) {
	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}
	sc.Subscription = s
	sc.AccountClient = storage.NewAccountsClient(sc.Subscription)
	sc.AccountClient.Authorizer = a
	sc.Ctx = context.Background()
	sc.getKey()
}

func (sc *StorageAccount) getKey() {

	k, err := sc.AccountClient.ListKeys(sc.Ctx, sc.ResourceGroup, sc.AccountName, "kerb")
	if err != nil {
		log.Fatalln(err)
	}

	sc.Key = *(*k.Keys)[0].Value
}
