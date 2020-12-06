package cli

import (
	"context"

	"github.com/foryouandyourcustomers/azapim/internal/apimclient"
	ucli "github.com/urfave/cli/v2"
)

var (
	apimClient apimclient.ApimClient

	// Collection contains alli api commands from cli package
	Collection = append(VersionedAPICli, DisasterRecoveryCli...)

	// GlobalFlags contains the definition of all global parameters
	GlobalFlags = []ucli.Flag{
		&ucli.StringFlag{
			Name:        "subscription",
			Usage:       "Azure Subscription `ID` of the API management service",
			Required:    true,
			EnvVars:     []string{"SUBSCRIPTION"},
			Destination: &apimClient.Subscription,
		},
		&ucli.StringFlag{
			Name:        "resourcegroup",
			Usage:       "`Name` of the resource group containing the API management service",
			Required:    true,
			EnvVars:     []string{"RESOURCEGROUP"},
			Destination: &apimClient.ResourceGroup,
		},
		&ucli.StringFlag{
			Name:        "servicename",
			Usage:       "`Name` of the API management service",
			Required:    true,
			EnvVars:     []string{"APIMGMT"},
			Destination: &apimClient.ServiceName,
		},
	}
	// BeforeFunction is executed prior to execution of any subcommand
	BeforeFunction = func(c *ucli.Context) error {
		apimClient.Ctx = context.Background()
		apimClient.Authenticate()
		return nil
	}
)
