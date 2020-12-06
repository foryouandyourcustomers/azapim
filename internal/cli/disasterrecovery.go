package cli

import (
	"fmt"
	"time"

	"github.com/foryouandyourcustomers/azapim/internal/apimclient"
	ucli "github.com/urfave/cli/v2"
)

var (
	dr apimclient.DisasterRecovery

	// DisasterRecoveryCli contains the disaster recovery cli
	DisasterRecoveryCli = []*ucli.Command{
		{
			Name:     "disasterrecovery",
			Category: "Management",
			Usage:    "Create APIM disaster recovery backups or restore from them",
			Flags: []ucli.Flag{
				&ucli.StringFlag{
					Name:        "storageaccount",
					Usage:       "the storage account where the disaster recovery backup is stored",
					Required:    true,
					EnvVars:     []string{"STORAGEACCOUNT"},
					Destination: &dr.Storage.AccountName,
				},
				&ucli.StringFlag{
					Name:        "storageaccountrg",
					Usage:       "the storage account resource group",
					Required:    true,
					EnvVars:     []string{"STORAGEACCOUNTRG"},
					Destination: &dr.Storage.ResourceGroup,
				},
				&ucli.StringFlag{
					Name:        "blobname",
					Usage:       "the blob container containing the api management backups",
					Required:    true,
					EnvVars:     []string{"BLOBNAME"},
					Destination: &dr.Storage.BlobName,
				},
				&ucli.StringFlag{
					Name:        "backupname",
					Usage:       "the name of the backup to create or restore from",
					Required:    false,
					EnvVars:     []string{"BACKUPNAME"},
					Destination: &dr.BackupName,
				},
			},
			Subcommands: []*ucli.Command{
				{
					Name:  "backup",
					Usage: "Backup the api management service",
					Action: func(c *ucli.Context) error {
						err := apimClient.Backup(apimClient.ResourceGroup, apimClient.ServiceName, dr.Parameters)
						if err != nil {
							return ucli.Exit("Unable to create backup for service", 1)
						}
						return nil
					},
				},
				{
					Name:  "restore",
					Usage: "Resttore the api management service",
					Action: func(c *ucli.Context) error {
						err := apimClient.Restore(apimClient.ResourceGroup, apimClient.ServiceName, dr.Parameters)
						if err != nil {
							return ucli.Exit("Unable to restore service", 1)
						}
						return nil
					},
				},
			},
			Before: func(c *ucli.Context) error {
				if len(dr.BackupName) == 0 {
					dr.BackupName = fmt.Sprintf("%s-%d", apimClient.ServiceName, time.Now().Unix())
				}
				dr.Initialize(apimClient.Subscription)
				return nil
			},
		},
	}
)
