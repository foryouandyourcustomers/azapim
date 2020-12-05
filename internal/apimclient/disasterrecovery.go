package apimclient

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/foryouandyourcustomers/azapim/internal/disasterrecovery"
	log "github.com/sirupsen/logrus"
)

// DisasterRecovery contains the configuration for the backupd and restore of the apim service
type DisasterRecovery struct {
	BackupName string
	Storage    disasterrecovery.StorageAccount
	Parameters apimanagement.ServiceBackupRestoreParameters
}

// Initialize inializes the storage account client
func (dr *DisasterRecovery) Initialize(s string) {
	dr.Storage.InitializeClient(s)
	dr.Parameters.AccessKey = &dr.Storage.Key
	dr.Parameters.ContainerName = &dr.Storage.BlobName
	dr.Parameters.BackupName = &dr.BackupName
	dr.Parameters.StorageAccount = &dr.Storage.AccountName
}

// Backup backups the specified api management service
func (apim *ApimClient) Backup(rg string, s string, p apimanagement.ServiceBackupRestoreParameters) error {
	log.Infof("Execute DR Backup for Service '%s' with name '%s' to storage account and blob '%s/%s'", s, p.BackupName, p.StorageAccount, p.StorageAccount)
	future, err := apim.ServiceClient.Backup(apim.Ctx, rg, s, p)
	if err != nil {
		return err
	}
	err = future.WaitForCompletionRef(apim.Ctx, apim.ServiceClient.Client)
	if err != nil {
		return err
	}
	return nil
}

// Restore disaster recovery backup for apim service
func (apim *ApimClient) Restore(rg string, s string, p apimanagement.ServiceBackupRestoreParameters) error {
	log.Infof("Execute DR Restore for Service '%s' with name '%s' to storage account and blob '%s/%s'", s, p.BackupName, p.StorageAccount, p.StorageAccount)
	future, err := apim.ServiceClient.Restore(apim.Ctx, rg, s, p)
	if err != nil {
		return err
	}
	err = future.WaitForCompletionRef(apim.Ctx, apim.ServiceClient.Client)
	if err != nil {
		return err
	}
	return nil
}
